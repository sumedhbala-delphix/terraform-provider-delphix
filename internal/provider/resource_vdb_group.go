package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	dctapi "github.com/delphix/dct-sdk-go/v25"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVdbGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Resource for managing VDB Groups.",

		CreateContext: resourceVdbGroupCreate,
		ReadContext:   resourceVdbGroupRead,
		UpdateContext: resourceVdbGroupUpdate,
		DeleteContext: resourceVdbGroupDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vdb_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceVdbGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := meta.(*apiClient).client

	vdbGroupCreateReq := *dctapi.NewCreateVDBGroupRequest(d.Get("name").(string))
	vdbGroupCreateReq.SetVdbIds(toStringArray(d.Get("vdb_ids")))
	if v, ok := d.GetOk("tags"); ok {
		tags := make([]dctapi.Tag, 0)
		for _, tag := range v.([]interface{}) {
			tagMap := tag.(map[string]interface{})
			tag := dctapi.NewTag(tagMap["key"].(string), tagMap["value"].(string))
			tags = append(tags, *tag)
		}
		vdbGroupCreateReq.SetTags(tags)
	}
	apiRes, httpRes, err := client.VDBGroupsAPI.CreateVdbGroup(ctx).CreateVDBGroupRequest(vdbGroupCreateReq).Execute()

	if diags := apiErrorResponseHelper(ctx, apiRes, httpRes, err); diags != nil {
		return diags
	}

	d.SetId(apiRes.VdbGroup.GetId())
	readDiags := resourceVdbGroupRead(ctx, d, meta)

	if readDiags.HasError() {
		return readDiags
	}
	return diags
}

func resourceVdbGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*apiClient).client

	var diags diag.Diagnostics

	vdbGroupId := d.Id()
	tflog.Info(ctx, DLPX+INFO+"VdbGroupId: "+vdbGroupId)
	apiRes, httpRes, err := client.VDBGroupsAPI.GetVdbGroup(ctx, vdbGroupId).Execute()

	if diags := apiErrorResponseHelper(ctx, apiRes, httpRes, err); diags != nil {
		d.SetId("")
		return diags
	}

	d.Set("name", apiRes.GetName())
	d.Set("vdb_ids", apiRes.GetVdbIds())
	if tags := apiRes.GetTags(); tags != nil {
		tagList := make([]map[string]string, 0)
		for _, tag := range tags {
			tagList = append(tagList, map[string]string{
				"key":   tag.GetKey(),
				"value": tag.GetValue(),
			})
		}
		d.Set("tags", tagList)
	}
	return diags
}

func resourceVdbGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).client
	var diags diag.Diagnostics
	vdbGroupId := d.Id()

	// Only handle tag updates
	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		oldTagList := oldTags.([]interface{})
		newTagList := newTags.([]interface{})

		// Delete old tags
		if len(oldTagList) > 0 {
			deleteTag := dctapi.NewDeleteTag()
			tags := make([]dctapi.Tag, 0)
			for _, tag := range oldTagList {
				tagMap := tag.(map[string]interface{})
				tag := dctapi.NewTag(tagMap["key"].(string), tagMap["value"].(string))
				tags = append(tags, *tag)
			}
			deleteTag.SetTags(tags)
			tagDelResp, tagDelErr := client.VDBGroupsAPI.DeleteVdbGroupTags(ctx, vdbGroupId).DeleteTag(*deleteTag).Execute()
			if diags := apiErrorResponseHelper(ctx, tagDelResp, nil, tagDelErr); diags != nil {
				return diags
			}
		}

		// Create new tags
		if len(newTagList) > 0 {
			tags := make([]dctapi.Tag, 0)
			for _, tag := range newTagList {
				tagMap := tag.(map[string]interface{})
				tag := dctapi.NewTag(tagMap["key"].(string), tagMap["value"].(string))
				tags = append(tags, *tag)
			}
			_, httpResp, tagCrtErr := client.VDBGroupsAPI.CreateVdbGroupsTags(ctx, vdbGroupId).TagsRequest(*dctapi.NewTagsRequest(tags)).Execute()
			if diags := apiErrorResponseHelper(ctx, nil, httpResp, tagCrtErr); diags != nil {
				return diags
			}
		}
	}

	// Read the updated resource to ensure state is in sync
	readDiags := resourceVdbGroupRead(ctx, d, meta)
	if readDiags.HasError() {
		return readDiags
	}

	return diags
}

func resourceVdbGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).client

	var diags diag.Diagnostics

	vdbGroupId := d.Id()

	deleteVdbParams := dctapi.NewDeleteVDBParametersWithDefaults()
	deleteVdbParams.SetForce(false)

	httpRes, err := client.VDBGroupsAPI.DeleteVdbGroup(ctx, vdbGroupId).Execute()

	if diags := apiErrorResponseHelper(ctx, nil, httpRes, err); diags != nil {
		return diags
	}
	if err != nil {
		resBody, err := ResponseBodyToString(ctx, httpRes.Body)
		if err != nil {
			return diag.FromErr(err)
		}
		return diag.Errorf(resBody)
	}

	return diags
}
