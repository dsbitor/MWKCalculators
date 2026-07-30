rpc
===

RPN (Reverse Polish Notation) stack-oriented scientific calculator.

**Converted from:** `RPC.C` (M. W. Klotz), `MWKC/Math/rpc.zip` **Go source:**
`MWKGo/rpc/rpc.go`

Purpose
-------

A four-register stack calculator (X, Y, Z, T, in the usual RPN sense) with a
memory register and a wide set of named operators: arithmetic, trigonometry
(degrees or radians), inverse trig, rectangular/polar conversion, exponentials
and logarithms (including base-2 and an arbitrary-base variant), bitwise
operators, number-base conversion (decimal/hex/octal/binary), and unit
conversion (temperature, weight, length). A full second copy of the calculator's
state can be swapped in for auxiliary calculations and swapped back.

The original program is driven entirely by mouse clicks on an on-screen keypad,
building up a pending numeric entry one keystroke at a time before committing it
to the stack on the next operator click. This conversion replaces that with a
single command line per step: a line that parses as a number (in the current
input base) pushes it onto the stack immediately, and any other recognised line
runs the matching named operator. Because a whole number is always typed and
committed in one step, there is no equivalent of the original's separate
pending-entry buffer to convert; every operator that would have auto-entered a
pending value now simply operates on values already on the stack.

Inputs
------

A number, entered on its own line, pushes it onto the stack (interpreted in the
current base: decimal, hex, octal, or binary). Everything else is a named
command; run `help` for the full list. The main categories:

| Commands                                                  | Effect                                               |
|-----------------------------------------------------------|------------------------------------------------------|
| `+ - * / x/y`                                             | basic arithmetic (`x/y` divides the y register by x) |
| `roll rolldn xy xm`                                       | stack rearrangement                                  |
| `e lastx 1/x chs`                                         | constants and single-value operators                 |
| `pi sqr sqrt fact`                                        | pi, x², √x, x!                                       |
| `store store+ store- store* store/`                       | store into memory                                    |
| `rcall rcall+ rcall- rcall* rcall/`                       | recall from memory                                   |
| `sin cos tan rss`                                         | trig; `rss` = √(x²+y²)                               |
| `asin acos atan unrss`                                    | inverse trig; `unrss` = √\|y²−x²\|                   |
| `topolar torect`                                          | rectangular ↔ polar conversion                       |
| `deg2rad rad2deg atan2`                                   | angle conversion, four-quadrant arctangent           |
| `frac int split ymodx`                                    | fractional/integer parts, y mod x                    |
| `floor ceil round gcd lcm`                                |                                                      |
| `e^x 10^x 2^x y^x y^(1/x)`                                | exponentials and roots                               |
| `ln log log2 ylogx`                                       | logarithms                                           |
| `and or 1comp xor`                                        | bitwise, 4-byte integer range                        |
| `temp weight length`                                      | unit conversion using the current from/to scales     |
| `tempfrom tempto weightfrom weightto lengthfrom lengthto` | cycle the from/to scale                              |
| `fix eng sci`                                             | display notation                                     |
| `dp=n` / `adj`                                            | fixed decimal places / self-adjusting decimals       |
| `dec hex bin oct`                                         | number input/display base                            |
| `deg rad`                                                 | angle mode                                           |
| `clearx clrstk clrmem clrall`                             | clearing                                             |
| `swap swapx`                                              | swap with / swap x with the secondary calculator     |
| `undo`                                                    | undo the last operation                              |
| `notes`                                                   | usage notes                                          |
| `quit`                                                    | exit                                                 |

Output
------

The stack (M, L, T, Z, Y, X), each formatted per the current decimal/notation
settings; the hex/octal/binary breakdown of the X register whenever it holds a
whole number within a 4-byte integer range; and a status line showing the
current unit-conversion scales, input base, angle mode, and display notation.

Method
------

Each binary operator follows RPN convention: `y op x`, where x is the top of
stack (most recently entered) and y is just below it, so `10 enter 3 -` computes
`10 - 3`. Every operator that mutates the stack, x register, or memory first
snapshots the pre-operation state for `undo`, and updates the "last x" register
to the x register's value just before the operation — both exactly matching the
original's own `save()` function. Two categories of operator are deliberately
**not** undoable, matching the original precisely: entering a number (it updates
"last x" but bypasses the undo snapshot), and the `store` family (memory-only
writes, unlike the `rcall` family which does snapshot first). An `undo` right
after one of these reverts to whichever earlier operation last took a snapshot,
not to the state just before the non-undoable step.

Trigonometric operators use the current angle mode (degrees by default) to
convert to and from radians for the underlying math library calls. The number of
decimal places shown, comma grouping, and engineering/scientific notation are
handled by the shared `internal/numdisplay` package (extracted from this same
display logic when converting `mix` in the same session), since both programs'
original C used the identical `vp`/`dplaces` display functions.

The bitwise operators (`and`, `or`, `xor`, `1comp`) and the hex/octal/binary
display operate on the original's 4-byte `long` integer range: a value is
converted through a signed 32-bit integer, so a negative x register displays as
the two's-complement bit pattern of that 4-byte integer, matching how the
original C compiler truncated a negative `double` into an `unsigned long` for
display.

The `ylogx` operator preserves a genuine quirk in the original: it computes a
logarithm of the y register, but to a base taken from `lastX` — the value that
was in the x register before the *previous* operation — not from the current x
register, even though the current x register is checked for positivity before
the operation runs. This looks like a mismatch between the guard and the
formula, but it is exactly what `RPC.C`'s `mbranch()` computes
(`log(t)/log(lastx)` where `t` is the y register), so it is preserved rather
than "corrected"; see `TestLogYX_UsesLastXAsBase` for a worked demonstration.

Worked Example
--------------

`RPC.TXT` describes the operators but includes no worked numeric transcript. As
independently verifiable checks, this conversion's tests confirm: `3 enter 4 +`
= 7; `5!` = 120; `gcd(12,18)` = 6 and `lcm(4,6)` = 12; `asin(sin(x))` reproduces
x for x in [-90°,90°]; converting rectangular to polar and back reproduces the
original (x,y); `ln(eˣ)` and `log10(10ˣ)` both reproduce x; `2^3` = 8; `12 and
10` = 8, `12 or 10` = 14, `12 xor 10` = 6 (all in binary); the one's complement
of 0 is -1 (all 4 bytes set); a store followed by an undo reverts to the last
real checkpoint rather than to the state just before the store, while a recall
followed by undo does revert the recall; and `ylogx` with lastX=8, y=100
computes log₈(100), not log_x(100) for whatever is in the x register.
