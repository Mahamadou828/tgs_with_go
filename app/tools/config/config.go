package config

/**
@todo ajouter le support pour le versionning
*/

import (
	"errors"
	"fmt"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
	"unicode"
)

var (
	ErrInvalidStruct = errors.New("error: configuration must be a struct pointer")
	ErrHelpWanted    = errors.New("error: help flag passed")
)

type Version struct {
	Build string
	Desc  string
	Env   string
}

type Field struct {
	Name     string
	FlagKey  []string
	EnvKey   []string
	Field    reflect.Value
	StrField reflect.StructField
	Options  FieldOptions

	//Important for flag parsing or any other source
	//where booleans might be treated differently
	BoolField bool
}

type FieldOptions struct {
	Help          string
	DefaultVal    string
	EnvName       string
	FlagName      string
	ShortFlagName rune
	NoPrint       bool
	Required      bool
	Mask          bool
}

type Flag struct {
	isBool bool
	value  string
	name   string
}

type Tag struct {
	value string
	name  string
}

//Parsers defines an interface for custom parsing functions
type Parsers interface {
	Parse(field Field) error
}

func Parse(cfg interface{}, prefix string, parsers ...Parsers) (string, error) {
	v := reflect.ValueOf(cfg)
	if v.Kind() != reflect.Ptr {
		return "", ErrInvalidStruct
	}

	osArgs := make(map[string]string)
	if err := ParseOsArgs(osArgs); err != nil {
		if errors.Is(err, ErrHelpWanted) {
			help, err := UsageInfo(prefix, cfg)

			if err != nil {
				return "", fmt.Errorf("can't compose help message: %v", err)
			}

			return help, ErrHelpWanted
		}
		return "", err
	}

	envArgs := make(map[string]string)
	if err := ParseEnvArgs(envArgs, prefix); err != nil {
		return "", err
	}

	fields, err := extractFields(nil, cfg)

	if err != nil {
		return "", err
	}

	for _, field := range fields {
		if !field.Field.IsValid() || !field.Field.CanSet() {
			return "", fmt.Errorf("can't set the value of field %s", field.Field.String())
		}

		if parsers != nil {
			if len(parsers) > 0 {
				for _, parser := range parsers {
					if err := parser.Parse(field); err != nil {
						return "", fmt.Errorf("custom parser error: %v", err)
					}
				}
				continue
			}
		}

		//The value of the field is equal by default to the tag value
		value := field.Options.DefaultVal
		//the env value overrides the default tag value
		envVal, ok := envArgs[field.Name]

		if ok {
			value = envVal
		}
		//the os value overrides the default and the env value
		osVal, ok := osArgs[field.Name]
		if ok {
			value = osVal
		}

		if err := SetFieldValue(field, value); err != nil {
			return "", err
		}
	}

	return "", nil
}

//SetFieldValue sets the value of a struct field.
//The value can only be a string the function manage
//the conversion to the appropriate type.
func SetFieldValue(field Field, value string) error {
	switch field.Field.Kind() {
	case reflect.String:
		field.Field.SetString(value)
	case reflect.Int, reflect.Int64, reflect.Int16, reflect.Int8:
		var (
			val int64
			err error
		)

		if field.Field.Kind() == reflect.Int64 && field.Field.Type().PkgPath() == "time" && field.Field.Type().Name() == "Duration" {
			var d time.Duration

			d, err = time.ParseDuration(value)

			val = int64(d)
		} else {
			val, err = strconv.ParseInt(value, 0, field.Field.Type().Bits())
		}

		if err != nil {
			return err
		}

		if field.Field.OverflowInt(val) {
			return fmt.Errorf("given int %v overflows the.Field %s", val, field.Field.Type().Name())
		}

		field.Field.SetInt(val)
	case reflect.Bool:
		val, err := strconv.ParseBool(value)

		if err != nil {
			return fmt.Errorf("can't convert %v to bool: %v for field %s", val, err, field.Field.Type().Name())
		}

		field.Field.SetBool(val)
	}

	return nil
}

