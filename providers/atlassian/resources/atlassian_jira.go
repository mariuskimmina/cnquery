package resources

import (
	"context"

	"github.com/rs/zerolog/log"
	"go.mondoo.com/cnquery/v9/llx"
	"go.mondoo.com/cnquery/v9/providers/atlassian/connection"
)

func (a *mqlAtlassianJira) id() (string, error) {
	return "wip", nil
}

func (a *mqlAtlassianJira) users() ([]interface{}, error) {
	conn := a.MqlRuntime.Connection.(*connection.AtlassianConnection)
	jira := conn.Jira()
	users, response, err := jira.User.Search.Do(context.Background(), "", " ", 0, 1000)
	if err != nil {
		log.Fatal().Err(err)
	}
	if response.Status != "200 OK" {
		log.Fatal().Msgf("Received response: %s\n", response.Status)
	}

	res := []interface{}{}
	for _, user := range users {
		mqlAtlassianJiraUser, err := CreateResource(a.MqlRuntime, "atlassian.jira.user",
			map[string]*llx.RawData{
				"id":      llx.StringData(user.AccountID),
				"name":    llx.StringData(user.DisplayName),
				"type":    llx.StringData(user.AccountType),
				"picture": llx.StringData(user.AvatarUrls.One6X16),
			})
		if err != nil {
			log.Fatal().Err(err)
		}
		res = append(res, mqlAtlassianJiraUser)
	}
	return res, nil
}

func (a *mqlAtlassianJiraUser) groups() ([]interface{}, error) {
	conn := a.MqlRuntime.Connection.(*connection.AtlassianConnection)
	jira := conn.Jira()
	groups, response, err := jira.Group.Bulk(context.Background(), nil, 0, 1000)
	if err != nil {
		log.Fatal().Err(err)
	}
	if response.Status != "200 OK" {
		log.Fatal().Msgf("Received response: %s\n", response.Status)
	}

	res := []interface{}{}
	for _, group := range groups.Values {
		mqlAtlassianJiraUserGroup, err := CreateResource(a.MqlRuntime, "atlassian.jira.user.group",
			map[string]*llx.RawData{
				"id": llx.StringData(group.GroupID),
			})
		if err != nil {
			log.Fatal().Err(err)
		}
		res = append(res, mqlAtlassianJiraUserGroup)
	}
	return res, nil
}

func (a *mqlAtlassianJira) projects() ([]interface{}, error) {
	conn := a.MqlRuntime.Connection.(*connection.AtlassianConnection)
	jira := conn.Jira()
	projects, response, err := jira.Project.Search(context.Background(), nil, 0, 1000)
	if err != nil {
		log.Fatal().Err(err)
	}
	if response.Status != "200 OK" {
		log.Fatal().Msgf("Received response: %s\n", response.Status)
	}

	res := []interface{}{}
	for _, project := range projects.Values {
		mqlAtlassianJiraProject, err := CreateResource(a.MqlRuntime, "atlassian.jira.project",
			map[string]*llx.RawData{
				"id":       llx.StringData(project.ID),
				"name":     llx.StringData(project.Name),
				"uuid":     llx.StringData(project.UUID),
				"key":      llx.StringData(project.Key),
				"url":      llx.StringData(project.URL),
				"email":    llx.StringData(project.Email),
				"private":  llx.BoolData(project.IsPrivate),
				"deleted":  llx.BoolData(project.Deleted),
				"archived": llx.BoolData(project.Archived),
			})
		if err != nil {
			log.Fatal().Err(err)
		}
		res = append(res, mqlAtlassianJiraProject)
	}
	return res, nil
}

func (a *mqlAtlassianJiraUser) id() (string, error) {
	return a.Id.Data, nil
}

func (a *mqlAtlassianJiraUserGroup) id() (string, error) {
	return a.Id.Data, nil
}

func (a *mqlAtlassianJiraProject) id() (string, error) {
	return a.Id.Data, nil
}
