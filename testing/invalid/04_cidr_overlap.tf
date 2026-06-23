# EXPECTED ERROR: Simulated deployment failure or overlap validation warning (if implemented by mock provider logic).
terraform { required_providers { mockazurerm = { source = "peterknowles/mockazurerm" } } }
provider "mockazurerm" {}

resource "mockazurerm_subnet" "sub1" {
  name                 = "sub1"
  resource_group_name  = "example-rg"
  virtual_network_name = "example-vnet"
  address_prefixes     = ["10.0.1.0/24"]
}

resource "mockazurerm_subnet" "sub2" {
  name                 = "sub2"
  resource_group_name  = "example-rg"
  virtual_network_name = "example-vnet"
  address_prefixes     = ["10.0.1.0/24"] # Overlaps with sub1
}
