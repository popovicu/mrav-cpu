package main

import (
	"fmt"
	"log/slog"
	"strings"
	"syscall/js"

	"mrav/software/asm"
	"mrav/software/format"
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

func assembleModule(this js.Value, args []js.Value) interface{} {
	src := args[0].String()

	logger := slog.Default()
	logger.Info("Got the source, moving on to assembling")

	program, err := asm.AssembleModules([]string{src})

	if err != nil {
		return wrapError(fmt.Errorf("unable to assemble: %w", err))
	}

	humanReadable, err := format.HumanReadable(program)

	if err != nil {
		return wrapError(fmt.Errorf("Cannot output the machine code: %w", err))
	}

	programOutput := strings.Join(humanReadable, "\n") + "\n"

	logger.Info("Successfully assembled")

	return map[string]interface{}{
		"data":  programOutput,
		"error": nil,
	}
}

func main() {
	c := make(chan struct{})
	js.Global().Set("assembleModule", js.FuncOf(assembleModule))
	<-c
}
