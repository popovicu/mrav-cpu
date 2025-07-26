package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"mrav/hardware/rtl/mravbus/generator"
)

type BusDevicesDescriptor struct {
	Descriptor []*generator.Peripheral
}

func main() {
	devicesFile := flag.String("devices_file", "", "JSON file describing the devices connecting to the bus")
	rtlFile := flag.String("rtl_file", "", "output RTL file containing Verilog")

	flag.Parse()

	devicesFileHandle, err := os.ReadFile(*devicesFile)

	if err != nil {
		log.Fatalf("cannot open the devices JSON file: %v", err)
	}

	var descriptor BusDevicesDescriptor

	err = json.Unmarshal(devicesFileHandle, &descriptor)

	if err != nil {
		log.Fatalf("cannot read the JSON descriptor: %v", err)
	}

	outputFile, err := os.Create(*rtlFile)

	if err != nil {
		log.Fatalf("cannot prepare the output RTL file: %v", err)
	}

	defer outputFile.Close()

	opts := &generator.BusGenOpts{
		Peripherals: descriptor.Descriptor,
		Writer:      outputFile,
	}

	if err := generator.GenerateBus(opts); err != nil {
		log.Fatalf("cannot generate the code: %v", err)
	}
}
