package main

import (
	"github.com/hashicorp/terraform/plugin"
	tutorial "github.com/naka-gawa/terraform-provider-tutorial/pkg"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: tutorial.Provider,
	})
}
