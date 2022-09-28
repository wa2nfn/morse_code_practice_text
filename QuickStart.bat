@echo OFF
SET PATH=%PATH%;%CD%
CMD /c cls
ECHO.
ECHO MCPT Quick Start 
ECHO.
ECHO A quick look at the use of SOME features, without explanation. See the QuickStart.html for more details.
ECHO Ctrl-C will EXIT the Quick Start.
ECHO.
ECHO Look at the command and options before you hit Enter. Look at the  generated output and you will 
ECHO likely see the direct correlation.
ECHO.
ECHO Hitting any key AFTER viewing the output runs the next command. (use scroll bars as needed)
ECHO.
PAUSE
CMD /c cls
ECHO Some Code Group Command Practice
ECHO.

ECHO mcpt -codeGroups -cgLen=3:6
ECHO.
PAUSE
CMD /c "echo. &&mcpt -codeGroups -cgLen=3:6&&echo.&&echo."
PAUSE
CMD /c cls

ECHO mcpt -codeGroups -cgLen=5 -num=50 -len=80
ECHO.
PAUSE
CMD /c "echo. &&mcpt -codeGroups -cgLen=5 -num=50 -len=80&&echo.&&echo."
PAUSE
CMD /c cls

ECHO mcpt -codeGroups -cgLen=3:6 -lesson=12 -tutor=lcwo
ECHO.
PAUSE
CMD /c "echo. &&mcpt -codeGroups -cgLen=3:6 -lesson=12 -tutor=lcwo&&echo.&&echo."
PAUSE
CMD /c cls

ECHO mcpt -permute=p -lesson=15 -tutor=lcwo
ECHO.
PAUSE
CMD /c "echo. &&mcpt -permute=p -lesson=15 -tutor=lcwo &&echo.&&echo."
PAUSE
CMD /c cls

ECHO mcpt -permute=b -lesson=15 -lesson=12 -tutor=g4fon
ECHO.
PAUSE
CMD /c "echo. &&mcpt -permute=b -lesson=12 -tutor=g4fon &&echo.&&echo."
PAUSE
CMD /c cls

ECHO.
ECHO Some Word Command Practice (some with code groups too)
ECHO.

ECHO mcpt -inFile=words.txt -inlen=3:10
ECHO.
PAUSE
CMD /c "echo. &&mcpt -inFile=words.txt -inlen=3:10&&echo.&&echo."
PAUSE
CMD /c cls

ECHO mcpt -inFile=words.txt -mixedMode=3 -cgLen=3:5
ECHO.
PAUSE
CMD /c "echo. &&mcpt -inFile=words.txt -mixedMode=3 -cgLen=3:5 &&echo.&&echo."
PAUSE
CMD /c cls

ECHO mcpt -inFile=words.txt -mixedMode=3 -sufLen=2
ECHO.
PAUSE
CMD /c "echo. &&mcpt -inFile=words.txt -mixedMode=3 -sufLen=2&&echo.&&echo."
PAUSE
CMD /c cls

ECHO mcpt -inFile=words.txt -preLen=1 -sufLen=2:2
ECHO.
PAUSE
CMD /c "echo. &&mcpt -inFile=words.txt -preLen=1 -sufLen=2:2&&echo.&&echo."
PAUSE
CMD /c cls

ECHO mcpt -inFile=words.txt -sufLen=1 -sufList="123=?"
ECHO.
PAUSE
CMD /c "echo. &&mcpt -inFile=words.txt -sufLen=1 -sufList=123&&echo.&&echo."
PAUSE
CMD /c cls

ECHO mcpt -inFile=words.txt -delimiter=1-5 -delimiterNum=2:2 -num=20
ECHO Note: delimiters are very flexible, they can do prosigns, and more - but not from a demo batch file.
ECHO.
PAUSE
CMD /c "echo.&&mcpt -inFile=words.txt -delimiter=1-5 -delimiterNum=2:2 -num=20&&echo.&&echo."
PAUSE
CMD /c cls

