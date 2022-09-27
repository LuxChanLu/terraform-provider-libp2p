package internal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

type libp2p struct {
	configured bool
}

func NewProvider() provider.Provider {
	return &libp2p{}
}

func (p *libp2p) GetSchema(context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{}, nil
}

func (p *libp2p) Configure(context.Context, provider.ConfigureRequest, *provider.ConfigureResponse) {
}

func (p *libp2p) DataSources(context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *libp2p) Resources(context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewKeyResource,
	}
}

func (p *libp2p) Metadata(_ context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "libp2p"
	resp.Version = BuildVersion
}
