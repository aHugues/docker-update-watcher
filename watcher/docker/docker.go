// Package docker contains everything related to docker images and containers parsing
package docker

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type SemVer struct {
	Major int
	Minor int
	Rev   int
}

func (v *SemVer) Equal(cmp *SemVer) bool {
	return v.Major == cmp.Major && v.Minor == cmp.Minor && v.Rev == cmp.Rev
}

func (v *SemVer) Lower(cmp *SemVer) bool {
	return (v.Major < cmp.Major) || (v.Major == cmp.Major && v.Minor < cmp.Minor) || (v.Major == cmp.Major && v.Minor == cmp.Minor && v.Rev < cmp.Rev)
}

func NewSemVer(raw string) (ret *SemVer, err error) {
	vals := strings.Split(raw, ".")
	if len(vals) == 0 {
		return nil, errors.New("invalid format")
	}
	ret.Major, err = strconv.Atoi(vals[0])
	if err != nil {
		return nil, fmt.Errorf("Invalid major version %s: %w", vals[0], err)
	}
	if len(vals) < 2 {
		return
	}
	ret.Minor, err = strconv.Atoi(vals[1])
	if err != nil {
		return nil, fmt.Errorf("Invalid minor version %s: %w", vals[1], err)
	}
	if len(vals) < 3 {
		return
	}
	if err != nil {
		return nil, fmt.Errorf("Invalid revision %s: %w", vals[1], err)
	}
	return
}

type VersionLatest string

type Versioner interface {
	Lower(cmp Versioner) bool
}

type Image struct {
	Name    string
	Version interface{}
	Digest  string
}

// NeedUpdate returns true if the given image is newer than the current one
func (i *Image) NeedUpdate(comp *Image) (bool, error) {
	switch currentVer := i.Version.(type) {
	case SemVer:
		return i.needUpdateSemVer(&currentVer, comp)
	case VersionLatest:
		return i.needUpdateLatest(comp)
	default:
		return nil
	}
}

func (i *Image) needUpdateLatest(cmp *Image) (bool, error) {
	versionStr, ok := cmp.Version.(string)
	if !ok {
		return false, errors.New("failed to compare version latest")
	}
	if versionStr != "latest" {
		return false, fmt.Errorf("unexpected version %s", versionStr)
	}
	return i.Digest != cmp.Digest, nil
}

func (i *Image) needUpdateSemVer(currentVersion *SemVer, cmp *Image) (bool, error) {
	switch cmpVer := cmp.Version.(type) {
	case SemVer:
		return currentVersion.Lower(&cmpVer), nil
	default:
		return false, errors.New("failed to compare semver versions")
	}
}
