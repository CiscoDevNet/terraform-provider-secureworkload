package secureworkload

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	// client "github.com/secureworkload-exchange/terraform-go-sdk"
	// secureworkload "github.com/secureworkload-exchange/terraform-go-sdk"
)

func resourceSecureWorkloadApplication() *schema.Resource {
	return &schema.Resource{
		Description: "Resource for creating a workspace in Secure Workload\n" +
			"\n" +
			"## Example\n" +
			"An example is shown below: \n" +
			"```hcl\n" +
			"resource \"secureworkload_workspace\" \"workspace1\" {\n" +
			"	 app_scope_id = data.secureworkload_scope.scope.id\n" +
			"    name = \"Product Service\"\n" +
			"    description = \"Demo description for workspace\"\n" +
			"    alternate_query_mode = true\n" +
			"    strict_validation = true\n" +
			"    primary = true \n" +
			"    cluster {\n" +
			"	 	 id = <ID_OF_Cluster>\n" +
			"    	 name = <NAME_OF_Cluster>\n" +
			"    	 description = <Cluster_Description>\n" +
			"        node {\n" +
			"            ip_address = \"1.2.3.4\"\n" +
			"        	 name = \"Product Service\"\n" +
			"        }\n" +
			"	 }\n" +
			"    filter {\n" +
			"	 	 id = <ID_OF_Cluster>\n" +
			"    	 name = <NAME_OF_Cluster>\n" +
			"    	 query = <<EOF\n" +
			"                {" +
			"        		 \"type\":\"eq\",\n" +
			"        		 \"field\": \"ip\",\n" +
			"        		 \"value\": \"10.0.0.1\"\n" +
			"        		 }\n" +
			"        		 EOF\n" +
			"	 }\n" +
			"    absolute_policy {\n" +
			"	 	 consumer_filter_id = <CONSUMER_FILTER_ID>\n" +
			"    	 provider_filter_id = <PROVIDER_FILTER_ID>\n" +
			"    	 action = \"ALLOW\"\n" +
			"        layer_4_network_policy {\n" +
			"            port_range = [80,80]\n" +
			"        	 protocol = 6\n" +
			"        }\n" +
			"	 }\n" +
			"    default_policy {\n" +
			"	 	 consumer_filter_id = <CONSUMER_FILTER_ID>\n" +
			"    	 provider_filter_id = <PROVIDER_FILTER_ID>\n" +
			"    	 action = \"DENY\"\n" +
			"        layer_4_network_policy {\n" +
			"            port_range = [80,80]\n" +
			"        	 protocol = 6\n" +
			"        }\n" +
			"	 }\n" +
			"    catch all action  = false \n" +
			"}\n" +
			"```\n" +
			"**Note:** If creating multiple resources for workspaces during a single `terraform apply`, you may have to use `depends_on` to chain the resources so that terraform creates it in the same order that you intended.\n",
		Create:        resourceSecureWorkloadApplicationCreate,
		Read:          resourceSecureWorkloadApplicationRead,
		Update:        resourceSecureWorkloadApplicationUpdate,
		Delete:        resourceSecureWorkloadApplicationDelete,
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"app_scope_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the scope assigned to the application.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "(Optional) User-specified name for the application. Updated in place (no replacement) via PUT /applications/{id}.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "(Optional) User-specified description of the application. Updated in place (no replacement) via PUT /applications/{id}.",
			},
			"alternate_query_mode": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "(Optional) Indicates if “dynamic mode” is used for the application. In dynamic mode, an ADM run creates one or more candidate queries for each cluster. Default value is true.",
			},
			"strict_validation": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "(Optional) Return an error if there are unknown keys/attributes in the uploaded data. Useful for catching misspelled keys. Default value is false.",
			},
			"primary": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "(Optional) Set to true to indicate this application is primary for the given scope. Default value is true. Updated in place (no replacement) via PUT /applications/{id}. The API enforces one primary application per scope and returns an error if this would conflict with an existing primary workspace.",
			},
			"cluster": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Cluster wraps a groups of nodes to be used to define policies.",
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Unique identifier to be used with policies.",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: " Cluster display name.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description of the cluster.",
						},
						"node": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Node represents an endpoint that is part of a cluster",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_address": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "IP address or subnet of the node; for example, 10.0.0.1/8 or 1.2.3.4.",
									},
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Displayed name of the node.",
									},
								},
							},
						},
						"consistent_uuid": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Must be unique to a given application. After an ADM run, the similar/same clusters in the next version will maintain the consistent_uuid.",
						},
					},
				},
			},
			"filter": {
				Type:        schema.TypeList,
				ForceNew:    true,
				Optional:    true,
				Description: "Filter wrap a collection of inventory filters on data center assets used to define an                application policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Unique identifier to be used with policies.",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Displayed name of the cluster.",
						},
						"query": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "JSON object representation of an inventory filter query.",
						},
					},
				},
			},
			"absolute_policy": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Ordered application policy to be created with the absolute rank.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"consumer_filter_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of a cluster, user inventory filter, or application scope.",
						},
						"consumer_filter_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Named filter. If more than one filter with the same name exists you must specify consumer_filter_id. ",
						},
						"provider_filter_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of a cluster, user inventory filter, or application scope.",
						},
						"provider_filter_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Named filter. If more than one filter with the same name exists you must specify consumer_filter_id. ",
						},
						"action": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "“ALLOW” or “DENY”",
						},
						"layer_4_network_policy": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Parameters for enforcing a layer 4 networking policy based off a flows                            protocol and ports.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Protocol integer value (NULL means all protocols).",
									},
									"port_range": {
										Type:        schema.TypeList,
										Required:    true,
										MinItems:    2,
										MaxItems:    2,
										Description: "Inclusive range of ports; for example, [80, 80] or [5000, 6000].",
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
									"approved": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "(Optional) Indicates whether the policy is approved. Default is false.",
									},
								},
							},
						},
					},
				},
			},
			"default_policy": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Ordered application policy to be created with the default rank.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"consumer_filter_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of a cluster, user inventory filter, or application scope.",
						},
						"consumer_filter_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Named filter. If more than one filter with the same name exists you must specify consumer_filter_id.",
						},
						"provider_filter_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of a cluster, user inventory filter, or application scope.",
						},
						"provider_filter_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Named filter. If more than one filter with the same name exists you must specify consumer_filter_id.",
						},
						"action": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "“ALLOW” or “DENY”",
						},
						"layer_4_network_policy": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Parameters for enforcing a layer 4 networking policy based off a flows protocol and ports.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:        schema.TypeInt,
										Optional:    true,
										Default:     nil,
										Description: "Protocol integer value (NULL means all protocols).",
									},
									"port_range": {
										Type:        schema.TypeList,
										Required:    true,
										MinItems:    2,
										MaxItems:    2,
										Description: "Inclusive range of ports; for example, [80, 80] or [5000, 6000].",
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
									"approved": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: "(Optional) Indicates whether the policy is approved. Default is false.",
									},
								},
							},
						},
					},
				},
			},
			"catch_all_action": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "DENY",
				ForceNew:    true,
				Description: "“ALLOW” or “DENY”",
			},
			"author": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "First and last name of the user who created the application.",
			},
			"created_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Unix timestamp indicating when the application was created.",
			},
			"latest_adm_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The latest adm (v*) version of the application.",
			},
			"enforcement_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if enforcement is enabled on the application.",
			},
			"enforced_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The enforced p* version of the application.",
			},
		},
	}
}

func resourceSecureWorkloadApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	isPrimaryApplication := d.Get("primary").(bool)
	tempAppScopeId := d.Get("app_scope_id").(string)
	if isPrimaryApplication {
		existingApplications, err := client.ListApplications(tempAppScopeId)
		if err != nil {
			return err
		}
		for _, existingApplication := range existingApplications {
			if existingApplication.Primary {
				return errors.New(fmt.Sprintf("Existing application '' %s '' exists for scope '' %s '' that is marked as primary. Please demote the workspace to secondary before continuing.", existingApplication.Name, existingApplication.AppScopeId))
			}
		}
	}
	createApplicationParams := CreateApplicationRequest{
		AppScopeId:         d.Get("app_scope_id").(string),
		Name:               d.Get("name").(string),
		Description:        d.Get("description").(string),
		AlternateQueryMode: d.Get("alternate_query_mode").(bool),
		StrictValidation:   d.Get("strict_validation").(bool),
		Primary:            isPrimaryApplication,
		CatchAllAction:     d.Get("catch_all_action").(string),
	}
	if value, ok := d.GetOk("cluster"); ok {
		var clusters []Cluster
		tfClusters := value.([]interface{})
		for _, tfCluster := range tfClusters {
			if tfCluster == nil {
				continue
			}
			cluster, err := clusterFromTerraform(tfCluster.(terraformObject))
			if err != nil {
				return err
			}
			clusters = append(clusters, cluster)
		}
		createApplicationParams.Clusters = clusters
	}
	if value, ok := d.GetOk("filter"); ok {
		var filters []PolicyFilter
		tfFilters := value.([]interface{})
		for _, tfFilter := range tfFilters {
			if tfFilter == nil {
				continue
			}
			filter, err := filterFromTerraform(tfFilter.(terraformObject))
			if err != nil {
				return err
			}
			filters = append(filters, filter)
		}
		createApplicationParams.Filters = filters
	}
	if value, ok := d.GetOk("absolute_policy"); ok {
		var absolutePolicies []Policy
		tfAbsolutePolicies := value.([]interface{})
		for _, tfAbsolutePolicy := range tfAbsolutePolicies {
			if tfAbsolutePolicy == nil {
				continue
			}
			abosolutePolicy, err := policyFromTerraform(client, tfAbsolutePolicy.(terraformObject))
			if err != nil {
				return err
			}
			absolutePolicies = append(absolutePolicies, abosolutePolicy)
		}
		createApplicationParams.AbsolutePolicies = absolutePolicies
	}
	if value, ok := d.GetOk("default_policy"); ok {
		var defaultPolicies []Policy
		tfDefaultPolicies := value.([]interface{})
		for _, tfDefaultPolicy := range tfDefaultPolicies {
			if tfDefaultPolicy == nil {
				continue
			}
			abosolutePolicy, err := policyFromTerraform(client, tfDefaultPolicy.(terraformObject))
			if err != nil {
				return err
			}
			defaultPolicies = append(defaultPolicies, abosolutePolicy)
		}
		createApplicationParams.DefaultPolicies = defaultPolicies
	}
	application, err := client.CreateApplication(createApplicationParams)
	if err != nil {
		return err
	}
	d.Set("author", application.Author)
	d.Set("created_at", application.CreatedAt)
	d.Set("latest_adm_version", application.LatestADMVersion)
	d.Set("enforcement_enabled", application.EnforcementEnabled)
	d.Set("enforced_version", application.EnforcedVersion)
	d.SetId(application.Id)
	return nil
}

