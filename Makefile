.PHONY: all clean dist

all: crosscompile dist

build:
	go build grvm.go

crosscompile:
	GOOS=linux   GOARCH=amd64 go build -o ./crosscompile/linux_amd64_grvm   grvm.go
	GOOS=linux   GOARCH=386   go build -o ./crosscompile/linux_386_grvm     grvm.go
	GOOS=linux   GOARCH=arm   go build -o ./crosscompile/linux_arm_grvm     grvm.go
	GOOS=darwin  GOARCH=amd64 go build -o ./crosscompile/darwin_amd64_grvm  grvm.go

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
