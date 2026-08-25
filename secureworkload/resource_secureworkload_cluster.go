package secureworkload

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	// client "github.com/secureworkload-exchange/terraform-go-sdk"
)

// jsonQueryDiffSuppress compares two JSON documents for semantic equality
// (ignoring key ordering / whitespace differences) so that a server-side
// re-serialization of a user-supplied query does not produce a permanent
// diff on plan. Used for query-like attributes (e.g. "query", "short_query")
// that are stored as raw JSON strings but round-trip through a server that
// may reformat them.
func jsonQueryDiffSuppress(k, oldValue, newValue string, d *schema.ResourceData) bool {
	if oldValue == newValue {
		return true
	}
	var oldParsed, newParsed interface{}
	if err := json.Unmarshal([]byte(oldValue), &oldParsed); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(newValue), &newParsed); err != nil {
		return false
	}
	return reflect.DeepEqual(oldParsed, newParsed)
}

func resourceSecureWorkloadCluster() *schema.Resource {
	return &schema.Resource{
		Description: "Resource for creating a new Cluster on Secure Workload\n" +
			"\n" +
			"## Example\n" +
			"An example is shown below: \n" +
			"```hcl\n" +
			"resource \"secureworkload_cluster\" \"cluster\" {\n" +
			"	 workspace_id = data.secureworkload_workspace.workspace.id\n" +
			"    name = \"New-Cluster\"\n" +
			"    description = \"New Cluster via TF\"\n" +
			"    query = <<EOF\n" +
			"                {" +
			"        		 \"type\":\"eq\",\n" +
			"        		 \"field\": \"ip\",\n" +
			"        		 \"value\": \"10.0.0.1\"\n" +
			"        		 }\n" +
			"        	EOF\n" +
			"    approved = false \n" +
			"}\n" +
			"resource \"secureworkload_cluster\" \"cluster2\" {\n" +
			"    depends_on = [secureworkload_cluster.cluster] \n" +
			"	 workspace_id = data.secureworkload_workspace.workspace.id\n" +
			"    name = \"New-Cluster2\"\n" +
			"    description = \"Second Cluster via TF\"\n" +
			"    query = <<EOF\n" +
			"                {" +
			"        		 \"type\":\"subnet\",\n" +
			"        		 \"field\": \"ip\",\n" +
			"        		 \"value\": \"10.0.2.0/24\"\n" +
			"        		 }\n" +
			"        	EOF\n" +
			"    approved = false \n" +
			"}\n" +
			"```\n" +
			"**Note:** If creating multiple clusters during a single `terraform apply`, remember to use `depends_on` to chain the filters so that terraform creates them in a specific order to avoid *429:too_many_request* error.\n",
		Create: resourceSecureWorkloadClusterCreate,
		Update: nil,
		Read:   resourceSecureWorkloadClusterRead,
		Delete: resourceSecureWorkloadClusterDelete,

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "User-specified name for the inventory cluster.",
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Indicates the version of the workspace the cluster will be added to.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "(optional) The description of the cluster.",
			},
			"query": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: jsonQueryDiffSuppress,
				Description:      "JSON object representation of an inventory filter query. *type* is operator, *field* is label key & *value* is label value. Operator can any of the following: [and, or, eq, subnet, contains, regex, gt, gte, lt, lte, in, range, ranges, not, all, none]",
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the workspace associated with the cluster.",
			},
			"approved": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "(optional) An approved cluster will not be updated during an automatic policy discovery run. Default false.",
			},
		},
	}
}

var requiredCreateClusterParams = []string{"name", "workspace_id", "query"}

func resourceSecureWorkloadClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	for _, param := range requiredCreateClusterParams {
		if d.Get(param) == "" {
			return fmt.Errorf("%s is required but was not provided", param)
		}
	}
	createClusterParams := CreateClusterRequest{
		Name:        d.Get("name").(string),
		Version:     d.Get("version").(string),
		Description: d.Get("description").(string),
		Query:       []byte(d.Get("query").(string)),
		Approved:    d.Get("approved").(bool),
	}
	cluster, err := client.CreateCluster(createClusterParams, d.Get("workspace_id").(string))
	if err != nil {
		return err
	}
	d.SetId(cluster.Id)
	return nil
}

func resourceSecureWorkloadClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	cluster, err := client.DescribeCluster(d.Id())
	if err != nil {
		if IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}
	d.Set("name", cluster.Name)
	d.Set("version", cluster.Version)
	d.Set("description", cluster.Description)
	d.Set("approved", cluster.Approved)
	// Refresh from short_query, NOT query. The API returns both: "query" is
	// the effective query, expanded with the parent scope's filters (it has
	// a vrf_id clause injected and is wrapped in an "and"), while
	// "short_query" is the query the user actually supplied. Setting the
	// expanded form here can never match the configuration, and because
	// query is ForceNew that means a destroy/recreate on every plan.
	// If the API did not return short_query, leave the stored value alone
	// rather than overwriting it with something that cannot round-trip.
	if len(cluster.ShortQueryJSON) > 0 {
		d.Set("query", string(cluster.ShortQueryJSON))
	}
	// NOTE: workspace_id is intentionally not refreshed here. The
	// Clusters struct (secureworkload_cluster.go) returned by
	// DescribeCluster does not carry a workspace/application id field,
	// so there is nothing to set it from without inventing an API
	// response field that doesn't exist.
	return nil
}

func resourceSecureWorkloadClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	return client.DeleteCluster(d.Get("workspace_id").(string), d.Id())
}
