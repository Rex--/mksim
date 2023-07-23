/ File: write_tape.p8
/
/ This program reads a tape and make a copy of it.

*10
MSGPTR, 0o700-1

*200
MAIN,   CLA CLL
        TAD I MSGPTR
        SNA
        JMP COPY
        JMS ECHO
        JMP MAIN

COPY,   JMS GTCHR
        SNA
        JMP EXIT
        JMS PTCHR
        CLA
        TAD DOT
        JMS ECHO
        JMP MAIN

EXIT,   HLT

PTCHR,  0
        PSF
        JMP .-1
        PLS
        JMP I PTCHR

GTCHR,  0
        RFC             / Fetch character from tape
        RSF             / Skip if reader flag == 1
        JMP .-1         / Jump back and test flag again
        CLA             / Clear AC
        RRB             / Load AC from reader buffer
        JMP I GTCHR     / Return from subroutine

ECHO,   0
        TSF             / Skip if teleprinter ready for character
        JMP .-1         / Else jump back and test again
        TLS             / Output the character in the AC to the teleprinter
        JMP I ECHO      / Return from subroutine
*700
MSG,    'C'
        'o'
        'p'
        'y'
        'i'
        'n'
        'g'
        ' '
        't'
        'a'
        'p'
        'e'
        0

DOT,    '.'

$
