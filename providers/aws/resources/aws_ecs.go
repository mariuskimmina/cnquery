// Copyright (c) Mondoo, Inc.
// SPDX-License-Identifier: BUSL-1.1

package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsservice "github.com/aws/aws-sdk-go-v2/service/ecs"
	ecstypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
	"go.mondoo.com/cnquery/v9/llx"
	"go.mondoo.com/cnquery/v9/providers-sdk/v1/plugin"
	"go.mondoo.com/cnquery/v9/providers-sdk/v1/util/convert"
	"go.mondoo.com/cnquery/v9/providers-sdk/v1/util/jobpool"
	"go.mondoo.com/cnquery/v9/providers/aws/connection"

	"go.mondoo.com/cnquery/v9/types"
	"go.mondoo.com/cnquery/v9/utils/stringx"
)

func (a *mqlAwsEcs) id() (string, error) {
	return "aws.ecs", nil
}

func (a *mqlAwsEcs) containers() ([]interface{}, error) {
	obj, err := CreateResource(a.MqlRuntime, "aws.ecs", map[string]*llx.RawData{})
	if err != nil {
		return nil, err
	}
	ecs := obj.(*mqlAwsEcs)

	clusters, err := ecs.clusters()
	if err != nil {
		return nil, err
	}
	containers := []interface{}{}

	for i := range clusters {
		tasks, err := clusters[i].(*mqlAwsEcsCluster).tasks()
		if err != nil {
			return nil, err
		}
		for i := range tasks {
			c := tasks[i].(*mqlAwsEcsTask).Containers
			containers = append(containers, c.Data...)
		}
	}
	return containers, nil
}

func (a *mqlAwsEcs) containerInstances() ([]interface{}, error) {
	obj, err := CreateResource(a.MqlRuntime, "aws.ecs", map[string]*llx.RawData{})
	if err != nil {
		return nil, err
	}
	ecs := obj.(*mqlAwsEcs)

	clusters, err := ecs.clusters()
	if err != nil {
		return nil, err
	}
	containerInstances := []interface{}{}

	for i := range clusters {
		ci, err := clusters[i].(*mqlAwsEcsCluster).containerInstances()
		if err != nil {
			return nil, err
		}
		containerInstances = append(containerInstances, ci...)

	}
	return containerInstances, nil
}

func (a *mqlAwsEcsInstance) ec2Instance() (*mqlAwsEc2Instance, error) {
	return a.Ec2Instance.Data, nil
}

func (a *mqlAwsEcs) clusters() ([]interface{}, error) {
	conn := a.MqlRuntime.Connection.(*connection.AwsConnection)
	res := []interface{}{}
	poolOfJobs := jobpool.CreatePool(a.getECSClusters(conn), 5)
	poolOfJobs.Run()

	// check for errors
	if poolOfJobs.HasErrors() {
		return nil, poolOfJobs.GetErrors()
	}
	// get all the results
	for i := range poolOfJobs.Jobs {
		res = append(res, poolOfJobs.Jobs[i].Result.([]interface{})...)
	}

	return res, nil
}

func (a *mqlAwsEcs) getECSClusters(conn *connection.AwsConnection) []*jobpool.Job {
	tasks := make([]*jobpool.Job, 0)
	var err error
	regions, err := conn.Regions()
	if err != nil {
		return []*jobpool.Job{{Err: err}} // return the error
	}
	log.Debug().Msgf("regions being called for ecs clusters list are: %v", regions)
	for ri := range regions {
		region := regions[ri]
		f := func() (jobpool.JobResult, error) {
			svc := conn.Ecs(region)
			ctx := context.Background()
			res := []interface{}{}

			nextToken := aws.String("no_token_to_start_with")
			params := &ecsservice.ListClustersInput{}
			for nextToken != nil {
				resp, err := svc.ListClusters(ctx, params)
				if err != nil {
					if Is400AccessDeniedError(err) {
						log.Warn().Str("region", region).Msg("error accessing region for AWS API")
						return res, nil
					}
					return nil, errors.Wrap(err, "could not gather ecs cluster information")
				}
				nextToken = resp.NextToken
				if resp.NextToken != nil {
					params.NextToken = nextToken
				}
				for _, cluster := range resp.ClusterArns {
					mqlCluster, err := NewResource(a.MqlRuntime, "aws.ecs.cluster",
						map[string]*llx.RawData{
							"arn": llx.StringData(cluster),
						})
					if err != nil {
						return nil, err
					}
					res = append(res, mqlCluster)
				}
			}
			return jobpool.JobResult(res), nil
		}
		tasks = append(tasks, jobpool.NewJob(f))
	}
	return tasks
}

