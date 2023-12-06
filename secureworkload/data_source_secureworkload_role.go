package secureworkload

import (
	"strings"
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSecureWorkloadRole() *schema.Resource {
	return &schema.Resource{
		Description: "Data source for fetching roles of a root scope from secure-workload\n\n" +
			"An example is shown below: \n" +
			"```hcl\n" +
			"data \"secureworkload_scope\" \"scope\" {\n" +
			"	exact_name = \"RootScope:ChildScope\"\n" +
			"}\n" +
			"data \"secureworkload_role\" \"role\" {\n" +
			"	app_scope_id = data.secureworkload_scope.scope.id\n" +
			"}\n" +
			"```",
		ReadContext: dataSourceSecureWorkloadRoleRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the role",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of this role",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description(if-any) of the role",
			},
			"app_scope_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the root app scope", 
			},
		},
	}
}

func dataSourceSecureWorkloadRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// Logic for query param

	getUrl := "?"
	if val, ok := d.GetOk("app_scope_id"); ok {
		temp := strings.Replace(val.(string), " ", "%20", -1)
		getUrl = getUrl + "app_scope_id=" +  temp
	}
	resRole, err := c.GetRoleByParam( getUrl)
	if (err != nil) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to get the role",
			Detail:   err.Error(),
		})
		return diags
	}
	if len(resRole) > 0 {
		d.SetId(resRole[0].Id)

		if err := d.Set("app_scope_id", resRole[0].AppScopeId); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "unable to read role",
				Detail:   err.Error(),
			})
			return diags
		}

		if err := d.Set("name", resRole[0].Name); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "unable to read role",
				Detail:   err.Error(),
			})
			return diags
		}
		if err := d.Set("description", resRole[0].Description); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "unable to read role",
				Detail:   err.Error(),
			})
			return diags
		}

		return diags
	}
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "role with the given root scope id is not present",
		Detail:   "role with the given root scope id is not present",
	})
	return diags
	
}
