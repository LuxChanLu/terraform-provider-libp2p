package internal_test

import (
	"github.com/LuxChanLu/terraform-provider-libp2p/internal"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

func protoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"libp2p": providerserver.NewProtocol6WithError(internal.NewProvider()),
	}
}
