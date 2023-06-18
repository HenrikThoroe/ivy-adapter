// Package mgmt contains the logic for managing game engines.
// This includes downloading engines and fetch available engine configurations.
// It uses the conf package to connect to the EVC server and deterime the engine store location.
package mgmt

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/conf"
)

// Version is a struct that represents a version of a game engine.
type Version struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
	Patch int `json:"patch"`
}

type Flavour struct {
	Id           string   `json:"id"`
	Os           string   `json:"os"`
	Arch         string   `json:"arch"`
	Capabilities []string `json:"capabilities"`
}

type Variation struct {
	Version  Version   `json:"version"`
	Flavours []Flavour `json:"flavours"`
}

// Engine is a struct that represents a game engine.
type Engine struct {
	Name       string      `json:"name"`
	Variations []Variation `json:"variations"`
}

// EngineInstance is a struct that represents a specific instance of a game engine.
// This is used to represent a specific version of a game engine.
type EngineInstance struct {
	Engine  string
	Version Version
	Id      string
}

// VersionStyle is an enum that represents the style of a version string.
type VersionStyle int

const (
	DotVersionStyle     VersionStyle = iota // DotVersionStyle is the style of a version string that uses dots.
	UrlSaveVersionStyle                     // UrlSaveVersionStyle is the style of a version string that uses dashes.
)

// String returns a string representation of the version.
// If style is 'DotVersionStyle', the version will be returned in the format of "$major.$minor.$patch".
// If style is 'UrlSaveVersionStyle', the version will be returned in the format of "v$major-$minor-$patch".
func (v *Version) String(style VersionStyle) string {
	if style == DotVersionStyle {
		return strconv.Itoa(v.Major) + "." + strconv.Itoa(v.Minor) + "." + strconv.Itoa(v.Patch)
	}

	return "v" + strconv.Itoa(v.Major) + "-" + strconv.Itoa(v.Minor) + "-" + strconv.Itoa(v.Patch)
}

// GetInstances returns a slice of EngineInstance structs that represent all the versions of the engine.
func (e *Engine) GetInstances() []EngineInstance {
	instances := make([]EngineInstance, 0)

	for _, v := range e.Variations {
		for _, f := range v.Flavours {
			instances = append(instances, EngineInstance{
				Engine:  e.Name,
				Version: v.Version,
				Id:      f.Id,
			})
		}
	}

	return instances
}

// FileName returns the file name of the engine instance.
func (inst EngineInstance) FileName() string {
	return inst.Engine + "_" + inst.Version.String(UrlSaveVersionStyle) + "_" + inst.Id
}

// Path returns the path to the engine instance.
// The conf package is used to determine the engine store location.
func (inst EngineInstance) Path() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(conf.GetEngineStore(), inst.FileName()+".exe")
	}

	return filepath.Join(conf.GetEngineStore(), inst.FileName())
}

// IsInstalled returns whether the engine instance is installed.
func (inst EngineInstance) IsInstalled() bool {
	_, err := os.Stat(inst.Path())
	return err == nil
}

// URL returns the URL to the engine instance on the EVC server.
func (inst EngineInstance) URL() string {
	return conf.GetEVCConfig().GetURL() + "/engines/bin/" + inst.Engine + "/" + inst.Id
}

// ParseVersion parses a version string and returns a Version struct.
// If style is 'DotVersionStyle', the version string must be in the format of "$major.$minor.$patch".
// If style is 'UrlSaveVersionStyle', the version string must be in the format of "$major-$minor-$patch".
// A leading 'v' is ignored.
func ParseVersion(ver string, style VersionStyle) (*Version, error) {
	var version Version
	var separator string

	if style == DotVersionStyle {
		separator = "."
	} else {
		separator = "-"
	}

	clean := strings.Replace(ver, "v", "", 1)
	parts := strings.Split(clean, separator)

	if len(parts) != 3 {
		return nil, errors.New("Invalid version: " + ver)
	}

	if major, err := strconv.Atoi(parts[0]); err != nil {
		return nil, errors.New("Invalid major version: " + ver)
	} else {
		version.Major = major
	}

	if minor, err := strconv.Atoi(parts[1]); err != nil {
		return nil, errors.New("Invalid minor version: " + ver)
	} else {
		version.Minor = minor
	}

	if patch, err := strconv.Atoi(parts[2]); err != nil {
		return nil, errors.New("Invalid patch version: " + ver)
	} else {
		version.Patch = patch
	}

	return &version, nil
}

func (v Version) Equals(other Version) bool {
	return v.Major == other.Major && v.Minor == other.Minor && v.Patch == other.Patch
}
