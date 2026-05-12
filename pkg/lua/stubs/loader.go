package stubs

import (
	"io/fs"
	"maps"
	"os"
	"path/filepath"
)

type Loader struct {
	parser *Parser
}

func NewLoader() *Loader {
	return &Loader{
		parser: NewParser(),
	}
}

func (l *Loader) LoadDir(dir string) (*Schema, error) {
	schema := &Schema{
		Globals: map[string]Field{},
		Modules: map[string]*Module{},
	}

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		fileSchema, err := l.parser.Parse(f)
		if err != nil {
			return err
		}

		merge(schema, fileSchema)

		return nil
	})

	return schema, err
}

func merge(dst, src *Schema) {
	for k, v := range src.Modules {

		if _, ok := dst.Modules[k]; !ok {
			dst.Modules[k] = v
			continue
		}

		maps.Copy(dst.Modules[k].Fields, v.Fields)
	}
}
