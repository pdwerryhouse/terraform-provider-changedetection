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
	_ datasource.DataSource              = &tagsDataSource{}
	_ datasource.DataSourceWithConfigure = &tagsDataSource{}
)

func NewTagsDataSource() datasource.DataSource {
	return &tagsDataSource{}
}

type tagsDataSource struct {
	client *changedetectionclient.Client
}

type tagsDataSourceModel struct {
	Tag map[string]tagsModel `tfsdk:"tags"`
}

type tagsModel struct {
	ID                types.String `tfsdk:"id"`
	Title             types.String `tfsdk:"title"`
	NotificationMuted types.Bool   `tfsdk:"notification_muted"`
	DateCreated       types.Int64  `tfsdk:"date_created"`
}

func (d *tagsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tags"
}

func (d *tagsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state tagsDataSourceModel

	tags, err := d.client.ListTags(nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read changedetection.io Tag",
			err.Error(),
		)
		return
	}

	state.Tag = make(map[string]tagsModel)

	for id, tag := range tags {

		tagState := tagsModel{
			ID:                types.StringValue(id),
			Title:             types.StringValue(tag.Title),
			NotificationMuted: types.BoolValue(tag.NotificationMuted),
			DateCreated:       types.Int64Value(tag.DateCreated),
		}

		state.Tag[id] = tagState
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *tagsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

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

func (d *tagsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of tags.",
		Attributes: map[string]schema.Attribute{
			"tags": schema.MapNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Tag Id.",
							Computed:    true,
						},
						"title": schema.StringAttribute{
							Description: "The title of the tag.",
							Computed:    true,
						},
						"notification_muted": schema.BoolAttribute{
							Description: "True if notifications are muted for this tag.",
							Computed:    true,
						},
						"date_created": schema.Int64Attribute{
							Description: "The date at which the tag was created (seconds since Jan 1 1970).",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}
