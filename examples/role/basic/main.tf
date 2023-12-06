provider "secureworkload" {
  api_key                  = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  api_secret               = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  api_url                  = "https://acme.secureworkloadpreview.com"
  disable_tls_verification = false
}
data "secureworkload_scope" "scope" {
  exact_name = "ScopeName1"
}
data "secureworkload_scope" "scope2" {
  exact_name = "ScopeName2"
}
resource "secureworkload_role" "read_role" {
  name                = "read_role"
  app_scope_id        = data.secureworkload_scope.scope.id
  access_app_scope_id = data.secureworkload_scope.scope2.id
  access_type         = "scope_read"
  user_ids            = ["5eab4dd8497d4f2bec5c585f"]
  description         = "role which provides read-only access to role_your_own_application"
}
