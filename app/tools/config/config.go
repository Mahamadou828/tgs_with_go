package config

/**
Comment parser un struc avec des flags
comment lire les commands line argument
@todo ajouter le support pour les string, int, bool, time.Duration
@todo refacto pour improve readability
@todo ajouter les flags suivant: -, required
@todo ajouter le message help
*/

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
)

var (
	ErrInvalidStruct = errors.New("error: configuration must be a struct pointer")
	ErrHelpWanted    = errors.New("error: help flag passed")
)

type Flag struct {
	isBool bool
	value  any
	key    string
}

//@todo explorer ce que ca fait si on avait laisser *interface{}

func Parse(cfg interface{}, ssmSecrets map[string]string, prefix string) error {
	v := reflect.ValueOf(cfg)
	if v.Kind() != reflect.Ptr {
		return ErrInvalidStruct
	}

	flags := make(map[string]Flag)

	//Insert the aws ssm secret
	for key, secret := range ssmSecrets {
		flags[key] = Flag{
			isBool: false,
			value:  secret,
			key:    key,
		}
	}

	//Parsing the envArgs
	err := parseEnvArgs(flags, prefix)
	if err != nil {
		return err
	}

	//Parsing all the os args
	err = parseOsArgs(flags)
	if err != nil {
		if errors.Is(err, ErrHelpWanted) {
			//Start compose the help message
			return nil
		}
	}

	parseStruct(v.Elem(), flags)

	return nil
}

func parseOsArgs(flags map[string]Flag) error {
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

			if len(s) == '2' { // "--" terminates the flags
				continue
			}
		}

		name := s[numMinuses:]

		if len(name) == 0 || name[0] == '-' || name[0] == '=' {
			return fmt.Errorf("bad flag syntax: %s", s)
		}

		isBool, value := true, ""

		//Searching for the flag name and value
		for i := 1; i < len(name); i++ {
			if name[i] == '=' {
				value = name[i+1:]
				isBool = false
				name = strings.ToLower(name[0:i])
			}
		}

		_, ok := flags[name]

		//The flag was already parsed
		if ok {
			if name == "help" || name == "h" { //Check if this is the help flag
				return ErrHelpWanted
			}
		}

		r := strings.NewReplacer("_", "", "-", "")

		name = r.Replace(name)

		flags[name] = Flag{
			isBool: isBool,
			value:  value,
			key:    name,
		}

	}

	return nil
}

//parseEnvArgs parse environment variables
func parseEnvArgs(flags map[string]Flag, prefix string) error {
	uspace := fmt.Sprintf("$%s_", strings.ToUpper(prefix))
	env := os.Environ()

	//Loop and match each environment variable using the uppercase namespace.
	for _, s := range env {
		if !strings.HasPrefix(s, uspace) {
			continue
		}

		//Remove the $ sign
		s = s[1:]

		key, value := "", ""

		for i := 0; i < len(s); i++ {
			if s[i] == '=' {
				value = s[i+1:]
				key = strings.ToLower(s[0:i])
			}
		}

		r := strings.NewReplacer("_", "", "-", "")

		key = r.Replace(key)

		flags[key] = Flag{
			isBool: false,
			value:  value,
			key:    key,
		}
	}

	return nil
}

//parseStruct parses configuration into the specified struct
func parseStruct(v reflect.Value, flags map[string]Flag) error {
	for i := 0; i < v.NumField(); i++ {
		typ := v.Type().Field(i).Type

		//If the type is a struct, we should recursively
		//explore him
		if typ.Kind() == reflect.Struct {
			nv := v.Field(i)
			parseStruct(nv, flags)
			continue
		}

		sf := v.Type().Field(i)
		s, ok := sf.Tag.Lookup("conf")
		if !ok {
			return fmt.Errorf("no tag present for the field: %s", sf.Name)
		}

		name, defVal := "", ""
		for i := 0; i < len(s); i++ {
			if s[i] == ':' {
				name = s[0:i]
				defVal = s[i+1:]
			}
		}

		flgVal, ok := flags[name]
		if ok {
			field := v.Field(i)

			switch field.Kind() {
			case reflect.String:
				break
			case reflect.Bool:
				break
			case reflect.Int:
				break
			}

		}

	}

	return nil
}

func processStringField(field reflect.Value, flag Flag) error {
	name, defVal, err := processField(field)

	if err != nil {
		return err
	}
}

//@todo to refacto
func processField(field reflect.StructField) (name string, defVal string, err error) {
	s, ok := field.Tag.Lookup("conf")
	if !ok {
		return "", "", fmt.Errorf("no tag present for the field: %s", field.Name)
	}

	for i := 0; i < len(s); i++ {
		if s[i] == ':' {
			name = s[0:i]
			defVal = s[i+1:]
		}
	}

	return
}
