.PHONY: all clean dist version

VERSION := $(shell if [ -z "$$TRAVIS_TAG" ]; then echo "dev"; else echo "$$TRAVIS_TAG" | tail -c +2; fi)

all: deps crosscompile dist

deps:
	go get github.com/mitchellh/gox
	go get ./...

crosscompile:
	gox --osarch="linux/amd64" --output="crosscompile/linux_amd64_grvm" --ldflags='-X main.version=${VERSION}'
	gox --osarch="linux/386" --output="crosscompile/linux_386_grvm" --ldflags='-X main.version=${VERSION}'
	gox --osarch="linux/arm" --output="crosscompile/linux_arm_grvm" --ldflags='-X main.version=${VERSION}'
	gox --osarch="darwin/amd64" --output="crosscompile/darwin_amd64_grvm" --ldflags='-X main.version=${VERSION}'

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

devinstall:
	sh install.sh devinstall
