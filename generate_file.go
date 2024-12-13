package filetools

import (
	"bytes"
	"context"
	"fmt"
	"go/format"
	"os"
	"text/template"
)

type FileRequest struct {
	// full path of the file to generate
	OutputPath string

	// Template path
	TemplatePath string

	// Template bytes - has priority over TemplatePath
	TemplateBytes []byte

	// Data will be passed to the template
	Data any

	// Funcs extra functions to pass to the template
	Funcs template.FuncMap

	// DisableGoFormat Should disable goformat
	DisableGoFormat bool
}

func GenerateFile(ctx context.Context, req FileRequest) ([]byte, error) {
	funcs := template.FuncMap{}
	for n, f := range req.Funcs {
		funcs[n] = f
	}

	// load the template
	templateData := req.TemplateBytes
	if templateData == nil {
		var err error
		templateData, err = os.ReadFile(req.TemplatePath)
		if err != nil {
			fmt.Println(" error reading template", err)
			return nil, err
		}
	}

	// instantiate the template
	t, err := template.New("template").Funcs(funcs).Parse(string(templateData))
	if err != nil {
		fmt.Printf("Template Error: %v\n ", req.TemplatePath)
		fmt.Printf("Template Error: %v\n ", err)
		return nil, fmt.Errorf("error with provided template: %w", err)
	}

	// execute with data
	var buf bytes.Buffer
	err = t.Execute(&buf, req.Data)
	if err != nil {
		fmt.Printf("Execute Error: %v\n err:%v\n", req, err)
		return nil, fmt.Errorf("error executing template")
	}

	output := buf.Bytes()

	// format code
	if !req.DisableGoFormat {
		output, err = format.Source(output)
		if err != nil {
			output = buf.Bytes()
			fmt.Printf("error formating file: %v \n", req.OutputPath)
		}

	}

	// write the output
	err = Write(req.OutputPath, output)
	if err != nil {
		return nil, err
	}
	return output, nil
}
