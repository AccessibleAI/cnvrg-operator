package pg

import (
	"github.com/Masterminds/sprig"
	"github.com/markbates/pkger"
	"io/ioutil"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"text/template"
)

var (
	log = zap.New(zap.UseDevMode(true))
)

func Defaults() Pg {
	return Pg{
		Enabled:        "true",
		SecretName:     "cnvrg-pg-secret",
		Image:          "centos/postgresql-12-centos7",
		Port:           5432,
		StorageSize:    "80Gi",
		SvcName:        "postgres",
		Dbname:         "cnvrg_production",
		Pass:           "pg_pass",
		User:           "cnvrg",
		RunAsUser:      26,
		FsGroup:        26,
		StorageClass:   "use-default",
		CPURequest:     4,
		MemoryRequest:  "4Gi",
		MaxConnections: 100,
		SharedBuffers:  "64Mb",
		HugePages: HugePages{
			Enabled: "false",
			Size:    "2Mi",
			Memory:  "",
		},
	}
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
