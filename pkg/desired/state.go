package desired

import (
	"bytes"
	"github.com/Masterminds/sprig"
	mlopsv1 "github.com/cnvrg-operator/api/v1"
	"github.com/markbates/pkger"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"text/template"
)

var (
	log = zap.New(zap.UseDevMode(true))
)

func (s *State) GenerateDeployable(cnvrgApp *mlopsv1.CnvrgApp) error {
	var tpl bytes.Buffer
	f, err := pkger.Open(s.TemplatePath)
	if err != nil {
		log.Error(err, "error reading path", "path", s.TemplatePath)
		return err
	}
	b, err := ioutil.ReadAll(f)

	if err != nil {
		log.Error(err, "error reading file", "path", s.TemplatePath)
		return err
	}
	s.Template, err = template.New(s.Name).Funcs(sprig.TxtFuncMap()).Parse(string(b))
	if err != nil {
		log.Error(err, "parse error", "file", s.Name)
		return err
	}
	s.Obj.SetGroupVersionKind(s.GVR)
	if err := s.Template.Execute(&tpl, cnvrgApp); err != nil {
		log.Error(err, "rendering template error", "file", s.Name)
		return err
	}
	s.ParsedTemplate = tpl.String()
	log.Info("template: " + s.TemplatePath + "\n" + s.ParsedTemplate)
	dec := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	if _, _, err := dec.Decode([]byte(s.ParsedTemplate), nil, s.Obj); err != nil {
		log.Error(err, "parsing object", "template", s.ParsedTemplate)
		return err
	}
	s.Name = s.Obj.Object["metadata"].(map[string]interface{})["name"].(string)

	return nil
}
