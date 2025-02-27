// Copyright (c) Mondoo, Inc.
// SPDX-License-Identifier: BUSL-1.1

package resources

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
	"go.mondoo.com/cnquery/v9/providers-sdk/v1/inventory"
	"go.mondoo.com/cnquery/v9/providers-sdk/v1/vault"
	"go.mondoo.com/cnquery/v9/providers/aws/connection"
	"go.mondoo.com/cnquery/v9/providers/aws/connection/awsec2ebsconn"
	awsec2ebstypes "go.mondoo.com/cnquery/v9/providers/aws/connection/awsec2ebsconn/types"
	"go.mondoo.com/cnquery/v9/providers/os/id/awsec2"
	"go.mondoo.com/cnquery/v9/providers/os/id/containerid"
	"go.mondoo.com/cnquery/v9/providers/os/id/ids"
)

type mqlObject struct {
	name      string
	labels    map[string]string
	awsObject awsObject
}

type awsObject struct {
	account    string
	region     string
	id         string
	service    string
	objectType string
	arn        string
}

func MondooObjectID(awsObject awsObject) string {
	accountId := trimAwsAccountIdToJustId(awsObject.account)
	return "//platformid.api.mondoo.app/runtime/aws/" + awsObject.service + "/v1/accounts/" + accountId + "/regions/" + awsObject.region + "/" + awsObject.objectType + "/" + awsObject.id
}

func MqlObjectToAsset(account string, mqlObject mqlObject, conn *connection.AwsConnection) *inventory.Asset {
	if name := mqlObject.labels["Name"]; name != "" {
		mqlObject.name = name
	}
	if mqlObject.name == "" {
		mqlObject.name = mqlObject.awsObject.id
	}
	if err := validate(mqlObject); err != nil {
		log.Error().Err(err).Msg("missing values in mql object to asset translation")
		return nil
	}
	platformName := getPlatformName(mqlObject.awsObject)
	if platformName == "" {
		log.Error().Err(errors.New("could not fetch platform info for object")).Msg("missing runtime info")
		return nil
	}
	platformid := MondooObjectID(mqlObject.awsObject)
	t := conn.Conf
	t.PlatformId = platformid
	return &inventory.Asset{
		PlatformIds: []string{platformid, mqlObject.awsObject.arn},
		Name:        mqlObject.name,
		Platform:    connection.GetPlatformForObject(platformName),
		Labels:      mqlObject.labels,
		Connections: []*inventory.Config{cloneInventoryConf(conn.Conf)},
		Options:     conn.ConnectionOptions(),
	}
}

func cloneInventoryConf(invConf *inventory.Config) *inventory.Config {
	invConfClone := invConf.Clone()
	// We do not want to run discovery again for the already discovered assets
	invConfClone.Discover = &inventory.Discovery{}
	return invConfClone
}

func validate(m mqlObject) error {
	if m.name == "" {
		return errors.New("name required for mql aws object to asset translation")
	}
	if m.awsObject.id == "" {
		return errors.New("id required for mql aws object to asset translation")
	}
	if m.awsObject.region == "" {
		return errors.New("region required for mql aws object to asset translation")
	}
	if m.awsObject.account == "" {
		return errors.New("account required for mql aws object to asset translation")
	}
	if m.awsObject.arn == "" {
		return errors.New("arn required for mql aws object to asset translation")
	}
	return nil
}

