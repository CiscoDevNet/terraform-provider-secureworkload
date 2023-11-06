package main
import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	// "github.com/secureworkload-exchange/terraform-provider/secureworkload"
	"terraform-provider-secureworkload/secureworkload"
)
func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return secureworkload.Provider()
		},
	})
}
// Package secureworkload provides a golang SDK interface
// for managing a SecureWorkload installation via the the SecureWorkload HTTP API.
