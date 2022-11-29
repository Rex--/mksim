/ FILE: hello_loop.p8
/
/ This hello world program is based on the basic hello_world example.
/ The only difference is the addition of a loop that halts every time the
/ full string is printed. This allows you to continue and print the string
/ again or single step through the instructions.

TSF=6041
TLS=6046

*10                   / Set current assembly origin to address 10,
STPTR,  STRNG-1       / An auto-increment register (one of eight at 10-17)
*20
STRST,  STRNG-1       / Holds the start of string

*200                  / Set current assembly origin to program text area
HELLO,  CLA CLL       / Clear AC and Link again (needed when we loop back from tls)
        TAD I STPTR   / Get next character, indirect via PRE-auto-increment address from the zero page
        SNA           / Skip if non-zero (not end of string)
        JMP LOOP      / Else skip to loop procedure (end of string)
        TLS           / Output the character in the AC to the teleprinter
        TSF           / Skip if teleprinter ready for character
        JMP .-1       / Else jump back and try again
        JMP HELLO     / Jump back for the next character
    
LOOP,   HLT         / Halt and wait for continue
        CLA         / On continue - Clear AC
        TAD STRST   / Load start of string address into AC
        DCA STPTR   / Store address in string pointer
        JMP HELLO   / Jump to start of hello routine to print out string

STRNG,  110           / H
        145           / e
        154           / l
        154           / l
        157           / o
        54            / ,
        40            / (space)
        167           / w
        157           / o
        162           / r
        154           / l
        144           / d
        41            / !
        12            / /n
        0             / End of string
$                /DEFAULT TERMINATOR
