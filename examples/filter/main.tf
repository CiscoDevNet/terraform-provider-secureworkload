provider "secureworkload" {
  api_key                  = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  api_secret               = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  api_url                  = "https://acme.secureworkloadpreview.com"
  disable_tls_verification = false
}
data "secureworkload_scope" "scope" {
  exact_name = "RootScope:ChildScope"
}
resource "secureworkload_filter" "filter" {
  name         = "Terraform created filter"
  query        = <<EOF
                    {
                      "type": "eq",
                      "field": "ip",
                      "value": "10.0.0.1"
                    }
          EOF
  app_scope_id = data.secureworkload_scope.scope.id
  primary      = true
  public       = false
}