func initAwsEcsCluster(runtime *plugin.Runtime, args map[string]*llx.RawData) (map[string]*llx.RawData, plugin.Resource, error) {
	if len(args) > 2 {
		return args, nil, nil
	}

	if len(args) == 0 {
		if ids := getAssetIdentifier(runtime); ids != nil {
			args["arn"] = llx.StringData(ids.arn)
		}
	}

	if args["arn"] == nil {
		return nil, nil, errors.New("arn required to fetch ecs cluster")
	}
	a := args["arn"].Value.(string)
	conn := runtime.Connection.(*connection.AwsConnection)

	region := ""
	if arn.IsARN(a) {
		if val, err := arn.Parse(a); err == nil {
			region = val.Region
		}
	}
	svc := conn.Ecs(region)
	ctx := context.Background()
	clusterDetails, err := svc.DescribeClusters(ctx, &ecs.DescribeClustersInput{Clusters: []string{a}})
	if err != nil {
		return nil, nil, err
	}
	if len(clusterDetails.Clusters) != 1 {
		return nil, nil, errors.Newf("only expected one cluster, got %d", len(clusterDetails.Clusters))
	}
	c := clusterDetails.Clusters[0]
	configuration, err := convert.JsonToDict(c.Configuration)
	if err != nil {
		return nil, nil, err
	}
	args["name"] = llx.StringData(convert.ToString(c.ClusterName))
	args["tags"] = llx.MapData(ecsTags(c.Tags), types.String)
	args["runningTasksCount"] = llx.IntData(int64(c.RunningTasksCount))
	args["pendingTasksCount"] = llx.IntData(int64(c.PendingTasksCount))
	args["registeredContainerInstancesCount"] = llx.IntData(int64(c.RegisteredContainerInstancesCount))
	args["configuration"] = llx.MapData(configuration, types.String)
	args["status"] = llx.StringData(convert.ToString(c.Status))
	return args, nil, nil
}

func ecsTags(t []ecstypes.Tag) map[string]interface{} {
	res := map[string]interface{}{}
	for i := range t {
		tag := t[i]
		if tag.Key != nil && tag.Value != nil {
			res[*tag.Key] = *tag.Value
		}
	}
	return res
}

func (a *mqlAwsEcsCluster) containerInstances() ([]interface{}, error) {
	conn := a.MqlRuntime.Connection.(*connection.AwsConnection)
	clustera := a.Arn.Data
	region := ""
	if arn.IsARN(clustera) {
		if val, err := arn.Parse(clustera); err == nil {
			region = val.Region
		}
	}
	svc := conn.Ecs(region)
	ctx := context.Background()
	res := []interface{}{}

	params := &ecsservice.ListContainerInstancesInput{Cluster: &clustera}
	containerInstances, err := svc.ListContainerInstances(ctx, params)
	if err != nil {
		log.Error().Err(err).Msg("cannot list container instances") // no fail
	} else if len(containerInstances.ContainerInstanceArns) > 0 {
		containerInstancesDetail, err := svc.DescribeContainerInstances(ctx, &ecsservice.DescribeContainerInstancesInput{Cluster: &clustera, ContainerInstances: containerInstances.ContainerInstanceArns})
		if err == nil {
			for _, ci := range containerInstancesDetail.ContainerInstances {
				// container instance assets
				args := map[string]*llx.RawData{
					"arn":              llx.StringData(convert.ToString(ci.ContainerInstanceArn)),
					"agentConnected":   llx.BoolData(ci.AgentConnected),
					"id":               llx.StringData(convert.ToString(ci.Ec2InstanceId)),
					"capacityProvider": llx.StringData(convert.ToString(ci.CapacityProviderName)),
					"region":           llx.StringData(region),
				}
				if strings.HasPrefix(convert.ToString(ci.Ec2InstanceId), "i-") {
					mqlInstanceResource, err := CreateResource(a.MqlRuntime, "aws.ec2.instance",
						map[string]*llx.RawData{
							"arn": llx.StringData(fmt.Sprintf(ec2InstanceArnPattern, region, conn.AccountId(), convert.ToString(ci.Ec2InstanceId))),
						})
					if err == nil && mqlInstanceResource != nil {
						mqlInstance := mqlInstanceResource.(*mqlAwsEc2Instance)
						args["ec2Instance"] = llx.ResourceData(mqlInstance, mqlInstance.MqlName())
					}
				} else {
					args["ec2Instance"] = llx.NilData
				}

				mqlEcsInstance, err := CreateResource(a.MqlRuntime, "aws.ecs.instance", args)
				if err != nil {
					return nil, err
				}
				res = append(res, mqlEcsInstance)
			}
		} else {
			log.Error().Err(err).Msg("could not gather ecs container instances")
		}
	}
	return res, nil
}

