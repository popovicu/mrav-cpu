package generator

import (
	"embed"
	"io"
	"text/template"
)

type MemGenOpts struct {
	MemSize            uint32
	TopModule          string
	Payload            []string
	InternalAddressRtl string
	Writer             io.Writer
}

//go:embed memory.sv
var templateFS embed.FS

func GenerateMemory(memOpts *MemGenOpts) error {
	funcMap := template.FuncMap{
		"isLast": func(idx int, payload []string) bool {
			return idx == (len(payload) - 1)
		},
		"sub": func(a uint32, b uint32) uint32 {
			return a - b
		},
	}
	tpl := template.New("memory.sv").Funcs(funcMap)
	tmpl, err := tpl.ParseFS(templateFS, "memory.sv")

	if err != nil {
		return err
	}

	if err := tmpl.Execute(memOpts.Writer, *memOpts); err != nil {
		return err
	}

	return nil
}
