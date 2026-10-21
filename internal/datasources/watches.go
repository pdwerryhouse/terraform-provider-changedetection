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

package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	changedetectionclient "codeberg.org/pdwerryhouse/changedetection-client-go/changedetection"
)

var (
	_ datasource.DataSource              = &watchesDataSource{}
	_ datasource.DataSourceWithConfigure = &watchesDataSource{}
)

func NewWatchesDataSource() datasource.DataSource {
	return &watchesDataSource{}
}

type watchesDataSource struct {
	client *changedetectionclient.Client
}

type watchesDataSourceModel struct {
	Watch []watchesModel `tfsdk:"watches"`
}

type watchesModel struct {
	ID                types.String   `tfsdk:"id"`
	Title             types.String   `tfsdk:"title"`
	Url               types.String   `tfsdk:"url"`
	PageTitle         types.String   `tfsdk:"page_title"`
	Link              types.String   `tfsdk:"link"`
	OpenLink          types.String   `tfsdk:"open_link"`
	LinkToOpen        types.String   `tfsdk:"link_to_open"`
	Paused            types.Bool     `tfsdk:"paused"`
	NotificationMuted types.Bool     `tfsdk:"notification_muted"`
	LastChecked       types.Int64    `tfsdk:"last_checked"`
	Tags              []types.String `tfsdk:"tags"`
}

func (d *watchesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_watches"
}

func (d *watchesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state watchesDataSourceModel

	watches, err := d.client.ListWatches(nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read changedetection.io Watch",
			err.Error(),
		)
		return
	}

	for id, watch := range watches {

		var tags []types.String

		for _, tag := range watch.Tags {
			tags = append(tags, types.StringValue(tag))
		}

		watchState := watchesModel{
			ID:                types.StringValue(id),
			Title:             types.StringValue(watch.Title),
			PageTitle:         types.StringValue(watch.PageTitle),
			Url:               types.StringValue(watch.URL),
			Link:              types.StringValue(watch.Link),
			OpenLink:          types.StringValue(watch.OpenLink),
			LinkToOpen:        types.StringValue(watch.LinkToOpen),
			Paused:            types.BoolValue(watch.Paused),
			NotificationMuted: types.BoolValue(watch.NotificationMuted),
			LastChecked:       types.Int64Value(watch.LastChecked),
			Tags:              tags,
		}

		state.Watch = append(state.Watch, watchState)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *watchesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*changedetectionclient.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *changedetectionclient.Client, got %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *watchesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of watches.",
		Attributes: map[string]schema.Attribute{
			"watches": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Watch Id.",
							Computed:    true,
						},
						"title": schema.StringAttribute{
							Description: "The title of the watch.",
							Computed:    true,
						},
						"url": schema.StringAttribute{
							Description: "The watch's url.",
							Computed:    true,
						},
						"link": schema.StringAttribute{
							Computed: true,
						},
						"open_link": schema.StringAttribute{
							Computed: true,
						},
						"link_to_open": schema.StringAttribute{
							Computed: true,
						},
						"page_title": schema.StringAttribute{
							Description: "The watch's page title.",
							Computed:    true,
						},
						"paused": schema.BoolAttribute{
							Computed: true,
						},
						"notification_muted": schema.BoolAttribute{
							Computed: true,
						},
						"last_checked": schema.Int64Attribute{
							Computed: true,
						},
						"tags": schema.ListAttribute{
							Description: "The watch's tags.",
							ElementType: types.StringType,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}
