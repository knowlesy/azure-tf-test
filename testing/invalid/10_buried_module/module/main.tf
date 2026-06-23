variable "vnet_name" {
  type    = string
  default = "module-vnet"
}

resource "mockazurerm_virtual_network" "vnet" {
  name                = var.vnet_name
  resource_group_name = "example-rg"
  location            = "East US"
  address_space       = ["10.0.0.0/16"]
}
