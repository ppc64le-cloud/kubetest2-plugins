variable "vpc_api_key" {
  sensitive = true
}

variable "vpc_resource_group" {
  default = "default"
}

variable "vpc_ssh_key" {}

variable "vpc_name" {
  type        = string
  description = "(optional) Specify existing VPC name. If none is provided, it will create a new VPC named {cluster_name}-vpc"
  default     = ""
}

variable "vpc_subnet_name" {
  type        = string
  description = "(optional) Specify existing subnet name. If none is provided, it will create a new subnet named {cluster_name}-subnet. This must be provided if vpc_name has been set"
  default     = ""
}

variable "node_image" {
  default = "ibm-ubuntu-22-04-2-minimal-s390x-1"
}

variable "node_profile" {
  default = "bz2-2x8"
}

variable "vpc_region" {
  default = "eu-de"
}

variable "vpc_zone" {
  default = "eu-de-1"
}
