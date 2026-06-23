terraform { required_providers { mockazurerm = { source = "peterknowles/mockazurerm" } } }
provider "mockazurerm" {}

resource "mockazurerm_resource_group" "rg" {
  name     = "messaging-rg"
  location = "East US"
}

resource "mockazurerm_servicebus_namespace" "sb" {
  name                = "core-sb-ns"
  location            = mockazurerm_resource_group.rg.location
  resource_group_name = mockazurerm_resource_group.rg.name
  sku                 = "Standard"
}

resource "mockazurerm_servicebus_queue" "queue" {
  name         = "jobs-queue"
  namespace_id = mockazurerm_servicebus_namespace.sb.id
}
