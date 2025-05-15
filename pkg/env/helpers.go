package env

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

const (
	whiteSpace = 32
)

type StringerWriter interface {
	io.Writer
	fmt.Stringer
}

type indentWrapper struct {
	writer StringerWriter
	indent int
}

func newIndent(writer StringerWriter, indent int) *indentWrapper {
	return &indentWrapper{writer: writer, indent: indent}
}

func (i *indentWrapper) Write(p []byte) (n int, err error) {
	indenter := []byte{whiteSpace}
	cp := make([]byte, len(p))
	copy(cp, p)

	if bytes.Contains(cp, []byte("\n")) {
		return i.writer.Write(cp)
	}
	return i.writer.Write(append(cp, bytes.Repeat(indenter, i.indent)...))
}

func (i *indentWrapper) Bytes() []byte {
	return []byte(i.writer.String())
}

func genEnvConfig[Type any](cfg Type) []byte {

	_type := reflect.TypeOf(cfg)

	if _type.Kind() == reflect.Pointer {
		_type = _type.Elem()
	}

	if _type.Kind() != reflect.Struct {
		return nil
	}

	builder := new(strings.Builder)

	writer := newIndent(builder, 2)
	genEnvConfigRecursively(writer, _type)

	return writer.Bytes()
}

func genEnvConfigRecursively(
	writer io.Writer,
	_type reflect.Type,
) {

	for i := range _type.NumField() {

		field := _type.Field(i)

		if field.Type.Kind() == reflect.Struct {
			genEnvConfigRecursively(writer, field.Type)
			continue
		}

		writer.Write([]byte("# "))
		descTag, ok := field.Tag.Lookup("desc")
		if ok {
			writer.Write([]byte(descTag))
		}
		writer.Write([]byte(fmt.Sprintf(" (%s)\n", field.Type)))

		if envTag, ok := field.Tag.Lookup("env"); ok {
			writer.Write([]byte(envTag + "="))
		}

		if defaultTag, ok := field.Tag.Lookup("default"); ok {
			writer.Write([]byte(defaultTag))
		}

		writer.Write([]byte("\n"))

	}

}

func ConfigInfo[Type any](
	writer io.Writer,
) func(string) error {
	return func(string) error {
		defer os.Exit(0)
		t := new(Type)
		_, err := writer.Write(genEnvConfig(t))
		if err != nil {
			return err
		}
		return nil
	}
}

func ConfigDoc[Type any]() func(string) error {
	return func(filename string) error {

		if _, err := os.Stat(filename); err != nil {
			return err
		}

		file, err := os.Create(filename)
		if err != nil {
			return err
		}

		return ConfigInfo[Type](file)(filename)
	}
}
