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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &tagResource{}
	_ resource.ResourceWithConfigure   = &tagResource{}
	_ resource.ResourceWithImportState = &tagResource{}
)

func NewTagResource() resource.Resource {

	return &tagResource{}
}

type tagResource struct {
	client *changedetectionclient.Client
}

type tagResourceModel struct {
	ID                types.String `tfsdk:"id"`
	Title             types.String `tfsdk:"title"`
	NotificationMuted types.Bool   `tfsdk:"notification_muted"`
	DateCreated       types.Int64  `tfsdk:"date_created"`
}

func (r *tagResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tag"
}

func (r *tagResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create a changedetection.io Tag",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Description: "The tag's title.",
				Required:    true,
			},
			"notification_muted": schema.BoolAttribute{
				Description: "True if notifications are muted for this tag.",
				Optional:    true,
				Computed:    true,
			},
			"date_created": schema.Int64Attribute{
				Description: "The date at which the tag was created (seconds since Jan 1 1970).",
				Computed:    true,
			},
		},
	}
}

func (r *tagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan tagResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tag := changedetectionclient.CreateTagReq{
		Title:             plan.Title.ValueString(),
		NotificationMuted: plan.NotificationMuted.ValueBoolPointer(),
	}

	newTag, err := r.client.CreateTag(tag)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating tag",
			"Could not create tag, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(newTag.Uuid)

	fetchTag, err := r.client.GetTag(newTag.Uuid, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error fetching tag",
			"Could not fetch tag, unexpected error: "+err.Error(),
		)
		return
	}

	plan.DateCreated = types.Int64Value(fetchTag.DateCreated)
	plan.NotificationMuted = types.BoolValue(fetchTag.NotificationMuted)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *tagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state tagResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tag, err := r.client.GetTag(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Tag",
			"Could not read changedetection.io tag ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Title = types.StringValue(tag.Title)
	state.NotificationMuted = types.BoolValue(tag.NotificationMuted)
	state.DateCreated = types.Int64Value(tag.DateCreated)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *tagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan tagResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	uuid := plan.ID.ValueString()

	tag := changedetectionclient.UpdateTagReq{
		Title:             plan.Title.ValueStringPointer(),
		NotificationMuted: plan.NotificationMuted.ValueBoolPointer(),
	}

	err := r.client.UpdateTag(plan.ID.ValueString(), &tag)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating changedetection.io Tag",
			"Could not update order, unexpected error: "+err.Error(),
		)
		return
	}

	fetchTag, err := r.client.GetTag(uuid, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error fetching tag",
			"Could not fetch tag, unexpected error: "+err.Error(),
		)
		return
	}

	plan.DateCreated = types.Int64Value(fetchTag.DateCreated)
	plan.NotificationMuted = types.BoolValue(fetchTag.NotificationMuted)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *tagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state tagResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteTag(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting changedetection.io Tag",
			"Could not delete tag, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *tagResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *tagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