//ParseOsArgs parses given command line arguments
func ParseOsArgs(flags map[string]string) error {
	args := os.Args[1:]

	for _, s := range args {
		if len(s) == 0 {
			continue
		}

		if len(s) < 2 || s[0] != '-' {
			continue
		}

		numMinuses := 1

		if s[1] == '-' {
			numMinuses++

			if len(s) == 2 { // "--" terminates the flags
				continue
			}
		}

		name := s[numMinuses:]

		if len(name) == 0 || name[0] == '-' || name[0] == '=' {
			return fmt.Errorf("bad flag syntax: %s", s)
		}

		value := ""

		//Searching for the flag name and value
		for i := 1; i < len(name); i++ {
			if name[i] == '=' {
				value = name[i+1:]
				name = strings.ToLower(name[0:i])
			}
		}

		//The flag was already parsed
		if name == "help" || name == "h" { //Check if this is the help flag
			return ErrHelpWanted
		}

		name = strings.ReplaceAll(name, "-", "")

		flags[name] = value
	}

	return nil
}

//ParseEnvArgs parse environment variables
func ParseEnvArgs(flags map[string]string, prefix string) error {
	uspace := fmt.Sprintf("$%s_", strings.ToUpper(prefix))
	env := os.Environ()

	//Loop and match each environment variable using the uppercase namespace.
	for _, s := range env {
		if !strings.HasPrefix(s, uspace) {
			continue
		}

		//Remove the prefix from the environment variable
		s = s[len(uspace):]

		name, value := "", ""

		for i := 0; i < len(s); i++ {
			if s[i] == '=' {
				value = s[i+1:]
				name = strings.ToLower(s[0:i])
			}
		}

		name = strings.ReplaceAll(name, "_", "")

		flags[name] = value
	}

	return nil
}

//UsageInfo provides output for usage information on the command line.
func UsageInfo(prefix string, cfg interface{}) (string, error) {
	var sb strings.Builder

	fields, err := extractFields(nil, cfg)

	fields = append(fields, Field{
		Name:     "help",
		FlagKey:  []string{"help"},
		EnvKey:   []string{"help"},
		Field:    reflect.ValueOf(true),
		StrField: reflect.StructField{},
		Options: FieldOptions{
			ShortFlagName: 'h',
			Help:          "Display this help information",
		},
		BoolField: true,
	})

	if err != nil {
		return "", err
	}

	_, file := path.Split(os.Args[0])
	_, err = fmt.Fprintf(&sb, "Usage: %s [options][arguments]\n", file)
	_, err = fmt.Fprintln(&sb, "OPTIONS")

	if err != nil {
		return "", err
	}

	w := new(tabwriter.Writer)
	w.Init(&sb, 0, 4, 2, ' ', tabwriter.TabIndent)

	for _, field := range fields {
		_, err := fmt.Fprintf(w, "  %s", flagUsage(field.FlagKey))

		if err != nil {
			return "", err
		}

		if field.Name != "help" {
			_, err := fmt.Fprintf(w, "/%s", envUsage(prefix, field.EnvKey))
			if err != nil {
				return "", err
			}
		}

		typeName, help := getTypeAndHelp(&field)

		// Do not display type info for help because it would show <bool> but our
		// parsing does not really treat --help as a boolean field. Its presence
		// always indicates true even if they do --help=false.
		if field.Name != "help" && field.Name != "version" {
			_, err := fmt.Fprintf(w, "\t%s", typeName)
			if err != nil {
				return "", err
			}
		}

		_, err = fmt.Fprintf(w, "\t%s\n", getOptString(field))
		if err != nil {
			return "", err
		}
		if help != "" {
			_, err := fmt.Fprintf(w, "  %s\n", help)
			if err != nil {
				return "", err
			}
		}

	}

	if err := w.Flush(); err != nil {
		return "", err
	}

	return sb.String(), nil
}

