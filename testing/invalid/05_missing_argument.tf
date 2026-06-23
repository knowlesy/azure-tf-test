# EXPECTED ERROR: The argument "resource_group_name" is required, but no definition was found.
terraform { required_providers { mockazurerm = { source = "peterknowles/mockazurerm" } } }
provider "mockazurerm" {}

resource "mockazurerm_linux_virtual_machine" "vm" {
  name                  = "missing-arg-vm"
  location              = "East US"
  size                  = "Standard_B2s"
  admin_username        = "adminuser"
  network_interface_ids = ["/subscriptions/123/resourceGroups/rg/providers/Microsoft.Network/networkInterfaces/nic"]
  # MISSING resource_group_name
}
