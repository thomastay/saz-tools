all: sazdump sazserve

sazdump: $(wildcard cmd/sazdump/*.go pkg/dumper/*.go pkg/parser/*.go pkg/analyzer/*.go)
	go build cmd/sazdump/sazdump.go

sazserve: cmd/sazserve/assets.go $(wildcard cmd/sazserve/*.go pkg/parser/*.go pkg/analyzer/*.go)
	go build cmd/sazserve/sazserve.go cmd/sazserve/assets.go

cmd/sazserve/assets.go: $(wildcard cmd/sazserve/assets/* cmd/sazserve/assets/*/*)
	go-bindata -fs -o cmd/sazserve/assets.go \
		-prefix cmd/sazserve/assets cmd/sazserve/assets/...

debug-sazdump: $(wildcard cmd/sazdump/*.go pkg/dumper/*.go pkg/parser/*.go pkg/analyzer/*.go)
	go run cmd/sazdump/sazdump.go "$(SAZ)"

debug-sazserve: cmd/sazserve/assets.go $(wildcard cmd/sazserve/*.go pkg/parser/*.go pkg/analyzer/*.go)
	go run cmd/sazserve/sazserve.go cmd/sazserve/assets.go

debug-assets:
	go-bindata -debug -fs -o cmd/sazserve/assets.go \
		-prefix cmd/sazserve/assets cmd/sazserve/assets/...

minify: cmd/sazserve/assets/js/jquery.min.js cmd/sazserve/assets/js/datatables.min.js

cmd/sazserve/assets/js/jquery.min.js: cmd/sazserve/assets/js/jquery.js
	./node_modules/.bin/terser -o cmd/sazserve/assets/js/jquery.min.js -m -c \
		--comments false \
		-source-map filename=cmd/sazserve/assets/js/jquery.min.js.map \
		--source-map url=jquery.min.js.map cmd/sazserve/assets/js/jquery.js

cmd/sazserve/assets/js/datatables.min.js: cmd/sazserve/assets/js/datatables.js
	./node_modules/.bin/terser -o cmd/sazserve/assets/js/datatables.min.js -m -c \
		--comments false \
		-source-map filename=cmd/sazserve/assets/js/datatables.min.js.map \
		--source-map url=datatables.min.js.map cmd/sazserve/assets/js/datatables.js

clean:
	rm -f sazdump sazserve cmd/sazserve/assets.go

push:
	git push && git push heroku master

.PHONY: clean debug-sazdump debug-sazserve debug-assets push