type terraformObject = map[string]interface{}

func clusterFromTerraform(tf terraformObject) (Cluster, error) {
	cluster := Cluster{}
	cluster.Id = tf["id"].(string)
	cluster.Name = tf["name"].(string)
	cluster.Description = tf["description"].(string)
	if value := tf["node"]; len(value.([]interface{})) > 0 {
		nodes := []Node{}
		tfNodes := value.([]interface{})
		for _, tfNode := range tfNodes {
			if tfNode == nil {
				continue
			}
			nodes = append(nodes, nodeFromTerraform(tfNode.(terraformObject)))
		}
		cluster.Nodes = nodes
	}
	cluster.ConsistentUUID = tf["consistent_uuid"].(string)
	return cluster, nil
}

func nodeFromTerraform(tf terraformObject) Node {
	return Node{
		IPAddress: tf["ip_address"].(string),
		Name:      tf["name"].(string),
	}
}

func filterFromTerraform(tf terraformObject) (PolicyFilter, error) {
	return PolicyFilter{
		Id:    tf["id"].(string),
		Name:  tf["name"].(string),
		Query: []byte(tf["query"].(string)),
	}, nil
}

type policyFilterQuery struct {
	AbsoluteId string
	FilterName string
	ScopeName  string
}

func policyFilterIdForQuery(apiClient Client, query policyFilterQuery) (string, error) {
	if query.AbsoluteId == "" && query.FilterName == "" && query.ScopeName == "" {
		return "", errors.New("One  of policy filter id, filter name or scope name must be specified")
	}
	if query.AbsoluteId != "" && (query.FilterName != "" || query.ScopeName != "") {
		return "", errors.New("Only one of policy filter id, filter name or scope name can be specified")
	}
	if query.FilterName != "" && query.ScopeName != "" {
		return "", errors.New("Only one of policy filter id, filter name or scope name can be specified")
	}
	if query.AbsoluteId != "" {
		return query.AbsoluteId, nil
	}
	var secureworkloadPolicyFilterId string
	if query.FilterName != "" {
		inventoryFilters, err := apiClient.ListFilters()
		if err != nil {
			return "", err
		}
		var filtersWithMatchingName []Filter
		for _, inventoryFilter := range inventoryFilters {
			if inventoryFilter.Name == query.FilterName {
				filtersWithMatchingName = append(filtersWithMatchingName, inventoryFilter)
			}
		}
		if len(filtersWithMatchingName) > 1 {
			return "", errors.New(fmt.Sprintf("More than one filter exists with name %s, please use policy filter id to specify the exact one to use.", query.FilterName))
		}
		secureworkloadPolicyFilterId = filtersWithMatchingName[0].Id
	}
	if query.ScopeName != "" {
		scopes, err := apiClient.ListScopes()
		if err != nil {
			return "", err
		}
		var scopesWithMatchingName []Scope
		for _, scope := range scopes {
			if scope.ShortName == query.ScopeName {
				scopesWithMatchingName = append(scopesWithMatchingName, scope)
			}
		}
		if len(scopesWithMatchingName) > 1 {
			return "", errors.New(fmt.Sprintf("More than one scope exists with name %s, please use policy filter id to specify the exact one to use.", query.ScopeName))
		}
		secureworkloadPolicyFilterId = scopesWithMatchingName[0].Id
	}
	return secureworkloadPolicyFilterId, nil
}

