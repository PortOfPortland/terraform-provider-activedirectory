package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/portofportland/terraform-provider-activedirectory/activedirectory"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: activedirectory.Provider,
	})
}
