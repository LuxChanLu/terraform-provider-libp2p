package main

import (
	"github.com/LuxChanLu/terraform-provider-libp2p-keys/internal"

	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: internal.Provider,
	})
}
