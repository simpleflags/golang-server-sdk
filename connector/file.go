package connector

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/simpleflags/evaluation"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

type FileConnector struct {
	sdkKey string
	path   string
}

func NewFileConnector(sdkKey string, path string) (FileConnector, error) {
	envDir := path + "/" + sdkKey
	err := os.MkdirAll(envDir, 0777)
	if err != nil {
		return FileConnector{}, err
	}
	return FileConnector{
		sdkKey: sdkKey,
		path:   envDir,
	}, nil
}

func (f FileConnector) Configurations(ctx context.Context, identifiers ...string) (evaluation.Configurations, error) {
	flagsPath := path.Join(f.path, "flags")
	files, err := ioutil.ReadDir(flagsPath)
	if err != nil {
		return evaluation.Configurations{}, err
	}

	configurations := make(evaluation.Configurations, 0)
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}
		bytes, err := ioutil.ReadFile(path.Join(flagsPath, file.Name()))
		if err != nil {
			return evaluation.Configurations{}, err
		}
		var config evaluation.Configuration
		err = json.Unmarshal(bytes, &config)
		if err != nil {
			return evaluation.Configurations{}, err
		}
		configurations = append(configurations, config)
	}
	return configurations, nil
}

func (f FileConnector) Variables(ctx context.Context, identifiers ...string) ([]evaluation.Variable, error) {
	varsPath := path.Join(f.path, "variables")
	files, err := ioutil.ReadDir(varsPath)
	if err != nil {
		return []evaluation.Variable{}, err
	}

	variables := make([]evaluation.Variable, 0)
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}
		bytes, err := ioutil.ReadFile(path.Join(varsPath, file.Name()))
		if err != nil {
			return []evaluation.Variable{}, err
		}
		var variable evaluation.Variable
		err = json.Unmarshal(bytes, &variable)
		if err != nil {
			return []evaluation.Variable{}, err
		}
		variables = append(variables, variable)
	}
	return variables, nil
}

func (f FileConnector) Stream(ctx context.Context, updater Updater) error {
	return errors.New("not supported")
}

func (f FileConnector) Close() error {
	return nil
}
