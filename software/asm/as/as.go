package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"

	"mrav/software/asm"
	"mrav/software/format"
)

func readTextFiles(paths []string) ([]string, error) {
	contents := make([]string, 0, len(paths))

	for i, path := range paths {
		contentsBytes, err := os.ReadFile(path)

		if err != nil {
			return nil, fmt.Errorf("cannot read source for module %d: %w", i, err)
		}

		contents = append(contents, string(contentsBytes))
	}

	return contents, nil
}

func main() {
	debug := flag.Bool("debug", false, "enable debug output")
	outputFile := flag.String("output", "", "path to the output program file")
	outputFormat := flag.String("format", "human", "output format for the assembler")

	flag.Parse()

	inputFiles := flag.Args()
	logger := slog.Default()

	srcs, err := readTextFiles(inputFiles)

	if err != nil {
		log.Fatalf("unable to load source files: %v", err)
	}

	for moduleIdx, inputFile := range inputFiles {
		logger.Info("Module", "index", moduleIdx, "module_path", inputFile)
	}

	logger.Info("Finished reading source files, moving on to assembling")

	program, err := asm.AssembleModules(srcs)

	if err != nil {
		log.Fatalf("unable to assemble: %v", err)
	}

	if *debug {
		spew.Dump(program)
	}

	humanReadableOutput := func() {
		humanReadable, err := format.HumanReadable(program)

		if err != nil {
			log.Fatalf("Cannot output the machine code: %v", err)
		}

		programOutput := strings.Join(humanReadable, "\n") + "\n"

		if err := os.WriteFile(*outputFile, []byte(programOutput), 0644); err != nil {
			log.Fatalf("Cannot write the human readable output file: %v", err)
		}
	}

	binaryOutput := func() {
		binaryPayload, err := format.Binary(program)

		if err != nil {
			log.Fatalf("Cannot output the machine code: %v", err)
		}

		if err := os.WriteFile(*outputFile, binaryPayload, 0644); err != nil {
			log.Fatalf("Cannot write the binary output file: %v", err)
		}
	}

	outputProducers := map[string]func(){
		"human":  humanReadableOutput,
		"binary": binaryOutput,
	}

	outputProducer, found := outputProducers[*outputFormat]

	if !found {
		log.Fatalf("Unknown output format: %s", *outputFormat)
	}

	outputProducer()

	logger.Info("Successfully assembled")
}
