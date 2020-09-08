![Go](https://github.com/AlexeyGy/insert-date/workflows/Go/badge.svg)
# insert-date

A slim Go application for Windows that lets you insert the current date, e.g., 2020-01-20
 (ISO8601) into any text by pressing `CTRL+ALT+D`.
 
## Download
see the release tab on the right to get the latest release

## building

use build.bat this will generate an .exe file in /dist

## Autostart
Open your startup folder via `WIN+R`, type `shell:startup`, and drag and drop the `.exe` there. [see the Windows docs for more info](https://support.microsoft.com/en-us/help/4558286/windows-10-add-an-app-to-run-automatically-at-startup)

## under the hood

We use the Windows syscalls 
 - `RegisterHotKey` to register the hotkey `CTRL+ALT+D`
 - `GetMessageW` to listen for any input events
 - `SendInput` to send keyboard events
 
