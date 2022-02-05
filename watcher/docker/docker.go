// Package docker contains everything related to docker images and containers parsing
package docker

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ahugues/docker-update-watcher/version"
)

type jsonImage struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	OS           string `json:"os"`
	Architecture string `json:"architecture"`
	Digest       string `json:"digest"`
}

type Image struct {
	Namespace    string
	Name         string
	Version      version.Version
	Digest       string
	Architecture string
	OS           string
}

func (i *Image) UnmarshalJSON(data []byte) error {
	var jsonI jsonImage
	if err := json.Unmarshal(data, &jsonI); err != nil {
		return err
	}

	separatedName := strings.Split(jsonI.Name, "/")
	switch len(separatedName) {
	case 1:
		i.Name = separatedName[0]
	case 2:
		i.Name = separatedName[1]
		i.Namespace = separatedName[0]
	default:
		return fmt.Errorf("invalid name %q", jsonI.Name)
	}

	i.Digest = jsonI.Digest

	switch jsonI.Version {
	case "latest":
		i.Version = &version.VersionLatest{}
	default:
		var err error
		if i.Version, err = version.NewSemVer(jsonI.Version); err != nil {
			return fmt.Errorf("invalid version %s: %w", jsonI.Version, err)
		}
	}
	return nil
}

// NeedUpdate returns true if the given image is newer than the current one
func (i *Image) NeedUpdate(comp *Image) (bool, error) {
	switch currentVer := i.Version.(type) {
	case *version.SemVer:
		return i.needUpdateSemVer(currentVer, comp)
	case *version.VersionLatest:
		return i.needUpdateLatest(comp)
	default:
		return false, errors.New("invalid version type")
	}
}

func (i *Image) needUpdateLatest(cmp *Image) (bool, error) {
	_, ok := cmp.Version.(*version.VersionLatest)
	if !ok {
		return false, errors.New("failed to compare version latest")
	}
	return i.Digest != cmp.Digest, nil
}

func (i *Image) needUpdateSemVer(currentVersion *version.SemVer, cmp *Image) (bool, error) {
	switch cmpVer := cmp.Version.(type) {
	case *version.SemVer:
		return currentVersion.Older(cmpVer), nil
	default:
		return false, errors.New("failed to compare semver versions")
	}
}

type initialConf struct {
	Images []Image `json:"initial-images"`
}

func ReadInitialConfig(confLocation string) (*[]Image, error) {
	content, err := os.ReadFile(confLocation)
	if err != nil {
		return nil, errors.New("failed to read initial list of images")
	}

	var res initialConf
	if err := json.Unmarshal(content, &res); err != nil {
		return nil, errors.New("failed to parse initial list of images")
	}
	return &(res.Images), nil
}
