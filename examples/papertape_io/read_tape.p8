/ File: read_tape.p8
/
/ This simple program reads a papertape device character by character and
/ prints it to the teleprinter.

*200

MAIN,   JMS GTCHR       / Get next character and place in AC
        JMS ECHO        / Echo character to teleprinter
        SZA             / Skip if AC == 0 (EOF)
        JMP MAIN        / Else jump back and print the next character
        HLT             / Halt on EOF

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
$
