.PHONY: build
build:						## Builds for linux, windows, and darwin
	export CGO_ENABLED=0
	gox -osarch="linux/amd64 windows/amd64 darwin/amd64" \
	-output="pkg/{{.OS}}_{{.Arch}}/{{.OS}}-{{.Arch}}-terraform-provider-twilio" .

.PHONY: test
test:						## Run unit tests
	go test -v $(shell go list ./... | grep -v /vendor/) 

.PHONY: testacc
testacc:					## Run acceptance tests
	TF_ACC=1 go test -v ./plugin/providers/twilio -run="TestAcc"

.PHONY: install
install: clean build		## Build and reinstall the latest version of the plugin locally
	cp pkg/linux_amd64/linux-amd64-terraform-provider-twilio ~/.terraform.d/plugins/terraform-provider-twilio

.PHONY: tfplan
tfplan: install				## Build, install, and run terraform plan
	terraform init -upgrade && terraform plan	

.PHONY: tfplandebug
tfplandebug: install		## Build, install, and run terraform plan in debug mode
	TF_LOG=debug DEBUG=true terraform init -upgrade && terraform plan

.PHONY: tfapply
tfapply: install			## Build, install, and run terraform apply
	terraform init -upgrade && terraform apply

.PHONY: tfapplydebug
tfapplydebug: install		## Build, install, and run terraform apply in debug mode
	TF_LOG=debug DEBUG=true terraform init -upgrade && terraform apply

.PHONY: bump-packages
bump-packages:				## Updates dependencies in go.mod to latest
	go get -u ./...

.PHONY: clean
clean:						## Cleans build outputs
	rm -rf pkg/

.PHONY: help
help:           			## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
