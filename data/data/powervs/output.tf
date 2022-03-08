output "masters" {
  value = module.master.addresses[*][0].external_ip
  description = "k8s master node IP addresses"
}

output "workers" {
  value = module.workers.addresses[*][0].external_ip
  description = "k8s worker node IP addresses"
}

output "masters_private" {
  value = module.master.addresses[*][0].ip_address
  description = "k8s master nodes private IP addresses"
}

output "workers_private" {
  value = module.workers.addresses[*][0].ip_address
  description = "k8s worker nodes private IP addresses"
}

output "network" {
  value = ibm_pi_network.public_network
  description = "Network used for the deployment"
}
