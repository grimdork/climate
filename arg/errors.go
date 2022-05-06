package arg

import "errors"

var (
	// ErrNoArgs is returned when no arguments are provided.
	ErrNoArgs = errors.New("no arguments provided")
	// ErrMissingRequired is returned when a required option is missing.
	ErrMissingRequired = errors.New("missing required option")
	// ErrMissingParam is returned when an option is missing an argument.
	ErrMissingParam = errors.New("missing option parameter")
	// ErrMissingFunc is returned when a command is missing a function.
	ErrMissingFunc = errors.New("missing function")
	// ErrLongShort is returned when a short option is longer than one character.
	ErrLongShort = errors.New("short option must be one character")
	// ErrUnknownOption is returned when an undefined option is encountered.
	ErrUnknownOption = errors.New("unknown option")
	// ErrEmptyLong is returned when a long option is empty.
	ErrEmptyLong = errors.New("long option without a string")
	// ErrShortLong is returned when a long option is shorter than two characters.
	ErrShortLong = errors.New("long option must be at least two characters")
	// ErrUnknownType is returned when an unknown option variable type is encountered.
	ErrUnknownType = errors.New("unknown option type")
	// ErrNoPlaceholder is returned when a positional argument is missing a placeholder.
	ErrNoPlaceholder = errors.New("no placeholder")
)
