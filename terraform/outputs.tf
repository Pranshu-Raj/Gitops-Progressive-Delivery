output "public_ip" {
  value = oci_core_instance.node.public_ip
}

output "ssh" {
  value = "ssh ubuntu@${oci_core_instance.node.public_ip}"
}
