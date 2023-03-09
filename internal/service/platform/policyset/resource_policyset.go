package policyset

import (
	"context"
	"net/http"

	"github.com/antihax/optional"
	"github.com/harness/harness-go-sdk/harness/policymgmt"
	"github.com/harness/terraform-provider-harness/helpers"
	"github.com/harness/terraform-provider-harness/internal"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourcePolicyset() *schema.Resource {
	resource := &schema.Resource{
		Description: "Resource for creating a Harness Policy.",

		ReadContext:   resourcePolicysetRead,
		UpdateContext: resourcePolicysetCreateOrUpdate,
		DeleteContext: resourcePolicysetDelete,
		CreateContext: resourcePolicysetCreateOrUpdate,
		Importer:      helpers.OrgResourceImporter,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Name of the policy.",
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
			},
			"identifier": {
				Description: "Identifier of the policy.",
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
			},
		},
	}

	helpers.SetMultiLevelResourceSchema(resource.Schema)

	return resource
}

func resourcePolicysetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*internal.Session).GetPolicyManagementClient()

	id := d.Id()

	localVarOptionals := policymgmt.PolicysetsApiPolicysetsFindOpts{
		AccountIdentifier: optional.NewString(meta.(*internal.Session).AccountId),
		XApiKey:           optional.NewString(meta.(*internal.Session).PLClient.ApiKey),
	}
	// check for project and org
	if d.Get("project_id").(string) != "" {
		localVarOptionals.ProjectIdentifier = helpers.BuildField(d, "project_id")
	}
	if d.Get("org_id").(string) != "" {
		localVarOptionals.OrgIdentifier = helpers.BuildField(d, "org_id")
	}

	policy, httpResp, err := c.PolicysetsApi.PolicysetsFind(ctx, id, &localVarOptionals)

	if err != nil {
		return helpers.HandleApiError(err, d, httpResp)
	}

	readPolicyset(d, policy)

	return nil
}

func resourcePolicysetCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*internal.Session).GetPolicyManagementClient()
	var err error
	var responsePolicyset policymgmt.PolicySet2
	var httpResp *http.Response
	id := d.Id()

	if id == "" {
		body := policymgmt.CreateRequestBody2{
			Name:       d.Get("name").(string),
			Identifier: d.Get("identifier").(string),
		}
		localVarOptionals := policymgmt.PolicysetsApiPolicysetsCreateOpts{
			AccountIdentifier: optional.NewString(meta.(*internal.Session).AccountId),

			XApiKey: optional.NewString(meta.(*internal.Session).PLClient.ApiKey),
		}
		// check for project and org
		if d.Get("project_id").(string) != "" {
			localVarOptionals.ProjectIdentifier = helpers.BuildField(d, "project_id")
		}
		if d.Get("org_id").(string) != "" {
			localVarOptionals.OrgIdentifier = helpers.BuildField(d, "org_id")
		}
		responsePolicyset, httpResp, err = c.PolicysetsApi.PolicysetsCreate(ctx, body, &localVarOptionals)
	} else {
		body := policymgmt.UpdateRequestBody2{
			Name: d.Get("name").(string),
		}
		localVarOptionals := policymgmt.PolicysetsApiPolicysetsUpdateOpts{
			AccountIdentifier: optional.NewString(meta.(*internal.Session).AccountId),
			XApiKey:           optional.NewString(meta.(*internal.Session).PLClient.ApiKey),
		}
		if d.Get("project_id").(string) != "" {
			localVarOptionals.ProjectIdentifier = helpers.BuildField(d, "project_id")
		}
		if d.Get("org_id").(string) != "" {
			localVarOptionals.OrgIdentifier = helpers.BuildField(d, "org_id")
		}
		httpResp, err = c.PolicysetsApi.PolicysetsUpdate(ctx, body, id, &localVarOptionals)
		if err == nil && httpResp.StatusCode == http.StatusNoContent {
			// if we get a 204, we need to get the policy again to get the updated values
			findLocalVarOptionals := policymgmt.PolicysetsApiPolicysetsFindOpts{
				AccountIdentifier: optional.NewString(meta.(*internal.Session).AccountId),
				XApiKey:           optional.NewString(meta.(*internal.Session).PLClient.ApiKey),
			}
			// check for project and org
			if d.Get("project_id").(string) != "" {
				findLocalVarOptionals.ProjectIdentifier = helpers.BuildField(d, "project_id")
			}
			if d.Get("org_id").(string) != "" {
				findLocalVarOptionals.OrgIdentifier = helpers.BuildField(d, "org_id")
			}
			responsePolicyset, httpResp, err = c.PolicysetsApi.PolicysetsFind(ctx, id, &findLocalVarOptionals)
		}
	}
	if err != nil {
		return helpers.HandleApiError(err, d, httpResp)
	}

	readPolicyset(d, responsePolicyset)
	return nil
}

func resourcePolicysetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*internal.Session).GetPolicyManagementClient()

	localVarOptionals := policymgmt.PoliciesApiPoliciesDeleteOpts{
		AccountIdentifier: optional.NewString(meta.(*internal.Session).AccountId),
		XApiKey:           optional.NewString(meta.(*internal.Session).PLClient.ApiKey),
	}
	// check for project and org
	if d.Get("project_id").(string) != "" {
		localVarOptionals.ProjectIdentifier = helpers.BuildField(d, "project_id")
	}
	if d.Get("org_id").(string) != "" {
		localVarOptionals.OrgIdentifier = helpers.BuildField(d, "org_id")
	}
	httpResp, err := c.PoliciesApi.PoliciesDelete(ctx, d.Id(), &localVarOptionals)

	if err != nil {
		return helpers.HandleApiError(err, d, httpResp)
	}

	return nil
}

func readPolicyset(d *schema.ResourceData, policy policymgmt.PolicySet2) {
	d.SetId(policy.Identifier)
	_ = d.Set("identifier", policy.Identifier)
	_ = d.Set("org_id", policy.OrgId)
	_ = d.Set("account_id", policy.AccountId)
	_ = d.Set("project_id", policy.ProjectId)
	_ = d.Set("name", policy.Name)
}