func getPlatformName(awsObject awsObject) string {
	switch awsObject.service {
	case "s3":
		if awsObject.objectType == "bucket" {
			return "aws-s3-bucket"
		}
	case "cloudtrail":
		if awsObject.objectType == "trail" {
			return "aws-cloudtrail-trail"
		}
	case "rds":
		if awsObject.objectType == "dbinstance" {
			return "aws-rds-dbinstance"
		}
	case "dynamodb":
		if awsObject.objectType == "table" {
			return "aws-dynamodb-table"
		}
	case "redshift":
		if awsObject.objectType == "cluster" {
			return "aws-redshift-cluster"
		}
	case "vpc":
		if awsObject.objectType == "vpc" {
			return "aws-vpc"
		}
	case "ec2":
		switch awsObject.objectType {
		case "securitygroup":
			return "aws-security-group"
		case "volume":
			return "aws-ec2-volume"
		case "snapshot":
			return "aws-ec2-snapshot"
		case "instance":
			return "aws-ec2-instance"
		}
	case "iam":
		switch awsObject.objectType {
		case "user":
			return "aws-iam-user"

		case "group":
			return "aws-iam-group"
		}
	case "cloudwatch":
		if awsObject.objectType == "loggroup" {
			return "aws-cloudwatch-loggroup"
		}
	case "lambda":
		if awsObject.objectType == "function" {
			return "aws-lambda-function"
		}
	case "ecs":
		if awsObject.objectType == "container" {
			return "aws-ecs-container"
		}
		if awsObject.objectType == "instance" {
			return "aws-ecs-instance"
		}
	case "efs":
		if awsObject.objectType == "filesystem" {
			return "aws-efs-filesystem"
		}
	case "gateway":
		if awsObject.objectType == "restapi" {
			return "aws-gateway-restapi"
		}
	case "elb":
		if awsObject.objectType == "loadbalancer" {
			return "aws-elb-loadbalancer"
		}
	case "es":
		if awsObject.objectType == "domain" {
			return "aws-es-domain"
		}
	case "kms":
		if awsObject.objectType == "key" {
			return "aws-kms-key"
		}
	case "sagemaker":
		if awsObject.objectType == "notebookinstance" {
			return "aws-sagemaker-notebookinstance"
		}
	case "ssm":
		if awsObject.objectType == "instance" {
			return "aws-ssm-instance"
		}
	case "ecr":
		if awsObject.objectType == "image" {
			return "aws-ecr-image"
		}
	}
	return ""
}

func accountAsset(conn *connection.AwsConnection, awsAccount *mqlAwsAccount) *inventory.Asset {
	var alias string
	aliases := awsAccount.GetAliases()
	if len(aliases.Data) > 0 {
		alias = aliases.Data[0].(string)
	}
	accountId := trimAwsAccountIdToJustId(awsAccount.Id.Data)
	name := AssembleIntegrationName(alias, accountId)

	id := "//platformid.api.mondoo.app/runtime/aws/accounts/" + accountId

	return &inventory.Asset{
		PlatformIds: []string{id},
		Name:        name,
		Platform:    connection.GetPlatformForObject(""),
		Connections: []*inventory.Config{conn.Conf},
		Options:     conn.ConnectionOptions(),
	}
}

func trimAwsAccountIdToJustId(id string) string {
	return strings.TrimPrefix(id, "aws.account/")
}

func AssembleIntegrationName(alias string, id string) string {
	accountId := trimAwsAccountIdToJustId(id)
	if alias == "" {
		return fmt.Sprintf("AWS Account %s", accountId)
	}
	return fmt.Sprintf("AWS Account %s (%s)", alias, accountId)
}

func addConnectionInfoToEc2Asset(instance *mqlAwsEc2Instance, accountId string, conn *connection.AwsConnection) *inventory.Asset {
	asset := &inventory.Asset{}
	asset.PlatformIds = []string{awsec2.MondooInstanceID(accountId, instance.Region.Data, instance.InstanceId.Data)}
	asset.IdDetector = []string{"aws-ec2"}
	asset.Platform = &inventory.Platform{
		Kind:    "virtual_machine",
		Runtime: "aws_ec2",
	}
	asset.State = mapEc2InstanceStateCode(instance.State.Data)
	asset.Labels = mapStringInterfaceToStringString(instance.Tags.Data)
	name := instance.InstanceId.Data
	if labelName := asset.Labels["Name"]; name != "" {
		name = labelName
	}
	asset.Name = name
	asset.Options = conn.ConnectionOptions()
	// if there is a public ip & it is running, we assume ssh is an option
	if instance.PublicIp.Data != "" && instance.State.Data == string(types.InstanceStateNameRunning) {
		imageName := ""
		if instance.GetImage().Data != nil {
			imageName = instance.GetImage().Data.Name.Data
		}
		probableUsername := getProbableUsernameFromImageName(imageName)
		asset.Connections = []*inventory.Config{{
			Type:     "ssh",
			Host:     instance.PublicIp.Data,
			Insecure: true,
			Runtime:  "ssh",
			Credentials: []*vault.Credential{
				{
					Type: vault.CredentialType_aws_ec2_instance_connect,
					User: probableUsername,
				},
			},
			Options: map[string]string{
				"region":   instance.Region.Data,
				"profile":  conn.Profile(),
				"instance": instance.InstanceId.Data,
			},
		}}
		if len(instance.GetSsm().Data.(map[string]interface{})["InstanceInformationList"].([]interface{})) > 0 {
			if instance.GetSsm().Data.(map[string]interface{})["InstanceInformationList"].([]interface{})[0].(map[string]interface{})["PingStatus"] == "Online" {
				asset.Connections[0].Credentials = append(asset.Connections[0].Credentials, &vault.Credential{
					User: probableUsername,
					Type: vault.CredentialType_aws_ec2_ssm_session,
				})
			}
		}
	} else {
		log.Warn().Str("asset", asset.Name).Msg("no public ip address found")
		asset = MqlObjectToAsset(accountId,
			mqlObject{
				name: name, labels: mapStringInterfaceToStringString(instance.Tags.Data),
				awsObject: awsObject{
					account: accountId, region: instance.Region.Data, arn: instance.Arn.Data,
					id: instance.InstanceId.Data, service: "ec2", objectType: "instance",
				},
			}, conn)
	}
	return asset
}

