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
               "type": "and",
               "filters": [
                  {
                     "field": "vrf_id",
                     "type": "eq",
                     "value": 700056
                  },
                  {
                     "type": "or",
                     "filters": [
                        {
                           "field": "ip",
                           "type": "eq",
                           "value": "10.254.252.43"
                        },
                        {
                           "field": "ip",
                           "type": "eq",
                           "value": "10.254.252.51"
                        },
                        {
                           "field": "ip",
                           "type": "eq",
                           "value": "10.254.252.52"
                        }
                     ]
                  }
               ]
            }
          EOF
  app_scope_id = data.secureworkload_scope.scope.id
  primary      = true
  public       = false
}
