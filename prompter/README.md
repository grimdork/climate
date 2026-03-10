# climate/prompter
A simple, interactive tool for terminal-based user input.

`prompter` makes it easy to ask questions in the terminal, handle default answers when the user just hits "Enter," and securely hide sensitive input like passwords.

## Installation
```bash
go get github.com/grimdork/climate/prompter
```

## Usage

### Define questions and ask
```go
package main

import (
	"fmt"
	"github.com/grimdork/climate/prompter"
)

func main() {
	pr := prompter.New([]prompter.Question{
		{Question: "Username", Default: "admin"},
		{Question: "Password", Secret: true},
	})

	err := pr.Ask()
	if err != nil {
		panic(err)
	}

	fmt.Printf("User: %s\n", pr.Answers[0])
	fmt.Printf("Pass: %s\n", pr.Answers[1])
}
```

Output:
```
Username [admin]: alice
Password []:
User: alice
Pass: s3cret
```

If the user presses Enter without typing anything, the default value is kept.

### Secret input
Any question with `Secret: true` will suppress terminal echo, so typed characters are not visible. This uses raw terminal ioctl — no external dependencies.

### Testing with custom I/O
For testing or automation, use `NewWithReader` to inject a custom input source:

```go
input := strings.NewReader("testuser\n")
fakePass := func() ([]byte, error) {
	return []byte("testpass"), nil
}

pr := prompter.NewWithReader([]prompter.Question{
	{Question: "User", Default: "admin"},
	{Question: "Pass", Secret: true},
}, input, nil, fakePass)

pr.Ask()
// pr.Answers[0] == "testuser"
// pr.Answers[1] == "testpass"
```

Parameters: input reader, output writer (nil discards prompts), and an optional password reader function (nil falls back to reading a line from the input).

## Platform support
Works on macOS, Linux, FreeBSD, NetBSD, OpenBSD, and Dragonfly BSD.
