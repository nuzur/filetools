package filetools

import (
	"bytes"
	"context"
	"fmt"
	"go/format"
	"os"
	"path"
	"text/template"
)

type FileRequest struct {
	// full path of the file to generate
	OutputPath string

	// Template name to use
	TemplateName string

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

	// read the template
	templatesPath := resolveTemplatesPath()
	templateFileName := fmt.Sprintf("%s.tmpl", req.TemplateName)
	templateFilePath := path.Join(templatesPath, templateFileName)
	data, err := os.ReadFile(templateFilePath)
	if err != nil {
		fmt.Println("reading error", err)
		return nil, err
	}

	// instantiate the template
	t, err := template.New("template").Funcs(funcs).Parse(string(data))
	if err != nil {
		fmt.Printf("Template Error: %v\n ", req.TemplateName)
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
	rootDir := CurrentPath()
	finalOutput := path.Join(rootDir, req.OutputPath)
	err = Write(finalOutput, output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func resolveTemplatesPath() string {
	rootDir := CurrentPath()
	return path.Join(rootDir, "templates")
}
