PLUGIN_NAME := terraform-provider-tutorial

.PHONY: install
install:
	rm -rf test/${PLUGIN_NAME}
	go build -o test/${PLUGIN_NAME} main.go
