package main

import (
	"flag"
	"log"
	"log/slog"
	"os"

	"mrav/isa"
	"mrav/system"
	"mrav/system/easybus"
	"mrav/system/easybus/device"
	"mrav/system/easybus/device/memory"
	"mrav/system/easybus/device/timer"
)

func main() {
	softwareBinary := flag.String("software", "", "path to the software file")
	verbose := flag.Bool("verbose", false, "whether to produce verbose output")
	instructionsToSim := flag.Int("instructions_to_sim", 20, "number of instructions to simulate")
	coreStateOutput := flag.String("core_state_output", "", "path to the file where the state of the core should be output after the simulation")
	coreStateProtoOutput := flag.String("core_state_proto_output", "", "path to the file where the state of the core should be output after the simulation (proto format)")

	flag.Parse()

	softwareBytes, err := os.ReadFile(*softwareBinary)

	if err != nil {
		log.Fatalf("cannot load the software binary: %v", err)
	}

	logger := slog.Default()
	opts := &system.SystemOpts{
		Logger:  logger,
		Verbose: *verbose,
	}

	mem, err := memory.NewMem(1024, softwareBytes)

	if err != nil {
		log.Fatalf("cannot create the memory devices: %v", err)
	}

	tim := &timer.Timer{}

	sys, err := easybus.NewEasyBusSystem(opts, []device.Device{mem, tim})

	if err != nil {
		log.Fatalf("cannot create a system: %v", err)
	}

	for i := 0; i < *instructionsToSim; i++ {
		if err := sys.RunInstruction(); err != nil {
			log.Fatalf("cannot run a system instruction: %v", err)
		}

		snap, err := sys.CoreDebug([]isa.RegisterId{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15})

		if err != nil {
			log.Fatalf("cannot snapshot the core: %v", err)
		}

		if *verbose {
			logger.Info("[Core] Snapshot", "state", snap)
		}
	}

	if *coreStateOutput != "" {
		coreState, err := sys.CoreDebug([]isa.RegisterId{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15})

		if err != nil {
			log.Fatalf("unable to generate the core state string: %w", err)
		}

		if err := os.WriteFile(*coreStateOutput, []byte(coreState), 0644); err != nil {
			log.Fatalf("unable to dump the core state string: %w", err)
		}
	}

	if *coreStateProtoOutput != "" {
		if err := sys.ProtoCoreDebugFile(*coreStateProtoOutput); err != nil {
			log.Fatalf("unable to dump core proto: %w", err)
		}
	}
}
