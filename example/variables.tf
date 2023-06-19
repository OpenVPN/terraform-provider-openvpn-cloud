variable "users" {
  type = map(
    object({
      username = string
      email    = string
      group    = string
      role     = string
    })
  )
  default = {
    "Denis_Arslanbekov" = {
      username = "Arslanbekov_admin"
      email    = "admin@arslanbekov.com"
      group    = "Default"
      role     = "ADMIN"
    }
    "Vladimir_Kozyrev" = {
      username = "Arslanbekov_developer"
      email    = "developer@arslanbekov.com"
      group    = "Developer"
      role     = "MEMBER"
    }
    "Antonio_Graziano" = {
      username = "Arslanbekov_support"
      email    = "support@arslanbekov.com"
      group    = "Support"
      role     = "MEMBER"
    }
  }
}

variable "groups" {
  type = map(string)
  default = {
    "Default"   = "12312312-1234-1234-1234-123123123123"
    "Developer" = "12312312-1234-1234-1234-123123123123"
    "Support"   = "12312312-1234-1234-1234-123123123123"
  }
}

variable "networks" {
  type = map(string)
  default = {
    "example-network" = "12312312-1234-1234-1234-123123123123"
  }
}

variable "example-terraform_ipv4_routes" {
  type = list(map(string))
  default = [
    {
      value       = "10.0.0.0/24"
      description = "Example route 1"
    },
    {
      value       = "10.10.0.0/24"
      description = "Example route 2"
    },
    {
      value       = "10.20.0.0/24"
      description = "Example route 3"
    },
  ]
}
