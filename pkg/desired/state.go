package desired

import (
	"bytes"
	"github.com/Masterminds/sprig"
	mlopsv1 "github.com/cnvrg-operator/api/v1"
	"github.com/markbates/pkger"
	"go.uber.org/zap"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"strings"
	"text/template"
)

func cnvrgTemplateFuncs() map[string]interface{} {
	return map[string]interface{}{
		"httpScheme": func(httpsEnabled string) string {
			if httpsEnabled == "true" {
				return "https://"
			}
			return "http://"
		},
	}
}

func (s *State) GenerateDeployable(cnvrgApp *mlopsv1.CnvrgApp) error {
	var tpl bytes.Buffer
	f, err := pkger.Open(s.TemplatePath)
	if err != nil {
		zap.S().Error(err, "error reading path", "path", s.TemplatePath)
		return err
	}
	b, err := ioutil.ReadAll(f)

	if err != nil {
		zap.S().Errorf("%v, error reading file: %v", err, s.TemplatePath)
		return err
	}

	s.Template, err = template.New(strings.ReplaceAll(s.TemplatePath, "/", "-")).
		Funcs(sprig.TxtFuncMap()).
		Funcs(cnvrgTemplateFuncs()).
		Parse(string(b))
	if err != nil {
		zap.S().Errorf("%v, template: %v", err, s.TemplatePath)
		return err
	}
	s.Obj.SetGroupVersionKind(s.GVR)
	if err := s.Template.Execute(&tpl, cnvrgApp); err != nil {
		zap.S().Error(err, "rendering template error", "file", s.TemplatePath)
		return err
	}
	s.ParsedTemplate = tpl.String()
	zap.S().Infof("parsing: %v ", s.TemplatePath)
	zap.S().Debug("template: " + s.TemplatePath + "\n" + s.ParsedTemplate)
	dec := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	if _, _, err := dec.Decode([]byte(s.ParsedTemplate), nil, s.Obj); err != nil {
		zap.S().Errorf("%v, template: %v", err, s.ParsedTemplate)
		return err
	}
	s.Name = s.Obj.Object["metadata"].(map[string]interface{})["name"].(string)

	return nil
}