func policyFromTerraform(apiClient Client, tf terraformObject) (Policy, error) {
	policy := Policy{}
	// Allow users to specify a consumer or provider filter via
	// absolute id OR scope name OR filter name
	// returning an error if either more than one scope/filter exists
	// with the same name or if both an absolute id and name was provided
	consumingPolicyFilterQuery := policyFilterQuery{
		AbsoluteId: tf["consumer_filter_id"].(string),
		FilterName: tf["consumer_filter_name"].(string),
	}
	filterId, err := policyFilterIdForQuery(apiClient, consumingPolicyFilterQuery)
	if err != nil {
		return policy, err
	}
	policy.ConsumerFilterId = filterId
	providingPolicyFilterQuery := policyFilterQuery{
		AbsoluteId: tf["provider_filter_id"].(string),
		FilterName: tf["provider_filter_name"].(string),
	}
	filterId, err = policyFilterIdForQuery(apiClient, providingPolicyFilterQuery)
	if err != nil {
		return policy, err
	}
	policy.ProviderFilterId = filterId
	policy.Action = tf["action"].(string)
	if value := tf["layer_4_network_policy"]; len(value.([]interface{})) > 0 {
		layer4NetworkPolicies := []Layer4NetworkPolicy{}
		tfLayer4NetworkPolicies := value.([]interface{})
		for _, tfLayer4NetworkPolicy := range tfLayer4NetworkPolicies {
			if tfLayer4NetworkPolicy == nil {
				continue
			}
			layer4NetworkPolicies = append(layer4NetworkPolicies, layer4NetworkPolicyFromTerraform(tfLayer4NetworkPolicy.(terraformObject)))
		}
		policy.Layer4NetworkPolicies = layer4NetworkPolicies
	}
	return policy, nil
}