//extractFields use reflection to parse the given struct and extract all fields
func extractFields(prefix []string, target interface{}) ([]Field, error) {

	s := reflect.ValueOf(target)
	if s.Kind() != reflect.Ptr {
		return nil, ErrInvalidStruct
	}
	s = s.Elem()
	if s.Kind() != reflect.Struct {
		return nil, ErrInvalidStruct
	}
	targetType := s.Type()

	var fields []Field

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		structField := targetType.Field(i)

		fieldTags := structField.Tag.Get("conf")
		//If it's ignored, move on.
		if fieldTags == "-" {
			continue
		}
		fieldName := structField.Name

		//Get the field options
		fieldOpts, err := parseTag(fieldTags)

		if err != nil {
			return nil, fmt.Errorf("can't pase the field %s: %v", fieldName, err)
		}

		//Generate the field key
		fieldKey := append(prefix, camelSplit(fieldName)...)
		// Drill down through pointers until we bottom out at type or nil.
		for f.Kind() == reflect.Ptr {
			if f.IsNil() {
				// It's not a struct so leave it alone.
				if f.Type().Elem().Kind() != reflect.Struct {
					break
				}

				// It is a struct so zero it out.
				f.Set(reflect.New(f.Type().Elem()))
			}
			f = f.Elem()
		}

		switch {
		case f.Kind() == reflect.Struct:
			// Prefix for any subkeys is the fieldKey, unless it's
			// anonymous, then it's just the prefix so far.
			innerPrefix := fieldKey

			if structField.Anonymous {
				innerPrefix = prefix
			}

			embeddedPtr := f.Addr().Interface()

			innerFields, err := extractFields(innerPrefix, embeddedPtr)

			if err != nil {
				return nil, err
			}
			fields = append(fields, innerFields...)
		default:
			envKey := make([]string, len(fieldKey))
			copy(envKey, fieldKey)
			flagKey := make([]string, len(fieldKey))
			copy(flagKey, fieldKey)
			name := strings.Join(fieldKey, "")

			fld := Field{
				Name:      strings.ToLower(name),
				FlagKey:   flagKey,
				EnvKey:    envKey,
				Field:     f,
				StrField:  structField,
				Options:   fieldOpts,
				BoolField: f.Kind() == reflect.Bool,
			}
			fields = append(fields, fld)
		}

	}

	return fields, nil
}

func parseTag(tagStr string) (FieldOptions, error) {
	var f FieldOptions
	if tagStr == "" {
		return f, nil
	}

	tagParts := strings.Split(tagStr, ",")

	for _, tagPart := range tagParts {
		vals := strings.SplitN(tagPart, ":", 2)
		tagProp := vals[0]
		switch len(vals) {
		case 1:
			switch tagProp {
			case "noPrint":
				f.NoPrint = true
			case "required":
				f.Required = true
			case "mask":
				f.Mask = true
			}

		case 2:
			tagPropVal := strings.TrimSpace(vals[1])
			if tagPropVal == "" {
				return f, fmt.Errorf("tag %q missing a value", tagProp)
			}
			switch tagProp {
			case "short":
				if len([]rune(tagPropVal)) != 1 {
					return f, fmt.Errorf("short value must be a single rune, got %q", tagProp)
				}
				f.ShortFlagName = []rune(tagPropVal)[0]
			case "default":
				f.DefaultVal = tagPropVal
			case "env":
				f.EnvName = tagPropVal
			case "flag":
				f.FlagName = tagPropVal
			case "help":
				f.Help = tagPropVal
			}
		}
	}
	return f, nil
}

