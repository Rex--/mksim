/ RIM Loader as specified in DEC-08-LRAA-D

*7756
BEG,    KCC
        KSF
        JMP .-1
        KRB
        CLL RTL
        RTL
        SPA
        JMP BEG+1
        RTL
        KSF
        JMP .-1
        KRS
        SNL
        DCA I TEMP
        DCA TEMP
        JMP BEG
TEMP,   0
$
