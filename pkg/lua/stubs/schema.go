package stubs

type Schema struct {
	Globals map[string]Field
	Modules map[string]*Module
}

type Module struct {
	Name   string
	Fields map[string]Field
}

type Field struct {
	Name string
	Type string
	Path string
}
