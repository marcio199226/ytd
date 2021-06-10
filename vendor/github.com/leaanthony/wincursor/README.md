# WinCursor

This is a simple library to show/hide the cursor in the windows terminal. Works in both Command Prompt and PowerShell.

## Installation

```
go get -u github.com/leaanthony/wincursor
```

## Usage

It's pretty simple really...

```
import github.com/leaanthony/wincursor"

func main() {
    wincursor.Hide() // This will hide the cursor
    wincursor.Show() // This will show the cursor
}
```

## Rationale

I was making a cross-platform spinner and wanted to hide the cursor on Windows. This turned out to be a lot harder than it should have been. Once I got it working, I decided to pull out the relevant functions into a new package so others don't need to go through the same pain 😀

## With a little help from my friends

Original idea from this awesome [stackoverflow answer].

Some code on how to implement that from [go-ansiterm].

[stackoverflow answer]: https://stackoverflow.com/a/10455937
[go-ansiterm]: https://github.com/Azure/go-ansiterm

## go-ansiterm license

The MIT License (MIT)

Copyright (c) 2015 Microsoft Corporation

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.