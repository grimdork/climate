package cfmt

const (
	// Foreground
	Reset   = "\x1b[0m"
	Black   = "\x1b[30;1m"
	Red     = "\x1b[31;1m"
	Green   = "\x1b[32;1m"
	Yellow  = "\x1b[33;1m"
	Blue    = "\x1b[34;1m"
	Magenta = "\x1b[35;1m"
	Cyan    = "\x1b[36;1m"
	White   = "\x1b[37;1m"

	Grey         = "\x1b[90;1m"
	LightRed     = "\x1b[91;1m"
	LightGreen   = "\x1b[92;1m"
	LightYellow  = "\x1b[93;1m"
	LightBlue    = "\x1b[94;1m"
	LightMagenta = "\x1b[95;1m"
	LightCyan    = "\x1b[96;1m"
	LightWhite   = "\x1b[97;1m"

	// Background
	BGBlack   = "\x1b[40;5m"
	BGRed     = "\x1b[41;5m"
	BGGreen   = "\x1b[42;5m"
	BGYellow  = "\x1b[43;5m"
	BGBlue    = "\x1b[44;5m"
	BGMagenta = "\x1b[45;5m"
	BGCyan    = "\x1b[46;5m"
	BGWhite   = "\x1b[47;5m"

	BGGrey         = "\x1b[100;5m"
	BGLightRed     = "\x1b[101;5m"
	BGLightGreen   = "\x1b[102;5m"
	BGLightYellow  = "\x1b[103;5m"
	BGLightBlue    = "\x1b[104;5m"
	BGLightMagenta = "\x1b[105;5m"
	BGLightCyan    = "\x1b[106;5m"
	BGLightWhite   = "\x1b[107;5m"

	// Other options
	Bold          = "\x1b[1;1m"
	Fuzzy         = "\x1b[2;1m"
	Italic        = "\x1b[3;1m"
	Underscore    = "\x1b[4;1m"
	Blink         = "\x1b[5;1m"
	FastBlink     = "\x1b[6;1m"
	Reverse       = "\x1b[7;1m"
	Concealed     = "\x1b[8;1m"
	Strikethrough = "\x1b[9;1m"
)
