package connection

import (
	"errors"
	"os"

	"github.com/ctreminiom/go-atlassian/admin"
	"github.com/ctreminiom/go-atlassian/confluence"
	"github.com/ctreminiom/go-atlassian/jira/v2"
	_ "github.com/ctreminiom/go-atlassian/jira/v3"
	"github.com/rs/zerolog/log"
	"go.mondoo.com/cnquery/v9/providers-sdk/v1/inventory"
)

type AtlassianConnection struct {
	id        uint32
	Conf      *inventory.Config
	asset     *inventory.Asset
	admin     *admin.Client
	jira      *v2.Client
	confluece *confluence.Client
	Host      string
	// Add custom connection fields here
}

func NewAtlassianConnection(id uint32, asset *inventory.Asset, conf *inventory.Config) (*AtlassianConnection, error) {
	apiKey := os.Getenv("ATLASSIAN_ADMIN_TOKEN")
	if apiKey == "" {
		return nil, errors.New("you need to provide atlassian admin token via ATLASSIAN_ADMIN_TOKEN env")
	}

	host := os.Getenv("ATLASSIAN_HOST")
	if host == "" {
		log.Warn().Msg("ATLASSIAN_HOST not set")
	}

	mail := os.Getenv("ATLASSIAN_USER")
	if mail == "" {
		log.Warn().Msg("ATLASSIAN_USER not set")
	}

	token := os.Getenv("ATLASSIAN_TOKEN")
	admin, err := admin.New(nil)
	if err != nil {
		log.Fatal().Err(err)
	}
	admin.Auth.SetBearerToken(apiKey)
	admin.Auth.SetUserAgent("curl/7.54.0")

	jira, err := v2.New(nil, host)
	if err != nil {
		log.Fatal().Err(err)
	}

	//jira.Auth.SetBearerToken(apiKey)
	jira.Auth.SetBasicAuth(mail, token)
	jira.Auth.SetUserAgent("curl/7.54.0")

	confluence, err := confluence.New(nil, host)
	if err != nil {
		log.Fatal().Err(err)
	}

	//confluence.Auth.SetBearerToken(apiKey)
	confluence.Auth.SetBasicAuth(mail, token)
	confluence.Auth.SetUserAgent("curl/7.54.0")

	conn := &AtlassianConnection{
		Conf:      conf,
		Host:      host,
		id:        id,
		asset:     asset,
		admin:     admin,
		jira:      jira,
		confluece: confluence,
	}

	return conn, nil
}

func (c *AtlassianConnection) Name() string {
	return "atlassian"
}

func (c *AtlassianConnection) ID() uint32 {
	return c.id
}

func (c *AtlassianConnection) Asset() *inventory.Asset {
	return c.asset
}

func (c *AtlassianConnection) Admin() *admin.Client {
	return c.admin
}

func (c *AtlassianConnection) Jira() *v2.Client {
	return c.jira
}

func (c *AtlassianConnection) Confluence() *confluence.Client {
	return c.confluece
}
