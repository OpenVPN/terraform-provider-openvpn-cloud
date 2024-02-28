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
    "Default"   = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
    "Developer" = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
    "Support"   = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  }
}

variable "networks" {
  type = map(string)
  default = {
    "example-network" = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
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
