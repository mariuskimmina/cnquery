// Copyright (c) Mondoo, Inc.
// SPDX-License-Identifier: BUSL-1.1

package config

import (
	"go.mondoo.com/cnquery/v9/providers-sdk/v1/plugin"
	"go.mondoo.com/cnquery/v9/providers/gitlab/provider"
)

var Config = plugin.Provider{
	Name:    "gitlab",
	ID:      "go.mondoo.com/cnquery/v9/providers/gitlab",
	Version: "9.0.8",
	ConnectionTypes: []string{
		provider.ConnectionType,
		provider.GitlabGroupConnection,
		provider.GitlabProjectConnection,
	},
	Connectors: []plugin.Connector{
		{
			Name:  "gitlab",
			Use:   "gitlab",
			Short: "GitLab",
			Discovery: []string{
				provider.DiscoveryAuto,
				provider.DiscoveryGroup,
				provider.DiscoveryProject,
				provider.DiscoveryTerraform,
			},
			Flags: []plugin.Flag{
				{
					Long:    "token",
					Type:    plugin.FlagType_String,
					Default: "",
					Desc:    "Provide GitLab personal access token.",
				},
				{
					Long:    "group",
					Type:    plugin.FlagType_String,
					Option:  plugin.FlagOption_Required,
					Default: "",
					Desc:    "a GitLab group to scan",
				},
				{
					Long:    "project",
					Type:    plugin.FlagType_String,
					Default: "",
					Desc:    "a GitLab project to scan",
				},
				{
					Long:    "url",
					Type:    plugin.FlagType_String,
					Default: "",
					Desc:    "custom GitLab base url",
				},
			},
		},
	},
}
