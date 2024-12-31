output "vpc_id" {
  value = ibm_is_vpc.vpc.id
}

output "subnet_id" {
  value = ibm_is_subnet.primary.id
}

output "security_group_id" {
  value = ibm_is_vpc.vpc.default_security_group
}
