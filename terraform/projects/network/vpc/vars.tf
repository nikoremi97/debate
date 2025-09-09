variable "cidr_block" {
  type = string
}

variable "enable_dns_support" {
  default = true
}

variable "enable_dns_hostnames" {
  default = true
}

variable "tags" {
  type = map(string)
}

variable "public_subnets" {
  type = object({
    endpoint   = list(string)
    zone       = list(string)
    cidr_block = list(string)
  })
  default = {
    endpoint   = []
    zone       = []
    cidr_block = []
  }
}

variable "private_subnets" {
  type = object({
    endpoint   = list(string)
    zone       = list(string)
    cidr_block = list(string)
  })
  default = {
    endpoint   = []
    zone       = []
    cidr_block = []
  }
}

variable "protocol" {
  type = map(any)
  default = {
    ALL = "-1"
  }
}

variable "traffic_type" {
  type    = string
  default = "REJECT"
}

variable "use_eip" {
  type    = bool
  default = true
}

variable "map_public_ip_on_launch" {
  type    = bool
  default = true
}
