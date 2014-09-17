package main

import (
	"fmt"
	"io"
	"text/template"
)

const (
	MergeTemplate      = "{{template \"MergeT\" .}}\n"
	FanoutTemplate     = "{{template \"FanoutT\" .}}\n"
	ApplyTemplate      = "{{template \"ApplyT\"}}\n"
	PapplyTemplate     = "{{template \"PapplyT\"}}\n"
	MapTemplate        = "{{template \"MapT\"}}\n"
	PmapTemplate       = "{{template \"Pmap\"}}"
	FilterTemplate     = "{{template \"FilterT\"}}\n"
	TransformTemplate  = "{{template \"TransformT\"}}\n"
	PtransformTemplate = "{{template \"PtransformT\" .}}\n"
)

var (
	all  *template.Template
	tmap map[string]string
)

func init() {
	all = template.Must(template.ParseFiles("chan.tmpl"))
	tmap = map[string]string{
		"Merge":      MergeTemplate,
		"Fanout":     FanoutTemplate,
		"Apply":      ApplyTemplate,
		"Papply":     PapplyTemplate,
		"Map":        MapTemplate,
		"Pmap":       PmapTemplate,
		"Filter":     FilterTemplate,
		"Transform":  TransformTemplate,
		"Ptransform": PtransformTemplate,
	}
}

func WriteMethod(w io.Writer, name string, v interface{}) error {
	tl, ok := tmap[name]
	if !ok {
		return fmt.Errorf("no method %q", name)
	}

	local, err := all.Clone()
	if err != nil {
		return err
	}
	local, err = local.Parse(tl)
	if err != nil {
		return err
	}

	return local.Execute(w, v)
}
