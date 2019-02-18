package tote

import (
    "fmt"
    "os"
    "reflect"
    "strconv"
    "strings"
    "unicode"

    "github.com/pkg/errors"
)

var defaultEnvironmentPrefix = "TOTE"

// NOTE: This variable exists to inject a mock into the code for tests, it
// might be a better idea to use something like:
//   https://github.com/adammck/venv
var lookupEnv = os.LookupEnv

func readConfigFromEnvironment(
    config interface{}, prefixes ...string,
) error {
    finalPrefix := strings.Join(prefixes, "_")
    return reflectEnvToConfig(config, finalPrefix)
}

// reflectEnvToConfig reads all the attributes of an
// object and assigns them based on environment variables.
func reflectEnvToConfig(
    object interface{}, startPrefix string,
) error {
    // I had troubles with addressability when using recursion so
    // this is no longer recursive.

    type prefixedValueType struct {
        prefix string
        value  reflect.Value
        typ    reflect.Type
        fieldName string
    }

    field := make([]prefixedValueType, 1)

    field[0] = prefixedValueType{
        startPrefix,
        reflect.ValueOf(object),
        reflect.TypeOf(object),
        "",
    }

    for len(field) != 0 {
        var nextItem prefixedValueType
        nextItem, field = field[0], field[1:]

        nextValue := nextItem.value
        nextType := nextItem.typ
        prefix := nextItem.prefix

        // Resolve all pointers.
        for nextType.Kind() == reflect.Ptr {
            nextType = nextType.Elem()
            nextValue = nextValue.Elem()
        }

        if nextType.Kind() == reflect.Struct {
            for i := 0; i < nextType.NumField(); i++ {
                structField := nextType.Field(i)
                structFieldName := structField.Name

                if !unicode.IsUpper(rune(structFieldName[0])) {
                    continue
                }

                structFieldValue := nextValue.Field(i)
                structFieldType := structField.Type
                structFieldNameUpper := strings.ToUpper(structFieldName)

                newPrefix := fmt.Sprintf(
                    "%s_%s",
                    prefix,
                    structFieldNameUpper,
                )
                nextObject := prefixedValueType{
                    newPrefix,
                    structFieldValue,
                    structFieldType,
                    structFieldName,
                }

                field = append(field, nextObject)
            }
        } else {
            envVariable := prefix
            envVariableValue, exists := lookupEnv(envVariable)
            if exists {
                err := setField(
                    nextValue, nextType.Kind(), envVariableValue,
                )
                if err != nil {
                    return errors.Wrapf(
                        err,
                        "Error while reading environment variable '%s' into " +
                            "field '%s'",
                        envVariable,
                        nextItem.fieldName,
                    )
                }
            }
        }
    }

    return nil
}

func setField(
    fieldValue reflect.Value,
    kind reflect.Kind,
    value string,
) error {
    switch kind {
    case
        reflect.Int,
        reflect.Int8,
        reflect.Int16,
        reflect.Int32,
        reflect.Int64:
        intValue, err := strconv.ParseInt(value, 10, 64)
        if err != nil {
            return errors.Wrapf(
                err, "Expected int value but found value '%s' instead", value,
            )
        }
        fieldValue.SetInt(intValue)
    case
        reflect.Float32,
        reflect.Float64:
        floatValue, err := strconv.ParseFloat(value, 64)
        if err != nil {
            return errors.Wrapf(
                err,
                "Expected float value but found value '%s' instead",
                value,
            )
        }
        fieldValue.SetFloat(floatValue)
    case reflect.String:
        fieldValue.SetString(value)
    case reflect.Bool:
        boolValue, err := strconv.ParseBool(value)
        if err != nil {
            return errors.Wrapf(
                err,
                "Expected bool value but found value '%s' instead",
                value,
            )
        }
        fieldValue.SetBool(boolValue)
    default:
        return errors.Errorf(
            "Field expects unknown kind '%s'", kind,
        )
    }

    return nil
}