ECHO mcpt -inFile=words.txt -repeat=2 -num=20
ECHO.
PAUSE
CMD /c "echo.&&mcpt -inFile=words.txt -repeat=2 -num=20&&echo.&&echo."
PAUSE
CMD /c cls

ECHO mcpt -inFile=words.txt -reverse -lesson=15 -tutor=lcwo -num=10 -inlen=3:10
ECHO.
PAUSE
CMD /c "echo.&&mcpt -inFile=words.txt -reverse -lesson=10 -tutor=lcwo -inlen=5&&echo.&&echo."
PAUSE
CMD /c cls

ECHO mcpt -inFile=words.txt -prosignFile=prosign.txt -num=40
ECHO.
PAUSE
CMD /c "echo.&&mcpt -inFile=words.txt -prosignFile=prosign.txt -num=40&&echo.&&echo."
PAUSE
CMD /c cls

ECHO mcpt -inFile=words.txt -prosignFile=prosign.txt -num=40 -mixedMode=3 
ECHO.
PAUSE
CMD  /c "echo.&&mcpt -inFile=words.txt -prosignFile=prosign.txt -num=40 -mixedMode=3 &&echo.&&echo."
PAUSE
CMD /c cls

ECHO mcpt -inFile=words.txt -prosignFile=prosign.txt -num=40 -mixedMode=3
ECHO.
PAUSE
CMD  /c "echo.&&mcpt -inFile=words.txt -prosignFile=prosign.txt -num=40 -mixedMode=3&&echo.&&echo."
PAUSE
CMD /c cls

ECHO mcpt -codeGroups -prosignFile=prosign.txt -num=40
ECHO.
PAUSE
CMD /c "echo.&&mcpt -codeGroups -prosignFile=prosign.txt -num=40&&echo.&&echo."
PAUSE
CMD /c cls

ECHO mcpt -codeGroups -preLen=1 -preList=WK -cgLen=1 -cglist=0-9 -sufLen=2:3 -sufList=A-Z
ECHO.
PAUSE
CMD /c  "echo.&&mcpt -codeGroups -preLen=1 -preList=WK -cgLen=1 -cglist=0-9 -sufLen=2:3 -sufList=A-Z&&echo.&&echo."
PAUSE
CMD /c cls

ECHO mcpt -callSigns
ECHO.
PAUSE
CMD /c  "echo.&&mcpt -callSigns&&echo."
ECHO.
PAUSE
CMD /c cls

ECHO mcpt -send=1,6 -num=30 -displayFormat=LF
ECHO.
PAUSE
CMD /c  "echo.&&mcpt -send=1,6 -num=30 -displayFormat=LF&&echo."
PAUSE
CMD /c cls

ECHO.
ECHO Some Help Command Topics
ECHO.
ECHO These show some HELP information available without needing to refer to the User Guide.
ECHO You may need to scroll to the top of the page for some.
ECHO .

ECHO mcpt -help
ECHO.
PAUSE
CMD /c "echo.&&mcpt -help&&echo.&&echo."
PAUSE
CMD /c cls

ECHO mcpt -help=files
ECHO.
PAUSE
CMD /c  "echo.&&mcpt -help=files&&echo.&&echo."
PAUSE
CMD /c cls

ECHO mcpt -help=tutors
ECHO.
PAUSE
cmd /c  "echo.&&mcpt -help=tutors&&echo.&&echo."
PAUSE
CMD /c cls

ECHO mcpt -help=options (plus other options from a functioning command, shows ALL options and defaults)
ECHO. 
PAUSE
CMD /c "echo.&&mcpt -help=options -codeGroups&&echo.&&echo."
PAUSE
CMD /c cls

ECHO *** LCWO Users: The advanced, on-the-fly speed change features, controlled by the "LCWO_"
ECHO *** options have NOT been demonstrated. They are fully explained in the UserGuide.html
ECHO *** in the examples after #20 - you don't want miss these!
ECHO.


