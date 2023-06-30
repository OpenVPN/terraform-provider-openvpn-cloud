terraform {
  required_providers {
    openvpn-cloud = {
      version = "0.0.7"
      source  = "openvpncloud.dev/openvpn/openvpncloud"
    }
  }
}

provider "openvpn-cloud" {
  base_url = ""
}

variable "host_name" {
  default = "tf-autotest"
  type    = string
}

resource "openvpncloud_host" "host" {
  name            = "TEST_HOST_NAME"
  description     = "Terraform test description 2"
  internet_access = "LOCAL"

  connector {
    name          = "test"
    vpn_region_id = "us-west-1"
  }

  provider = openvpn-cloud
}

locals {
  connector_profile = [for connector in openvpncloud_host.host.connector : connector.profile][0]
}


output "host_id" {
  value = openvpncloud_host.host.id
}

output "connector_id" {
  value = [for connector in openvpncloud_host.host.connector : connector.id][0]
}
