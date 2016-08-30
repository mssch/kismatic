# Set the build version
ifeq ($(origin VERSION), undefined)
	VERSION := $(shell git rev-parse --short HEAD)
endif

# Setup some useful vars
HOST_GOOS = $(shell go env GOOS)
HOST_GOARCH = $(shell go env GOARCH)
GLIDE_VERSION = v0.11.1

build: vendor
	go build -o bin/kismatic -ldflags "-X main.version=$(VERSION)" ./cmd/kismatic

clean: 
	rm -rf bin
	rm -rf out
	rm -rf vendor
	rm -rf vendor-ansible/out
	rm -rf vendor-cfssl/out

test: 
	go test $(shell ./tools/glide nv) $(TEST_OPTS)

vendor: tools/glide
	./tools/glide install

tools/glide:
	mkdir -p tools
	curl -L https://github.com/Masterminds/glide/releases/download/$(GLIDE_VERSION)/glide-$(GLIDE_VERSION)-$(HOST_GOOS)-$(HOST_GOARCH).tar.gz | tar -xz -C tools
	mv tools/$(HOST_GOOS)-$(HOST_GOARCH)/glide tools/glide
	rm -r tools/$(HOST_GOOS)-$(HOST_GOARCH)

vendor-ansible/out:
	docker build -t apprenda/vendor-ansible -q vendor-ansible 
	docker run --rm -v $(shell pwd)/vendor-ansible/out:/ansible apprenda/vendor-ansible pip install --install-option="--prefix=/ansible" ansible

vendor-cfssl/out:
	mkdir -p vendor-cfssl/out/
	curl -L https://pkg.cfssl.org/R1.2/cfssl_linux-amd64 -o vendor-cfssl/out/cfssl_linux-amd64
	curl -L https://pkg.cfssl.org/R1.2/cfssljson_linux-amd64 -o vendor-cfssl/out/cfssljson_linux-amd64 
	curl -L https://pkg.cfssl.org/R1.2/cfssl_darwin-amd64 -o vendor-cfssl/out/cfssl_darwin-amd64
	curl -L https://pkg.cfssl.org/R1.2/cfssljson_darwin-amd64 -o vendor-cfssl/out/cfssljson_darwin-amd64

dist: vendor-ansible/out vendor-cfssl/out build
	mkdir -p out
	cp bin/kismatic out
	mkdir -p out/ansible
	cp -r vendor-ansible/out/* out/ansible
	cp -r ansible out/ansible/playbooks
	mkdir -p out/cfssl
	cp -r vendor-cfssl/out/* out/cfssl
	rm -f out/kismatic.tar.gz
	tar -cvzf kismatic.tar.gz -C out .
	mv kismatic.tar.gz out