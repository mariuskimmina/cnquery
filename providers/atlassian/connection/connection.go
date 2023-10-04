package connection

import (
	"os"

	"github.com/ctreminiom/go-atlassian/admin"
	"github.com/ctreminiom/go-atlassian/jira/v2"
	_ "github.com/ctreminiom/go-atlassian/jira/v3"
	"github.com/rs/zerolog/log"
	"go.mondoo.com/cnquery/providers-sdk/v1/inventory"
)

type AtlassianConnection struct {
	id    uint32
	Conf  *inventory.Config
	asset *inventory.Asset
	admin *admin.Client
	jira  *v2.Client
	// Add custom connection fields here
}

func NewAtlassianConnection(id uint32, asset *inventory.Asset, conf *inventory.Config) (*AtlassianConnection, error) {
	apiKey := os.Getenv("ATLASSIAN_KEY")
	token := os.Getenv("ATLASSIAN_TOKEN")
	admin, err := admin.New(nil)
	if err != nil {
		log.Fatal().Err(err)
	}
	admin.Auth.SetBearerToken(apiKey)
	admin.Auth.SetUserAgent("curl/7.54.0")

	jira, err := v2.New(nil, "lunalectric.atlassian.net")
	if err != nil {
		log.Fatal().Err(err)
	}

	jira.Auth.SetBasicAuth("marius@mondoo.com", token)

	conn := &AtlassianConnection{
		Conf:  conf,
		id:    id,
		asset: asset,
		admin: admin,
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