func layer4NetworkPolicyFromTerraform(tf terraformObject) Layer4NetworkPolicy {
	tfPortRange := tf["port_range"].([]interface{})
	return Layer4NetworkPolicy{
		Protocol:  tf["protocol"].(int),
		PortRange: [2]int{tfPortRange[0].(int), tfPortRange[1].(int)},
		Approved:  tf["approved"].(bool),
	}
}

func resourceSecureWorkloadApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	updateApplicationParams := UpdateApplicationRequest{}
	changed := false
	if d.HasChange("name") {
		updateApplicationParams.Name = d.Get("name").(string)
		changed = true
	}
	if d.HasChange("description") {
		updateApplicationParams.Description = d.Get("description").(string)
		changed = true
	}
	if d.HasChange("primary") {
		primary := d.Get("primary").(bool)
		updateApplicationParams.Primary = &primary
		changed = true
	}
	if changed {
		_, err := client.UpdateApplication(updateApplicationParams, d.Id())
		if err != nil {
			// Surface the error unmodified: a 422 here is typically the
			// one-primary-per-scope conflict returned by the API, which is
			// genuinely useful information for the user.
			return err
		}
	}
	return resourceSecureWorkloadApplicationRead(d, meta)
}

func resourceSecureWorkloadApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	describeApplicatioParams := DescribeApplicationRequest{
		ApplicationId: d.Id(),
	}
	application, err := client.DescribeApplication(describeApplicatioParams)
	if err != nil {
		if IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}
	if err := d.Set("name", application.Name); err != nil {
		return err
	}
	if err := d.Set("description", application.Description); err != nil {
		return err
	}
	if err := d.Set("primary", application.Primary); err != nil {
		return err
	}
	// app_scope_id was previously never refreshed here, even though it is
	// part of the API response. Setting it keeps state in sync in case it
	// ever drifts (e.g. imported resources), even though the field itself
	// remains ForceNew since the API silently ignores attempts to change
	// it via PUT.
	if err := d.Set("app_scope_id", application.AppScopeId); err != nil {
		return err
	}
	// alternate_query_mode is intentionally NOT refreshed here. The
	// Application struct declares an `alternate_query_mode` JSON tag, but
	// live GET/PUT responses from the API do not actually include this
	// field (see live-verified sample response in the #36 follow-up
	// investigation). Since the field is absent from the response body,
	// json.Unmarshal leaves it at its zero value (false), and setting it
	// unconditionally here would clobber the user's configured value with
	// `false` on every read/refresh -- a permanent diff on a ForceNew
	// attribute, i.e. exactly the bug class this fix addresses.
	//
	// strict_validation is a create-only request parameter (it only
	// exists on CreateApplicationRequest, not on the Application response
	// type returned by the API), so there is nothing to refresh for it
	// either.
	//
	// cluster / filter / absolute_policy / default_policy /
	// catch_all_action are intentionally NOT refreshed here either. The
	// Application response struct returned by DescribeApplication cannot
	// represent any of these nested structures, and inventing a
	// flattening from some other endpoint would risk introducing new
	// permanent diffs rather than fixing this one. These attributes
	// remain ForceNew, so Terraform will always plan a full replacement
	// if the user changes them, and there's no correctness gap in leaving
	// them un-refreshed on every plain Read.
	if err := d.Set("latest_adm_version", application.LatestADMVersion); err != nil {
		return err
	}
	if err := d.Set("enforcement_enabled", application.EnforcementEnabled); err != nil {
		return err
	}
	if err := d.Set("enforced_version", application.EnforcedVersion); err != nil {
		return err
	}
	return nil
}

func resourceSecureWorkloadApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	return client.DeleteApplication(d.Id())
}
