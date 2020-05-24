# Run with "CGO_ENABLED=0 GOOS=linux" in the environment for Docker.
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

VERSION=$(shell sed 's/.*"version": "\(.*\)".*/\1/;t;d' package.json)

all: sazdump sazserve

sazdump: cmd/sazdump/version.go $(wildcard cmd/sazdump/*.go pkg/dumper/*.go pkg/parser/*.go pkg/analyzer/*.go)
	cd cmd/sazdump && go build $(GOFLAGS) -o ../../sazdump

sazserve: cmd/sazserve/version.go $(ASSET_BIN) $(wildcard cmd/sazserve/*.go pkg/parser/*.go pkg/analyzer/*.go internal/cache/*.go)
	cd cmd/sazserve && go build $(GOFLAGS) -o ../../sazserve

$(ASSET_BIN): $(ASSET_DIR)/js/index.min.js $(wildcard $(ASSET_DIR)/* $(ASSET_DIR)/*/*)
	$(BINDATA) -fs -o $@ -prefix $(ASSET_DIR) $(ASSET_DIR)/...
	go run _tools/move-generated-comments/move-generated-comments.go -- $(ASSET_BIN)
	gofmt -s -w $@

$(ASSET_DIR)/js/index.min.js: node_modules/datatables.net/js/jquery.dataTables.js.vendor node_modules/datatables.net-buttons/js/dataTables.buttons.js.vendor cmd/sazserve/sources/js/mime-type-icons.js cmd/sazserve/assets/json/help-page.json cmd/sazserve/assets/json/help-table.json $(wildcard $(SOURCE_DIR)/*/*)
	mkdir -p $(ASSET_DIR)/css $(ASSET_DIR)/js
	$(ESBUILD) --outfile=$(ASSET_DIR)/js/index.min.js --format=iife --sourcemap \
		--bundle --minify cmd/sazserve/sources/js/index.js
	sed -i '1s/^\xEF\xBB\xBF//' node_modules/chardin.js/chardinjs.css
	$(MINIFY) -o $(ASSET_DIR)/css/index.min.css \
		node_modules/datatables.net-bs4/css/dataTables.bootstrap4.css \
		node_modules/datatables.net-buttons-bs4/css/buttons.bootstrap4.css \
		node_modules/datatables.net-colreorder-bs4/css/colReorder.bootstrap4.css \
		node_modules/datatables.net-fixedheader-bs4/css/fixedHeader.bootstrap4.css \
		node_modules/chardin.js/chardinjs.css $(SOURCE_DIR)/css/index.css
	$(MINIFY) -o $(ASSET_DIR)/css/bootstrap.flatly.min.css $(SOURCE_DIR)/css/bootstrap.flatly.css
	$(MINIFY) -o $(ASSET_DIR)/css/bootstrap.darkly.min.css $(SOURCE_DIR)/css/bootstrap.darkly.css
	$(MINIFY) -o $(ASSET_DIR)/css/overrides.darkly.min.css $(SOURCE_DIR)/css/overrides.darkly.css
	$(MINIFY) -o $(ASSET_DIR)/index.html $(SOURCE_DIR)/index.html

cmd/sazserve/assets/json/help-page.json: cmd/sazserve/sources/yml/help-page.yml
	mkdir -p $(ASSET_DIR)/json
	./node_modules/.bin/js-yaml -c $? > $@

cmd/sazserve/assets/json/help-table.json: cmd/sazserve/sources/yml/help-table.yml
	mkdir -p $(ASSET_DIR)/json
	./node_modules/.bin/js-yaml -c $? > $@

cmd/sazserve/sources/js/mime-type-icons.js: $(wildcard $(ASSET_DIR)/png/*)
	go run _tools/list-known-mime-types/list-known-mime-types.go -- $@

node_modules/datatables.net/js/jquery.dataTables.js.vendor: cmd/sazserve/sources/js/jquery.dataTables.js.diff
ifeq (,$(wildcard node_modules/datatables.net/js/jquery.dataTables.js.vendor))
	cp node_modules/datatables.net/js/jquery.dataTables.js $@
	patch -p0 < $?
endif

node_modules/datatables.net-buttons/js/dataTables.buttons.js.vendor: cmd/sazserve/sources/js/dataTables.buttons.js.diff
ifeq (,$(wildcard node_modules/datatables.net-buttons/js/dataTables.buttons.js.vendor))
	cp node_modules/datatables.net-buttons/js/dataTables.buttons.js $@
	patch -p0 < $?
endif

restore-datatables :: node_modules/datatables.net/js/jquery.dataTables.js node_modules/datatables.net-buttons/js/dataTables.buttons.js
ifneq (,$(wildcard node_modules/datatables.net/js/jquery.dataTables.js.vendor))
	mv node_modules/datatables.net/js/jquery.dataTables.js.vendor node_modules/datatables.net/js/jquery.dataTables.js
endif
ifneq (,$(wildcard node_modules/datatables.net-buttons/js/dataTables.buttons.js.vendor))
	mv node_modules/datatables.net-buttons/js/dataTables.buttons.js.vendor node_modules/datatables.net-buttons/js/dataTables.buttons.js
endif

cmd/sazdump/version.go: package.json
	echo "package main\n\nconst version = \"$(VERSION)\"" > $@

cmd/sazserve/version.go: package.json
	echo "package main\n\nconst version = \"$(VERSION)\"" > $@

generate ::
ifeq (,$(wildcard $(ASSET_BIN)))
	$(BINDATA) -fs -o $(ASSET_BIN) -prefix $(ASSET_DIR) $(ASSET_DIR)/...
	go run _tools/move-generated-comments/move-generated-comments.go -- $(ASSET_BIN)
endif

lint ::
	./node_modules/.bin/standard --verbose --fix scripts/*.js cmd/sazserve/sources/js/*.js
	golangci-lint run _tools/list-known-mime-types _tools/move-generated-comments \
		cmd/... pkg/... internal/...

run-dump :: cmd/sazdump/version.go $(wildcard cmd/sazdump/*.go pkg/dumper/*.go pkg/parser/*.go pkg/analyzer/*.go)
	go run cmd/sazdump/sazdump.go cmd/sazdump/version.go "$(SAZ)"

run-serve :: cmd/sazserve/version.go $(ASSET_BIN) $(wildcard cmd/sazserve/*.go pkg/parser/*.go pkg/analyzer/*.go internal/cache/*.go)
	go run cmd/sazserve/sazserve.go cmd/sazserve/api.go cmd/sazserve/version.go $(ASSET_BIN)

debug-serve :: debug-data cmd/sazserve/version.go $(wildcard cmd/sazserve/*.go pkg/parser/*.go pkg/analyzer/*.go internal/cache/*.go)
	go run cmd/sazserve/sazserve.go cmd/sazserve/api.go cmd/sazserve/version.go $(ASSET_BIN)

debug-data :: debug-assets $(wildcard $(ASSET_DIR)/* $(ASSET_DIR)/*/*)
	go-bindata -debug -fs -o $(ASSET_BIN) -prefix $(ASSET_DIR) $(ASSET_DIR)/...

