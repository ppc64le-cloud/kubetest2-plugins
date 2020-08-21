output "addresses" {
    value = ibm_pi_instance.pvminstance.*.addresses
}
