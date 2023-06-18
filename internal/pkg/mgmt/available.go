package mgmt

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/conf"
	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/sys"
	"golang.org/x/exp/slices"
)

type versionCache map[Version]*EngineInstance
type engineCache map[string]versionCache

var cache engineCache = make(engineCache)

func BestMatch(name string, version Version) (*EngineInstance, error) {
	if cache[name] == nil {
		cache[name] = make(versionCache)
	}

	if cache[name][version] != nil {
		return cache[name][version], nil
	}

	engines, err := GetAvailableEngines(name)
	device, _ := sys.DeviceInfo()
	var capabilities []string

	if len(device.Cpu) > 0 {
		capabilities = device.Cpu[0].Capabilities
	}

	if err != nil {
		return nil, err
	}

	for _, engine := range *engines {
		for _, vari := range engine.Variations {
			if !version.Equals(vari.Version) {
				continue
			}

			max := -1
			inst := &EngineInstance{}

			for _, flav := range vari.Flavours {
				if flav.Os != device.OS || flav.Arch != device.Arch {
					continue
				}

				counter := 0

				for _, cap := range capabilities {
					if slices.Contains(flav.Capabilities, strings.ToLower(cap)) {
						counter++
					}
				}

				//? Require that all capabilities of the engine are met
				if counter < len(flav.Capabilities) {
					continue
				}

				//? Select the engine with the most capabilities
				if counter > max {
					max = counter
					inst = &EngineInstance{
						Engine:  engine.Name,
						Version: vari.Version,
						Id:      flav.Id,
					}
				}
			}

			if max > -1 {
				cache[engine.Name][inst.Version] = inst
				return inst, nil
			}
		}
	}

	return nil, errors.New("could not find engine")
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

	var engines []Engine

	data, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if hasDefaultEngine {
		var engine Engine
		json.Unmarshal(data, &engine)
		engines = append(engines, engine)
	} else {
		json.Unmarshal(data, &engines)
	}

	return &engines, nil
}
