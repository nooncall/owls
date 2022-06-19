test: fmt
#	go test -race  ./controller/test/...
#	go test -race  ./service/checker/...
#	go test -race  ./service/sql_util/...
#	go test -race  ./util/...

config:
	mkdir -p ./bin/resource
	cp server/config.yaml ./bin/
	cp server/resource/rbac_model.conf ./bin/resource/rbac_model.conf

build: fmt config
	cd server && \
	go build -o ../bin/owls ./cmd/owls/  &&\
	cd ..

build-linux: config
	cd server && \
	CGO_ENABLED=0 GOOS=linux go build -o ../bin/owls -a -ldflags '-extldflags "-static"' ./cmd/owls/
	cd ..
fmt:
	go fmt ./...

run: config build-front build build-docs
	cd ./bin && ./owls

.ONESHELL:
build-front:
	mkdir -p bin
	rm -rf ./bin/static
	cd web/ && npm run build && cp -r ./dist ../bin/static
	cd ..

build-docs:
	mkdir -p bin
	rm -rf ./bin/docs-static
	cd docs/ && npm run build && cp -r ./build ../bin/docs-static
	cd ..

build-docker: build-front build-linux build-docs
	docker build -t mingbai/owls:v0.1.0 .

run-docker: build-docker
	docker run -p 80:80 -d  mingbai/owls:v0.1.0
