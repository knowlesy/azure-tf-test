# EXPECTED ERROR: This object has no argument, nested block, or exported attribute named "ip_address".
terraform { required_providers { mockazurerm = { source = "peterknowles/mockazurerm" } } }
provider "mockazurerm" {}

resource "mockazurerm_resource_group" "rg" {
  name     = "attr-rg"
  location = "East US"
}

output "rg_ip" {
  value = mockazurerm_resource_group.rg.ip_address # Resource Groups do not have IPs
}