func addSSMConnectionInfoToEc2Asset(instance *mqlAwsEc2Instance, accountId string, conn *connection.AwsConnection) *inventory.Asset {
	asset := &inventory.Asset{}
	asset.PlatformIds = []string{awsec2.MondooInstanceID(accountId, instance.Region.Data, instance.InstanceId.Data)}
	asset.IdDetector = []string{"aws-ec2"}
	asset.Platform = &inventory.Platform{
		Kind:    "virtual_machine",
		Runtime: "aws_ec2",
	}
	ssm := ""
	if s := instance.GetSsm().Data.(map[string]interface{})["InstanceInformationList"]; s != nil {
		if len(s.([]interface{})) > 0 {
			ssm = s.([]interface{})[0].(map[string]interface{})["PingStatus"].(string)
		}
	}
	asset.State = mapSmmManagedPingStateCode(ssm)
	asset.Options = conn.ConnectionOptions()
	asset.Labels = mapStringInterfaceToStringString(instance.Tags.Data)
	name := instance.InstanceId.Data
	if lname := asset.Labels["Name"]; name != "" {
		name = lname
	}
	asset.Name = name
	imageName := ""
	if instance.GetImage().Data != nil {
		imageName = instance.GetImage().Data.Name.Data
	}
	creds := []*vault.Credential{
		{
			User: getProbableUsernameFromImageName(imageName),
			Type: vault.CredentialType_aws_ec2_ssm_session,
		},
	}
	host := instance.InstanceId.Data
	if instance.PublicIp.Data != "" {
		host = instance.PublicIp.Data
	}
	if ssm == string(ssmtypes.PingStatusOnline) {
		asset.Connections = []*inventory.Config{{
			Host:        host,
			Insecure:    true,
			Runtime:     "aws_ec2",
			Credentials: creds,
			Options: map[string]string{
				"region":   instance.Region.Data,
				"profile":  conn.Profile(),
				"instance": instance.InstanceId.Data,
			},
		}}
	} else {
		asset = MqlObjectToAsset(accountId,
			mqlObject{
				name: name, labels: mapStringInterfaceToStringString(instance.Tags.Data),
				awsObject: awsObject{
					account: accountId, region: instance.Region.Data, arn: instance.Arn.Data,
					id: instance.InstanceId.Data, service: "ec2", objectType: "instance",
				},
			}, conn)
	}
	return asset
}

func mapEc2InstanceStateCode(state string) inventory.State {
	switch state {
	case string(types.InstanceStateNameRunning):
		return inventory.State_STATE_RUNNING
	case string(types.InstanceStateNamePending):
		return inventory.State_STATE_PENDING
	case string(types.InstanceStateNameShuttingDown): // 32 is shutting down, which is the step before terminated, assume terminated if we get shutting down
		return inventory.State_STATE_TERMINATED
	case string(types.InstanceStateNameStopping):
		return inventory.State_STATE_STOPPING
	case string(types.InstanceStateNameStopped):
		return inventory.State_STATE_STOPPED
	case string(types.InstanceStateNameTerminated):
		return inventory.State_STATE_TERMINATED
	default:
		log.Warn().Str("state", string(state)).Msg("unknown ec2 state")
		return inventory.State_STATE_UNKNOWN
	}
}

