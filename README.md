MKSIM
=====
A PDP-8 emulator with a console UI written in Go.


Usage
-----
To run a program, pass the compiled program as the first argument:

    mksim hello.po

Compiled programs must be in either Pobj or RIM format. To aquire programs in
this format, use either [pdpnasm](http://people.csail.mit.edu/ebakke/pdp8/)
(Pobj) or [mkasm](https://github.com/Rex--/mkasm.git) (Pobj or RIM).

### IOT Devices
Programs can take advantage of a Teletype IOT device that uses device addresses
`03` (keyboard) and `04` (printer).

### Help
```
Usage: ./mksim [options] <in_file>

Options:
  -F_CPU speed
        simulated clock speed (default 8000000)
  -exit
        Exit the simulator on HALT
  -halt
        HALT the machine before first instruction cycle
  -help
        Print this message and exit
  -lock page
        Lock memory viewer to page (default -1)
  -no-gui
        Do not display curses ui
  -print-return
        Print return code (AC) upon exiting
```


Build
-----
Building this project requires `go`, which can be downloaded online.

First, clone the repo to get the latest version:

    git clone https://github.com/Rex--/mksim.git

Next, build the `mksim` executable. In the repo directory run:

    go build .


Copying
-------
Copyright (c) 2022-2024 Rex McKinnon \
This software is available for free under the permissive University of
Illinois/NCSA Open Source License. See the LICENSE file for full details.