debug-assets :: node_modules/datatables.net/js/jquery.dataTables.js.vendor node_modules/datatables.net-buttons/js/dataTables.buttons.js.vendor cmd/sazserve/sources/js/mime-type-icons.js cmd/sazserve/assets/json/help-page.json cmd/sazserve/assets/json/help-table.json $(wildcard $(SOURCE_DIR)/*/*)
	mkdir -p $(ASSET_DIR)/css $(ASSET_DIR)/js
	$(ESBUILD) --outfile=$(ASSET_DIR)/js/index.min.js --format=iife --sourcemap \
		--bundle cmd/sazserve/sources/js/index.js
	sed -i '1s/^\xEF\xBB\xBF//' node_modules/chardin.js/chardinjs.css
	cat node_modules/datatables.net-bs4/css/dataTables.bootstrap4.css \
		node_modules/datatables.net-buttons-bs4/css/buttons.bootstrap4.css \
		node_modules/datatables.net-colreorder-bs4/css/colReorder.bootstrap4.css \
		node_modules/datatables.net-fixedheader-bs4/css/fixedHeader.bootstrap4.css \
		node_modules/chardin.js/chardinjs.css $(SOURCE_DIR)/css/index.css \
		> $(ASSET_DIR)/css/index.min.css
	cp $(SOURCE_DIR)/css/bootstrap.flatly.css $(ASSET_DIR)/css/bootstrap.flatly.min.css
	cp $(SOURCE_DIR)/css/bootstrap.darkly.css $(ASSET_DIR)/css/bootstrap.darkly.min.css
	cp $(SOURCE_DIR)/css/overrides.darkly.css $(ASSET_DIR)/css/overrides.darkly.min.css
	cp $(SOURCE_DIR)/index.html $(ASSET_DIR)/index.html

prepare :: npm-prepare go-prepare

npm-prepare ::
	npm ci

go-prepare ::
	go get -u github.com/go-bindata/go-bindata/v3/...
	go get -u github.com/evanw/esbuild/...
	go get -u github.com/tdewolff/minify/v2/...

clean ::
	rm -rf sazdump sazserve $(ASSET_BIN) $(ASSET_DIR)/css $(ASSET_DIR)/js \
		$(ASSET_DIR)/json cmd/sazserve/sources/js/mime-type-icons.js \
		cmd/sazdump/version.go cmd/sazserve/version.go $(ASSET_DIR)/index.html dist

push ::
	git push heroku master
	git push
	git push --tags

publish ::
	GITLAB_TOKEN= goreleaser --rm-dist
	npm publish
	cd ../homebrew-tap && git pull && make push-info
	cd ../scoop-bucket && git pull

snapshot ::
	GITLAB_TOKEN= goreleaser --skip-validate --skip-publish --snapshot --rm-dist

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
