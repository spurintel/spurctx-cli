# Spurctx

Spurctx is a command-line application written in Go that queries the Spur API for information about IP addresses. It supports multiple forms of input and can process multiple IP addresses in parallel.

## Usage

To use the application, you can provide IP addresses as command-line arguments, in a file, or through standard input. The application also supports a "garbage" mode where it extracts IP addresses from any text input.

```bash
Usage of ./target/spurctx:
  -f string
    	File with IP addresses
  -g	Garbage input
  -gf string
    	Garbage file input
  -ip string
    	IP addresses separated by comma
  -n int
    	Override the parallelism (default 20)
```

Here are some example usages:

- Command-line arguments: `spurctx -ip 192.168.1.1,192.168.1.2`
- File input: `spurctx -f ips.txt`
- Garbage file input (does regex search for ip on each line): `spurctx -gf ips.txt`
- Standard input: `echo 1.1.1.1 | spurctx`
- Garbage input: `echo "The IP address is 1.1.1.1." | spurctx -g`

## Building

You can build the application using the provided Makefile:

- `make fmt`: Formats the Go code.
- `make lint`: Checks the Go code for potential errors.
- `make bin`: Compiles the Go code into a binary.
- `make all`: Runs `fmt`, `lint`, and `bin` in that order.
- `make crosscompile`: Compiles the Go code into binaries for Linux, macOS, and Windows, both for AMD64 and ARM64 architectures.

The binaries are placed in the `target` directory.

## Installation

### Using Homebrew
```
brew tap spurintel/spurintel
brew install spurintel/spurintel/spurctx
```

### Using `go install`

If you have Go installed, you can install `spurctx` directly using `go install`:

```bash
go install github.com/spurintel/spurctx-cli@latest
```

This will install the `spurctx-cli` binary to your `$GOPATH/bin` directory.

### Downloading the Binary

You can also download the pre-compiled binaries from the GitHub releases page. Choose the binary that matches your operating system and architecture, download it, and move it to a directory in your `PATH`.

For example, to download the binary for Linux on AMD64:
```bash
wget https://github.com/spurintel/spurctx/releases/download/vX.Y.Z/spurctx-linux-amd64
chmod +x spurctx-linux-amd64
mv spurctx-linux-amd64 /usr/local/bin/spurctx
```

Replace `vX.Y.Z` with the version you want to download.