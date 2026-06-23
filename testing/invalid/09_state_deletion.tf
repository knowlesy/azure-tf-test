# EXPECTED ERROR: Delete Blocked: Resource is still referenced by Network Interface in Mock DB.
terraform { required_providers { mockazurerm = { source = "peterknowles/mockazurerm" } } }
provider "mockazurerm" {}

# Pretend a user is trying to apply a destruction of a subnet that still has NICs attached in the backend
resource "mockazurerm_subnet" "sub" {
  name                 = "in-use-subnet"
  resource_group_name  = "example-rg"
  virtual_network_name = "example-vnet"
  address_prefixes     = ["10.0.1.0/24"]
  # If a user runs terraform destroy on this while mockazurerm_network_interface exists in state, it blocks.
}
