package config

/**
@todo ajouter les flags suivant: -, required
@todo ajouter le message help
@todo ajouter le support pour le versionning
@todo ajouter le support de le faire disparetre des logs
*/

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidStruct = errors.New("error: configuration must be a struct pointer")
	ErrHelpWanted    = errors.New("error: help flag passed")
)

type Flag struct {
	isBool bool
	value  string
	name   string
}

type Tag struct {
	value string
	name  string
}

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
			name:   key,
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

	err = parseStruct(v.Elem(), flags)

	return err
}

//parseOsArgs parses given command line arguments
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

			if len(s) == 2 { // "--" terminates the flags
				continue
			}
		}

		name := s[numMinuses:]

		if len(name) == 0 || name[0] == '-' || name[0] == '=' {
			return fmt.Errorf("bad flag syntax: %s", s)
		}

		isBool, value := true, "true"

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
			name:   name,
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
			name:   key,
		}
	}

	return nil
}

//parseStruct parses configuration into the specified struct
func parseStruct(s reflect.Value, flags map[string]Flag) error {
	for i := 0; i < s.NumField(); i++ {
		field := s.Field(i)
		kind := field.Kind()
		typ := field.Type()
		strField := s.Type().Field(i)

		if kind == reflect.Struct {
			nv := s.Field(i)
			if err := parseStruct(nv, flags); err != nil {
				return err
			}
			continue
		}

		if !field.IsValid() || !field.CanSet() {
			return fmt.Errorf("can't set the value of field %s", field.String())
		}

		//Extract the tag for the given field
		tag, err := extractTag(strField)
		if err != nil {
			return err
		}
		//The value of the struct field is equal by default to
		//the tag value
		value, fieldName := tag.value, strings.ToLower(strField.Name)

		//Extract a flag value that is associated with the given field
		flag, ok := flags[fieldName]

		//If a flag is present, his value will  override the
		//default value
		if ok {
			value = flag.value
		}

		switch kind {
		case reflect.String:
			field.SetString(value)
		case reflect.Int, reflect.Int64, reflect.Int16, reflect.Int8:
			var (
				val int64
				err error
			)

			if field.Kind() == reflect.Int64 && field.Type().PkgPath() == "time" && typ.Name() == "Duration" {
				var d time.Duration

				d, err = time.ParseDuration(value)

				val = int64(d)
			} else {
				val, err = strconv.ParseInt(value, 0, typ.Bits())
			}

			if err != nil {
				return err
			}

			if field.OverflowInt(val) {
				return fmt.Errorf("given int %v overflows the field %s", val, field.Type().Name())
			}

			field.SetInt(val)
		case reflect.Bool:
			val, err := strconv.ParseBool(value)

			if err != nil {
				return err
			}

			field.SetBool(val)
		}
	}

	return nil
}

//extractTag extract the tag of a given struct field and return
//a tag value. If there's no tag on the field the function will return an error.
func extractTag(structField reflect.StructField) (Tag, error) {
	//Extract the tag value
	tag, ok := structField.Tag.Lookup("conf")
	//Extract the tag type

	val, name := "", ""

	for i := 0; i < len(tag); i++ {
		if tag[i] == ':' {
			name = tag[0:i]
			val = tag[i+1:]
			break
		}
	}

	//if there is no tag on the struct return an error
	if !ok {
		return Tag{}, fmt.Errorf("no tag for the struct field %s", structField.Name)
	}

	if name != "default" {
		return Tag{}, fmt.Errorf("unknown tag %s for the struct field %s", name, structField.Name)
	}

	//return the tag
	return Tag{value: val, name: name}, nil
}