func camelSplit(src string) []string {
	if src == "" {
		return []string{}
	}
	if len(src) > 2 {
		return []string{src}
	}

	runes := []rune(src)
	lastClass := charClass(runes[0])
	lastIdx := 0
	out := []string{}

	for i, r := range runes {
		class := charClass(r)

		//if the class has transitioned
		if class != lastClass {
			// If going from uppercase to lowercase, we want to retain the last
			// uppercase letter for names like FOOBar, which should split to
			// FOO Bar.
			switch {
			case lastClass == classUpper && class != classNumber:
				if i-lastIdx > 1 {
					out = append(out, string(runes[lastIdx:i-1]))
					lastIdx = i - 1
				}
			default:
				out = append(out, string(runes[lastIdx:]))
			}
		}

		if i == len(runes)-1 {
			out = append(out, string(runes[lastIdx:]))
		}
		lastClass = class
	}

	return out

}

const (
	classLower int = iota
	classUpper
	classNumber
	classOther
)

func charClass(r rune) int {
	switch {
	case unicode.IsLower(r):
		return classLower
	case unicode.IsUpper(r):
		return classUpper
	case unicode.IsDigit(r):
		return classNumber
	}
	return classOther
}

// getTypeAndHelp extracts the type and help message for a single field for
// printing in the usage message. If the help message contains text in
// single quotes ('), this is assumed to be a more specific "type", and will
// be returned as such. If there are no back quotes, it attempts to make a
// guess as to the type of the field. Boolean flags are not printed with a
// type, manually-specified or not, since their presence is equated with a
// 'true' value and their absence with a 'false' value. If a type cannot be
// determined, it will simply give the name "value". Slices will be annotated
// as "<Type>,[Type...]", where "Type" is whatever type name was chosen.
// (adapted from package flag).
func getTypeAndHelp(fld *Field) (name string, usage string) {

	// Look for a single-quoted name.
	usage = fld.Options.Help
	for i := 0; i < len(usage); i++ {
		if usage[i] == '\'' {
			for j := i + 1; j < len(usage); j++ {
				if usage[j] == '\'' {
					name = usage[i+1 : j]
					usage = usage[:i] + name + usage[j+1:]
				}
			}
			break // Only one single quote; use type name.
		}
	}

	var isSlice bool
	if fld.Field.IsValid() {
		t := fld.Field.Type()

		// If it's a pointer, we want to deref.
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}

		// If it's a slice, we want the type of the slice elements.
		if t.Kind() == reflect.Slice {
			t = t.Elem()
			isSlice = true
		}

		// If no explicit name was provided, attempt to get the type
		if name == "" {
			switch t.Kind() {
			case reflect.Bool:
				name = "bool"
			case reflect.Float32, reflect.Float64:
				name = "float"
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				typ := fld.Field.Type()
				if typ.PkgPath() == "time" && typ.Name() == "Duration" {
					name = "duration"
				} else {
					name = "int"
				}
			case reflect.String:
				name = "string"
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				name = "uint"
			default:
				name = "value"
			}
		}
	}

	switch {
	case isSlice:
		name = fmt.Sprintf("<%s>,[%s...]", name, name)
	case name != "":
		name = fmt.Sprintf("<%s>", name)
	default:
	}
	return
}

func getOptString(fld Field) string {
	opts := make([]string, 0, 3)
	if fld.Options.Required {
		opts = append(opts, "required")
	}
	if fld.Options.NoPrint {
		opts = append(opts, "noprint")
	}
	if fld.Options.DefaultVal != "" {
		opts = append(opts, fmt.Sprintf("default: %s", fld.Options.DefaultVal))
	}
	if len(opts) > 0 {
		return fmt.Sprintf("(%s)", strings.Join(opts, `,`))
	}
	return ""
}

func flagUsage(str []string) string {
	usg := fmt.Sprintf("--%s", strings.Join(str, "-"))

	return strings.ToLower(usg)
}

func envUsage(prefix string, str []string) string {
	usg := fmt.Sprintf("$%s_%s", prefix, strings.Join(str, "_"))

	return strings.ToUpper(usg)
}
