# climate/paths
Cross-platform directory resolution for configuration and data.

`paths` ensures your application stores its data in the correct place depending on the operating system, following platform-specific conventions.

## Installation
```bash
go get github.com/grimdork/climate/paths
```

## Usage

### Configuration paths
```go
package main

import (
	"fmt"
	"github.com/grimdork/climate/paths"
)

func main() {
	p, err := paths.New("myapp")
	if err != nil {
		panic(err)
	}

	fmt.Println("User config:", p.UserBase)
	fmt.Println("Server config:", p.ServerBase)
}
```

Default paths by platform:

| OS | UserBase | ServerBase |
| :--- | :--- | :--- |
| macOS | `~/Library/Application Support/myapp` | `~/Library/Application Support/myapp` |
| Linux/Unix | `$XDG_CONFIG_HOME/myapp` or `~/.config/myapp` | `/etc/myapp` |
| Windows | `%AppData%\myapp` | `%ProgramData%\myapp` |

### Custom base paths
Override the defaults if needed:
```go
p, _ := paths.New("myapp")
p.SetBase("/opt/configs")      // UserBase becomes /opt/configs/myapp
p.SetServerBase("/srv/data")   // ServerBase becomes /srv/data/myapp
```

Pass an empty string to reset to the platform default.

### File and directory checks
Utility functions for checking existence:
```go
// Returns (exists, isDir)
exists, isDir := paths.Exists("/some/path")

// Convenience wrappers
if paths.DirExists("/etc/myapp") { ... }
if paths.FileExists("/etc/myapp/config.json") { ... }
```

## Zero dependencies
Uses only `os`, `os/user`, and `path/filepath` from the standard library.
