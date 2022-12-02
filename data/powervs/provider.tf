terraform {
  required_providers {
    ibm = {
      source = "IBM-Cloud/ibm"
      version = ">= 1.12.0"
    }
  }
}

provider "ibm" {
  ibmcloud_api_key = var.powervs_api_key
  region = var.powervs_region
  zone = var.powervs_zone
}
