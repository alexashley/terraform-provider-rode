package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rode/rode/proto/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func resourcePolicy() *schema.Resource {
	return &schema.Resource{
		Description:   "An Open Policy Agent Rego policy.",
		CreateContext: resourcePolicyCreate,
		ReadContext:   resourcePolicyRead,
		UpdateContext: resourcePolicyUpdate,
		DeleteContext: resourcePolicyDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Policy name",
				Required:    true,
				Type:        schema.TypeString,
			},
			"description": {
				Description: "A brief summary of the policy",
				Required:    true,
				Type:        schema.TypeString,
			},
			"current_version": {
				Description: "Current version of the policy",
				Computed:    true,
				Type:        schema.TypeInt,
			},
			"policy_version_id": {
				Computed:    true,
				Description: "Policy version id",
				Type:        schema.TypeString,
			},
			"message": {
				Description: "A summary of changes since the last version",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"rego_content": {
				Description: "The Rego code",
				Type:        schema.TypeString,
				Required:    true,
			},
			"created": {
				Description: "Creation timestamp",
				Computed:    true,
				Type:        schema.TypeString,
			},
			"updated": {
				Description: "Last updated timestamp",
				Computed:    true,
				Type:        schema.TypeString,
			},
			"deleted": {
				Description: "",
				Computed:    true,
				Type:        schema.TypeBool,
			},
		},
	}
}

func resourcePolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rode := meta.(v1alpha1.RodeClient)
	policy := &v1alpha1.Policy{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Policy: &v1alpha1.PolicyEntity{
			Message:     d.Get("message").(string),
			RegoContent: d.Get("rego_content").(string),
		},
	}

	response, err := rode.CreatePolicy(ctx, policy)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(response.Id)

	return resourcePolicyRead(ctx, d, meta)
}

func resourcePolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rode := meta.(v1alpha1.RodeClient)

	policy, err := rode.GetPolicy(ctx, &v1alpha1.GetPolicyRequest{Id: d.Id()})
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", policy.Name)
	d.Set("description", policy.Description)
	d.Set("current_version", policy.CurrentVersion)
	d.Set("policy_version_id", policy.Policy.Id)
	d.Set("message", policy.Policy.Message)
	d.Set("rego_content", policy.Policy.RegoContent)
	d.Set("created", formatProtoTimestamp(policy.Created))
	d.Set("updated", formatProtoTimestamp(policy.Updated))
	d.Set("deleted", policy.Deleted)

	return nil
}

func resourcePolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Errorf("unimplemented")
}

func resourcePolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rode := meta.(v1alpha1.RodeClient)

	_, err := rode.DeletePolicy(ctx, &v1alpha1.DeletePolicyRequest{
		Id: d.Id(),
	})

	if status.Code(err) == codes.NotFound {
		d.SetId("")
		return nil
	}

	return diag.FromErr(err)
}
