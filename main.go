package main

import (
	"context"

	"github.com/LuxChanLu/terraform-provider-libp2p/internal"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	providerserver.Serve(context.Background(), internal.NewProvider, providerserver.ServeOpts{
		Address: "registry.terraform.io/hashicorp/random",
	})
}
