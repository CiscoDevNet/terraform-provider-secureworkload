package secureworkload

import (
	"strings"
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSecureWorkloadScope() *schema.Resource {
	return &schema.Resource{
		Description: "Data source for fetching scopes from secure-workload\n\n" +
			"An example is shown below: \n" +
			"```hcl\n" +
			"data \"secureworkload_scope\" \"scope\" {\n" +
			"	exact_name = \"RootScope:ChildScope\"\n" +
			"}\n" +
			"```",
		ReadContext: dataSourceSecureWorkloadScopeRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of this resource",
			},
			"vrf_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The VRF ID of this resource",
			},
			"exact_short_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Short name of the scope", 
			},
			"exact_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the scope",
			},
		},
	}
}

func dataSourceSecureWorkloadScopeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// Logic for query param

	getUrl := "?"
	if val, ok := d.GetOk("exact_name"); ok {
		temp := strings.Replace(val.(string), " ", "%20", -1)
		getUrl = getUrl + "exact_name=" +  temp
	}
	if val, ok := d.GetOk("exact_short_name"); ok {
		temp := strings.Replace(val.(string), " ", "%20", -1)
		getUrl = getUrl + "&exact_short_name=" + temp
	}
	if val, ok := d.GetOk("vrf_id"); ok {
		getUrl = getUrl + "&vrf_id=" + val.(string)
	}
	resScope, err := c.GetScopeByParam( getUrl)
	if (err != nil) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to get the scope",
			Detail:   err.Error(),
		})
		return diags
	}
	if len(resScope) > 0 {
		d.SetId(resScope[0].Id)

		if err := d.Set("exact_short_name", resScope[0].ExactShortName); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "unable to read scope",
				Detail:   err.Error(),
			})
			return diags
		}

		if err := d.Set("exact_name", resScope[0].ExactName); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "unable to read scope",
				Detail:   err.Error(),
			})
			return diags
		}
		if err := d.Set("vrf_id", resScope[0].VRFId); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "unable to read scope",
				Detail:   err.Error(),
			})
			return diags
		}

		return diags
	}
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "scope with the given name is not present",
		Detail:   "scope with the given name is not present",
	})
	return diags
	
}
