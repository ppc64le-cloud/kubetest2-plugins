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

module "master" {
  source = "./instance"

  ibmcloud_api_key = var.powervs_api_key
  image_name = var.powervs_image_name
  memory = var.controlplane_powervs_memory
  network = var.powervs_network_name == "" ? ibm_pi_network.public_network[0].network_id : data.ibm_pi_network.existing_net[0].id
  powervs_service_instance_id = var.powervs_service_id
  processors = var.controlplane_powervs_processors
  ssh_key_name = var.powervs_ssh_key
  storage_tier = var.powervs_storage_tier
  vm_name = "${var.cluster_name}-master"
  ibmcloud_region = var.powervs_region
  ibmcloud_zone = var.powervs_zone
}

module "workers" {
  source = "./instance"
  instance_count = var.workers_count

  ibmcloud_api_key = var.powervs_api_key
  image_name = var.powervs_image_name
  memory = var.powervs_memory
  network = var.powervs_network_name == "" ? ibm_pi_network.public_network[0].network_id : data.ibm_pi_network.existing_net[0].id
  powervs_service_instance_id = var.powervs_service_id
  processors = var.powervs_processors
  ssh_key_name = var.powervs_ssh_key
  storage_tier = var.powervs_storage_tier
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
