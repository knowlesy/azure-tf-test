terraform { required_providers { mockazurerm = { source = "peterknowles/mockazurerm" } } }
provider "mockazurerm" {}

resource "mockazurerm_resource_group" "rg" {
  name     = "identity-rg"
  location = "East US"
}

resource "mockazurerm_user_assigned_identity" "uami" {
  name                = "core-uami"
  location            = mockazurerm_resource_group.rg.location
  resource_group_name = mockazurerm_resource_group.rg.name
}

resource "mockazurerm_role_assignment" "role" {
  scope                = mockazurerm_resource_group.rg.id
  role_definition_id   = "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/acdd72a7-3385-48ef-bd42-f606fba81ae7"
  principal_id         = mockazurerm_user_assigned_identity.uami.principal_id
}
