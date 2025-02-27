// Copyright (c) Mondoo, Inc.
// SPDX-License-Identifier: BUSL-1.1

package config

import (
	"go.mondoo.com/cnquery/v9/providers-sdk/v1/plugin"
	"go.mondoo.com/cnquery/v9/providers/slack/provider"
)

var Config = plugin.Provider{
	Name:            "slack",
	ID:              "go.mondoo.com/cnquery/v9/providers/slack",
	Version:         "9.0.6",
	ConnectionTypes: []string{provider.ConnectionType},
	Connectors: []plugin.Connector{
		{
			Name:      "slack",
			Use:       "slack",
			Short:     "slack team",
			MinArgs:   0,
			MaxArgs:   0,
			Discovery: []string{},
			Flags: []plugin.Flag{
				{
					Long:    "token",
					Type:    plugin.FlagType_String,
					Default: "",
					Desc:    "Slack API token",
				},
			},
		},
	},
}
