terraform {
  required_providers {
    ibm = {
      source  = "IBM-Cloud/ibm"
      version = "~> 1.50.0"
    }
  }
}

provider "ibm" {
  ibmcloud_api_key = var.vpc_api_key
  region           = var.vpc_region
}
