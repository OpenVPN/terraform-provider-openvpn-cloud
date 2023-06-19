provider "openvpncloud" {
  base_url      = "https://${COMPANY_NAME_IN_OPENVPN_CLOUD}.api.openvpn.com"
  client_id     = COMPANY_NAME_IN_OPENVPN_CLOUD
  client_secret = API_TOKEN
}
