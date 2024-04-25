package secureworkload

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	// "github.com/hashicorp/terraform/terraform"
	// client "github.com/secureworkload-exchange/terraform-go-sdk"
)

// Provider returns a terraform resource provider for managing secureworkload resources.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SECUREWORKLOAD_API_KEY", nil),
				Description: "API key for calculating request signatures for SecureWorkload API calls.",
			},
			"api_secret": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SECUREWORKLOAD_API_SECRET", nil),
				Description: "API secret for calculating request signatures for SecureWorkload API calls.",
			},
			"api_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SECUREWORKLOAD_API_URL", nil),
				Description: "URL for a SecureWorkload API.",
			},
			"disable_tls_verification": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SECUREWORKLOAD_DISABLE_TLS_VERIFICATION", false),
				Description: "Allow connections to SecureWorkload endpoints without validating their TLS certificate.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"secureworkload_filter":    resourceSecureWorkloadFilter(),
			"secureworkload_scope":     resourceSecureWorkloadScope(),
			"secureworkload_label":     resourceSecureWorkloadLabel(),
			"secureworkload_user":      resourceSecureWorkloadUser(),
			"secureworkload_workspace": resourceSecureWorkloadApplication(),
			"secureworkload_role":      resourceSecureWorkloadRole(),
			"secureworkload_cluster":   resourceSecureWorkloadCluster(),
			"secureworkload_policies":  resourceSecureWorkloadPolicy(),
			"secureworkload_port":      resourceSecureWorkloadPort(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"secureworkload_scope":     dataSourceSecureWorkloadScope(),
			"secureworkload_workspace": dataSourceSecureWorkloadApplication(),
			"secureworkload_role":      dataSourceSecureWorkloadRole(),
		},
		ConfigureFunc: configureClient,
	}
}

func configureClient(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		APIKey:                 d.Get("api_key").(string),
		APISecret:              d.Get("api_secret").(string),
		APIURL:                 d.Get("api_url").(string),
		DisableTLSVerification: d.Get("disable_tls_verification").(bool),
	}
	if err := validate(config); err != nil {
		return nil, err
	}
	client, err := New(config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// validate validates the config needed to initialize a secureworkload client,
// returning a single error with all validation errors, or nil if no error.
func validate(config Config) error {
	var err *multierror.Error
	if config.APIKey == "" {
		err = multierror.Append(err, fmt.Errorf("API Key must be configured for the Secure Workload provider"))
	}
	if config.APISecret == "" {
		err = multierror.Append(err, fmt.Errorf("API Secret must be configured for the Secure Workload provider"))
	}
	if config.APIURL == "" {
		err = multierror.Append(err, fmt.Errorf("API URL must be configured for the Secure Workload provider"))
	}
	return err.ErrorOrNil()
}
