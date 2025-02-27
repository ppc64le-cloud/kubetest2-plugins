data "ibm_resource_group" "default_group" {
  name = var.vpc_resource_group
}

data "ibm_is_image" "node_image" {
  name = var.node_image
}

data "ibm_is_ssh_key" "ssh_key" {
  name = var.vpc_ssh_key
}

module "vpc" {
  source         = "./vpc-instance"
  vpc_name       = var.vpc_name != "" ? var.vpc_name : "${var.cluster_name}-vpc"
  cluster_name   = var.cluster_name
  zone           = var.vpc_zone
  resource_group = data.ibm_resource_group.default_group.id
}

locals {
  vpc_id            = module.vpc.vpc_id
  subnet_id         = module.vpc.subnet_id
  security_group_id = module.vpc.security_group_id
}

resource "ibm_is_instance_template" "node_template" {
  name           = "${var.cluster_name}-node-template"
  image          = data.ibm_is_image.node_image.id
  profile        = var.node_profile
  vpc            = local.vpc_id
  zone           = var.vpc_zone
  resource_group = data.ibm_resource_group.default_group.id
  keys           = [data.ibm_is_ssh_key.ssh_key.id]

  primary_network_interface {
    subnet          = local.subnet_id
    security_groups = [local.security_group_id]
  }
}

module "master" {
  source                    = "./node"
  node_name                 = "${var.cluster_name}-master"
  node_instance_template_id = ibm_is_instance_template.node_template.id
  resource_group            = data.ibm_resource_group.default_group.id
}

module "workers" {
  source                    = "./node"
  count                     = var.workers_count
  node_name                 = "${var.cluster_name}-worker-${count.index}"
  node_instance_template_id = ibm_is_instance_template.node_template.id
  resource_group            = data.ibm_resource_group.default_group.id
}

resource "null_resource" "wait-for-master-completes" {
  connection {
    type        = "ssh"
    user        = "root"
    host        = module.master.public_ip
    private_key = file(var.ssh_private_key)
    timeout     = "20m"
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
    type        = "ssh"
    user        = "root"
    host        = module.workers[count.index].public_ip
    private_key = file(var.ssh_private_key)
    timeout     = "15m"
  }
  provisioner "remote-exec" {
    inline = [
      "cloud-init status -w"
    ]
  }
}
