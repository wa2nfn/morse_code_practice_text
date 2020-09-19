@echo OFF
CMD /c cls
ECHO.
ECHO MCPT Quick Start 
ECHO.
ECHO A quick look at the use of some features, without explanation. See the QuickStart.html for details.
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

ECHO mcpt -codeGroups -cglen=3:6
ECHO.
PAUSE
CMD /c "echo. &&mcpt -codeGroups -cglen=3:6&&echo.&&echo."
PAUSE
CMD /c cls

ECHO mcpt -codeGroups -cglen=5 -num=50 -len=80
ECHO.
PAUSE
CMD /c "echo. &&mcpt -codeGroups -cglen=5 -num=50 -len=80&&echo.&&echo."
PAUSE
CMD /c cls

ECHO mcpt -codeGroups -cglen=3:6 -lesson=12 -tutor=lcwo
ECHO.
PAUSE
CMD /c "echo. &&mcpt -codeGroups -cglen=3:6 -lesson=12 -tutor=lcwo&&echo.&&echo."
PAUSE
CMD /c cls

ECHO mcpt -in=words.txt -permute=p -lesson=15 -tutor=lcwo
ECHO.
PAUSE
CMD /c "echo. &&mcpt -in=words.txt -permute=p -lesson=15 -tutor=lcwo &&echo.&&echo."
PAUSE
CMD /c cls

ECHO mcpt -in=words.txt -permute=b -lesson=15 -lesson=12 -tutor=g4fon
ECHO.
PAUSE
CMD /c "echo. &&mcpt -in=words.txt -permute=b -lesson=12 -tutor=g4fon &&echo.&&echo."
PAUSE
CMD /c cls

ECHO.
ECHO Some Word Command Practice (some with code groups too)
ECHO.

ECHO mcpt -in=words.txt -inlen=3:10
ECHO.
PAUSE
CMD /c "echo. &&mcpt -in=words.txt -inlen=3:10&&echo.&&echo."
PAUSE
CMD /c cls

ECHO mcpt -in=words.txt -mixedMode=3 -cglen=3:5
ECHO.
PAUSE
CMD /c "echo. &&mcpt -in=words.txt -mixedMode=3 -cglen=3:5 &&echo.&&echo."
PAUSE
CMD /c cls

ECHO mcpt -in=words.txt -mixedMode=3 -suflen=2
ECHO.
PAUSE
CMD /c "echo. &&mcpt -in=words.txt -mixedMode=3 -suflen=2&&echo.&&echo."
PAUSE
CMD /c cls

ECHO mcpt -in=words.txt -prelen=1 -suflen=2:2
ECHO.
PAUSE
CMD /c "echo. &&mcpt -in=words.txt -prelen=1 -suflen=2:2&&echo.&&echo."
PAUSE
CMD /c cls

ECHO mcpt -in=words.txt -suflen=1 -suflist="123=?"
ECHO.
PAUSE
CMD /c "echo. &&mcpt -in=words.txt -suflen=1 -suflist=123&&echo.&&echo."
PAUSE
CMD /c cls

ECHO mcpt -in=words.txt -delimiter=1-5 -DM=2:2 -num=20
ECHO Note: delimiters are very flexible, they can do prosigns, and more - but not from a batch file.
ECHO.
PAUSE
CMD /c "echo.&&mcpt -in=words.txt -delimiter=1-5 -DM=2:2 -num=20&&echo.&&echo."
PAUSE
CMD /c cls

ECHO mcpt -in=words.txt -repeat=2 -num=20
ECHO.
PAUSE
CMD /c "echo.&&mcpt -in=words.txt -repeat=2 -num=20&&echo.&&echo."
PAUSE
CMD /c cls

ECHO mcpt -in=words.txt -reverse -lesson=15 -tutor=lcwo -num=10 -inlen=3:10
ECHO.
PAUSE
CMD /c "echo.&&mcpt -in=words.txt -reverse -lesson=10 -tutor=lcwo -inlen=5&&echo.&&echo."
PAUSE
CMD /c cls

ECHO mcpt -in=words.txt -prosign=prosign.txt -num=40
ECHO.
PAUSE
CMD /c "echo.&&mcpt -in=words.txt -prosign=prosign.txt -num=40&&echo.&&echo."
PAUSE
CMD /c cls

ECHO mcpt -in=words.txt -prosign=prosign.txt -num=40 -mixedMode=3 
ECHO.
PAUSE
CMD  /c "echo.&&mcpt -in=words.txt -prosign=prosign.txt -num=40 -mixedMode=3 &&echo.&&echo."
PAUSE
CMD /c cls

ECHO mcpt -in=words.txt -prosign=prosign.txt -num=40 -mixedMode=3
ECHO.
PAUSE
CMD  /c "echo.&&mcpt -in=words.txt -prosign=prosign.txt -num=40 -mixedMode=3&&echo.&&echo."
PAUSE
CMD /c cls

ECHO mcpt -codeGroups -prosign=prosign.txt -num=40
ECHO.
PAUSE
CMD /c "echo.&&mcpt -codeGroups -prosign=prosign.txt -num=40&&echo.&&echo."
PAUSE
CMD /c cls

ECHO mcpt -codeGroups -prelen=1 -prelist=WK -cglen=1 -cglist=0-9 -suflen=2:3 -suflist=A-Z
ECHO.
PAUSE
CMD /c  "echo.&&mcpt -codeGroups -prelen=1 -prelist=WK -cglen=1 -cglist=0-9 -suflen=2:3 -suflist=A-Z&&echo.&&echo."
PAUSE
CMD /c cls

ECHO mcpt -callSigns
ECHO.
PAUSE
CMD /c  "echo.&&mcpt -callSigns&&echo."
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


