# Run with "CGO_ENABLED=0 GOOS=linux" in tyhe enfironment for Docker.
ifeq ($(DOCKER),1)
	export CGO_ENABLED=0
	export GOOS=linux
	GOFLAGS=-a -installsuffix cgo -ldflags '-extldflags "-static"'
endif

ASSET_DIR=cmd/sazserve/assets
ASSET_BIN=cmd/sazserve/assets.go

all: sazdump sazserve

sazdump: $(wildcard cmd/sazdump/*.go pkg/dumper/*.go pkg/parser/*.go pkg/analyzer/*.go)
	go build $(GOFLAGS) cmd/sazdump/sazdump.go

sazserve: $(ASSET_BIN) $(wildcard cmd/sazserve/*.go pkg/parser/*.go pkg/analyzer/*.go)
	go build $(GOFLAGS) cmd/sazserve/sazserve.go $(ASSET_BIN)

$(ASSET_BIN): $(wildcard $(ASSET_DIR)/* $(ASSET_DIR)/*/*)
	go-bindata -fs -o $(ASSET_BIN) -prefix $(ASSET_DIR) $(ASSET_DIR)/...

run-sazdump: $(wildcard cmd/sazdump/*.go pkg/dumper/*.go pkg/parser/*.go pkg/analyzer/*.go)
	go run cmd/sazdump/sazdump.go "$(SAZ)"

run-sazserve: $(ASSET_BIN) $(wildcard cmd/sazserve/*.go pkg/parser/*.go pkg/analyzer/*.go)
	go run cmd/sazserve/sazserve.go $(ASSET_BIN)

debug-sazdump: $(wildcard cmd/sazdump/*.go pkg/dumper/*.go pkg/parser/*.go pkg/analyzer/*.go)
	go run cmd/sazdump/sazdump.go "$(SAZ)"

debug-sazserve: debug-assets $(wildcard cmd/sazserve/*.go pkg/parser/*.go pkg/analyzer/*.go)
	go run cmd/sazserve/sazserve.go $(ASSET_BIN) cmd/sazserve/debug.go

debug-assets:
	go-bindata -debug -fs -o $(ASSET_BIN) -prefix $(ASSET_DIR) $(ASSET_DIR)/...

minify: $(ASSET_DIR)/js/jquery.min.js $(ASSET_DIR)/js/datatables.min.js

$(ASSET_DIR)/js/jquery.min.js: $(ASSET_DIR)/js/jquery.js
	./node_modules/.bin/terser -o $(ASSET_DIR)/js/jquery.min.js -m -c \
		--comments false \
		-source-map filename=$(ASSET_DIR)/js/jquery.min.js.map \
		--source-map url=jquery.min.js.map $(ASSET_DIR)/js/jquery.js

$(ASSET_DIR)/js/datatables.min.js: $(ASSET_DIR)/js/datatables.js
	./node_modules/.bin/terser -o $(ASSET_DIR)/js/datatables.min.js -m -c \
		--comments false \
		-source-map filename=$(ASSET_DIR)/js/datatables.min.js.map \
		--source-map url=datatables.min.js.map $(ASSET_DIR)/js/datatables.js

prepare:
	go get -u github.com/go-bindata/go-bindata/...

clean:
	rm -f sazdump sazserve $(ASSET_BIN)

push:
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
	docker run --rm -it sazserve

docker-dump-example ::
	docker run --rm -it -v ${PWD}:/work -w /work sazdump examples/test.saz

docker-serve-example ::
	docker run --rm -it -v ${PWD}:/work -w /work sazdump examples/test.saz
	docker run --rm -it sazserve

docker-tag ::
	docker tag sazdump prantlf/sazdump:latest
	docker tag sazserve prantlf/sazserve:latest

docker-login ::
	docker login --username=prantlf

docker-push ::
	docker push prantlf/sazdump:latest
	docker push prantlf/sazserve:latest

.PHONY: clean debug-sazdump debug-sazserve debug-assets push prepare
