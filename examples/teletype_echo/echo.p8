/ File: echo.pa
/
/ This program echos incoming characters from the keyboard to the teleprinter.
/ A press of the 'Enter' key breaks from the loop and halts the computer.

*200
ECHO,   KSF             / Skip if character ready
        JMP .-1         / Jump back and wait if not ready
        KRB             / Read character into AC
        TSF             / Skip if teleprinter ready for character
        JMP .-1         / Else jump back and wait
        TLS             / Print character in AC
        TAD NNL         / Add negated new line
        SZA             / Skip if newline character (ac is zero)
        JMP ECHO        / Jump back for the next character
        JMP EXIT        / Break from loop
NNL,    7766            / Two's complement of "\n" (ascii 10)

EXIT,   HLT
        JMP ECHO
$