func (s *mqlAwsEcsInstance) id() (string, error) {
	return s.Arn.Data, nil
}

func (s *mqlAwsEcsCluster) id() (string, error) {
	return s.Arn.Data, nil
}

func (a *mqlAwsEcsCluster) tasks() ([]interface{}, error) {
	conn := a.MqlRuntime.Connection.(*connection.AwsConnection)
	clustera := a.Arn.Data
	name := a.Arn.Data

	region := ""
	if arn.IsARN(clustera) {
		if val, err := arn.Parse(clustera); err == nil {
			region = val.Region
		}
	}
	svc := conn.Ecs(region)
	ctx := context.Background()
	res := []interface{}{}

	nextToken := aws.String("no_token_to_start_with")
	params := &ecsservice.ListTasksInput{Cluster: &clustera}
	for nextToken != nil {
		resp, err := svc.ListTasks(ctx, params)
		if err != nil {
			return nil, errors.Wrap(err, "could not gather ecs tasks information")
		}
		nextToken = resp.NextToken
		if resp.NextToken != nil {
			params.NextToken = nextToken
		}
		for _, task := range resp.TaskArns {
			mqlTask, err := CreateResource(a.MqlRuntime, "aws.ecs.task",
				map[string]*llx.RawData{
					"arn":         llx.StringData(task),
					"clusterName": llx.StringData(name),
				})
			if err != nil {
				return nil, err
			}
			res = append(res, mqlTask)
		}
	}
	return res, nil
}

func (s *mqlAwsEcsTask) id() (string, error) {
	return s.Arn.Data, nil
}

func initAwsEcsTask(runtime *plugin.Runtime, args map[string]*llx.RawData) (map[string]*llx.RawData, plugin.Resource, error) {
	if len(args) > 2 {
		return args, nil, nil
	}

	if len(args) == 0 {
		if ids := getAssetIdentifier(runtime); ids != nil {
			args["arn"] = llx.StringData(ids.arn)
		}
	}

	if args["arn"] == nil {
		return nil, nil, errors.New("arn required to fetch ecs task")
	}
	a := args["arn"].Value.(string)

	conn := runtime.Connection.(*connection.AwsConnection)

	region := ""
	clusterName := ""
	if arn.IsARN(a) {
		if val, err := arn.Parse(a); err == nil {
			region = val.Region
			if res := strings.Split(val.Resource, "/"); len(res) == 3 {
				clusterName = res[1]
			}
		}
	}
	svc := conn.Ecs(region)
	ctx := context.Background()
	params := &ecs.DescribeTasksInput{Tasks: []string{a}}
	params.Cluster = &clusterName
	taskDetails, err := svc.DescribeTasks(ctx, params)
	if err != nil {
		return nil, nil, err
	}
	if len(taskDetails.Tasks) != 1 {
		return nil, nil, errors.Newf("only expected one task, got %d", len(taskDetails.Tasks))
	}
	t := taskDetails.Tasks[0]
	taskDefinitionArn := t.TaskDefinitionArn
	definition, err := svc.DescribeTaskDefinition(ctx, &ecs.DescribeTaskDefinitionInput{TaskDefinition: taskDefinitionArn})
	if err != nil {
		return nil, nil, err
	}
	containerLogDriverMap := make(map[string]string)
	containerCommandMap := make(map[string][]string)

	for i := range definition.TaskDefinition.ContainerDefinitions {
		cd := definition.TaskDefinition.ContainerDefinitions[i]
		if cd.Name != nil {
			containerCommandMap[*cd.Name] = cd.Command
			if cd.LogConfiguration != nil {
				containerLogDriverMap[*cd.Name] = string(cd.LogConfiguration.LogDriver)
			} else {
				containerLogDriverMap[*cd.Name] = "none"
			}
		}
	}

	args["connectivity"] = llx.StringData(string(t.Connectivity))
	args["lastStatus"] = llx.StringData(convert.ToString(t.LastStatus))
	args["platformFamily"] = llx.StringData(convert.ToString(t.PlatformFamily))
	args["platformVersion"] = llx.StringData(convert.ToString(t.PlatformVersion))
	args["tags"] = llx.MapData(ecsTags(t.Tags), types.String)

	containers := []interface{}{}
	// pf, _ := a.platformFamily()
	// pv, _ := a.PlatformVersion()

	for _, c := range t.Containers {
		cmds := []interface{}{}
		for i := range containerCommandMap[convert.ToString(c.Name)] {
			cmds = append(cmds, containerCommandMap[convert.ToString(c.Name)][i])
		}
		publicIp := getContainerIP(ctx, conn, t.Attachments, c, region)
		name := convert.ToString(c.Name)
		if publicIp != "" {
			name = name + "-" + publicIp
		}

		mqlContainer, err := CreateResource(runtime, "aws.ecs.container",
			map[string]*llx.RawData{
				"name": llx.StringData(name),
				// "platformFamily", pf,
				// "platformVersion", pv,
				"status":            llx.StringData(convert.ToString(c.LastStatus)),
				"publicIp":          llx.StringData(publicIp),
				"arn":               llx.StringData(convert.ToString(c.ContainerArn)),
				"logDriver":         llx.StringData(containerLogDriverMap[convert.ToString(c.Name)]),
				"image":             llx.StringData(convert.ToString(c.Image)),
				"clusterName":       llx.StringData(clusterName),
				"taskDefinitionArn": llx.StringData(convert.ToString(taskDefinitionArn)),
				"region":            llx.StringData(region),
				"command":           llx.ArrayData(cmds, types.Any),
				"taskArn":           llx.StringData(convert.ToString(t.TaskArn)),
				"runtimeId":         llx.StringData(convert.ToString(c.RuntimeId)),
				"containerName":     llx.StringData(convert.ToString(c.Name)),
			})
		if err != nil {
			return args, nil, err
		}
		containers = append(containers, mqlContainer)
	}
	args["containers"] = llx.ArrayData(containers, types.Resource("aws.ecs.container"))

	return args, nil, nil
}

