PLUGIN_NAME := terraform-provider-tutorial_v0.1.0_x4

.PHONY: install
install:
	rm -rf test/${PLUGIN_NAME}
	go build -o test/${PLUGIN_NAME} main.go
