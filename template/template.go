package template

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/steinfletcher/t8/domain"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

// TODO test this

// fsRenderer implements TemplateRenderer
type fsRenderer struct{}

func NewTemplateRenderer() domain.TemplateRenderer {
	return fsRenderer{}
}

// Render walks the targetDir on the filesystem and renders files with the given config.
// It treats every file apart from t8's known config as a template.
func (f fsRenderer) Render(targetDir string, conf domain.Config) {
	err := filepath.Walk(targetDir, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return nil
		}

		if conf.ShouldExcludePath(path) {
			return os.Remove(path)
		}

		file, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		data, err := render(string(file), createTemplateModel(conf))
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(path, data, 0644)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func createTemplateModel(config domain.Config) map[string]interface{} {
	model := map[string]interface{}{
		"Name": config.Name,
	}

	params := map[string]interface{}{}

	for _, param := range config.Parameters {
		if param.Actual != nil {
			params[param.Name] = param.Actual
			continue
		}
		params[param.Name] = param.Default
	}

	if len(params) > 0 {
		model["Parameter"] = params
	}

	return model
}

func render(tpl string, model interface{}) ([]byte, error) {
	t, err := template.New("template").Parse(string(tpl))
	if err != nil {
		return nil, errors.Wrap(err, "error parsing template")
	}

	buffer := bytes.NewBuffer([]byte{})
	err = t.Execute(buffer, model)
	if err != nil {
		return nil, errors.Wrap(err, "error rendering template")
	}

	return buffer.Bytes(), nil
}
