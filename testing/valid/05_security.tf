terraform { required_providers { mockazurerm = { source = "peterknowles/mockazurerm" } } }
provider "mockazurerm" {}

resource "mockazurerm_resource_group" "rg" {
  name     = "security-rg"
  location = "East US"
}

resource "mockazurerm_key_vault" "kv" {
  name                = "corekv12345"
  location            = mockazurerm_resource_group.rg.location
  resource_group_name = mockazurerm_resource_group.rg.name
  tenant_id           = "00000000-0000-0000-0000-000000000000"
  sku_name            = "standard"
}

resource "mockazurerm_key_vault_secret" "secret" {
  name         = "database-password"
  value        = "SuperSecret123!"
  key_vault_id = mockazurerm_key_vault.kv.id
}
