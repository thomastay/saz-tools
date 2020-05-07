# Run with "CGO_ENABLED=0 GOOS=linux" in the enfironment for Docker.
ifeq ($(DOCKER),1)
	export CGO_ENABLED=0
	export GOOS=linux
	GOFLAGS=-a -installsuffix cgo -ldflags '-extldflags "-static"'
endif

# Run on Heroku which does not include GOBIN in PATH.
ifdef GOBIN
	BINDATA="$(GOBIN)/go-bindata"
  MINIFY="$(GOBIN)/minify"
else
	BINDATA=go-bindata
  MINIFY=minify
endif

SOURCE_DIR=cmd/sazserve/sources
ASSET_DIR=cmd/sazserve/assets
ASSET_BIN=cmd/sazserve/assets.go

all: sazdump sazserve

sazdump: $(wildcard cmd/sazdump/*.go pkg/dumper/*.go pkg/parser/*.go pkg/analyzer/*.go)
	go build $(GOFLAGS) cmd/sazdump/sazdump.go

sazserve: $(ASSET_BIN) $(wildcard cmd/sazserve/*.go pkg/parser/*.go pkg/analyzer/*.go)
	go build $(GOFLAGS) cmd/sazserve/sazserve.go $(ASSET_BIN)

$(ASSET_BIN): $(ASSET_DIR)/js/all.min.js $(wildcard $(ASSET_DIR)/* $(ASSET_DIR)/*/*)
	$(BINDATA) -fs -o $(ASSET_BIN) -prefix $(ASSET_DIR) $(ASSET_DIR)/...
	go run internal/move-generated-comments/move-generated-comments.go -- $(ASSET_BIN)

run-dump :: $(wildcard cmd/sazdump/*.go pkg/dumper/*.go pkg/parser/*.go pkg/analyzer/*.go)
	go run cmd/sazdump/sazdump.go "$(SAZ)"

run-serve :: $(ASSET_BIN) $(wildcard cmd/sazserve/*.go pkg/parser/*.go pkg/analyzer/*.go)
	go run cmd/sazserve/sazserve.go $(ASSET_BIN)

debug-serve :: debug-assets $(wildcard cmd/sazserve/*.go pkg/parser/*.go pkg/analyzer/*.go)
	go run cmd/sazserve/sazserve.go $(ASSET_BIN)

debug-assets :: concatenate
	go-bindata -debug -fs -o $(ASSET_BIN) -prefix $(ASSET_DIR) $(ASSET_DIR)/...

lint ::
	golint ./...

generate :: $(ASSET_BIN)

concatenate :: $(wildcard $(SOURCE_DIR)/*/*)
	mkdir -p $(ASSET_DIR)/css $(ASSET_DIR)/js
	cat $(SOURCE_DIR)/js/jquery.js $(SOURCE_DIR)/js/bootstrap.bundle.js \
		$(SOURCE_DIR)/js/datatables.js $(SOURCE_DIR)/js/saz.js > $(ASSET_DIR)/js/all.min.js
	cat  $(SOURCE_DIR)/css/datatables.css $(SOURCE_DIR)/css/saz.css \
		> $(ASSET_DIR)/css/common.min.css
	cp $(SOURCE_DIR)/css/bootstrap.flatly.css $(ASSET_DIR)/css/bootstrap.flatly.min.css
	cp $(SOURCE_DIR)/css/bootstrap.darkly.css $(ASSET_DIR)/css/bootstrap.darkly.min.css
	cp $(SOURCE_DIR)/css/saz.darkly.css $(ASSET_DIR)/css/saz.darkly.min.css

$(ASSET_DIR)/js/all.min.js: $(wildcard $(SOURCE_DIR)/*/*)
	mkdir -p $(ASSET_DIR)/css $(ASSET_DIR)/js
	$(MINIFY) -o $(ASSET_DIR)/js/all.min.js $(SOURCE_DIR)/js/jquery.js \
		$(SOURCE_DIR)/js/bootstrap.bundle.js $(SOURCE_DIR)/js/datatables.js $(SOURCE_DIR)/js/saz.js
	$(MINIFY) -o $(ASSET_DIR)/css/common.min.css $(SOURCE_DIR)/css/datatables.css \
		$(SOURCE_DIR)/css/saz.css
	$(MINIFY) -o $(ASSET_DIR)/css/bootstrap.flatly.min.css $(SOURCE_DIR)/css/bootstrap.flatly.css
	$(MINIFY) -o $(ASSET_DIR)/css/bootstrap.darkly.min.css $(SOURCE_DIR)/css/bootstrap.darkly.css
	$(MINIFY) -o $(ASSET_DIR)/css/saz.darkly.min.css $(SOURCE_DIR)/css/saz.darkly.css

prepare ::
	go get -u github.com/go-bindata/go-bindata/v3/...
	go get -u github.com/tdewolff/minify/v2/...

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
