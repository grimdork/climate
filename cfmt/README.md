# climate/cfmt
Simple ANSI color formatting using printf-style templates.

`cfmt` allows you to add color and style to your terminal output using a clean, readable template syntax. It handles the ANSI escape codes for you, keeping your Go code free of messy control characters.

## Installation
```bash
go get [github.com/grimdork/climate/cfmt](https://github.com/grimdork/climate/cfmt)
```

## Usage
The package uses a simple `{{.Color}}` syntax. Always remember to include {{.Reset}} to prevent color bleeding into the rest of the terminal.

### Basic Example
```go
package main

import "[github.com/grimdork/climate/cfmt](https://github.com/grimdork/climate/cfmt)"

func main() {
	// Simple colored output
	cfmt.Printf("{{.Red}}Error:{{.Reset}} file not found\n")
	cfmt.Printf("{{.Green}}Success:{{.Reset}} Configuration loaded\n")

	// Using styles
	cfmt.Printf("{{.Bold}}{{.Yellow}}Warning:{{.Reset}} This action is permanent.\n")
}
```

### Available Tags
| Category|	Tags |
| :---- | :---- |
| Reset | {{.reset}} (Clears all formatting) |
| Text | {{.black}}, {{.red}}, {{.green}}, {{.yellow}}, {{.blue}}, {{.magenta}}, {{.cyan}}, {{.white}} |
| Light text | {{.grey}}, {{.gray}}, {{.lred}}, {{.lgreen}}, {{.Lyellow}}, {{.lblue}}, {{.lmagenta}}, {{.lcyan}}, {{.lwhite}} |
| Background | {{.bgblack}}, {{.bgred}}, {{.bggreen}}, {{.bgyellow}}, {{.bgblue}}, {{.bgmagenta}}, {{.bgcyan}}, {{.bgwhite}}
| Light background | {{.bggrey}}, {{.bggray}}, {{.bglred}}, {{.bglgreen}}, {{.bglyellow}}, {{.bglblue}}, {{.bglmagenta}}, {{.bglcyan}}, {{.bglwhite}}
| Styles | {{.bold}}, {{.fuzzy}}, {{.italic}}, {{.under}}, {{.blink}}, {{.fast}}, {{.reverse}} {{.conceal}}, {{.strike}} |
