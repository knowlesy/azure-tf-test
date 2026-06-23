# EXPECTED ERROR: Inappropriate value for attribute "address_prefixes": list of string required.
terraform { required_providers { mockazurerm = { source = "peterknowles/mockazurerm" } } }
provider "mockazurerm" {}

resource "mockazurerm_subnet" "sub" {
  name                 = "type-sub"
  resource_group_name  = "example-rg"
  virtual_network_name = "example-vnet"
  address_prefixes     = "10.0.1.0/24" # Should be a list
}
