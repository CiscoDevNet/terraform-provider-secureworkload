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
  short_query = <<EOF
  {
    "type": "or",
    "filters": [
      {
        "type": "and",
        "filters": [
          {
            "type": "contains",
            "field": "user_orchestrator_system/name",
            "value": "Random"
          },
          {
            "type": "eq",
            "field": "ip",
            "value": "10.0.1.1"
          }
        ]
      },
      {
        "type": "gt",
        "field": "host_tags_cvss3",
        "value": 2
      }
    ]
  }
  EOF
  sub_type = "DNS_SERVERS"
  parent_app_scope_id = data.secureworkload_scope.scope.id
}

