data "ibm_pi_network" "power_network" {
    pi_network_name      = var.network
    pi_cloud_instance_id = var.powervs_service_instance_id
}

data "ibm_pi_image" "power_images" {
    pi_image_name        = var.image_name
    pi_cloud_instance_id = var.powervs_service_instance_id
}

resource "ibm_pi_instance" "pvminstance" {
    count = var.instance_count
    pi_memory             = var.memory
    pi_processors         = var.processors
    pi_instance_name      = var.instance_count == 1 ? var.vm_name : "${var.vm_name}-${count.index}"
    pi_proc_type          = var.proc_type
    pi_image_id           = data.ibm_pi_image.power_images.id
    pi_key_pair_name      = var.ssh_key_name
    pi_sys_type           = var.system_type
    pi_storage_type       = var.storage_tier
    pi_cloud_instance_id = var.powervs_service_instance_id
    pi_user_data          = var.user_data
    # Wait for the WARNING state instead of OK state to save some time because we aren't performing any DLPAR operations
    # on this LPARS and later in the flow we also have ssh connectivity check to confirm deployed vms are up and running.
    pi_health_status      = "WARNING"

    pi_network {
      network_id = data.ibm_pi_network.power_network.id
    }
    timeouts {
      create = "30m"
      delete = "30m"
    }
}
