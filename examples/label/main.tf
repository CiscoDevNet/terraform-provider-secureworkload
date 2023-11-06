provider "secureworkload" {
  api_key                  = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  api_secret               = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  api_url                  = "https://acme.secureworkloadpreview.com/"
  disable_tls_verification = false
}

resource "secureworkload_label" "tag" {
  tenant_name = "acme"
  ip          = "10.0.0.1"
  attributes = {
    Environment = "test"
    Datacenter  = "aws"
    app_name    = "product-service"
  }
}
