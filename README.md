![Go](https://github.com/AlexeyGy/insert-date/workflows/Go/badge.svg)
# insert-date

A slim Go application for Windows that lets you insert the current date (ISO8601) into any text. E.g., 2020-01-20

## building

use build.bat this will generate an .exe file in /dist

## Autostart
open your startup folder via `shell:startup` and drag and drop the `.exe` there. [see the Windows docs for more info](https://support.microsoft.com/en-us/help/4558286/windows-10-add-an-app-to-run-automatically-at-startup)

## under the hood

We use the Windows syscalls SendInput
