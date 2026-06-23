# EXPECTED ERROR: Unsupported argument: An argument named "invalid_variable" is not expected here.
terraform { required_providers { mockazurerm = { source = "peterknowles/mockazurerm" } } }
provider "mockazurerm" {}

module "network" {
  source           = "./module"
  invalid_variable = "this will fail"
}
