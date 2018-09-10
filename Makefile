go_files := $(filter-out vendor/% _playground/%, $(wildcard **/*.go))
test_files := $(filter %_test.go, $(go_files))
src_files := $(filter-out $(test_files), $(go_files))
go_path := $(GOPATH)
pwd := $(shell pwd)
package := $(subst $(go_path)/src/, , $(pwd))
version := $(shell node -p 'require("./package.json").version')

server_cmd := $(package)/cmd/log.go-server
harvester_cmd := $(package)/cmd/log.go-harvester
provide_version := -ldflags='-X whitetrefoil.com/log-go.Version=$(version)'

.PHONY: dev all zip build test clean doc_server

all : test zip
	@echo All done...

zip : build
	@echo Compressing built binaries...
	@tar -cJf out/darwin_amd64.tar.xz -C out/darwin_amd64 log.go-harvester log.go-server
	@tar -cJf out/linux_amd64.tar.xz -C out/linux_amd64 log.go-harvester log.go-server
	@tar -cJf out/linux_386.tar.xz -C out/linux_386 log.go-harvester log.go-server

build : \
	out/darwin_amd64/log.go-harvester \
	out/darwin_amd64/log.go-server \
	out/linux_amd64/log.go-harvester \
	out/linux_amd64/log.go-server \
	out/linux_386/log.go-harvester \
	out/linux_386/log.go-server

out/darwin_amd64/log.go-server : $(src_files)
	@echo Building log.go-server for darwin/amd64
	@GOOS=darwin GOARCH=amd64 go build -o out/darwin_amd64/log.go-server $(provide_version) $(server_cmd)

out/darwin_amd64/log.go-harvester : $(src_files)
	@echo Building log.go-harvester for darwin/amd64
	@GOOS=darwin GOARCH=amd64 go build -o out/darwin_amd64/log.go-harvester $(provide_version) $(harvester_cmd)

out/linux_amd64/log.go-server : $(src_files)
	@echo Building log.go-server for linux/amd64
	@GOOS=linux GOARCH=amd64 go build -o out/linux_amd64/log.go-server $(provide_version) $(server_cmd)

out/linux_amd64/log.go-harvester : $(src_files)
	@echo Building log.go-harvester for linux/amd64
	@GOOS=linux GOARCH=amd64 go build -o out/linux_amd64/log.go-harvester $(provide_version) $(harvester_cmd)

out/linux_386/log.go-server : $(src_files)
	@echo Building log.go-server for linux/386
	@GOOS=linux GOARCH=386 go build -o out/linux_386/log.go-server $(provide_version) $(server_cmd)

out/linux_386/log.go-harvester : $(src_files)
	@echo Building log.go-harvester for linux/386
	@GOOS=linux GOARCH=386 go build -o out/linux_386/log.go-harvester $(provide_version) $(harvester_cmd)

test : $(test_files)
	@echo Testing...
	@go test -cover ./...

clean :
	rm -rf out

version :
	@echo $(version)

dev :
	@echo $(provide_version)

doc_server :
	@godoc -http=:4444 -index
