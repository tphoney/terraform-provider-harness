package policyset

import (
	"context"
	"errors"
	"net/http"

	"github.com/antihax/optional"
	"github.com/harness/harness-go-sdk/harness/policymgmt"
	"github.com/harness/terraform-provider-harness/helpers"
	"github.com/harness/terraform-provider-harness/internal"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourcePolicyset() *schema.Resource {
	resource := &schema.Resource{
		Description: "Data source for retrieving a Harness policyset.",

		ReadContext: dataSourceProjectRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Name of the policyset.",
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
			},
			"identifier": {
				Description: "Identifier of the policyset.",
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
			},
		},
	}

	helpers.SetMultiLevelDatasourceSchema(resource.Schema)

	return resource
}

func dataSourceProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*internal.Session).GetPolicyManagementClient()
	id := d.Get("identifier").(string)

	var err error
	var policyset policymgmt.PolicySet2
	var httpResp *http.Response

	if id != "" {
		policyset, _, _ = c.PolicysetsApi.PolicysetsFind(ctx, id, &policymgmt.PolicysetsApiPolicysetsFindOpts{
			AccountIdentifier: optional.NewString(meta.(*internal.Session).AccountId),
			XApiKey:           optional.NewString(meta.(*internal.Session).PLClient.ApiKey),
		})
	} else {
		return diag.FromErr(errors.New("identifier must be specified"))
	}
	if err != nil {
		return helpers.HandleApiError(err, d, httpResp)
	}

	bla := policymgmt.PolicySet2{}
	if policyset.Identifier == bla.Identifier {
		d.SetId("")
		d.MarkNewResource()
		return nil
	}

	readPolicyset(d, policyset)
	return nil
}
