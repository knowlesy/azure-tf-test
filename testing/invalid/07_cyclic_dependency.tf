# EXPECTED ERROR: Cycle: mockazurerm_subnet.sub2, mockazurerm_subnet.sub1
terraform { required_providers { mockazurerm = { source = "peterknowles/mockazurerm" } } }
provider "mockazurerm" {}

resource "mockazurerm_subnet" "sub1" {
  name                 = mockazurerm_subnet.sub2.id # Cycle!
  resource_group_name  = "example-rg"
  virtual_network_name = "example-vnet"
  address_prefixes     = ["10.0.1.0/24"]
}

resource "mockazurerm_subnet" "sub2" {
  name                 = mockazurerm_subnet.sub1.id # Cycle!
  resource_group_name  = "example-rg"
  virtual_network_name = "example-vnet"
  address_prefixes     = ["10.0.2.0/24"]
}
