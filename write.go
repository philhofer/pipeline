package main

import (
	"fmt"
	"io"
	"os"
	"text/template"
)

const (
	MergeTemplate      = "{{template \"MergeT\" .}}\n\n"
	FanoutTemplate     = "{{template \"FanoutT\" .}}\n\n"
	ApplyTemplate      = "{{template \"ApplyT\" .}}\n\n"
	PapplyTemplate     = "{{template \"PapplyT\" .}}\n\n"
	MapTemplate        = "{{template \"MapT\" .}}\n\n"
	PmapTemplate       = "{{template \"PmapT\" .}}\n\n"
	FilterTemplate     = "{{template \"FilterT\" .}}\n\n"
	SendAllTemplate    = "{{template \"SendAllT\" .}}"
	RecvAllTemplate    = "{{template \"RecvAllT\" .}}\n\n"
	RecvNTemplate      = "{{template \"RecvNT\" .}}\n\n"
	BufferTemplate     = "{{template \"BufferT\" .}}\n\n"
	TransformTemplate  = "{{template \"TransformT\" .}}\n\n"
	PtransformTemplate = "{{template \"PtransformT\" .}}\n\n"
)

var (
	all  *template.Template
	tmap map[string]string
)

func init() {
	gopath := os.Getenv("GOPATH")
	tloc := gopath + "/src/github.com/philhofer/pipeline/chan.tmpl"
	all = template.Must(template.ParseFiles(tloc))
	tmap = map[string]string{
		"Merge":      MergeTemplate,
		"Fanout":     FanoutTemplate,
		"Apply":      ApplyTemplate,
		"Papply":     PapplyTemplate,
		"Map":        MapTemplate,
		"Pmap":       PmapTemplate,
		"Filter":     FilterTemplate,
		"RecvAll":    RecvAllTemplate,
		"RecvN":      RecvNTemplate,
		"SendAll":    SendAllTemplate,
		"Buffer":     BufferTemplate,
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
