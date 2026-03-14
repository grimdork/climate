package ini

import "fmt"

// Field contains a variable and its data.
type Field struct {
	// Value will be stripped of surrounding whitespace when loaded.
	Value string
	// Type lets the user choose which Get* method to use when loading unknown files.
	Type   byte
	boolV  bool
	intV   int64
	floatV float64
}

// GetBool returns a field as a bool.
func (f *Field) GetBool() bool {
	return f.boolV
}

// SetBool sets a field to a bool.
func (f *Field) SetBool(value bool) {
	f.boolV = value
	f.Type = Bool
	f.Value = fmt.Sprintf("%t", value)
}

// GetInt returns a field as an int.
func (f *Field) GetInt() int64 {
	return f.intV
}

// SetInt sets a field as an int.
func (f *Field) SetInt(value int64) {
	f.intV = value
	f.Type = Int
	f.Value = fmt.Sprintf("%d", value)
}

// GetFloat returns a field as a float64.
func (f *Field) GetFloat() float64 {
	return f.floatV
}

// SetFloat sets a field to a float64.
func (f *Field) SetFloat(value float64) {
	f.floatV = value
	f.Type = Float
	f.Value = fmt.Sprintf("%f", value)
}

// SetString sets a field as a string.
func (f *Field) SetString(value string) {
	f.Value = value
	f.Type = String
}
