package desired

import (
	"bytes"
	"context"
	"embed"
	v1mlops "github.com/AccessibleAI/cnvrg-operator/api/v1"
	"github.com/Masterminds/sprig"
	"github.com/go-logr/logr"
	"github.com/imdario/mergo"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	"regexp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"text/template"
)

var defaultDataLoaderRegex, _ = regexp.Compile(`mlops.cnvrg.io\/default-loader:.*"true"`)
var ownRegex, _ = regexp.Compile(`mlops.cnvrg.io\/own:.*"true"`)
var updatableRegex, _ = regexp.Compile(`mlops.cnvrg.io\/updatable:.*"true"`)
var ocpRouteRegex, _ = regexp.Compile(`kind:.*Route`)
var nginxIngressRegex, _ = regexp.Compile(`kind:.*Ingress`)
var istioVsRegex, _ = regexp.Compile(`kind:.*VirtualService`)

type LoadFilter struct {
	Ingress       *v1mlops.IngressType
	AssetName     []string
	DefaultLoader bool
}

type AssetsGroup struct {
	Assets   []*Asset
	Fs       embed.FS
	BasePath string
	Log      logr.Logger
	Filter   *LoadFilter
}

type Asset struct {
	Name           string
	RawTemplate    string
	ParsedTemplate string
	Obj            *unstructured.Unstructured
	Own            bool
	Updatable      bool
	Error          error
}

func NewAsset(name, rawTemplate string) *Asset {
	asset := &Asset{
		Name:        name,
		Obj:         &unstructured.Unstructured{},
		RawTemplate: rawTemplate,
	}
	asset.Own = asset.own()
	asset.Updatable = asset.updatable()
	return asset
}

func (a *Asset) Default() bool {
	return defaultDataLoaderRegex.MatchString(a.RawTemplate)
}

func (a *Asset) updatable() bool {
	return updatableRegex.MatchString(a.RawTemplate)
}

func (a *Asset) own() bool {
	return ownRegex.MatchString(a.RawTemplate)
}

func (a *Asset) IngressType() v1mlops.IngressType {
	if istioVsRegex.MatchString(a.RawTemplate) {
		return v1mlops.IstioIngress
	}
	if nginxIngressRegex.MatchString(a.RawTemplate) {
		return v1mlops.NginxIngress
	}
	if ocpRouteRegex.MatchString(a.RawTemplate) {
		return v1mlops.OpenShiftIngress
	}
	return ""
}

func (a *Asset) Render(data interface{}) error {
	var tpl bytes.Buffer
	tmpl, err := template.New(a.Name).
		Funcs(sprig.TxtFuncMap()).
		Funcs(cnvrgTemplateFuncs()).
		Option("missingkey=error").
		Parse(a.RawTemplate)
	if err != nil {
		a.Error = err
		return err
	}
	if err := tmpl.Execute(&tpl, data); err != nil {
		a.Error = err
		return a.Error
	}
	parsedTemplateBytes := tpl.Bytes()
	a.ParsedTemplate = string(parsedTemplateBytes)

	dec := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)

	if _, _, err := dec.Decode(parsedTemplateBytes, nil, a.Obj); err != nil {
		a.Error = err
		return a.Error
	}

	return a.Error
}

func NewAssetsGroup(fs embed.FS, basePath string, log logr.Logger, filter *LoadFilter) *AssetsGroup {

	return &AssetsGroup{
		Assets:   []*Asset{},
		Fs:       fs,
		BasePath: basePath,
		Log:      log.WithValues("basePath", basePath),
		Filter:   filter,
	}
}

func (g *AssetsGroup) LoadAssets() error {
	dirEntries, err := g.Fs.ReadDir(g.BasePath)
	if err != nil {
		return err
	}
	for _, e := range dirEntries {
		if e.IsDir() {
			continue
		}
		f, err := g.Fs.ReadFile(g.BasePath + "/" + e.Name())
		if err != nil {
			g.Log.WithValues("assetName", e.Name()).Error(err, "error reading asset file")
			continue
		}
		a := NewAsset(e.Name(), string(f))
		if g.pass(a) {
			g.Assets = append(g.Assets, a)
		}
	}
	return nil
}

func (g *AssetsGroup) pass(a *Asset) bool {

	if g.Filter == nil {
		return true
	}

	if len(g.Filter.AssetName) > 0 {
		for _, asset := range g.Filter.AssetName {
			if asset == a.Name {
				return true
			}
		}
		return false
	}

	if a.Default() != g.Filter.DefaultLoader {
		return false
	}

	if g.Filter.Ingress != nil && a.IngressType() != "" {
		if a.IngressType() != *g.Filter.Ingress {
			return false
		}
	}

	return true
}

func (g *AssetsGroup) Render(data interface{}) error {
	var e error
	for _, a := range g.Assets {
		if err := a.Render(data); err != nil {
			g.Log.WithValues("assetName", a.Name).Error(err, "error parsing asset")
			e = err
			continue
		}
	}

	return e
}

func (g *AssetsGroup) Apply(spec v1.Object, c client.Client, s *runtime.Scheme, log logr.Logger) (e error) {

	for _, a := range g.Assets {
		g.Log.Info("applying", "assetName", a.Name)
		if a.Own {
			if err := ctrl.SetControllerReference(spec, a.Obj, s); err != nil {
				log.Error(err, "error setting controller reference", "name", a.Obj.GetName())
				e = err
			}
		}

		ctx := context.Background()
		existingObj := &unstructured.Unstructured{}
		existingObj.SetGroupVersionKind(a.Obj.GroupVersionKind())

		err := c.Get(ctx, types.NamespacedName{Name: a.Obj.GetName(), Namespace: a.Obj.GetNamespace()}, existingObj)

		// if object not found, all good, just create it and continue
		if errors.IsNotFound(err) {
			if err := c.Create(ctx, a.Obj); err != nil {
				g.Log.Error(err, "failed create object")
			}
			continue
		}

		// if failed to get the object
		// and not b/c the object is not just not found, return error
		if err != nil {
			g.Log.Error(err, "failed retrieve the object")
			continue
		}

		// if object is not updatable, skip it and continue
		if !a.Updatable {
			log.Info("asset is not updatable", "assetName", a.Name)
			continue
		}

		// need to merge the object before applying the update
		if err = mergo.Merge(existingObj, a.Obj, mergo.WithOverride); err != nil {
			log.Error(err, "can't merge")
			continue
		}

		if err = c.Update(ctx, existingObj); err != nil {
			log.Info("error updating object")
			if updateErr, ok := err.(*errors.StatusError); ok && updateErr.Status().Code == http.StatusUnprocessableEntity {
				if deleteErr := c.Delete(ctx, existingObj); deleteErr == nil {
					if createErr := c.Create(ctx, a.Obj); createErr != nil {
						log.Info("error recreated object")
					} else {
						log.Info("recreated object")
					}
				}
			}
		}

	}
	return
}
