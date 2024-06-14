package secureworkload

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSecureWorkloadCluster() *schema.Resource {
	return &schema.Resource{
		Description: "Data source for fetching cluster from a workspace\n\n" +
			"An example is shown below: \n" +
			"```hcl\n" +
			"data \"secureworkload_cluster\" \"clusters\" {\n" +
			"	workspace_id = \"data.secureworkload_workspace.workspace1.id\"\n" +
			"	name = \"cluster1\" \n" +
			"}\n" +
			"**Note:** If cluster with the given name is not found, this data source block will respond with the first cluster in the list.\n",
		ReadContext: dataSourceSecureWorkloadClusterRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of this resource",
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the workspace where this cluster is located",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the cluster you need to find",
			},
		},
	}
}

func dataSourceSecureWorkloadClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	var index int
	// Logic for query param
	workspaceId := d.Get("workspace_id").(string)
	name := d.Get("name").(string)
	getUrl := workspaceId + "/clusters"

	resCluster, err := c.GetClusterByParam(getUrl, name)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to get the cluster",
			Detail:   err.Error(),
		})
		return diags
	}
	if len(resCluster) > 0 {
		for i, j := range resCluster {
			if j.Name == d.Get("name").(string) {
				index = i
				break
			}
		}
		d.SetId(resCluster[index].Id)
		if err := d.Set("name", resCluster[index].Name); err != nil {
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
		Summary:  "cluster with the given name is not present",
		Detail:   "cluster with the given name is not present",
	})
	return diags

}
