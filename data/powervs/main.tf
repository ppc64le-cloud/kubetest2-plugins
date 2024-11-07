data "ibm_pi_network" "existing_net" {
  count                = var.powervs_network_name == "" ? 0 : 1
  pi_network_name      = var.powervs_network_name
  pi_cloud_instance_id = var.powervs_service_id
}

resource "ibm_pi_network" "public_network" {
  count                   = var.powervs_network_name == "" ? 1 : 0
  pi_network_name           = "${var.cluster_name}-pub-net"
  pi_cloud_instance_id      = var.powervs_service_id
  pi_network_type           = "pub-vlan"
  pi_dns = [ "8.8.4.4", "8.8.8.8"]
}

data "ibm_pi_catalog_images" "catalog_images" {
  pi_cloud_instance_id = var.powervs_service_id
}
data "ibm_pi_images" "project_images" {
  pi_cloud_instance_id = var.powervs_service_id
}

locals {
  catalog_power_image = [for x in data.ibm_pi_catalog_images.catalog_images.images : x if x.name == var.powervs_image_name]
  project_power_image = [for x in data.ibm_pi_images.project_images.image_info : x if x.name == var.powervs_image_name]
  invalid_power_image = length(local.project_power_image) == 0 && length(local.catalog_power_image) == 0
  # If invalid then use name to fail in ibm_pi_instance resource; else if not found in project then import using ibm_pi_image; else use the power image id
  power_image_id = (
    local.invalid_power_image ? var.powervs_image_name : (
      length(local.project_power_image) == 0 ? ibm_pi_image.power[0].image_id : local.project_power_image[0].id
    )
  )
}

# Copy image from catalog if not in the project and present in catalog
resource "ibm_pi_image" "power" {
  count                = length(local.project_power_image) == 0 && length(local.catalog_power_image) == 1 ? 1 : 0
  pi_image_name        = var.powervs_image_name
  pi_image_id          = local.catalog_power_image[0].image_id
  pi_cloud_instance_id = var.powervs_service_id
}

module "master" {
  depends_on = [ibm_pi_image.power]
  source = "./instance"

  ibmcloud_api_key = var.powervs_api_key
  image_name = var.powervs_image_name
  memory = var.controlplane_powervs_memory
  network = var.powervs_network_name == "" ? ibm_pi_network.public_network[0].network_id : data.ibm_pi_network.existing_net[0].id
  powervs_service_instance_id = var.powervs_service_id
  processors = var.controlplane_powervs_processors
  ssh_key_name = var.powervs_ssh_key
  vm_name = "${var.cluster_name}-master"
  ibmcloud_region = var.powervs_region
  ibmcloud_zone = var.powervs_zone
}

module "workers" {
  depends_on = [ibm_pi_image.power]
  source = "./instance"
  instance_count = var.workers_count

  ibmcloud_api_key = var.powervs_api_key
  image_name = var.powervs_image_name
  memory = var.powervs_memory
  network = var.powervs_network_name == "" ? ibm_pi_network.public_network[0].network_id : data.ibm_pi_network.existing_net[0].id
  powervs_service_instance_id = var.powervs_service_id
  processors = var.powervs_processors
  ssh_key_name = var.powervs_ssh_key
  vm_name = "${var.cluster_name}-worker"
  ibmcloud_region = var.powervs_region
  ibmcloud_zone = var.powervs_zone
}

resource "null_resource" "wait-for-master-completes" {
  connection {
    type = "ssh"
    user = "root"
    host = module.master.addresses[0][0].external_ip
    private_key = file(var.ssh_private_key)
    timeout = "20m"
  }
  provisioner "remote-exec" {
    inline = [
      "cloud-init status -w"
    ]
  }
}

resource "null_resource" "wait-for-workers-completes" {
  count = var.workers_count
  connection {
    type = "ssh"
    user = "root"
    host = module.workers.addresses[count.index][0].external_ip
    private_key = file(var.ssh_private_key)
    timeout = "15m"
  }
  provisioner "remote-exec" {
    inline = [
      "cloud-init status -w"
    ]
  }
}
