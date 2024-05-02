terraform {
  backend "local" {}
  required_providers {
    openvpncloud = {
      source  = "OpenVPN/cloudconnexa"
      version = "0.0.12"
    }
  }
}
