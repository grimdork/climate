package prompter

import (
	"bufio"
	"fmt"
	"os"
	"syscall"

	"golang.org/x/term"
)

// Prompter is the interface for prompting the user for input.
type Prompter struct {
	*bufio.Reader
	// Questions to ask.
	Questions []string
	// Secret sets the corresponding question to be a secret, hiding input.
	Secret []bool
	// Answers from the user.
	Answers []string
	// Defaults for each question.
	Defaults []string
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

// New creates a new Prompter.
func New(q []Question) *Prompter {
	pr := &Prompter{}
	pr.Reader = bufio.NewReader(os.Stdin)
	for _, q := range q {
		pr.Questions = append(pr.Questions, q.Question)
		pr.Secret = append(pr.Secret, q.Secret)
		pr.Answers = append(pr.Answers, q.Default)
		pr.Defaults = append(pr.Defaults, q.Default)
	}

	return pr
}

// Ask the user for input.
func (pr *Prompter) Ask() error {
	var t string
	var err error
	for i, q := range pr.Questions {
		fmt.Printf("%s [%s]: ", q, pr.Answers[i])
		if pr.Secret[i] {
			sec, err := term.ReadPassword(int(syscall.Stdin))
			println()
			if err != nil {
				return err
			}

			t = string(sec)
		} else {
			t, err = pr.ReadString('\n')
			if err != nil {
				return err
			}
		}

		if t == "" {
			pr.Answers[i] = pr.Defaults[i]
		} else {
			pr.Answers[i] = t
		}
	}

	return nil
}
