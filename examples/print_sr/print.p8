/ File: print.p8

/ This program prints the contents of the lower 8-bits of the switch register
/ to the teletype as an ascii character. Key in the wanted ascii code to the
/ SR and press Continue to print the char.

PRNTSR,     HLT         / Wait for user to continue
            LAS         / Load SR into AC
            TLS         / Print the character and clear flag
            TSF         / Skip if flag is set (means printer is ready)
            JMP .-1     / Jump back and wait till ready
            JMP PRNTSR  / Jump to beginning of loop if printer is ready
$
