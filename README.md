pipeline
==========

This package is a code generator for a number of generic channel pipelining methods.


### Usage

First, install using the usual `go get`.

This tool is targeted at the `go generate` tool. In your source code, you will
need to include one or more directives like the following:

```go
//go:generate pipeline -methods=Merge,Transform,Apply type=*int,string
```

(The above would generate the `Merge` and `Apply` methods for `*int`s and `Transform` for `*int`-to-`string`.) The default output 
file is the name of the file the directive was called from with `_gen.go` added as a suffix. You can override 
the default output with the `-o` flag.

### Methods

 - Merge (merge channels)
 - Fanout (fanout channels)
 - Apply (sequential modification)
 - Papply (parallel modification)
 - Map (sequential modification)
 - Pmap (parallel modification)
 - Filter
 - Transform (sequential transformation)
 - Ptransform (parallel transformation)
 

Prototypes and descriptions of the methods are available at [godoc](http://godoc.org/github.com/philhofer/pipeline).