data "ibm_is_vpc" "vpc" {
  count = var.vpc_name == "" ? 0 : 1
  name  = var.vpc_name
}

data "ibm_is_subnet" "subnet" {
  count = var.vpc_subnet_name == "" ? 0 : 1
  name  = var.vpc_subnet_name
}

data "ibm_resource_group" "default_group" {
  name = var.vpc_resource_group
}

module "vpc" {
  # Create new vpc and subnet only if vpc_name is not set
  count          = var.vpc_name == "" ? 1 : 0
  source         = "./vpc-instance"
  cluster_name   = var.cluster_name
  zone           = var.vpc_zone
  resource_group = data.ibm_resource_group.default_group.id
}

locals {
  vpc_id            = var.vpc_name == "" ? module.vpc[0].vpc_id : data.ibm_is_vpc.vpc[0].id
  subnet_id         = var.vpc_name == "" ? module.vpc[0].subnet_id : data.ibm_is_subnet.subnet[0].id
  security_group_id = var.vpc_name == "" ? module.vpc[0].security_group_id : data.ibm_is_vpc.vpc[0].default_security_group
}

data "ibm_is_image" "node_image" {
  name = var.node_image
}

data "ibm_is_ssh_key" "ssh_key" {
  name = var.vpc_ssh_key
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