func getProbableUsernameFromImageName(name string) string {
	if strings.Contains(name, "centos") {
		return "centos"
	}
	if strings.Contains(name, "ubuntu") {
		return "ubuntu"
	}
	return "ec2-user"
}

func addConnectionInfoToSSMAsset(instance *mqlAwsSsmInstance, accountId string, conn *connection.AwsConnection) *inventory.Asset {
	asset := &inventory.Asset{}
	asset.Labels = mapStringInterfaceToStringString(instance.Tags.Data)
	name := instance.InstanceId.Data
	if labelName := asset.Labels["Name"]; name != "" {
		name = labelName
	}
	asset.Name = name
	creds := []*vault.Credential{
		{
			User: getProbableUsernameFromSSMPlatformName(strings.ToLower(instance.PlatformName.Data)),
		},
	}

	host := instance.InstanceId.Data
	if instance.IpAddress.Data != "" {
		host = instance.IpAddress.Data
	}
	asset.Options = conn.ConnectionOptions()
	asset.PlatformIds = []string{awsec2.MondooInstanceID(accountId, instance.Region.Data, instance.InstanceId.Data)}
	asset.Platform = &inventory.Platform{
		Kind:    "virtual_machine",
		Runtime: "ssm_managed",
	}
	asset.State = mapSmmManagedPingStateCode(instance.PingStatus.Data)

	if strings.HasPrefix(instance.InstanceId.Data, "i-") && instance.PingStatus.Data == string(ssmtypes.PingStatusOnline) {
		creds[0].Type = vault.CredentialType_aws_ec2_ssm_session // this will only work for ec2 instances
		asset.Connections = []*inventory.Config{{
			Host:        host,
			Insecure:    true,
			Runtime:     "aws_ec2",
			Credentials: creds,
			Options: map[string]string{
				"region":   instance.Region.Data,
				"profile":  conn.Profile(),
				"instance": instance.InstanceId.Data,
			},
		}}
	} else {
		log.Warn().Str("asset", asset.Name).Str("id", instance.InstanceId.Data).Msg("cannot use ssm session credentials for connection")
		asset = MqlObjectToAsset(accountId,
			mqlObject{
				name: name, labels: mapStringInterfaceToStringString(instance.Tags.Data),
				awsObject: awsObject{
					account: accountId, region: instance.Region.Data, arn: instance.Arn.Data,
					id: instance.InstanceId.Data, service: "ssm", objectType: "instance",
				},
			}, conn)
	}
	return asset
}

func getProbableUsernameFromSSMPlatformName(name string) string {
	if strings.HasPrefix(name, "centos") {
		return "centos"
	}
	if strings.HasPrefix(name, "ubuntu") {
		return "ubuntu"
	}
	return "ec2-user"
}

func mapSmmManagedPingStateCode(pingStatus string) inventory.State {
	switch pingStatus {
	case string(ssmtypes.PingStatusOnline):
		return inventory.State_STATE_RUNNING
	case string(ssmtypes.PingStatusConnectionLost):
		return inventory.State_STATE_PENDING
	case string(ssmtypes.PingStatusInactive):
		return inventory.State_STATE_STOPPED
	default:
		return inventory.State_STATE_UNKNOWN
	}
}

func MondooImageRegistryID(id string) string {
	return "//platformid.api.mondoo.app/runtime/docker/registry/" + id
}

