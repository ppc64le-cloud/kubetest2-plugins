variable "powervs_resource_group" {
  description = "IBM Cloud Resource Group"
}

variable "powervs_api_key" {
  description = "Denotes the IBM Cloud API key to use"
}

variable "powervs_dns" {
  description = "IBM Cloud VPC Private DNS name"
}

variable "powervs_dns_zone" {
  description = "IBM Cloud VPC Private DNS Zone Name"
  default = "k8s.test"
}

variable "powervs_image_name" {
  description = "PowerVS image name to be used for the deployment"
}

variable "powervs_memory" {
  description = "PowerVS memory in GB"
}

variable "powervs_processors" {
  description = "PowerVS processor units"
}

variable "powervs_network_name" {
  description = "PowerVS Network name to be used for the deployment"
}

variable "powervs_service_id" {
  description = "PowerVS service ID"
}

variable "powervs_ssh_key" {
  description = "PowerVS SSH Key ID"
}

variable "powervs_region" {
  description = "PowerVS Region"
}

variable "powervs_zone" {
  description = "PowerVS Zone"
}

variable "apiserver_port" {
  description = "Kubernetes API Server Port"
}

