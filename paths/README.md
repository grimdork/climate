# climate/paths
Cross-platform directory resolution for configuration and data.

`paths` ensures your application stores its data in the "correct" place depending on the operating system. It helps avoid cluttering the user's home directory by following platform-specific conventions.

## Installation
```bash
go get github.com/grimdork/climate/paths
```

## Why use paths?
Different operating systems have different standards for where local application data should live.
- On macOS: Users expect data in ~/Library/Application Support/.
- On Linux/Unix: Users expect hidden "dot-directories" in ~/ or XDG-compliant paths.

## Usage

### Get a configuration directory
This function returns the path to a directory named after your app. If the directory doesn't exist, it creates it with 0700 permissions.

```Go
package main

import (
	"fmt"
	"github.com/grimdork/climate/paths"
)

func main() {
	// Returns:
	// - ~/Library/Application Support/mytool (macOS)
	// - ~/.mytool (Linux)
	// - $XDG_CONFIG_HOME/.mytool (Linux with XDG_CONFIG_HOME set)
	configPath, err := paths.GetConfigDir("mytool")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Store your config at: %s/config.json\n", configPath)
}
```

### Standard home resolution
A simple helper to get the current user's home directory without manually checking environment variables.

```go
home := paths.Home()
```

### Technical Details
Platform Mapping

| OS | Resulting Path |
| :-- | :-- |
| macOS | /Users/\<user>/Library/Application Support/\<appname> |
| Linux/Unix | $XDG_CONFIG_HOME/\<appname> or ~/.config/\<appname> |
| Linux server | /etc/\<appname> |

### Zero dependencies
Like all climate packages, this relies strictly on the Go standard library (os, runtime, and path/filepath).

### Best Practices
When using GetConfigDir, it is recommended to use your application's binary name as the argument. This keeps your data isolated and easy for the user to find if they need to manually edit a config file.
