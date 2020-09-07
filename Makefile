PLUGIN_NAME := terraform-provider-tutorial_v0.1.0_x4

.PHONY: install
install:
	rm -rf test/${PLUGIN_NAME} test/.terraform test/terraform.tfstate
	go build -o test/${PLUGIN_NAME} main.go
