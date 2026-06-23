terraform { required_providers { mockazurerm = { source = "peterknowles/mockazurerm" } } }
provider "mockazurerm" {}

resource "mockazurerm_resource_group" "rg" {
  name     = "compute-rg"
  location = "East US"
}

resource "mockazurerm_virtual_network" "vnet" {
  name                = "compute-vnet"
  resource_group_name = mockazurerm_resource_group.rg.name
  location            = mockazurerm_resource_group.rg.location
  address_space       = ["10.0.0.0/16"]
}

resource "mockazurerm_subnet" "sub" {
  name                 = "compute-subnet"
  resource_group_name  = mockazurerm_resource_group.rg.name
  virtual_network_name = mockazurerm_virtual_network.vnet.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "mockazurerm_network_security_group" "nsg" {
  name                = "compute-nsg"
  location            = mockazurerm_resource_group.rg.location
  resource_group_name = mockazurerm_resource_group.rg.name
}

resource "mockazurerm_network_interface" "nic" {
  name                = "compute-nic"
  location            = mockazurerm_resource_group.rg.location
  resource_group_name = mockazurerm_resource_group.rg.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = mockazurerm_subnet.sub.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "mockazurerm_linux_virtual_machine" "vm" {
  name                  = "compute-vm"
  resource_group_name   = mockazurerm_resource_group.rg.name
  location              = mockazurerm_resource_group.rg.location
  size                  = "Standard_B2s"
  admin_username        = "adminuser"
  network_interface_ids = [mockazurerm_network_interface.nic.id]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-LTS"
    version   = "latest"
  }
}
