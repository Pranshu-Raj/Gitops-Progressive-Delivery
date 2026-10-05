variable "region" {
  type = string
}

variable "compartment_ocid" {
  type = string
}

variable "ssh_public_key" {
  type = string
}

variable "admin_cidr" {
  type        = string
  description = "Source range allowed to reach SSH and the Kubernetes API"
}

variable "shape" {
  type    = string
  default = "VM.Standard.A1.Flex"
}

variable "ocpus" {
  type    = number
  default = 4
}

variable "memory_gb" {
  type    = number
  default = 12
}

variable "k3s_version" {
  type    = string
  default = "v1.37.1+k3s1"
}