func addConnectionInfoToEcrAsset(image *mqlAwsEcrImage, conn *connection.AwsConnection) *inventory.Asset {
	a := &inventory.Asset{}
	a.PlatformIds = []string{containerid.MondooContainerImageID(image.Digest.Data)}
	a.Platform = &inventory.Platform{
		Kind:    "container_image",
		Runtime: "aws_ecr",
	}
	a.Options = conn.ConnectionOptions()
	a.Name = ecrImageName(image.RepoName.Data, image.Digest.Data)
	a.State = inventory.State_STATE_ONLINE
	imageTags := []string{}
	for i := range image.Tags.Data {
		tag := image.Tags.Data[i].(string)
		imageTags = append(imageTags, tag)
		a.Connections = append(a.Connections, &inventory.Config{
			Type: "registry-image",
			Host: image.Uri.Data + ":" + tag,
			Options: map[string]string{
				"region":  image.Region.Data,
				"profile": conn.Profile(),
			},
		})

	}
	a.Labels = make(map[string]string)
	// store digest
	a.Labels[fmt.Sprintf("ecr.%s.amazonaws.com/digest", image.Region.Data)] = image.Digest.Data
	a.Labels[fmt.Sprintf("ecr.%s.amazonaws.com/tags", image.Region.Data)] = strings.Join(imageTags, ",")

	// store repo digest
	repoDigests := []string{image.Uri.Data + "@" + image.Digest.Data}
	a.Labels[fmt.Sprintf("ecr.%s.amazonaws.com/repo-digests", image.Region.Data)] = strings.Join(repoDigests, ",")

	return a
}

func ecrImageName(repoName string, digest string) string {
	return repoName + "@" + digest
}

func mapContainerInstanceState(status *string) inventory.State {
	if status == nil {
		return inventory.State_STATE_UNKNOWN
	}
	switch *status {
	case "REGISTERING":
		return inventory.State_STATE_PENDING
	case "REGISTRATION_FAILED":
		return inventory.State_STATE_ERROR
	case "ACTIVE":
		return inventory.State_STATE_ONLINE
	case "INACTIVE":
		return inventory.State_STATE_OFFLINE
	case "DEREGISTERING":
		return inventory.State_STATE_STOPPING
	case "DRAINING":
		return inventory.State_STATE_STOPPING
	default:
		return inventory.State_STATE_UNKNOWN
	}
}

func addConnectionInfoToECSContainerAsset(container *mqlAwsEcsContainer, accountId string, conn *connection.AwsConnection) *inventory.Asset {
	a := &inventory.Asset{}

	runtimeId := container.RuntimeId.Data
	if runtimeId == "" {
		return nil
	}
	state := container.Status.Data
	containerArn := container.Arn.Data
	taskArn := container.TaskArn.Data
	publicIp := container.GetPublicIp().Data
	region := container.Region.Data

	a.Name = container.Name.Data
	a.Options = conn.ConnectionOptions()
	a.PlatformIds = []string{containerid.MondooContainerID(runtimeId), MondooECSContainerID(containerArn)}
	a.Platform = &inventory.Platform{
		Kind:    "container",
		Runtime: "aws_ecs",
	}
	a.State = mapContainerState(state)
	taskId := ""
	if arn.IsARN(taskArn) {
		if parsed, err := arn.Parse(taskArn); err == nil {
			if taskIds := strings.Split(parsed.Resource, "/"); len(taskIds) > 1 {
				taskId = taskIds[len(taskIds)-1]
			}
		}
	}

	if publicIp != "" {
		a.Connections = []*inventory.Config{{
			Host: publicIp,
			Options: map[string]string{
				"region":         region,
				"container_name": container.Name.Data,
				"task_id":        taskId,
			},
		}}
	} else {
		log.Warn().Str("asset", a.Name).Msg("no public ip address found")
		a = MqlObjectToAsset(accountId,
			mqlObject{
				name: container.Name.Data, labels: make(map[string]string),
				awsObject: awsObject{
					account: accountId, region: container.Region.Data, arn: container.Arn.Data,
					id: container.Arn.Data, service: "ecs", objectType: "container",
				},
			}, conn)
	}

	return a
}

func addConnectionInfoToECSContainerInstanceAsset(inst *mqlAwsEcsInstance, accountId string, conn *connection.AwsConnection) *inventory.Asset {
	m := mqlObject{
		name: inst.Id.Data, labels: map[string]string{},
		awsObject: awsObject{
			account: accountId, region: inst.Region.Data, arn: inst.Arn.Data,
			id: inst.Id.Data, service: "ecs", objectType: "instance",
		},
	}
	a := MqlObjectToAsset(accountId, m, conn)
	a.Connections = append(a.Connections, &inventory.Config{
		Type: "ssh", // fallback to ssh
		Options: map[string]string{
			"region": inst.Region.Data,
		},
	})
	return a
}

