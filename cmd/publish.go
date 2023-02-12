package cmd

import (
	"errors"
	"flag"
	"os"

	"github.com/brer/brer-cli/api"
	"gopkg.in/yaml.v3"
)

func Publish() error {
	var err error = nil

	if len(os.Args) < 3 {
		return errors.New("filepath is missing")
	}
	filepath := os.Args[2]

	// TODO: validate manifest
	var manifest Manifest
	err = ParseManifest(filepath, &manifest)
	if err != nil {
		return err
	}
	if manifest.Version != 0 {
		return errors.New("unsupported manifest version")
	}
	if len(manifest.Functions) <= 0 {
		return errors.New("nothing to publish")
	}

	var tag string
	var token string
	var url string

	f := flag.NewFlagSet("publish", flag.ExitOnError)
	f.StringVar(&tag, "tag", "", "container image tag")
	f.StringVar(&token, "token", "", "authorization header")
	f.StringVar(&url, "url", "", "brer server url")
	f.Parse(os.Args[3:])

	if len(url) <= 0 {
		url = os.Getenv("BRER_URL")
	}
	if len(url) <= 0 {
		return errors.New("url option is required")
	}

	if len(token) <= 0 {
		token = os.Getenv("BRER_TOKEN")
	}
	if len(token) <= 0 {
		return errors.New("token option is required")
	}

	// TODO: this is hacky I suppose
	if len(tag) <= 0 {
		tag = os.Getenv("BRER_TAG")
	}
	if len(tag) > 0 {
		manifest.Image.Tag = tag
	}

	for _, fn := range manifest.Functions {
		var data api.Function
		data.Image = manifest.Image.Repository + ":" + manifest.Image.Tag
		data.SecretName = fn.SecretName
		data.Env = mapEnvs(fn.Env)

		err = api.UpdateFunction(url, token, &data)
		if err != nil {
			return err
		}
	}

	return err
}

func mapEnvs(items []ManifestFunctionEnv) []api.FunctionEnv {
	mapped := make([]api.FunctionEnv, len(items))

	for i, item := range items {
		mapped[i] = api.FunctionEnv{
			Name:      item.Name,
			SecretKey: item.SecretKey,
			Value:     item.Value,
		}
	}

	return mapped
}

type Manifest struct {
	Version   uint8              `yaml:"version"`
	Image     ManifestImage      `yaml:"image"`
	Functions []ManifestFunction `yaml:"functions"`
}

type ManifestImage struct {
	Repository string `yaml:"repository"`
	Tag        string `yaml:"tag"`
}

type ManifestFunction struct {
	Name       string                `yaml:"name"`
	SecretName string                `yaml:"secretName"`
	Env        []ManifestFunctionEnv `yaml:"env"`
}

type ManifestFunctionEnv struct {
	Name      string `yaml:"name"`
	Value     string `yaml:"value"`
	SecretKey string `yaml:"secretKey"`
}

func ParseManifest(filepath string, manifest *Manifest) error {
	file, err := os.Open(filepath)
	if err == nil {
		defer file.Close()
		decoder := yaml.NewDecoder(file)
		err = decoder.Decode(manifest)
	}
	return err
}
