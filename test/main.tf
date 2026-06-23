terraform {
  required_providers {
    mockazurerm = {
      source  = "peterknowles/mockazurerm"
      version = "1.0.0"
    }
  }
}

provider "mockazurerm" {
  subscription_id = "00000000-0000-0000-0000-000000000000"
}

resource "mockazurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
  tags = {
    environment = "dev"
  }
}

resource "mockazurerm_virtual_network" "example" {
  name                = "example-network"
  resource_group_name = mockazurerm_resource_group.example.name
  location            = mockazurerm_resource_group.example.location
  address_space       = ["10.0.0.0/16"]
  dns_servers         = ["10.0.0.4", "10.0.0.5"]
  tags = {
    environment = "dev"
  }
}

resource "mockazurerm_subnet" "example" {
  name                 = "internal"
  resource_group_name  = mockazurerm_resource_group.example.name
  virtual_network_name = mockazurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "mockazurerm_network_interface" "example" {
  name                = "example-nic"
  location            = mockazurerm_resource_group.example.location
  resource_group_name = mockazurerm_resource_group.example.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = mockazurerm_subnet.example.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "mockazurerm_linux_virtual_machine" "example" {
  name                  = "example-vm"
  resource_group_name   = mockazurerm_resource_group.example.name
  location              = mockazurerm_resource_group.example.location
  size                  = "Standard_B1s"
  admin_username        = "adminuser"
  network_interface_ids = [mockazurerm_network_interface.example.id]

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

resource "mockazurerm_storage_account" "example" {
  name                     = "examplestoracc"
  resource_group_name      = mockazurerm_resource_group.example.name
  location                 = mockazurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  network_rules {
    default_action = "Deny"
  }
}

resource "mockazurerm_servicebus_namespace" "example" {
  name                = "example-sb-ns"
  location            = mockazurerm_resource_group.example.location
  resource_group_name = mockazurerm_resource_group.example.name
  sku                 = "Standard"
}

resource "mockazurerm_user_assigned_identity" "example" {
  name                = "example-identity"
  location            = mockazurerm_resource_group.example.location
  resource_group_name = mockazurerm_resource_group.example.name
}
