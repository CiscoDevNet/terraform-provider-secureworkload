package secureworkload

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSecureWorkloadFilter() *schema.Resource {
	return &schema.Resource{
		Description: "Data source for fetching filter from secure-workload\n\n" +
			"An example is shown below: \n" +
			"```hcl\n" +
			"data \"secureworkload_filter\" \"filters\" {\n" +
			"	workspace_id = \"data.secureworkload_workspace.workspace1.id\"\n" +
			"	name = \"filter1\" \n" +
			"}\n" +
			"```" +
			"**Note:** If filter with the given name is not found, this data source block will respond with the first cluster in the list.\n" +
			"```",
		ReadContext: dataSourceSecureWorkloadFilterRead,
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
				Description: "Exact name of the filter",
			},
		},
	}
}

func dataSourceSecureWorkloadFilterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	var index int

	// Logic for query param

	getUrl := "?"
	if val, ok := d.GetOk("name"); ok {
		temp := strings.Replace(val.(string), " ", "%20", -1)
		getUrl = getUrl + "exact_name=" + temp
	}
	resFilter, err := c.GetFilterByParam(getUrl)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to read filter",
			Detail:   err.Error(),
		})
		return diags
	}
	if len(resFilter) > 0 {
		for i, j := range resFilter {
			if j.Name == d.Get("name").(string) {
				index = i
				break
			}
		}
		d.SetId(resFilter[index].Id)
		if err := d.Set("name", resFilter[index].Name); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "unable to read cluster",
				Detail:   err.Error(),
			})
			return diags
		}
		return diags
	}
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "filter with the given name is not present",
		Detail:   "filter with the given name is not present",
	})
	return diags

}
