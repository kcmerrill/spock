package info

import (
	"bytes"
	"encoding/json"
	"html/template"
)

// New returns information passed along via our check
func New(encoded []byte) *Info {
	i := &Info{}
	json.Unmarshal(encoded, i)
	i.Clean()
	i.CreateTemplate()
	return i
}

// Info holds all of our check information
type Info struct {
	ID         string
	Attempts   int
	Module     string
	Error      string
	Output     string
	Name       string
	Template   string
	Properties struct {
		Attempts struct{ Value int }
		Module   struct{ Value string }
		Error    struct{ Value string }
		Output   struct{ Value string }
		Name     struct{ Value string }
		Template struct{ Value string }
	}
}

// Clean makes the template bit more convienant
func (i *Info) Clean() {
	i.Attempts = i.Properties.Attempts.Value
	i.Module = i.Properties.Module.Value
	i.Error = i.Properties.Error.Value
	i.Output = i.Properties.Output.Value
	i.Name = i.Properties.Name.Value
	i.Template = i.Properties.Template.Value
}

// CreateTemplate takes in the template param and tries to create a template off it.
func (i *Info) CreateTemplate() {
	if i.Template != "" {
		template := template.Must(template.New("translate").Parse(i.Template))
		b := new(bytes.Buffer)
		err := template.Execute(b, i)
		if err == nil {
			i.Template = b.String()
		} else {
			// busted ... don't use it
			i.Template = ""
		}
	}
}
