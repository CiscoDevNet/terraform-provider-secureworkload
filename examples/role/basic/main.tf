provider "secureworkload" {
  api_key                  = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  api_secret               = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  api_url                  = "https://acme.secureworkloadpreview.com"
  disable_tls_verification = false
}

resource "secureworkload_role" "read_role" {
  name                = "read_role"
  app_scope_id        = "5ce71503497d4f2c23af85b7"
  access_app_scope_id = "5ceea87b497d4f753baf85bc"
  access_type         = "scope_read"
  user_ids            = ["5eab4dd8497d4f2bec5c585f"]
  description         = "role which provides read-only access to role_your_own_application"
}
