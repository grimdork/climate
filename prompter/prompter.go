package prompter

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// readPasswordFunc is the function signature for reading a secret from the terminal.
type readPasswordFunc func() ([]byte, error)

// Prompter is the interface for prompting the user for input.
type Prompter struct {
	reader   *bufio.Reader
	output   io.Writer
	readPass readPasswordFunc
	// Questions to ask.
	Questions []string
	// Secret sets the corresponding question to be a secret, hiding input.
	Secret []bool
	// Answers from the user.
	Answers []string
}

// Question is a question to ask the user.
type Question struct {
	// Question to present. Include a question mark if necessary.
	Question string
	// Secret sets the question to be a secret, hiding input.
	Secret bool
	// Default answer, if any.
	Default string
}

// New creates a new Prompter that reads from stdin and writes prompts to stdout.
func New(q []Question) *Prompter {
	pr := &Prompter{
		reader:   bufio.NewReader(os.Stdin),
		output:   os.Stdout,
		readPass: readPassword,
	}
	for _, q := range q {
		pr.Questions = append(pr.Questions, q.Question)
		pr.Secret = append(pr.Secret, q.Secret)
		pr.Answers = append(pr.Answers, q.Default)
	}

	return pr
}

// NewWithReader creates a Prompter with a custom input reader, output writer, and password reader.
// Useful for testing. If output is nil, prompts are discarded. If readPass is nil, it falls back
// to reading a line from the input reader (without echo suppression).
func NewWithReader(q []Question, input io.Reader, output io.Writer, readPass readPasswordFunc) *Prompter {
	if output == nil {
		output = io.Discard
	}
	if readPass == nil {
		r := bufio.NewReader(input)
		readPass = func() ([]byte, error) {
			line, err := r.ReadString('\n')
			if len(line) > 0 && line[len(line)-1] == '\n' {
				line = line[:len(line)-1]
			}
			return []byte(line), err
		}
	}
	pr := &Prompter{
		reader:   bufio.NewReader(input),
		output:   output,
		readPass: readPass,
	}
	for _, q := range q {
		pr.Questions = append(pr.Questions, q.Question)
		pr.Secret = append(pr.Secret, q.Secret)
		pr.Answers = append(pr.Answers, q.Default)
	}

	return pr
}

// Ask the user for input.
func (pr *Prompter) Ask() error {
	var t string
	var err error
	for i, q := range pr.Questions {
		fmt.Fprintf(pr.output, "%s [%s]: ", q, pr.Answers[i])
		if pr.Secret[i] {
			sec, err := pr.readPass()
			fmt.Fprintln(pr.output)
			if err != nil {
				return err
			}

			t = string(sec)
		} else {
			t, err = pr.reader.ReadString('\n')
			if err != nil {
				return err
			}
		}

		if len(t) > 0 && !pr.Secret[i] {
			t = t[:len(t)-1]
		}

		if t != "" {
			pr.Answers[i] = t
		}
	}

	return nil
}
