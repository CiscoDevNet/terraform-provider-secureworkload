provider "secureworkload" {
  api_key                  = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  api_secret               = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  api_url                  = "https://acme.secureworkloadpreview.com"
  disable_tls_verification = false
}
data "secureworkload_scope" "scope" {
  exact_name = "RootScope:ChildScope"
}
resource "secureworkload_scope" "scope" {
  short_name          = "Terraform created scope"
  short_query_type    = "eq"
  short_query_field   = "ip"
  short_query_value   = "10.0.0.1"
  parent_app_scope_id = data.secureworkload_scope.scope.id
}
