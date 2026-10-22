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
	_ resource.Resource                = &watchResource{}
	_ resource.ResourceWithConfigure   = &watchResource{}
	_ resource.ResourceWithImportState = &watchResource{}
)

func NewWatchResource() resource.Resource {

	return &watchResource{}
}

type watchResource struct {
	client *changedetectionclient.Client
}

type watchResourceModel struct {
	ID         types.String `tfsdk:"id"`
	Title      types.String `tfsdk:"title"`
	Url        types.String `tfsdk:"url"`
	LinkToOpen types.String `tfsdk:"link_to_open"`
}

func (r *watchResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_watch"
}

func (r *watchResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create an changedetection.io Watch",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Description: "The watch's title.",
				Optional:    true,
			},
			"url": schema.StringAttribute{
				Description: "The watch url.",
				Required:    true,
			},
			"link_to_open": schema.StringAttribute{
				Description: "Optional human facing URL to open in the browser.",
				Optional:    true,
			},
		},
	}
}

func (r *watchResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan watchResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	watch := changedetectionclient.CreateWatchReq{
		Title:      plan.Title.ValueStringPointer(),
		URL:        plan.Url.ValueString(),
		LinkToOpen: plan.LinkToOpen.ValueStringPointer(),
	}

	newWatch, err := r.client.CreateWatch(watch)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating watch",
			"Could not create watch, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(newWatch.Uuid)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *watchResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state watchResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	watch, err := r.client.GetWatch(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Watch",
			"Could not read changedetection.io watch ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Title = types.StringValue(watch.Title)
	state.Url = types.StringValue(watch.URL)
	state.LinkToOpen = types.StringValue(watch.LinkToOpen)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *watchResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan watchResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	watch := changedetectionclient.UpdateWatchReq{
		Title:      plan.Title.ValueStringPointer(),
		URL:        plan.Url.ValueStringPointer(),
		LinkToOpen: plan.LinkToOpen.ValueStringPointer(),
	}

	err := r.client.UpdateWatch(plan.ID.ValueString(), &watch)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating changedetection.io Watch",
			"Could not update order, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *watchResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state watchResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteWatch(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting changedetection.io Watch",
			"Could not delete watch, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *watchResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *watchResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
