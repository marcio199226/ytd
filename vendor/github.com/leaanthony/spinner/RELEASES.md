# Releases

## 0.5.0 - WIP
* Trim messages to fit console width
  * Attempts to handle console resizes
* Fixed race conditions: Dynamic updates safe

## 0.4 - 17 July 2018
Added UpdateMessage(string)
Handle ctrl-c interrupts
Added SetAbortMessage(string) to set message shown on ctrl-c
Updated documentation

## 0.3 - 15 July 2018
Made message for New/NewSpinner optional
Start() now takes an optional message
Removed Restart()
Use '>' as defualt success symbol on Windows
Use '!' as defualt error symbol on Windows

## 0.2.1 - 11 July 2018
Issue warning if attempting to stop a non-running spinner

## 0.2 - 10 July 2018
Added custom Sussess/Error functions: Successf() and Error()

## 0.1 - 6 July 2018
Initial Release
