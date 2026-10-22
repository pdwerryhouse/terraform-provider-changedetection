// Copyright (C) 2026 Paul Dwerryhouse <paul@dwerryhouse.com.au>
//
// This file is part of terraform-provider-changedetection.
//
// terraform-provider-changedetection is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// terraform-provider-changedetection is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with terraform-provider-changedetection.  If not, see <https://www.gnu.org/licenses/>.

package provider

import (
	"context"
	"os"
	"terraform-provider-changedetection/internal/datasources"
	"terraform-provider-changedetection/internal/resources"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	changedetectionclient "codeberg.org/pdwerryhouse/changedetection-client-go/changedetection"
)

var (
	_ provider.Provider = &changedetectionProvider{}
)

type changedetectionProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	ApiKey   types.String `tfsdk:"apikey"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &changedetectionProvider{
			version: version,
		}
	}
}

type changedetectionProvider struct {
	version string
}

func (p *changedetectionProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "changedetection"
	resp.Version = p.version
}

func (p *changedetectionProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Terraform provider for changedetection.io",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Description: "The changedetection.io API Endpoint. May also be provided via CHANGEDETECTION_ENDPOINT.",
				Optional:    true,
			},
			"apikey": schema.StringAttribute{
				Description: "The changedetection.io API Key. May also be provided via CHANGEDETECTION_API_KEY.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *changedetectionProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config changedetectionProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// XXX improve error messages
	if config.Endpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown changedetection.io host",
			"The changedetection.io host is unknown.",
		)
	}

	if config.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("apikey"),
			"Unknown changedetection.io apikey",
			"The changedetection.io apikey is unknown.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := os.Getenv("CHANGEDETECTION_ENDPOINT")
	apikey := os.Getenv("CHANGEDETECTION_API_KEY")

	if !config.Endpoint.IsNull() {
		endpoint = config.Endpoint.ValueString()
	}

	if !config.ApiKey.IsNull() {
		apikey = config.ApiKey.ValueString()
	}

	// XXX improve these
	if endpoint == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Unknown changedetection.io endpoint",
			"The changedetection.io endpoint is unknown.",
		)
	}

	if apikey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("apikey"),
			"Unknown changedetection.io apikey",
			"The changedetection.io apikey is unknown.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	c, err := changedetectionclient.NewClient(&endpoint, &apikey)

	// XXX improve this
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create changedetection.io API Client",
			"An Unexpected error occurred.",
		)
	}

	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *changedetectionProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		datasources.NewWatchesDataSource,
		datasources.NewTagsDataSource,
	}
}

func (p *changedetectionProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewWatchResource,
		resources.NewTagResource,
	}
}
