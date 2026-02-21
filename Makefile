PKGS = freetype2
FLAGS_FILE = compile_flags.txt

.PHONY: all clean flags

all: flags

flags:
	@echo "Generating $(FLAGS_FILE) for packages $(PKGS)"
	@pkg-config --cflags-only-I $(PKGS) | tr ' ' '\n' > $(FLAGS_FILE)
	@echo "Done."
