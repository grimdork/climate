# climate/tab
Clean, elastic tabbed columns for terminal output.

`tab` takes a string of whitespace-separated text and formats it into aligned columns using Go's `tabwriter`.

## Installation
```bash
go get github.com/grimdork/climate/tab
```

## Usage

### Full column mode
Every whitespace boundary becomes a column break:
```go
input := "Name Role Status\nAlice Admin Active\nBob User Offline"
output, _ := tab.Tabulate(input, false)
fmt.Print(output)
// Name   Role   Status
// Alice  Admin  Active
// Bob    User   Offline
```

### Two-column mode
The first word is column one, everything else is column two:
```go
input := "serve Start the web server\nconfig Show current configuration\nhelp Display this message"
output, _ := tab.Tabulate(input, true)
fmt.Print(output)
// serve   Start the web server
// config  Show current configuration
// help    Display this message
```

Useful for command lists, key-description pairs, or any label+text layout.

### SplitColumns
If you need the parsed rows without formatting:
```go
rows, _ := tab.SplitColumns(input, false)
// rows is [][]string
```

### CSV input
Parse CSV with proper quoting support. The first row becomes a header with a separator line:
```go
csv := "Name,Role,Status\nAlice,Admin,Active\n\"Eve, Jr\",Guest,Pending"
output, _ := tab.TabulateCSV(csv)
fmt.Print(output)
// Name     Role   Status
// -------  -----  -------
// Alice    Admin  Active
// Eve, Jr  Guest  Pending
```

Handles quoted fields, commas inside values, and all standard CSV rules via Go's `encoding/csv`.