func mapContainerState(state string) inventory.State {
	switch strings.ToLower(state) {
	case "running":
		return inventory.State_STATE_RUNNING
	case "created":
		return inventory.State_STATE_PENDING
	case "paused":
		return inventory.State_STATE_STOPPED
	case "exited":
		return inventory.State_STATE_TERMINATED
	case "restarting":
		return inventory.State_STATE_PENDING
	case "dead":
		return inventory.State_STATE_ERROR
	default:
		log.Warn().Str("state", state).Msg("unknown container state")
		return inventory.State_STATE_UNKNOWN
	}
}

func MondooECSContainerID(containerArn string) string {
	var account, region, id string
	if arn.IsARN(containerArn) {
		if p, err := arn.Parse(containerArn); err == nil {
			account = p.AccountID
			region = p.Region
			id = p.Resource
		}
	}
	return "//platformid.api.mondoo.app/runtime/aws/ecs/v1/accounts/" + account + "/regions/" + region + "/" + id
}

func SSMConnectAsset(args []string, opts map[string]string) *inventory.Asset {
	var user, id string
	if len(args) == 3 {
		if args[0] == "ec2" && args[1] == "ssm" {
			if targets := strings.Split(args[2], "@"); len(targets) == 2 {
				user = targets[0]
				id = targets[1]
			}
		}
	}
	asset := &inventory.Asset{}
	opts["instance"] = id
	asset.IdDetector = []string{ids.IdDetector_Hostname, ids.IdDetector_CloudDetect, ids.IdDetector_SshHostkey}
	asset.Connections = []*inventory.Config{{
		Type:     "ssh",
		Host:     id,
		Insecure: true,
		Runtime:  "ssh",
		Credentials: []*vault.Credential{
			{
				Type: vault.CredentialType_aws_ec2_ssm_session,
				User: user,
			},
		},
		Options: opts,
	}}
	return asset
}

func InstanceConnectAsset(args []string, opts map[string]string) *inventory.Asset {
	var user, id string
	if len(args) == 3 {
		if args[0] == "ec2" && args[1] == "instance-connect" {
			if targets := strings.Split(args[2], "@"); len(targets) == 2 {
				user = targets[0]
				id = targets[1]
			}
		}
	}
	asset := &inventory.Asset{}
	asset.IdDetector = []string{ids.IdDetector_Hostname, ids.IdDetector_CloudDetect, ids.IdDetector_SshHostkey}
	opts["instance"] = id
	asset.Connections = []*inventory.Config{{
		Type:     "ssh",
		Host:     id,
		Insecure: true,
		Runtime:  "ssh",
		Credentials: []*vault.Credential{
			{
				Type: vault.CredentialType_aws_ec2_instance_connect,
				User: user,
			},
		},
		Options: opts,
	}}
	return asset
}

func EbsConnectAsset(args []string, opts map[string]string) *inventory.Asset {
	var target, targetType string
	if len(args) >= 3 {
		if args[0] == "ec2" && args[1] == "ebs" {
			// parse for target type: instance, volume, snapshot
			switch args[2] {
			case awsec2ebstypes.EBSTargetVolume:
				target = args[3]
				targetType = awsec2ebstypes.EBSTargetVolume
			case awsec2ebstypes.EBSTargetSnapshot:
				target = args[3]
				targetType = awsec2ebstypes.EBSTargetSnapshot
			default:
				// in the case of an instance target, this is the instance id
				target = args[2]
				targetType = awsec2ebstypes.EBSTargetInstance
			}
		}
	}
	asset := &inventory.Asset{}
	opts["type"] = targetType
	opts["id"] = target
	asset.Name = target
	asset.IdDetector = []string{ids.IdDetector_Hostname} // do not use cloud detect or host key here
	asset.Connections = []*inventory.Config{{
		Type:     string(awsec2ebsconn.EBSConnectionType),
		Host:     target,
		Insecure: true,
		Runtime:  "aws-ebs",
		Options:  opts,
	}}
	return asset
}
