package cfmt

const (
	//
	// Foreground
	//

	// Reset all attributes
	Reset = "\x1b[0m"
	// Black text colour
	Black = "\x1b[30;1m"
	// Red text colour
	Red = "\x1b[31;1m"
	// Green text colour
	Green = "\x1b[32;1m"
	// Yellow text colour
	Yellow = "\x1b[33;1m"
	// Blue text colour
	Blue = "\x1b[34;1m"
	// Magenta text colour
	Magenta = "\x1b[35;1m"
	// Cyan text colour
	Cyan = "\x1b[36;1m"
	// White text colour
	White = "\x1b[37;1m"

	// Grey text colour
	Grey = "\x1b[90;1m"
	// LightRed text colour
	LightRed = "\x1b[91;1m"
	// LightGreen text colour
	LightGreen = "\x1b[92;1m"
	// LightYellow text colour
	LightYellow = "\x1b[93;1m"
	// LightBlue text colour
	LightBlue = "\x1b[94;1m"
	// LightMagenta text colour
	LightMagenta = "\x1b[95;1m"
	// LightCyan text colour
	LightCyan = "\x1b[96;1m"
	// LightWhite text colour
	LightWhite = "\x1b[97;1m"

	//
	// Background colours
	//

	// BGBlack is the background colour black.
	BGBlack = "\x1b[40;5m"
	// BGRed is the background colour red.
	BGRed = "\x1b[41;5m"
	// BGGreen is the background colour green.
	BGGreen = "\x1b[42;5m"
	// BGYellow is the background colour yellow.
	BGYellow = "\x1b[43;5m"
	// BGBlue is the background colour blue.
	BGBlue = "\x1b[44;5m"
	// BGMagenta is the background colour magenta.
	BGMagenta = "\x1b[45;5m"
	// BGCyan is the background colour cyan.
	BGCyan = "\x1b[46;5m"
	// BGWhite is the background colour white.
	BGWhite = "\x1b[47;5m"

	// BGGrey is the background colour grey.
	BGGrey = "\x1b[100;5m"
	// BGLightRed is the background colour light red.
	BGLightRed = "\x1b[101;5m"
	// BGLightGreen is the background colour light green.
	BGLightGreen = "\x1b[102;5m"
	// BGLightYellow is the background colour light yellow.
	BGLightYellow = "\x1b[103;5m"
	// BGLightBlue is the background colour light blue.
	BGLightBlue = "\x1b[104;5m"
	// BGLightMagenta is the background colour light magenta.
	BGLightMagenta = "\x1b[105;5m"
	// BGLightCyan is the background colour light cyan.
	BGLightCyan = "\x1b[106;5m"
	// BGLightWhite is the background colour light white.
	BGLightWhite = "\x1b[107;5m"

	//
	// Style options
	//

	// Bold text style.
	Bold = "\x1b[1;1m"
	// Fuzzy text style (dim or faint).
	Fuzzy = "\x1b[2;1m"
	// Italic text style.
	Italic = "\x1b[3;1m"
	// Underscore (underline) text style.
	Underscore = "\x1b[4;1m"
	// Blink text style. Use sparingly.
	Blink = "\x1b[5;1m"
	// FastBlink text style. Use even more sparingly.
	FastBlink = "\x1b[6;1m"
	// Reverse text style.
	Reverse = "\x1b[7;1m"
	// Concealed text style (hidden). Useful for passwords or spoilers.
	Concealed = "\x1b[8;1m"
	// Strikethrough text style.
	Strikethrough = "\x1b[9;1m"
)