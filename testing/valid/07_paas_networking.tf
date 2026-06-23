terraform { required_providers { mockazurerm = { source = "peterknowles/mockazurerm" } } }
provider "mockazurerm" {}

resource "mockazurerm_resource_group" "rg" {
  name     = "paas-rg"
  location = "East US"
}

resource "mockazurerm_virtual_network" "vnet" {
  name                = "paas-vnet"
  resource_group_name = mockazurerm_resource_group.rg.name
  location            = mockazurerm_resource_group.rg.location
  address_space       = ["10.1.0.0/16"]
}

resource "mockazurerm_subnet" "sub" {
  name                 = "paas-subnet"
  resource_group_name  = mockazurerm_resource_group.rg.name
  virtual_network_name = mockazurerm_virtual_network.vnet.name
  address_prefixes     = ["10.1.1.0/24"]
}

resource "mockazurerm_storage_account" "sa" {
  name                     = "paassa123"
  resource_group_name      = mockazurerm_resource_group.rg.name
  location                 = mockazurerm_resource_group.rg.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "mockazurerm_private_endpoint" "pe" {
  name                = "paas-pe"
  location            = mockazurerm_resource_group.rg.location
  resource_group_name = mockazurerm_resource_group.rg.name
  subnet_id           = mockazurerm_subnet.sub.id
}
