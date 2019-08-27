build:
	gox -osarch="linux/amd64 windows/amd64 darwin/amd64" \
	-output="pkg/{{.OS}}_{{.Arch}}/{{.OS}}-{{.Arch}}-terraform-provider-twilio" .

test:
	go test -v $(shell go list ./... | grep -v /vendor/) 

testacc:
	TF_ACC=1 go test -v ./plugin/providers/twilio -run="TestAcc"

install: clean build
	cp pkg/linux_amd64/linux-amd64-terraform-provider-twilio ~/.terraform.d/plugins/terraform-provider-twilio

tfplan: install
	terraform init -upgrade && terraform plan	

tfplandebug: install
	TF_LOG=debug DEBUG=true terraform init -upgrade && terraform plan

tfapply: install
	terraform init -upgrade && terraform apply

tfapplydebug: install
	TF_LOG=debug DEBUG=true terraform init -upgrade && terraform apply

release: release_bump release_build

release_bump:
	scripts/release_bump.sh

release_build:
	scripts/release_build.sh

deps:
	dep ensure -vendor-only
	
clean:
	rm -rf pkg/
