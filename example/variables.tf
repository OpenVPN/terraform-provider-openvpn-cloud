variable "company_name" {
  type        = string
  description = "Company name in CloudConnexa"
  # default = ""
}

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
    "Username1" = {
      username = "Username1"
      email    = "username1@company.com"
      group    = "Default"
      role     = "ADMIN"
    }
    "Username2" = {
      username = "Username2"
      email    = "username2@company.com"
      group    = "Developer"
      role     = "MEMBER"
    }
    "Username3" = {
      username = "Username3"
      email    = "username3@company.com"
      group    = "Support"
      role     = "MEMBER"
    }
  }
}

variable "groups" {
  type = map(string)
  default = {
    "Default"   = "11111111-1111-1111-1111-111111111111"
    "Developer" = "22222222-1111-1111-1111-111111111111"
    "Support"   = "33333333-1111-1111-1111-111111111111"
  }
}

variable "networks" {
  type = map(string)
  default = {
    "example-network" = "11111111-2222-3333-4444-555555555555"
  }
}

variable "routes" {
  type = list(map(string))
  default = [
    {
      value       = "10.0.0.0/18"
      description = "Example Route with subnet /18"
    },
    {
      value       = "10.10.0.0/20"
      description = "Example Route with subnet /20"
    },
    {
      value       = "10.20.0.0/24"
      description = "Example Route with subnet /24"
    },
  ]
}
