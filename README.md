```
NNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNMTMMNNMTMMNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNN
NNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNg  gNNg  gNNM  MNNNNNNNNNNNNNNNNNNNNNNNNNNNNN
NNNMMMMMMMNNNNNNMMMMMMNNNNNMMMMMMNNMMMMNNMMMMNNM  MMMMNNNNNMMMMMNNNNNNMMMMMMNNNN
NNNN  y    MNNM   yy mNNNM    y mNNM  NNNK  NNNM   yyyNNM        MNNNm   y  NNNN
NNNNMMMMM  mNNm  MMMNNNNM  mNNNNNNNM  NNNK  NNNM  MNNNNN   NNNNm  MNNm  NNNNNNNN
NNM   ggg  mNNNNgy   MNNM  MNNNNNNNM  NNNK  NNNM  MNNNNN  MNNNNM  MNNm  NNNNNNNN
NNK  MMMM  mNNMMMMMM  NNN   MMMMMNNM  NNNK  NNNN  MMMMNNm  MMMM   NNNm  NNNNNNNN
NNNgy     ygNNy     ggNNNNgg    yNNM  NNNK  NNNNNy    gNNNgy   ygNNNNm  NNNNNNNN
NNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNN
```

## installation
### pre-requirements
- gcc compiler
- [freetype library](https://download.savannah.gnu.org/releases/freetype/) (or `brew install freetype`)
- go >=v1.25

### build (cli)
```bash
make flags      # generate compile_flags for freetype
make build:cli  # build binary to build/asciitor
```
or run directly:
```bash
make run:cli
```

## options (cli)
### structure
```
asciitor [INPUT] ...flags
```
- `input` - input writer (file or pipe), default to stdin
- `--output -o` - output writer (file or pipe), default to stdout
- `--width -w` - output width (px/characters/tw/original)
- `--charset -c` - charset preset (alphanumeric/ascii/braille) or file
- `--height -h` - output height (px/lines/th/original)
- `--font -f` - font from the bundled (mono/symbols/inter) or font file
- `--face` - font face index
- `--font-size -s` - font size in pixels
- `--block-size -b` - block size (bigger values better represent colors but lack details)
- `--fill` - fill output with no respect to original proportions
- `--inverse` - inverse image colors

## todo
* ~~basic library functionality~~
* ~~cli tool~~
* examples
* text recognition
* better README
* xkcd to ascii tool (more examples)
* frontend
