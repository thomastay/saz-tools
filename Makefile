# Run with "CGO_ENABLED=0 GOOS=linux" in the enfironment for Docker.
ifeq ($(DOCKER),1)
	export CGO_ENABLED=0
	export GOOS=linux
	GOFLAGS=-a -installsuffix cgo -ldflags '-extldflags "-static"'
endif

# Run on Heroku which does not include GOBIN in PATH.
ifdef GOBIN
	BINDATA="$(GOBIN)/go-bindata"
	ESBUILD="$(GOBIN)/esbuild"
  MINIFY="$(GOBIN)/minify"
else
	BINDATA=go-bindata
	ESBUILD=esbuild
  MINIFY=minify
endif

SOURCE_DIR=cmd/sazserve/sources
ASSET_DIR=cmd/sazserve/assets
ASSET_BIN=cmd/sazserve/assets.go

all: sazdump sazserve

sazdump: $(wildcard cmd/sazdump/*.go pkg/dumper/*.go pkg/parser/*.go pkg/analyzer/*.go)
	go build $(GOFLAGS) cmd/sazdump/sazdump.go

sazserve: $(ASSET_BIN) $(wildcard cmd/sazserve/*.go pkg/parser/*.go pkg/analyzer/*.go internal/cache/*.go)
	go build $(GOFLAGS) cmd/sazserve/sazserve.go cmd/sazserve/api.go $(ASSET_BIN)

$(ASSET_BIN): $(ASSET_DIR)/js/saz.min.js $(wildcard $(ASSET_DIR)/* $(ASSET_DIR)/*/*)
	$(BINDATA) -fs -o $(ASSET_BIN) -prefix $(ASSET_DIR) $(ASSET_DIR)/...
	go run _tools/move-generated-comments/move-generated-comments.go -- $(ASSET_BIN)

generate ::
ifeq (,$(wildcard $(ASSET_BIN)))
	$(BINDATA) -fs -o $(ASSET_BIN) -prefix $(ASSET_DIR) $(ASSET_DIR)/...
	go run _tools/move-generated-comments/move-generated-comments.go -- $(ASSET_BIN)
endif

lint ::
	golint ./...

run-dump :: $(wildcard cmd/sazdump/*.go pkg/dumper/*.go pkg/parser/*.go pkg/analyzer/*.go)
	go run cmd/sazdump/sazdump.go "$(SAZ)"

run-serve :: $(ASSET_BIN) $(wildcard cmd/sazserve/*.go pkg/parser/*.go pkg/analyzer/*.go internal/cache/*.go)
	go run cmd/sazserve/sazserve.go cmd/sazserve/api.go $(ASSET_BIN)

debug-serve :: debug-assets $(wildcard cmd/sazserve/*.go pkg/parser/*.go pkg/analyzer/*.go internal/cache/*.go)
	go run cmd/sazserve/sazserve.go cmd/sazserve/api.go $(ASSET_BIN)

debug-assets :: $(ASSET_DIR)/js/saz.min.js $(wildcard $(ASSET_DIR)/* $(ASSET_DIR)/*/*)
	go-bindata -debug -fs -o $(ASSET_BIN) -prefix $(ASSET_DIR) $(ASSET_DIR)/...

$(ASSET_DIR)/js/saz.min.js: node_modules/datatables.net/js/jquery.dataTables.js.vendor $(wildcard $(SOURCE_DIR)/*/*)
	mkdir -p $(ASSET_DIR)/css $(ASSET_DIR)/js
	$(ESBUILD) --outfile=$(ASSET_DIR)/js/saz.min.js --format=iife --sourcemap \
		--bundle --minify cmd/sazserve/sources/js/saz.js
	$(MINIFY) -o $(ASSET_DIR)/css/common.min.css \
		node_modules/datatables.net-bs4/css/dataTables.bootstrap4.css \
		node_modules/datatables.net-buttons-bs4/css/buttons.bootstrap4.css \
		node_modules/datatables.net-colreorder-bs4/css/colReorder.bootstrap4.css \
		node_modules/datatables.net-fixedcolumns-bs4/css/fixedColumns.bootstrap4.css \
		node_modules/datatables.net-scroller-bs4/css/scroller.bootstrap4.css \
		$(SOURCE_DIR)/css/saz.css
	$(MINIFY) -o $(ASSET_DIR)/css/bootstrap.flatly.min.css $(SOURCE_DIR)/css/bootstrap.flatly.css
	$(MINIFY) -o $(ASSET_DIR)/css/bootstrap.darkly.min.css $(SOURCE_DIR)/css/bootstrap.darkly.css
	$(MINIFY) -o $(ASSET_DIR)/css/saz.darkly.min.css $(SOURCE_DIR)/css/saz.darkly.css

prepare ::
	go get -u github.com/go-bindata/go-bindata/v3/...
	go get -u github.com/evanw/esbuild/...
	go get -u github.com/tdewolff/minify/v2/...
	npm ci

node_modules/datatables.net/js/jquery.dataTables.js.vendor: cmd/sazserve/sources/js/jquery.dataTables.js.diff
ifeq (,$(wildcard node_modules/datatables.net/js/jquery.dataTables.js.vendor))
	cp node_modules/datatables.net/js/jquery.dataTables.js \
		node_modules/datatables.net/js/jquery.dataTables.js.vendor
	patch -p0 < cmd/sazserve/sources/js/jquery.dataTables.js.diff
endif

restore-datatables :: node_modules/datatables.net/js/jquery.dataTables.js
ifneq (,$(wildcard node_modules/datatables.net/js/jquery.dataTables.js.vendor))
	mv node_modules/datatables.net/js/jquery.dataTables.js.vendor \
		node_modules/datatables.net/js/jquery.dataTables.js
endif

clean ::
	rm -rf sazdump sazserve $(ASSET_BIN) $(ASSET_DIR)/css $(ASSET_DIR)/js

push ::
	git push && git push heroku master

docker-clean ::
	docker image rm sazdump
	docker image rm sazserve

docker-lint ::
	docker run --rm -i \
		-v ${PWD}/.hadolint.yaml:/bin/hadolint.yaml \
		-e XDG_CONFIG_HOME=/bin hadolint/hadolint \
		< Dockerfile.sazdump
	docker run --rm -i \
		-v ${PWD}/.hadolint.yaml:/bin/hadolint.yaml \
		-e XDG_CONFIG_HOME=/bin hadolint/hadolint \
		< Dockerfile.sazserve

docker-build ::
	docker build -f Dockerfile.sazdump -t sazdump .
	docker build -f Dockerfile.sazserve -t sazserve .

docker-run-help ::
	docker run --rm -it sazdump

docker-dump-example ::
	docker run --rm -it -v ${PWD}:/work -w /work sazdump examples/test.saz

docker-serve-example ::
	docker run --rm -it sazserve

docker-tag ::
	docker tag sazdump prantlf/sazdump:latest
	docker tag sazserve prantlf/sazserve:latest

docker-login ::
	docker login --username=prantlf

docker-push ::
	docker push prantlf/sazdump:latest
	docker push prantlf/sazserve:latest
