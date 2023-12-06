package secureworkload

import (
	"strings"
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSecureWorkloadApplication() *schema.Resource {
	return &schema.Resource{
		Description: "Data source for fetching workspace from secure-workload\n\n" +
			"An example is shown below: \n" +
			"```hcl\n" +
			"data \"secureworkload_workspace\" \"workspace\" {\n" +
			"	name = \"MyWorkspace\"\n" +
			"}\n" +
			"```",
		ReadContext: dataSourceSecureWorkloadApplicationRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of this resource",
			},
			"app_scope_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the parent scope",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Short name of the workspace",
			},
			"author": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Author of the workspace",
			},
			"primary": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Determines whether the workspace is a primary one or not",
			},
		},
	}
}

func dataSourceSecureWorkloadApplicationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// Logic for query param

	getUrl := "?"
	if val, ok := d.GetOk("name"); ok {
		temp := strings.Replace(val.(string), " ", "%20", -1)
		getUrl = getUrl + "exact_name=" +  temp
	}
	resApp, err := c.GetApplicationByParam( getUrl)
	if (err != nil) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to read workspace",
			Detail:   err.Error(),
		})
		return diags
	}
	if len(resApp) > 0 {
		d.SetId(resApp[0].Id)
			if err := d.Set("name", resApp[0].Name); err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "unable to read workspace",
					Detail:   err.Error(),
				})
				return diags
			}
			if err := d.Set("app_scope_id", resApp[0].AppScopeId); err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "unable to read workspace",
					Detail:   err.Error(),
				})
				return diags
			}	
			if err := d.Set("primary", resApp[0].Primary); err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "unable to read workspace",
					Detail:   err.Error(),
				})
				return diags
			}	
			if err := d.Set("author", resApp[0].Author); err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "unable to read workspace",
					Detail:   err.Error(),
				})
				return diags
			}	

		return diags
	}
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "workspace with the given name is not present",
		Detail:   "workspace with the given name is not present",
	})
	return diags
	
}
