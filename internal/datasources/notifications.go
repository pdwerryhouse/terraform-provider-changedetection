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
	_ datasource.DataSource              = &notificationsDataSource{}
	_ datasource.DataSourceWithConfigure = &notificationsDataSource{}
)

func NewNotificationsDataSource() datasource.DataSource {
	return &notificationsDataSource{}
}

type notificationsDataSource struct {
	client *changedetectionclient.Client
}

type notificationsDataSourceModel struct {
	NotificationUrls []types.String `tfsdk:"notification_urls"`
}

func (d *notificationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notifications"
}

func (d *notificationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state notificationsDataSourceModel

	notifications, err := d.client.ListNotifications()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read changedetection.io Notification",
			err.Error(),
		)
		return
	}

	for _, notification := range notifications.NotificationUrls {

		notificationUrl := types.StringValue(notification)
		state.NotificationUrls = append(state.NotificationUrls, notificationUrl)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *notificationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

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

func (d *notificationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of notifications.",
		Attributes: map[string]schema.Attribute{
			"notification_urls": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}
