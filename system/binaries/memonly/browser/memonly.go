package main

import (
	"fmt"
	"log/slog"
	"syscall/js"

	"mrav/isa"
	"mrav/system"
	"mrav/system/easybus"
	"mrav/system/easybus/device"
	"mrav/system/easybus/device/memory"
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
		Details: fmt.Sprintf("%T: %s", err, err.Error()),
	}

	return js.ValueOf(map[string]interface{}{
		"error":   errWrapper.Error,
		"details": errWrapper.Details,
	})
}

func simulateSystem(this js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return wrapError(fmt.Errorf("simulateSystem requires binary data and instruction count"))
	}

	// Convert JavaScript array to byte slice
	jsArray := args[0]
	arrayLength := jsArray.Get("length").Int()
	softwareBytes := make([]byte, arrayLength)

	for i := 0; i < arrayLength; i++ {
		softwareBytes[i] = byte(jsArray.Index(i).Int())
	}

	instructionsToSim := args[1].Int()

	logger := slog.Default()
	opts := &system.SystemOpts{
		Logger:  logger,
		Verbose: false,
	}

	mem, err := memory.NewMem(1024, softwareBytes)
	if err != nil {
		return wrapError(fmt.Errorf("cannot create memory device: %w", err))
	}

	// Note: Timer is intentionally omitted for browser version
	sys, err := easybus.NewEasyBusSystem(opts, []device.Device{mem})
	if err != nil {
		return wrapError(fmt.Errorf("cannot create system: %w", err))
	}

	instructionCount := 0
	for i := 0; i < instructionsToSim; i++ {
		if err := sys.RunInstruction(); err != nil {
			return wrapError(fmt.Errorf("cannot run instruction %d: %w", i, err))
		}
		instructionCount++
	}

	// Get the core and build a completely flat map
	core := sys.GetCore()
	result := make(map[string]interface{})
	result["pc"] = fmt.Sprintf("0x%04x", core.Pc)
	for i := 0; i < int(isa.RegsNumber); i++ {
		result[fmt.Sprintf("r%d", i)] = fmt.Sprintf("0x%04x", core.Registers[i])
	}
	result["instructions"] = instructionCount
	result["error"] = nil

	// Add memory contents to result
	memBytes := mem.GetMemoryBytes()
	for i, b := range memBytes {
		result[fmt.Sprintf("mem_%d", i)] = fmt.Sprintf("0x%02X", b)
	}

	return result
}

func main() {
	c := make(chan struct{})
	js.Global().Set("simulateSystem", js.FuncOf(simulateSystem))
	<-c
}
