//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/Enflick/smithy-go/ptr"
)

func main() {
	types := ptr.GetScalars()

	for filename, tmplName := range map[string]string{
		"to_ptr.go":   "scalar to pointer",
		"from_ptr.go": "scalar from pointer",
	} {
		if err := generateFile(filename, tmplName, types); err != nil {
			log.Fatalf("%s file generation failed, %v", filename, err)
		}
	}
}

func generateFile(filename string, tmplName string, types ptr.Scalars) (err error) {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create %s file, %v", filename, err)
	}

	defer func() {
		closeErr := f.Close()
		if err == nil {
			err = closeErr
		} else if closeErr != nil {
			err = fmt.Errorf("close error: %v, original error: %w", closeErr, err)
		}
	}()

	if err = ptrTmpl.ExecuteTemplate(f, tmplName, types); err != nil {
		return fmt.Errorf("failed to generate %s file, %v", filename, err)
	}

	return nil
}

var ptrTmpl = template.Must(template.New("ptrTmpl").Parse(`
{{- define "header" }}
	// Code generated by smithy-go/ptr/generate.go DO NOT EDIT.
	package ptr

	import (
		{{- range $_, $import := $.Imports }}
			"{{ $import.Path }}"
		{{- end }}
	)
{{- end }}

{{- define "scalar from pointer" }}
	{{ template "header" $ }}

	{{ range $_, $type := $ }}
		{{ template "from pointer func" $type }}
		{{ template "from pointers func" $type }}
	{{- end }}
{{- end }}

{{- define "scalar to pointer" }}
	{{ template "header" $ }}

	{{ range $_, $type := $ }}
		{{ template "to pointer func" $type }}
		{{ template "to pointers func" $type }}
	{{- end }}
{{- end }}

{{- define "to pointer func" }}
	// {{ $.Name }} returns a pointer value for the {{ $.Symbol }} value passed in.
	func {{ $.Name }}(v {{ $.Symbol }}) *{{ $.Symbol }} {
		return &v
	}
{{- end }}

{{- define "to pointers func" }}
	// {{ $.Name }}Slice returns a slice of {{ $.Symbol }} pointers from the values
	// passed in.
	func {{ $.Name }}Slice(vs []{{ $.Symbol }}) []*{{ $.Symbol }} {
		ps := make([]*{{ $.Symbol }}, len(vs))
		for i, v := range vs {
			vv := v
			ps[i] = &vv
		}

		return ps
	}

	// {{ $.Name }}Map returns a map of {{ $.Symbol }} pointers from the values
	// passed in.
	func {{ $.Name }}Map(vs map[string]{{ $.Symbol }}) map[string]*{{ $.Symbol }} {
		ps := make(map[string]*{{ $.Symbol }}, len(vs))
		for k, v := range vs {
			vv := v
			ps[k] = &vv
		}

		return ps
	}
{{- end }}

{{- define "from pointer func" }}
	// To{{ $.Name }} returns {{ $.Symbol }} value dereferenced if the passed
	// in pointer was not nil. Returns a {{ $.Symbol }} zero value if the
	// pointer was nil.
	func To{{ $.Name }}(p *{{ $.Symbol }}) (v {{ $.Symbol }}) {
		if p == nil {
			return v
		}
			
		return *p
	}
{{- end }}

{{- define "from pointers func" }}
	// To{{ $.Name }}Slice returns a slice of {{ $.Symbol }} values, that are
	// dereferenced if the passed in pointer was not nil. Returns a {{ $.Symbol }}
	// zero value if the pointer was nil.
	func To{{ $.Name }}Slice(vs []*{{ $.Symbol }}) []{{ $.Symbol }} {
		ps := make([]{{ $.Symbol }}, len(vs))
		for i, v := range vs {
			ps[i] = To{{ $.Name }}(v)
		}

		return ps
	}

	// To{{ $.Name }}Map returns a map of {{ $.Symbol }} values, that are
	// dereferenced if the passed in pointer was not nil. The {{ $.Symbol }}
	// zero value is used if the pointer was nil.
	func To{{ $.Name }}Map(vs map[string]*{{ $.Symbol }}) map[string]{{ $.Symbol }} {
		ps := make(map[string]{{ $.Symbol }}, len(vs))
		for k, v := range vs {
			ps[k] = To{{ $.Name }}(v)
		}

		return ps
	}
{{- end }}
`))
