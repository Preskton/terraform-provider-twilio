.PHONY: build plugin clean
DIRS = ~/.terraform.d/plugins
$(shell mkdir -p $(DIRS) )

build:
	go build

plugin: build
	mv terraform-provider-twilio ~/.terraform.d/plugins

clean:
	rm -rf terraform-provider-twilio
