package main

import (
	"fmt"
	"log/slog"
	"strings"
	"syscall/js"

	"mrav/software/asm"
	"mrav/software/format"
	"mrav/software/model"
)

type ErrorWrapper struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

func wrapError(err error) js.Value {
	if err == nil {
		return js.Null()
	}

	errWrapper := ErrorWrapper{
		Error:   err.Error(),
		Details: fmt.Sprintf("%T: %s", err, err.Error()), // Include the error type
	}

	return js.ValueOf(map[string]interface{}{
		"error":   errWrapper.Error,
		"details": errWrapper.Details,
	})
}

func assembleModule(src string, formatFunc func(*model.MravModule) (interface{}, error)) interface{} {
	logger := slog.Default()
	logger.Info("Got the source, moving on to assembling")

	program, err := asm.AssembleModules([]string{src})

	if err != nil {
		return wrapError(fmt.Errorf("unable to assemble: %w", err))
	}

	output, err := formatFunc(program)

	if err != nil {
		return wrapError(fmt.Errorf("Cannot output the machine code: %w", err))
	}

	logger.Info("Successfully assembled")

	return map[string]interface{}{
		"data":  output,
		"error": nil,
	}
}

func assembleModuleHumanReadable(this js.Value, args []js.Value) interface{} {
	src := args[0].String()
	return assembleModule(src, func(program *model.MravModule) (interface{}, error) {
		humanReadable, err := format.HumanReadable(program)
		if err != nil {
			return nil, err
		}
		return strings.Join(humanReadable, "\n") + "\n", nil
	})
}

func assembleModuleBinary(this js.Value, args []js.Value) interface{} {
	src := args[0].String()
	return assembleModule(src, func(program *model.MravModule) (interface{}, error) {
		binaryData, err := format.Binary(program)
		if err != nil {
			return nil, err
		}
		// Convert []byte to []interface{} so js.ValueOf can handle it
		jsArray := make([]interface{}, len(binaryData))
		for i, b := range binaryData {
			jsArray[i] = int(b)
		}
		return jsArray, nil
	})
}

func main() {
	c := make(chan struct{})
	js.Global().Set("assembleModuleHumanReadable", js.FuncOf(assembleModuleHumanReadable))
	js.Global().Set("assembleModuleBinary", js.FuncOf(assembleModuleBinary))
	<-c
}
