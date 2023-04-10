package mgmt

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/conf"
)

// version represents an engine version as sent by the EVC backend.
type version struct {
	ID   string `json:"id"`
	Path string `json:"path"`
}

// engine represents an engine as sent by the EVC backend.
type engine struct {
	Name     string    `json:"name"`
	Versions []version `json:"versions"`
}

// GetEngineInstance returns an EngineInstance for the given engine and version.
func GetEngineInstance(engine string, version string) (*EngineInstance, error) {
	ver, err := ParseVersion(version, DotVersionStyle)

	if err != nil {
		return nil, err
	}

	inst := EngineInstance{
		Engine:  engine,
		Version: *ver,
	}

	return &inst, nil
}

// GetAvailableEngines returns a list of available engines.
// If defaultEngine is not empty, only the default engine will be returned.
func GetAvailableEngines(defaultEngine string) (*[]Engine, error) {
	hasDefaultEngine := defaultEngine != ""
	url := conf.GetEVCConfig().GetURL() + "/engines"

	if hasDefaultEngine {
		url += "/" + defaultEngine
	}

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("Could not download file: " + resp.Status)
	}

	var body []engine

	data, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if hasDefaultEngine {
		var engine engine
		json.Unmarshal(data, &engine)
		body = append(body, engine)
	} else {
		json.Unmarshal(data, &body)
	}

	engines := make([]Engine, len(body))

	for i, e := range body {
		err := parseEngine(&e, &engines[i])

		if err != nil {
			return nil, err
		}
	}

	return &engines, nil
}

// parseEngine parses an engine from the EVC backend.
// It returns an error if the engine could not be parsed.
// The data is read from the given engine and written to the given target.
func parseEngine(eng *engine, target *Engine) error {
	target.Versions = make([]Version, len(eng.Versions))
	target.Name = eng.Name

	err := parseVersions(eng.Versions, target)

	if err != nil {
		return err
	}

	return nil
}

// parseVersions parses a list of versions of an engine.
func parseVersions(vers []version, target *Engine) error {
	for i, v := range vers {
		parsed, err := ParseVersion(v.ID, DotVersionStyle)

		if err != nil {
			return err
		}

		target.Versions[i] = *parsed
	}

	return nil
}
