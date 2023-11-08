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
		Description: "Resource for creating application in Secure Workload\n" +
			"\n" +
			"## Example\n" +
			"An example is shown below: \n" +
			"```hcl\n" +
			"resource \"secureworkload_application\" \"application1\" {\n" +
			"	 app_scope_id = data.secureworkload_scope.scope.id\n" +
			"    name = \"Product Service\"\n" +
			"    description = \"Demo description for application\"\n" +
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
			"**Note:** If creating multiple rules during a single `terraform apply`, remember to use `depends_on` to chain the rules so that terraform creates it in the same order that you intended.\n" ,
		Create:        resourceSecureWorkloadApplicationCreate,
		Read:          resourceSecureWorkloadApplicationRead,
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
				ForceNew:    true,
				Computed:    true,
				Description: "(Optional) User-specified name for the application.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "(Optional) User-specified description of the application.",
			},
			"alternate_query_mode": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     true,
				Description: "(Optional) Indicates if “dynamic mode” is used for the application. In dynamic mode, an ADM run creates one or more candidate queries for each cluster. Default value is true.",
			},
			"strict_validation": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "(Optional) Return an error if there are unknown keys/attributes in the uploaded data. Useful for catching misspelled keys. Default value is false.",
			},
			"primary": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     true,
				Description: "(Optional) Set to true to indicate this application is primary for the given scope. Default value is true.",
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
							Description: "Named filter. If more than one filter with the same name exists you must specify consumer_filter_id. Only one of consumer_filter_id, consumer_filter_name or consumer_scope_name can be specified.",
						},
						"consumer_scope_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Named application scope. If more than one application scope with the same name exists you must specify consumer_filter_id. Only one of consumer_filter_id, consumer_filter_name or consumer_scope_name can be specified.",
						},
						"provider_filter_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of a cluster, user inventory filter, or application scope.",
						},
						"provider_filter_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Named filter. If more than one filter with the same name exists you must specify consumer_filter_id. Only one of consumer_filter_id, consumer_filter_name or consumer_scope_name can be specified.",
						},
						"provider_scope_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Named application scope. If more than one application scope with the same name exists you must specify consumer_filter_id. Only one of consumer_filter_id, consumer_filter_name or consumer_scope_name can be specified.",
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
										Description: "Inclusive range of ports; for example, [80, 80] or [5000, 6000].",
										Elem: &schema.Schema{
											Type:     schema.TypeInt,
											MinItems: 2,
											MaxItems: 2,
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
							Description: "Named filter. If more than one filter with the same name exists you must specify consumer_filter_id. Only one of consumer_filter_id, consumer_filter_name or consumer_scope_name can be specified.",
						},
						"consumer_scope_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Named application scope. If more than one application scope with the same name exists you must specify consumer_filter_id. Only one of consumer_filter_id, consumer_filter_name or consumer_scope_name can be specified.",
						},
						"provider_filter_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of a cluster, user inventory filter, or application scope.",
						},
						"provider_filter_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Named filter. If more than one filter with the same name exists you must specify consumer_filter_id. Only one of consumer_filter_id, consumer_filter_name or consumer_scope_name can be specified.",
						},
						"provider_scope_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Named application scope. If more than one application scope with the same name exists you must specify consumer_filter_id. Only one of consumer_filter_id, consumer_filter_name or consumer_scope_name can be specified.",
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
										Description: "Inclusive range of ports; for example, [80, 80] or [5000, 6000].",
										Elem: &schema.Schema{
											Type:     schema.TypeInt,
											MinItems: 2,
											MaxItems: 2,
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
				Required:    true,
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
	if isPrimaryApplication {
		existingApplications, err := client.ListApplications()
		if err != nil {
			return err
		}
		for _, existingApplication := range existingApplications {
			if existingApplication.Primary {
				return errors.New(fmt.Sprintf("Existing application %s exists for scope %s that is marked as primary. Please demote the workspace to secondary before continuing.", existingApplication.Name, existingApplication.AppScopeId))
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
		ScopeName:  tf["consumer_scope_name"].(string),
	}
	filterId, err := policyFilterIdForQuery(apiClient, consumingPolicyFilterQuery)
	if err != nil {
		return policy, err
	}
	policy.ConsumerFilterId = filterId
	providingPolicyFilterQuery := policyFilterQuery{
		AbsoluteId: tf["provider_filter_id"].(string),
		FilterName: tf["provider_filter_name"].(string),
		ScopeName:  tf["provider_scope_name"].(string),
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

func resourceSecureWorkloadApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	describeApplicatioParams := DescribeApplicationRequest{
		ApplicationId: d.Id(),
	}
	application, err := client.DescribeApplication(describeApplicatioParams)
	if err != nil {
		return err
	}
	d.Set("name", application.Name)
	d.Set("description", application.Description)
	d.Set("primary", application.Primary)
	d.Set("alternate_query_mode", application.AlternateQueryMode)
	d.Set("latest_adm_version", application.LatestADMVersion)
	d.Set("enforcement_enabled", application.EnforcementEnabled)
	d.Set("enforced_version", application.EnforcedVersion)
	return nil
}

func resourceSecureWorkloadApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	return client.DeleteApplication(d.Id())
}
