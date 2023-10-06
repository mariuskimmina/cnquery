package main

import (
	"os"

	"go.mondoo.com/cnquery/v9/providers-sdk/v1/plugin"
	"go.mondoo.com/cnquery/v9/providers/atlassian/provider"
)

func main() {
	plugin.Start(os.Args, provider.Init())
}
