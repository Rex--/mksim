/ FILE: hello.p8
/
/ This hello world program was taken from the pdp-8 wikipedia page:
/ https://en.wikipedia.org/wiki/PDP-8
/ It is licensed under the terms of CC BY-SA 3.0, available here:
/ https://creativecommons.org/licenses/by-sa/3.0/
/
/ MODIFIED: 2022-11-28 Rex McKinnon
/    Modified the source to compile with pdpnasm.

TSF=6041
TLS=6046

*10                   / Set current assembly origin to address 10,
STPTR,    STRNG-1     / An auto-increment register (one of eight at 10-17)

*200                  / Set current assembly origin to program text area
HELLO,  CLA CLL       / Clear AC and Link again (needed when we loop back from tls)
        TAD I STPTR   / Get next character, indirect via PRE-auto-increment address from the zero page
        SNA           / Skip if non-zero (not end of string)
        HLT           / Else halt on zero (end of string)
        TLS           / Output the character in the AC to the teleprinter
        TSF           / Skip if teleprinter ready for character
        JMP .-1       / Else jump back and try again
        JMP HELLO     / Jump back for the next character

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
