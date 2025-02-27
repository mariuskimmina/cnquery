// Copyright (c) Mondoo, Inc.
// SPDX-License-Identifier: BUSL-1.1

package connection

import (
	"context"
	"net/http"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/google/go-github/v55/github"
	"github.com/rs/zerolog/log"
	"go.mondoo.com/cnquery/v9/providers-sdk/v1/inventory"
	"go.mondoo.com/cnquery/v9/providers-sdk/v1/vault"
	"golang.org/x/oauth2"
)

type GithubConnection struct {
	id     uint32
	Conf   *inventory.Config
	asset  *inventory.Asset
	client *github.Client
}

func NewGithubConnection(id uint32, asset *inventory.Asset, conf *inventory.Config) (*GithubConnection, error) {
	token := conf.Options["token"]

	// if no token was provided, lets read the env variable
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}

	// if a secret was provided, it always overrides the env variable since it has precedence
	if len(conf.Credentials) > 0 {
		for i := range conf.Credentials {
			cred := conf.Credentials[i]
			if cred.Type == vault.CredentialType_password {
				token = string(cred.Secret)
			} else {
				log.Warn().Str("credential-type", cred.Type.String()).Msg("unsupported credential type for GitHub provider")
			}
		}
	}

	if token == "" {
		return nil, errors.New("a valid GitHub token is required, pass --token '<yourtoken>' or set GITHUB_TOKEN environment variable")
	}

	var oauthClient *http.Client
	if token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		ctx := context.Background()
		oauthClient = oauth2.NewClient(ctx, ts)
	}

	client := github.NewClient(oauthClient)
	// perform a quick call to verify the token's validity.
	_, resp, err := client.Zen(context.Background())
	if err != nil {
		if resp != nil && resp.StatusCode == 401 {
			return nil, errors.New("invalid GitHub token provided. check the value passed with the --token flag or the GITHUB_TOKEN environment variable")
		}
		return nil, err
	}
	conn := &GithubConnection{
		Conf:   conf,
		id:     id,
		asset:  asset,
		client: client,
	}

	return conn, nil
}

func (c *GithubConnection) Name() string {
	return "github"
}

func (c *GithubConnection) ID() uint32 {
	return c.id
}

func (c *GithubConnection) Asset() *inventory.Asset {
	return c.asset
}

func (c *GithubConnection) Client() *github.Client {
	return c.client
}
