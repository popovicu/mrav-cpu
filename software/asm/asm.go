package asm

import (
	"fmt"

	"mrav/software/asm/firstpass"
	"mrav/software/asm/parsing"
	"mrav/software/asm/secondpass"
	"mrav/software/linker"
	"mrav/software/model"
)

func AssembleModules(modules []string) (*model.MravModule, error) {
	objects := make([]*secondpass.MravObject, 0, len(modules))

	for i, m := range modules {
		parsedModule, err := parsing.ParseModuleString(m)

		if err != nil {
			return nil, fmt.Errorf("unable to parse module %d: %w", i, err)
		}

		firstPassModule, err := firstpass.ModuleFirstPass(parsedModule)

		if err != nil {
			return nil, fmt.Errorf("unable to do the first pass on the module %d: %w", i, err)
		}

		object, err := secondpass.BuildObject(firstPassModule)

		if err != nil {
			return nil, fmt.Errorf("unable to do the second pass on the module %d: %w", i, err)
		}

		objects = append(objects, &object)
	}

	program, err := linker.Link(objects)

	if err != nil {
		return nil, err
	}

	return program, nil
}
