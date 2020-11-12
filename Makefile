all: build
	@echo done

build:
ifeq (,$(shell which gocat))
	$(install_gocat)
endif
	@mkdir -p dist
	go vet .
	gocat -n -p main *.go | sed -e s/__USER__/${USER}/ > dist/response.go
	

define install_gocat
	GO111MODULE=off go get github.com/naegelejd/gocat
endef