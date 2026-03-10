# climate/prompter
A simple, interactive tool for terminal-based user input.

`prompter` makes it easy to ask questions in the terminal, handle default answers when the user just hits "Enter," and securely hide sensitive input like passwords.

## Installation
```bash
go get github.com/grimdork/climate/prompter
```

## Core Features

### Basic questions
Ask a question and provide a fallback value if the user provides no input.

```Go
name := prompter.Ask("What is your name?", "Guest", false)
fmt.Printf("Hello, %s!\n", name)
```

### Secure input (passwords)
Mask the user's keystrokes for sensitive information.

```go
// The third argument 'true' enables hidden input
dbPass := prompter.Ask("Database Password:", "", true)
```

### Smart Defaults
If a default value is provided, it is typically displayed in brackets:
```
What is your name? [Guest]:
```

If the user simply presses Enter, the default value is returned.

## Technical Details
### Handling "Hidden" Input

Under the hood, prompter uses terminal escape codes or raw mode (depending on the OS) to ensure that characters typed by the user do not echo back to the screen. This is essential for:

- Database credentials
- API keys
- Authentication tokens

---

Unlike larger libraries that pull in heavy terminal-control dependencies, prompter stays lean. It uses standard os.Stdin and bufio.Scanner patterns, making it highly compatible with TinyGo and small-footprint binaries.
