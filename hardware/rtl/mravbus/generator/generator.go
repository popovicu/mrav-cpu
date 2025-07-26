package generator

import (
	"embed"
	"io"
	"text/template"
)

type Peripheral struct {
	DeviceId string `json:"device_id"`
	AddrLo   int    `json:"addr_lo"`
	AddrHi   int    `json:"addr_hi"`
}

type BusGenOpts struct {
	Peripherals []*Peripheral
	Writer      io.Writer
}

//go:embed bus_tpl.sv
var templateFS embed.FS

func GenerateBus(busOpts *BusGenOpts) error {
	funcMap := template.FuncMap{
		"isLast": func(idx int, periphs []*Peripheral) bool {
			return idx == (len(periphs) - 1)
		},
	}

	tpl := template.New("bus_tpl.sv").Funcs(funcMap)
	tmpl, err := tpl.ParseFS(templateFS, "bus_tpl.sv")

	if err != nil {
		return err
	}

	if err := tmpl.Execute(busOpts.Writer, *busOpts); err != nil {
		return err
	}

	return nil
}
