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

package resources

import (
	"context"
	"fmt"

	changedetectionclient "codeberg.org/pdwerryhouse/changedetection-client-go/changedetection"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &notificationResource{}
	_ resource.ResourceWithConfigure   = &notificationResource{}
	_ resource.ResourceWithImportState = &notificationResource{}
)

func NewNotificationResource() resource.Resource {

	return &notificationResource{}
}

type notificationResource struct {
	client *changedetectionclient.Client
}

type notificationResourceModel struct {
	NotificationUrls []types.String `tfsdk:"notification_urls"`
}

func (r *notificationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification"
}

func (r *notificationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create changedetection.io notifications",
		Attributes: map[string]schema.Attribute{
			"notification_urls": schema.ListAttribute{
				Description: "The list of notification urls.",
				ElementType: types.StringType,
				Required:    true,
			},
		},
	}
}

func (r *notificationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan notificationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var notification_urls []string

	for _, notification := range plan.NotificationUrls {
		notification_url := notification.ValueString()
		notification_urls = append(notification_urls, notification_url)
	}

	notification := changedetectionclient.Notifications{
		NotificationUrls: notification_urls,
	}

	_, err := r.client.ReplaceNotifications(notification)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating notification",
			"Could not create notification, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *notificationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state notificationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	notifications, err := r.client.ListNotifications()
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

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *notificationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan notificationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var notification_urls []string

	for _, notification := range plan.NotificationUrls {
		notification_url := notification.ValueString()
		notification_urls = append(notification_urls, notification_url)
	}

	notification := changedetectionclient.Notifications{
		NotificationUrls: notification_urls,
	}

	_, err := r.client.ReplaceNotifications(notification)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating notification",
			"Could not create notification, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *notificationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state notificationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	notification := changedetectionclient.Notifications{
		NotificationUrls: []string{},
	}

	_, err := r.client.ReplaceNotifications(notification)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting changedetection.io Notification",
			"Could not delete notification, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *notificationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*changedetectionclient.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *changedetectionclient.Client, got %T. Please report this issue to the provider developer.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *notificationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
