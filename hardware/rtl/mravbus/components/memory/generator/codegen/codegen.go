package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"mrav/hardware/rtl/mravbus/components/memory/generator"
)

func main() {
	memSize := flag.Uint("mem_size", 0, "size of the generated block")
	rtlFile := flag.String("rtl_file", "", "output RTL file containing Verilog")
	topModule := flag.String("top_module", "", "name of the top RTL module")
	binPayload := flag.String("bin_payload", "", "file containing the binary payload for the memory image")
	internalAddress := flag.String("internal_address", "", "RTL slice generating the internal address")

	flag.Parse()

	binContents, err := os.ReadFile(*binPayload)

	if err != nil {
		log.Fatalf("cannot read the binary payload for file '%s': %v", *binPayload, err)
	}

	verilogBytes := make([]string, 0, len(binContents))

	for _, payloadByte := range binContents {
		verilogBytes = append(verilogBytes, fmt.Sprintf("h%02X", payloadByte))
	}

	outputFile, err := os.Create(*rtlFile)

	if err != nil {
		log.Fatalf("cannot prepare the output RTL file: %v", err)
	}

	defer outputFile.Close()

	opts := &generator.MemGenOpts{
		MemSize:            uint32(*memSize),
		TopModule:          *topModule,
		Payload:            verilogBytes,
		InternalAddressRtl: *internalAddress,
		Writer:             outputFile,
	}

	if err := generator.GenerateMemory(opts); err != nil {
		log.Fatalf("cannot generate the code: %v", err)
	}
}
