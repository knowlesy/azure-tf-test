terraform { required_providers { mockazurerm = { source = "peterknowles/mockazurerm" } } }
provider "mockazurerm" {}

resource "mockazurerm_resource_group" "rg" {
  name     = "routing-rg"
  location = "East US"
}

resource "mockazurerm_route_table" "rt" {
  name                = "core-rt"
  location            = mockazurerm_resource_group.rg.location
  resource_group_name = mockazurerm_resource_group.rg.name
}

resource "mockazurerm_route" "route" {
  name                   = "to-firewall"
  resource_group_name    = mockazurerm_resource_group.rg.name
  route_table_name       = mockazurerm_route_table.rt.name
  address_prefix         = "0.0.0.0/0"
  next_hop_type          = "VirtualAppliance"
  next_hop_in_ip_address = "10.0.1.4"
}
