provider "openvpncloud" {
  base_url      = "https://${var.company_name}.api.openvpn.com"
  client_id     = var.client_id
  client_secret = var.client_secret
}

## Use ENV variables:
# export TF_VAR_client_id=''
# export TF_VAR_client_secret=''
