provider "ibm" {
  version = "1.10.0"
  ibmcloud_api_key = var.powervs_api_key
  region = var.powervs_region
  zone = var.powervs_zone
}
