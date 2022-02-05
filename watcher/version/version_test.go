package version_test

import (
	"testing"

	"github.com/ahugues/docker-update-watcher/version"
)

func TestCreateSemver(t *testing.T) {
	t.Run("OK full", func(t *testing.T) {
		v, err := version.NewSemVer("1.2.3")
		if err != nil {
			t.Fatalf("Unexpected error %s", err.Error())
		}
		if v.Major != 1 || v.Minor != 2 || v.Rev != 3 {
			t.Fatalf("Unexpected version %+v", v)
		}
		if v.String() != "1.2.3" {
			t.Fatalf("Unexpected string value %v", v.String())
		}
	})

	t.Run("OK no rev", func(t *testing.T) {
		v, err := version.NewSemVer("1.2")
		if err != nil {
			t.Fatalf("Unexpected error %s", err.Error())
		}
		if v.Major != 1 || v.Minor != 2 || v.Rev != 0 {
			t.Fatalf("Unexpected version %+v", v)
		}
		if v.String() != "1.2.0" {
			t.Fatalf("Unexpected string value %v", v.String())
		}
	})

	t.Run("OK only major", func(t *testing.T) {
		v, err := version.NewSemVer("1")
		if err != nil {
			t.Fatalf("Unexpected error %s", err.Error())
		}
		if v.Major != 1 || v.Minor != 0 || v.Rev != 0 {
			t.Fatalf("Unexpected version %+v", v)
		}
		if v.String() != "1.0.0" {
			t.Fatalf("Unexpected string value %v", v.String())
		}
	})

	t.Run("not a semver", func(t *testing.T) {
		_, err := version.NewSemVer("latest")
		if err == nil {
			t.Fatal("Unexpected nil error")
		} else if err.Error() != `invalid major version latest: strconv.ParseInt: parsing "latest": invalid syntax` {
			t.Fatalf("Unexpected error message %s", err.Error())
		}
	})

	t.Run("invalid major", func(t *testing.T) {
		_, err := version.NewSemVer("toto.1.2")
		if err == nil {
			t.Fatal("Unexpected nil error")
		} else if err.Error() != `invalid major version toto: strconv.ParseInt: parsing "toto": invalid syntax` {
			t.Fatalf("Unexpected error message %s", err.Error())
		}
	})

	t.Run("invalid minor", func(t *testing.T) {
		_, err := version.NewSemVer("1.toto.2")
		if err == nil {
			t.Fatal("Unexpected nil error")
		} else if err.Error() != `invalid minor version toto: strconv.ParseInt: parsing "toto": invalid syntax` {
			t.Fatalf("Unexpected error message %s", err.Error())
		}
	})

	t.Run("invalid rev", func(t *testing.T) {
		_, err := version.NewSemVer("1.2.toto")
		if err == nil {
			t.Fatal("Unexpected nil error")
		} else if err.Error() != `invalid revision toto: strconv.ParseInt: parsing "toto": invalid syntax` {
			t.Fatalf("Unexpected error message %s", err.Error())
		}
	})
}

func TestCompareSemver(t *testing.T) {
	t.Run("Equal", func(t *testing.T) {
		v1 := version.SemVer{1, 2, 3}
		v2 := version.SemVer{1, 2, 3}
		v3 := version.SemVer{1, 2, 4}
		v4 := version.VersionLatest{}
		if !v1.Equal(&v2) {
			t.Fatal("Should be equal")
		}

		if v1.Equal(&v3) {
			t.Fatal("Should not be equal")
		}
		if v1.Equal(&v4) {
			t.Fatal("Should not be equal to latest")
		}
	})

	t.Run("Older", func(t *testing.T) {
		v1 := version.SemVer{1, 2, 3}
		v2 := version.SemVer{1, 2, 4}
		v3 := version.SemVer{1, 3, 1}
		v4 := version.SemVer{2, 0, 0}

		for _, v := range []*version.SemVer{&v2, &v3, &v4} {
			if !v1.Older(v) {
				t.Errorf("%v should be older than %v", &v1, v)
			}
		}
	})

	t.Run("Compare with latest", func(t *testing.T) {
		v1 := version.SemVer{1, 2, 3}
		v2 := version.VersionLatest{}

		if !v1.Older(&v2) {
			t.Fatal("Should be older than latest")
		}
	})
}

func TestLatest(t *testing.T) {
	v1 := version.SemVer{1, 2, 3}
	v2 := version.VersionLatest{}
	v3 := version.VersionLatest{}

	t.Run("String", func(t *testing.T) {
		if v2.String() != "latest" {
			t.Fatalf("Unexpected string value %s", v2.String())
		}
	})

	t.Run("Equal", func(t *testing.T) {
		if v2.Equal(&v1) {
			t.Fatal("Should not be equal")
		}
		if !v2.Equal(&v3) {
			t.Fatal("Should be equal")
		}
	})

	t.Run("Older", func(t *testing.T) {
		if v2.Older(&v1) {
			t.Fatal("Should not be older")
		}
		if v2.Older(&v3) {
			t.Fatal("Should not be older")
		}
	})
}
