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

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
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
	ID                types.String `tfsdk:"id"`
	Body              types.String `tfsdk:"body"`
	FetchBackend      types.String `tfsdk:"fetch_backend"`
	Headers           types.Map    `tfsdk:"headers"`
	IgnoreStatusCodes types.Bool   `tfsdk:"ignore_status_codes"`
	LinkToOpen        types.String `tfsdk:"link_to_open"`
	Method            types.String `tfsdk:"method"`
	NotificationMuted types.Bool   `tfsdk:"notification_muted"`
	Paused            types.Bool   `tfsdk:"paused"`
	Processor         types.String `tfsdk:"processor"`
	Proxy             types.String `tfsdk:"proxy"`
	Tags              types.List   `tfsdk:"tags"`
	Title             types.String `tfsdk:"title"`
	WebdriverDelay    types.Int64  `tfsdk:"webdriver_delay"`
	Url               types.String `tfsdk:"url"`
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
			"body": schema.StringAttribute{
				Description: "HTTP Request Body.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"fetch_backend": schema.StringAttribute{
				Description: "Backend to use for fetching content: system, html_requests, html_webdriver, extra_browser_* . Default is system.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("system"),
			},
			"headers": schema.MapAttribute{
				Description: "Map of headers.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default: mapdefault.StaticValue(
					types.MapValueMust(types.StringType, map[string]attr.Value{}),
				),
			},
			"ignore_status_codes": schema.BoolAttribute{
				Description: "Ignore HTTP status code errors. Default is false.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"link_to_open": schema.StringAttribute{
				Description: "Optional human facing URL to open in the browser.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"method": schema.StringAttribute{
				Description: "HTTP method to use (GET, POST, DELETE, PUT).",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("GET"),
			},
			"notification_muted": schema.BoolAttribute{
				Description: "Whether notifications are muted. Default is false.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"paused": schema.BoolAttribute{
				Description: "Whether the web page change monitor (watch) is paused.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"processor": schema.StringAttribute{
				Description: "Processor mode to be used for change detection. text_json_diff or restock_diff.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("text_json_diff"),
			},
			"proxy": schema.StringAttribute{
				Description: "Proxy configuration.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"tags": schema.ListAttribute{
				Description: "List of tag uuids.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default: listdefault.StaticValue(
					types.ListValueMust(types.StringType, []attr.Value{}),
				),
			},
			"title": schema.StringAttribute{
				Description: "The watch's title.",
				Optional:    true,
			},
			"webdriver_delay": schema.Int64Attribute{
				Description: "Delay in seconds for the webdriver.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
			"url": schema.StringAttribute{
				Description: "The watch url.",
				Required:    true,
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

	var tags []string
	diags = plan.Tags.ElementsAs(ctx, &tags, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var headers map[string]string
	diags = plan.Headers.ElementsAs(ctx, &headers, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	watch := changedetectionclient.CreateWatchReq{
		Body:              plan.Body.ValueStringPointer(),
		FetchBackend:      plan.FetchBackend.ValueStringPointer(),
		IgnoreStatusCodes: plan.IgnoreStatusCodes.ValueBoolPointer(),
		LinkToOpen:        plan.LinkToOpen.ValueStringPointer(),
		Method:            plan.Method.ValueStringPointer(),
		NotificationMuted: plan.NotificationMuted.ValueBoolPointer(),
		Paused:            plan.Paused.ValueBoolPointer(),
		Processor:         plan.Processor.ValueStringPointer(),
		Proxy:             plan.Proxy.ValueStringPointer(),
		Tags:              &tags,
		Title:             plan.Title.ValueStringPointer(),
		URL:               plan.Url.ValueString(),
		WebdriverDelay:    plan.WebdriverDelay.ValueInt64Pointer(),
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

	tags, diags := types.ListValueFrom(ctx, types.StringType, watch.Tags)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	headers, diags := types.MapValueFrom(ctx, types.StringType, watch.Headers)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Body = types.StringValue(watch.Body)
	state.FetchBackend = types.StringValue(watch.FetchBackend)
	state.Headers = headers
	state.IgnoreStatusCodes = types.BoolValue(watch.IgnoreStatusCodes)
	state.LinkToOpen = types.StringValue(watch.LinkToOpen)
	state.Method = types.StringValue(watch.Method)
	state.NotificationMuted = types.BoolValue(watch.NotificationMuted)
	state.Paused = types.BoolValue(watch.Paused)
	state.Processor = types.StringValue(watch.Processor)
	state.Proxy = types.StringValue(watch.Proxy)
	state.Tags = tags
	state.Title = types.StringValue(watch.Title)
	state.WebdriverDelay = types.Int64Value(watch.WebdriverDelay)
	state.Url = types.StringValue(watch.URL)

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

	var tags []string
	diags = plan.Tags.ElementsAs(ctx, &tags, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var headers map[string]string
	diags = plan.Headers.ElementsAs(ctx, &headers, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	watch := changedetectionclient.UpdateWatchReq{
		Body:              plan.Body.ValueStringPointer(),
		FetchBackend:      plan.FetchBackend.ValueStringPointer(),
		Headers:           &headers,
		IgnoreStatusCodes: plan.IgnoreStatusCodes.ValueBoolPointer(),
		LinkToOpen:        plan.LinkToOpen.ValueStringPointer(),
		Method:            plan.Method.ValueStringPointer(),
		NotificationMuted: plan.NotificationMuted.ValueBoolPointer(),
		Paused:            plan.Paused.ValueBoolPointer(),
		Processor:         plan.Processor.ValueStringPointer(),
		Proxy:             plan.Proxy.ValueStringPointer(),
		Tags:              &tags,
		Title:             plan.Title.ValueStringPointer(),
		WebdriverDelay:    plan.WebdriverDelay.ValueInt64Pointer(),
		URL:               plan.Url.ValueStringPointer(),
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
