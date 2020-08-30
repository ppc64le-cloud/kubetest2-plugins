data "ibm_resource_group" "rg" {
  name = var.powervs_resource_group
}

data "ibm_resource_instance" "test-pdns-instance" {
  name = var.powervs_dns
  resource_group_id = data.ibm_resource_group.rg.id
}

data "ibm_dns_zones" "ds_pdnszone" {
  instance_id = data.ibm_resource_instance.test-pdns-instance.guid
}

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
  pi_dns = [ "9.9.9.9", "8.8.8.8"]
}

module "master" {
  source = "./instance"

  ibmcloud_api_key = var.powervs_api_key
  image_name = var.powervs_image_name
  memory = var.powervs_memory
  networks = [var.powervs_network_name == "" ? ibm_pi_network.public_network[0].pi_network_name : data.ibm_pi_network.existing_net[0].pi_network_name]
  powervs_service_instance_id = var.powervs_service_id
  processors = var.powervs_processors
  ssh_key_name = var.powervs_ssh_key
  vm_name = "${var.cluster_name}-master"
  user_data = base64encode(templatefile("${path.module}/user_data.tmpl",{port=var.apiserver_port, extra_domain="${var.cluster_name}-master.${var.powervs_dns_zone}", release_marker=var.release_marker, build_version=var.build_version}))
  ibmcloud_region = var.powervs_region
  ibmcloud_zone = var.powervs_zone
}

module "workers" {
  source = "./instance"
  instance_count = var.workers_count

  ibmcloud_api_key = var.powervs_api_key
  image_name = var.powervs_image_name
  memory = var.powervs_memory
  networks = [var.powervs_network_name == "" ? ibm_pi_network.public_network[0].pi_network_name : data.ibm_pi_network.existing_net[0].pi_network_name]
  powervs_service_instance_id = var.powervs_service_id
  processors = var.powervs_processors
  ssh_key_name = var.powervs_ssh_key
  vm_name = "${var.cluster_name}-worker"
  user_data = base64encode(templatefile("${path.module}/user_data.tmpl",{port=var.apiserver_port, extra_domain="${var.cluster_name}-master.${var.powervs_dns_zone}", release_marker=var.release_marker, build_version=var.build_version}))
  ibmcloud_region = var.powervs_region
  ibmcloud_zone = var.powervs_zone
}

resource "ibm_dns_resource_record" "test-pdns-resource-record-a" {
  instance_id = data.ibm_resource_instance.test-pdns-instance.guid
  zone_id     = local.zoneids[var.powervs_dns_zone]
  type        = "A"
  name        = "${var.cluster_name}-master"
  rdata       = module.master.addresses[0][0].external_ip
  ttl         = 3600
}

resource "null_resource" "wait-for-master-completes" {
  connection {
    type = "ssh"
    user = "root"
    host = module.master.addresses[0][0].external_ip
    private_key = file(var.ssh_private_key)
    timeout = "15m"
  }
  provisioner "remote-exec" {
    inline = [
      "cloud-init status -w"
    ]
  }
}

resource "null_resource" "kubeadm-init" {
  depends_on = [null_resource.wait-for-master-completes]
  connection {
    type = "ssh"
    user = "root"
    host = module.master.addresses[0][0].external_ip
    private_key = file(var.ssh_private_key)
    timeout = "15m"
  }
  provisioner "remote-exec" {
    inline = [
      "kubeadm init --apiserver-bind-port=${var.apiserver_port} --apiserver-cert-extra-sans ${var.cluster_name}-master.${var.powervs_dns_zone} --pod-network-cidr=172.20.0.0/16 --kubernetes-version ${var.release_marker} --token ${var.bootstrap_token}",
      "curl https://docs.projectcalico.org/manifests/calico.yaml -O",
      "sed -i 's/veth_mtu\\:.*/veth_mtu: \"8940\"/' calico.yaml",
      "KUBECONFIG=/etc/kubernetes/admin.conf kubectl create -f calico.yaml"
    ]
  }
  provisioner "local-exec" {
    command = "scp -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -i ${var.ssh_private_key} root@${module.master.addresses[0][0].external_ip}:/etc/kubernetes/admin.conf ${var.kubeconfig_path}"
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

resource "null_resource" "kubeadm-join-worker" {
  depends_on = [null_resource.wait-for-workers-completes, null_resource.kubeadm-init]
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
      "kubeadm join ${module.master.addresses[0][0].ip}:${var.apiserver_port} --token ${var.bootstrap_token} --discovery-token-unsafe-skip-ca-verification"
    ]
  }
}


locals {
  zoneids = zipmap(data.ibm_dns_zones.ds_pdnszone.dns_zones[*]["name"], data.ibm_dns_zones.ds_pdnszone.dns_zones[*]["zone_id"])
}
