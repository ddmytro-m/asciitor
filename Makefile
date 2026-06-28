PKGS = freetype2
FLAGS_FILE = compile_flags.txt

.PHONY: all run\:cli flags clean

all: run\:cli

run\:cli:
	go run cmd/asciitor/main.go

build\:cli:
	go build -o build/asciitor cmd/asciitor/main.go

test:
	go test -v ./...

flags:
	@echo "Generating $(FLAGS_FILE) for packages $(PKGS)"
	@pkg-config --cflags-only-I $(PKGS) | tr ' ' '\n' > $(FLAGS_FILE)
	@echo "Done."

clean:
	rm $(FLAGS_FILE)