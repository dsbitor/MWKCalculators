# mwkGo

mwkGo converts Marv Klotz's DOS shop-calculator programs (each
originally a standalone DOS C program) into equivalent Go
programs that build and run on macOS, Linux, and Windows 10 or
later. Each original program becomes one Go program under
`MWKGo/<name>/`, keeping the original program's calculation,
worked examples, and known quirks intact rather than redesigning
the tool. Conversion is complete: `MWKGo/` contains 124 Go
programs (a handful of the 118 original archives bundled more
than one C program, for example `cuts.zip`'s `CUTS.C`,
`CUTLIST.C`, and `REMNANT.C`, each converted separately). Three
original programs have no Go equivalent; see "Excluded" in
`docs/calculators.md` for why. The full per-tier list of every
converted program, its purpose, and its worked example is
`docs/calculators.md`.

**Published at [github.com/dsbitor/MWKCalculators](https://github.com/dsbitor/MWKCalculators).**
This repository is a deliberate subset of a larger project
(this README, the Go source under `MWKGo/`, and `docs/`) whose
primary development happens in a separate Fossil-managed
repository, not published here. That upstream repository also
holds the original DOS C source these programs were converted
from, the conversion plan and its recorded decisions, and this
project's own engineering standards — none of which this GitHub
copy needs in order to build, run, or be understood on its own.

## Build and install

Go 1.26 or later and Git are required.

```bash
git clone https://github.com/dsbitor/MWKCalculators.git
cd MWKCalculators
task build
```

compiles every program under `MWKGo/` (this only checks that
everything builds; `go build ./...` across many `main` packages
does not retain the resulting binaries). To build one program's
binary:

```bash
cd MWKGo/<name>
go build -o ../../bin/<name> .
```

`bin/` is reserved for built binaries (not tracked in version
control) but nothing currently installs into it automatically;
there is no separate install step yet.

Built independently (one binary per program, no shared
executable), the 124 binaries range from 2.43 MB (`combi`) to
9.49 MB (`diam`), averaging 3.37 MB; 418 MB in total. The
largest binaries are the ones that embed a reference database
(see "Configuration" below) and so also link in
`modernc.org/sqlite`, a pure-Go SQLite implementation that
accounts for most of the size difference, not the embedded data
itself.

### Windows

No prebuilt Windows binaries are distributed yet, since doing so
without code-signing them triggers a SmartScreen warning on
first run. Building from source avoids this: a binary compiled
locally never acquires the Mark of the Web that SmartScreen
checks for, since that mark is applied by the browser or email
client a file was downloaded through, not by a compiler. Building
from source is also otherwise straightforward here, since nothing
in this project uses cgo (including `modernc.org/sqlite`, a
pure-Go SQLite implementation): no C/C++ toolchain (MSVC,
MinGW, or otherwise) needs to be installed alongside Go.

1. Install Go 1.26 or later from
   [go.dev/dl](https://go.dev/dl/) (the Windows installer sets
   up `PATH` automatically; a terminal restart may be needed
   for it to take effect).
2. Install Git from
   [git-scm.com](https://git-scm.com/downloads) and clone the
   source:

   ```powershell
   git clone https://github.com/dsbitor/MWKCalculators.git
   cd MWKCalculators
   ```

3. From the same PowerShell prompt:

   ```powershell
   cd MWKGo
   go build -o ..\bin\<name>.exe .\<name>
   ```

   substituting the calculator's own directory name for
   `<name>` (see `docs/calculators.md` for the full list). Go
   adds the `.exe` extension automatically even if it's left off
   the `-o` path.

To build every program at once instead of one at a time, install
[go-task](https://taskfile.dev/) (itself a single, dependency-free
binary, also cross-platform) and run `task build` from the
repository root, exactly as on macOS or Linux; it only verifies
that everything compiles; it does not produce or install
binaries into `bin/` on its own, on any platform (see above).

## Configuration

There is no shared runtime configuration file. Each program
takes its own input one of two ways, matching how the original
DOS program took input:

- Interactively, via numbered prompts with defaults shown in
  brackets (`internal/promptio`, the Go equivalent of the
  original programs' `vin`/`vpr` library calls). Most converted
  programs work this way; just run the binary and answer the
  prompts.
- Via a `-data <file>` flag pointing at a small text file in the
  original programs' own `STARTOFDATA`/`ENDOFDATA` format.
  Programs that read a job's own one-off data (for example
  `cuts`, `ogive`, `spline`) work this way; each program's own
  page under `docs/calculators/` documents its exact file
  format, and a `testdata/example.dat` alongside its source is a
  working example.

A small number of programs (`drill`, `unit`, and others reading a
shared reference table rather than one job's own data) instead
read from a SQLite database built once from the original
programs' shipped `.DAT` tables by `MWKGo/tools/build-refdb`.
Rebuilding it is a rare, explicit step covered in that tool's own
documentation, not something run per invocation.

Two SQLite databases exist, both stored per-user under
`os.UserConfigDir()/mwkgo` (`~/Library/Application Support/mwkgo`
on macOS, `$XDG_CONFIG_HOME/mwkgo` on Linux, `%AppData%\mwkgo` on
Windows):

- `reference.db`: universal reference tables (drill sizes,
  shaft/hole fits, cutting speeds, `unit`'s 246-unit database,
  and similar). The 88 KB committed at
  `MWKGo/internal/refdata/reference.db` is embedded into every
  binary that needs it and copied to the user's own
  `reference.db` the first time such a program runs on a given
  machine; that copy is never overwritten afterward, so a
  machine with an older seeded copy will not pick up newer
  embedded data until that file is removed.
- `user.db`: machine-specific data the program has no way to
  know on its own (a particular lathe's change gears, hole-circle
  presets, and similar), used by `change`, `ddh`, `diam`,
  `divhead`, `diffthrd`, `gearatio`, `spaceblk`, and `sun`. It is
  not shipped at all; it's created empty (schema only) the first
  time any of those eight programs runs, and grows only with what
  the user enters.

## Usage

```bash
# Interactive: preferred-numbers (Renard) series generator
go run ./MWKGo/ratio

# Data-file driven: cut a list of pipe lengths from standard
# bars, searching for zero-waste combinations first
go run ./MWKGo/cuts -data MWKGo/cuts/testdata/example.dat -zero-waste

# Mixed: interactive prompts for the cam's own parameters, plus
# an optional flag to render the resulting profile as SVG
go run ./MWKGo/cam -svg cam.svg
```

Every program also accepts `-h` for its own flag reference where
it takes flags at all.

## Testing

```bash
task test   # go test -race -count=1 ./... (required for every change)
task vet    # go vet ./...
task fmt    # gofmt -l -w .
```

`task build`/`test`/`vet`/`fmt`/`coverage` are thin wrappers around
plain `go` commands (see `Taskfile.yml`, not published to this
GitHub mirror). If [go-task](https://taskfile.dev/) isn't
installed, `compilingGo.md` gives the equivalent `go build`/
`go test`/`go vet`/`gofmt`/`go tool cover` commands directly, with
an explanation of what each flag does.

## Context file library

This project's engineering standards (Go style, testing,
logging, SQLite usage, Fossil workflow, documentation and
Markdown conventions) live in the upstream Fossil repository's
own `ai/context/` directory, not in this GitHub copy — see
"Published at" above.

## Acknowledgments

The conversion of all 124 Go programs, their tests and
documentation, and the setup of this GitHub publication were
carried out with Claude (Anthropic), working from the original
DOS C source under Marv Klotz's guidance-in-code rather than a
redesign brief. About two days of preparatory work (assembling a
clean copy of the original C sources, and setting up project
infrastructure and standards) preceded the first Fossil commit
on 2026-07-26; from there, 55 Fossil commits over the following
four days brought the conversion, its tests, its documentation,
and this GitHub publication to completion by 2026-07-30 — under
a week end to end.
