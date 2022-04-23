test: fmt
#	go test -race  ./controller/test/...
#	go test -race  ./service/checker/...
#	go test -race  ./service/sql_util/...
#	go test -race  ./util/...

config:
	mkdir -p ./bin/config
	cp server/config.yaml ./bin/
	cp server/config/config.yaml ./bin/config/config.yaml
	
build: fmt config
	cd server && \
	go build -o ../bin/owl ./cmd/owl/  &&\
	cd ..

build-linux: fmt config
	CGO_ENABLED=0 GOOS=linux go build -o bin/owl -a -ldflags '-extldflags "-static"' ./server/cmd/owl/

fmt:
	go fmt ./...

run: build-front build
	./bin/owl

.ONESHELL:
build-front:
	mkdir -p bin
	cd web/ && npm run build && mv ./dist ../bin/static
	cd ..

build-docker: build-front build
	docker build -t mingbai/owls:v0.1.0 .

run-docker: build-docker
	docker run -p 8787:8787 -d  palfish/owl:v0.1.0
