package desired

import (
	"bytes"
	"github.com/Masterminds/sprig"
	mlopsv1 "github.com/cnvrg-operator/api/v1"
	"github.com/markbates/pkger"
	"go.uber.org/zap"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"strings"
	"text/template"
)

func (s *State) GenerateDeployable(cnvrgApp *mlopsv1.CnvrgApp) error {
	var tpl bytes.Buffer
	f, err := pkger.Open(s.TemplatePath)
	if err != nil {
		zap.S().Error(err, "error reading path", "path", s.TemplatePath)
		return err
	}
	b, err := ioutil.ReadAll(f)

	if err != nil {
		zap.S().Error(err, "error reading file", "path", s.TemplatePath)
		return err
	}
	s.Template, err = template.New(s.Name).Funcs(sprig.TxtFuncMap()).Parse(string(b))
	if err != nil {
		zap.S().Error(err, "parse error", "file", s.Name)
		return err
	}
	s.Obj.SetGroupVersionKind(s.GVR)
	if err := s.Template.Execute(&tpl, cnvrgApp); err != nil {
		zap.S().Error(err, "rendering template error", "file", s.TemplatePath)
		return err
	}
	s.ParsedTemplate = tpl.String()
	zap.S().Debug("template: " + s.TemplatePath + "\n" + s.ParsedTemplate)
	dec := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	if _, _, err := dec.Decode([]byte(s.ParsedTemplate), nil, s.Obj); err != nil {
		zap.S().Error(err, "parsing object", "template", s.ParsedTemplate)
		return err
	}
	s.Name = s.Obj.Object["metadata"].(map[string]interface{})["name"].(string)
	return nil
}

func (s *State) Apply() error {
	var kubeconfig string

	kubeconfig = "/Users/Dima/.kube/config"

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		zap.S().Error(err, "can't create k8s config")
		return err
	}
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		zap.S().Error(err, "can't create dynamic client")
		return err
	}
	resourceSchema := schema.GroupVersionResource{
		Group:    s.GVR.Group,
		Version:  s.GVR.Version,
		Resource: strings.ToLower(s.GVR.Kind) + "s",
	}
	zap.S().Debug("creating resource", "file", s.TemplatePath)
	_, err = client.Resource(resourceSchema).Get(s.Name, metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		_, err = client.Resource(resourceSchema).Create(s.Obj, metav1.CreateOptions{})
		if err != nil {
			zap.S().Error(err, "fail to create crd", "file", s.TemplatePath)
			return err
		}
	}
	return nil
}
