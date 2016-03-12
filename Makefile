.PHONY: all clean dist

all: deps crosscompile dist

deps:
	go get github.com/mitchellh/gox
	go get ./...

build:
	go build grvm.go

crosscompile:
	gox --osarch="linux/amd64" --output="crosscompile/linux_amd64_grvm"
	gox --osarch="linux/386" --output="crosscompile/linux_386_grvm"
	gox --osarch="linux/arm" --output="crosscompile/linux_arm_grvm"
	gox --osarch="darwin/amd64" --output="crosscompile/darwin_amd64_grvm"

dist:
	@for target in crosscompile/*_grvm; do \
		rm -rf build; \
		mkdir -p build; \
		mkdir -p build/bin; \
		cp -R scripts build; \
		cp -R $$target build/bin/grvm; \
		tar -cvzf $$target.tar.gz --directory=build bin scripts; \
	done

clean:
	rm -rf crosscompile

localinstall:
	go build grvm.go
	rm -rf $$HOME/.grvm
	mkdir -p $$HOME/.grvm/bin
	mkdir -p $$HOME/.grvm/scripts
	cp scripts/grvm $$HOME/.grvm/scripts/grvm
	cp grvm $$HOME/.grvm/bin/grvm
	$$HOME/.grvm/bin/grvm doctor
