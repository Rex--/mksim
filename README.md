MKSIM
=========

A simple PDP-8 emulator written in Go.

**Status:** Most of the main instructions are implemented, the most notable ones
missing are the rotate instructions. \
There is only one IOT device implemented:
the teletype printer. This device implements identical instructions to the OG
teletype device, but instead of physically printing the character it
figuratively prints the character to stdout.

Usage
----------
Pass the compiled object file as the first argument:
```
./mksim examples/hello_world/hello.po
```
Input files must be in the random undocumented format produced by
[pdpnasm](http://people.csail.mit.edu/ebakke/pdp8/).

