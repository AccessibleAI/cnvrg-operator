package pg

import (
	"bytes"
	"github.com/Masterminds/sprig"
	mlopsv1 "github.com/cnvrg-operator/api/v1"
	"github.com/markbates/pkger"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"text/template"
)

var (
	log = zap.New(zap.UseDevMode(true))
)

var S = []*State{
	{
		Name:           "pg-pvc",
		TemplatePath:   "/pkg/pg/tmpl/pvc.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            PvcGVR,
	},
	{
		Name:           "pg-dep",
		TemplatePath:   "/pkg/pg/tmpl/dep.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            DeploymentGVR,
	},
	{
		Name:           "pg-secret",
		TemplatePath:   "/pkg/pg/tmpl/secret.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            SecretGVR,
	},
	{
		Name:           "pg-svc",
		TemplatePath:   "/pkg/pg/tmpl/svc.tpl",
		Template:       nil,
		ParsedTemplate: "",
		Obj:            &unstructured.Unstructured{},
		GVR:            PvcGVR,
	},
}

func (s *State) InitTemplate(cnvrgApp mlopsv1.CnvrgApp) error {
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
	return nil
}

func GetTemplates() map[string]*template.Template {
	manifests, err := ReadTemplatesFiles()
	if err != nil {

	}

	tmpls, err := LoadTemplates(manifests)
	if err != nil {

	}
	return tmpls

}

func ReadTemplatesFiles() (map[string]string, error) {
	var manifests = make(map[string]string)
	err := pkger.Walk("/pkg/pg/tmpl", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		f, err := pkger.Open(path)
		log.Info(path)
		if err != nil {
			log.Error(err, "error reading path", "path", path)
			return err
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			log.Error(err, "error reading file", "path", path)
			return err
		}
		manifests[info.Name()] = string(b)
		return nil
	})
	if err != nil {
		log.Error(err, "error walking dir", "dir", "/pkg/pg/tmpl")
		return nil, err
	}
	return manifests, nil
}

func LoadTemplates(tmplFiles map[string]string) (map[string]*template.Template, error) {
	var templates = make(map[string]*template.Template)
	for k, v := range tmplFiles {
		tmpl, err := template.New(k).Funcs(sprig.TxtFuncMap()).Parse(v)
		if err != nil {
			log.Error(err, "parse error", "file", k)
			return nil, err
		}
		templates[k] = tmpl
	}
	return templates, nil
}
