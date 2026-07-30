# Compiling and testing the Go source

The commands below are what `task build`, `task test`, `task vet`,
`task fmt`, and `task coverage` actually run (see `Taskfile.yml`,
which isn't published to this GitHub mirror). Anyone without the
[go-task](https://taskfile.dev/) tool installed can run these
directly instead. All of them are run from the `MWKGo` directory:

```bash
cd MWKGo
```

## Build

Compiles every program in the module. Since this targets many
`main` packages at once rather than a single one, `go build`
only checks that everything compiles; it does not write out any
binaries. To produce an actual runnable binary for one program,
build its own directory directly with an explicit `-o` path (see
the README's "Build and install" section).

```bash
go build ./...
```

## Test

Runs every package's test suite. `-race` enables Go's data race
detector; `-count=1` disables the test result cache, forcing every
test to actually run rather than reusing a cached pass/fail result
from an earlier invocation (a cached result did not re-run the race
detector, so it is not valid evidence of a clean run). Both flags
are required for every change to this project, not just optional
extras.

```bash
go test -race -count=1 ./...
```

## Vet

Runs Go's static analysis checker (suspicious constructs like
unreachable code, wrong `Printf` verbs, and similar) across every
package.

```bash
go vet ./...
```

## Format

Rewrites every `.go` file in place to match `gofmt`'s standard
formatting. `-l` lists which files were changed; `-w` writes the
changes back to disk instead of only printing the reformatted
source.

```bash
gofmt -l -w .
```

## Coverage

Runs the same test suite as above, but also records which lines
were exercised, then prints a per-function coverage percentage
from that recording.

```bash
go test -race -count=1 -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```
