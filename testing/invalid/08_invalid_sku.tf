# EXPECTED ERROR: Expected sku to be one of [Basic Standard Premium], got FakeSKU.
terraform { required_providers { mockazurerm = { source = "peterknowles/mockazurerm" } } }
provider "mockazurerm" {}

resource "mockazurerm_servicebus_namespace" "sb" {
  name                = "invalid-sku-sb"
  location            = "East US"
  resource_group_name = "example-rg"
  sku                 = "FakeSKU" # Fails enum validation
}
