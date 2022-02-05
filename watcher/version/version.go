// package version handles everything related to versions comparison
package version

import (
	"fmt"
	"strconv"
	"strings"
)

type Version interface {
	String() string
	Equal(cmp Version) bool
	Older(cmp Version) bool
}

// SemVer represents a version using semantic versioning (Vx.x.x)
type SemVer struct {
	Major int64
	Minor int64
	Rev   int64
}

func (v *SemVer) Equal(cmp Version) bool {
	switch casted := cmp.(type) {
	case *SemVer:
		return v.Major == casted.Major && v.Minor == casted.Minor && v.Rev == casted.Rev
	default:
		return false
	}
}

func (v *SemVer) Older(cmp Version) bool {
	switch casted := cmp.(type) {
	case *SemVer:
		return (v.Major < casted.Major) || (v.Major == casted.Major && v.Minor < casted.Minor) || (v.Major == casted.Major && v.Minor == casted.Minor && v.Rev < casted.Rev)
	default:
		return true
	}
}

func (v *SemVer) String() string {
	return strconv.FormatInt(v.Major, 10) + "." + strconv.FormatInt(v.Minor, 10) + "." + strconv.FormatInt(v.Rev, 10)
}

func NewSemVer(raw string) (ret *SemVer, err error) {
	var vers SemVer
	vals := strings.Split(raw, ".")
	vers.Major, err = strconv.ParseInt(vals[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid major version %s: %w", vals[0], err)
	}
	if len(vals) < 2 {
		return &vers, nil
	}
	vers.Minor, err = strconv.ParseInt(vals[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid minor version %s: %w", vals[1], err)
	}
	if len(vals) < 3 {
		return &vers, nil
	}
	vers.Rev, err = strconv.ParseInt(vals[2], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid revision %s: %w", vals[2], err)
	}
	return &vers, nil
}

// VersionLatest represents the `latest` tags, no version comparison can be done using this type of versioning, the push date is necessary
type VersionLatest struct {
}

func (vl *VersionLatest) String() string {
	return "latest"
}

func (vl *VersionLatest) Equal(cmp Version) bool {
	switch cmp.(type) {
	case *VersionLatest:
		return true
	default:
		return false
	}
}

// Older returns always false since the version number cannot determine the age of an image in this case
func (vl *VersionLatest) Older(cmp Version) bool {
	return false
}
