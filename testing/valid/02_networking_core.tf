terraform { required_providers { mockazurerm = { source = "peterknowles/mockazurerm" } } }
provider "mockazurerm" {}

resource "mockazurerm_resource_group" "rg" {
  name     = "networking-rg"
  location = "East US"
}

resource "mockazurerm_virtual_network" "vnet" {
  name                = "core-vnet"
  resource_group_name = mockazurerm_resource_group.rg.name
  location            = mockazurerm_resource_group.rg.location
  address_space       = ["10.0.0.0/16"]
}

resource "mockazurerm_subnet" "sub1" {
  name                 = "subnet1"
  resource_group_name  = mockazurerm_resource_group.rg.name
  virtual_network_name = mockazurerm_virtual_network.vnet.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "mockazurerm_subnet" "sub2" {
  name                 = "subnet2"
  resource_group_name  = mockazurerm_resource_group.rg.name
  virtual_network_name = mockazurerm_virtual_network.vnet.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "mockazurerm_subnet" "sub3" {
  name                 = "subnet3"
  resource_group_name  = mockazurerm_resource_group.rg.name
  virtual_network_name = mockazurerm_virtual_network.vnet.name
  address_prefixes     = ["10.0.3.0/24"]
}
