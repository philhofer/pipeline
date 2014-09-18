package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
)

// used for two-type methods
type Transfrm struct {
	Src string
	Dst string
}

func (t Transfrm) Esc() string {
	return strings.Replace(t.Src, "*", "ptr", -1) + strings.Replace(t.Dst, "*", "ptr", -1)
}

type Basic struct {
	s string
}

func (b Basic) String() string { return b.s }
func (b Basic) Esc() string    { return strings.Replace(b.s, "*", "ptr", -1) }

var (
	methodlist string
	tnames     string
	outf       string
)

func init() {
	flag.StringVar(&methodlist, "methods", "", "methods to generate")
	flag.StringVar(&tnames, "type", "", "the type to use")
	flag.StringVar(&outf, "o", "_gen.go", "out file")
}

func main() {
	flag.Parse()

	if methodlist == "" {
		fmt.Println("missing method list (-methods=\"Merge,Fanout,Apply...\"")
		os.Exit(1)
	}
	if tnames == "" {
		fmt.Println("missing type name(s) (-type=\"int,string...\"")
	}

	var buf bytes.Buffer
	ms := strings.Split(methodlist, ",")
	ts := strings.Split(tnames, ",")
	switch len(ts) {
	case 1, 2:
	default:
		fmt.Println("need either 1 or 2 type names")
		os.Exit(1)
	}

	pkg := os.Getenv("GOPACKAGE")
	if pkg == "" {
		pkg = "main"
	}

	// write package name and required imports (which
	// should just be "sync")
	buf.WriteString(fmt.Sprintf("package %s\n\n", pkg))
	buf.WriteString("import (\n\t\"sync\"\n)\n\n")

	if outf == "_gen.go" {
		if os.Getenv("GOFILE") != "" {
			outf = os.Getenv("GOFILE") + "_gen.go"
		}
	}

	for _, m := range ms {
		var err error
		switch m {
		// one-type methods
		case "Merge", "Fanout", "Apply", "Papply", "Map", "Pmap", "Filter", "SendAll", "RecvAll", "RecvN", "Buffer":
			err = WriteMethod(&buf, m, Basic{s: ts[0]})

		// two-type methods
		case "Transform", "Ptransform":
			if len(ts) != 2 {
				fmt.Println("need 2 type names for Transform methods")
				os.Exit(1)
			}
			err = WriteMethod(&buf, m, Transfrm{Src: ts[0], Dst: ts[1]})
		default:
			fmt.Printf("Unrecognized method name %q", m)
			os.Exit(1)
		}
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
	}

	file, err := os.Create(outf)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	defer file.Close()
	_, err = buf.WriteTo(file)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