func initAwsEcsContainer(runtime *plugin.Runtime, args map[string]*llx.RawData) (map[string]*llx.RawData, plugin.Resource, error) {
	if len(args) > 2 {
		return args, nil, nil
	}
	if len(args) == 0 {
		if ids := getAssetIdentifier(runtime); ids != nil {
			args["arn"] = llx.StringData(ids.arn)
		}
	}
	obj, err := CreateResource(runtime, "aws.ecs", map[string]*llx.RawData{})
	if err != nil {
		return nil, nil, err
	}
	ecs := obj.(*mqlAwsEcs)

	rawResources := ecs.GetContainers()

	if rawResources.Error != nil {
		return nil, nil, rawResources.Error
	}

	arnVal := args["arn"].Value.(string)
	for i := range rawResources.Data {
		container := rawResources.Data[i].(*mqlAwsEcsContainer)
		if container.Arn.Data == arnVal {
			return args, container, nil
		}
	}

	return nil, nil, errors.New("container does not exist")
}

func getContainerIP(ctx context.Context, conn *connection.AwsConnection, attachments []ecstypes.Attachment, c ecstypes.Container, region string) string {
	containerAttachmentIds := []string{}
	for _, ca := range c.NetworkInterfaces {
		containerAttachmentIds = append(containerAttachmentIds, *ca.AttachmentId)
	}
	var publicIp string
	for _, a := range attachments {
		if stringx.Contains(containerAttachmentIds, *a.Id) {
			for _, detail := range a.Details {
				if *detail.Name == "networkInterfaceId" {
					publicIp = getPublicIpForContainer(ctx, conn, *detail.Value, region)
				}
			}
		}
	}
	return publicIp
}

func getPublicIpForContainer(ctx context.Context, conn *connection.AwsConnection, nii string, region string) string {
	svc := conn.Ec2(region)
	ni, err := svc.DescribeNetworkInterfaces(ctx, &ec2.DescribeNetworkInterfacesInput{NetworkInterfaceIds: []string{nii}})
	if err == nil {
		if len(ni.NetworkInterfaces) == 1 {
			if ni.NetworkInterfaces[0].Association != nil {
				return *ni.NetworkInterfaces[0].Association.PublicIp
			}
		}
	}
	return ""
}

func (s *mqlAwsEcsContainer) id() (string, error) {
	return s.Arn.Data, nil
}
