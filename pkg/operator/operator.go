/*
 * Copyright (c) 2019, 2020 Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

// The operator package contains types and functions used directly by the Operator main
package operator

import (
	"encoding/json"
	"github.com/oracle/coherence-operator/pkg/data"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
)

const (
	// The environment variable holding the default Coherence image name
	EnvCoherenceImage = "COHERENCE_IMAGE"
	// The environment variable holding the default Coherence Utils image name
	EnvUtilsImage = "UTILS_IMAGE"
)

var config *Config

// The Operator configuration loaded from the data.json file.
type Config struct {
	CoherenceImage string `json:"CoherenceImage,omitempty"`
	UtilsImage string `json:"UtilsImage,omitempty"`
}

// GetDefaultCoherenceImage returns the name of the Coherence image to use.
func (in Config) GetDefaultCoherenceImage() string {
	return in.CoherenceImage
}

// GetDefaultUtilsImage returns the name of the operator utils image to use.
func (in Config) GetDefaultUtilsImage() string {
	return in.UtilsImage
}

func GetOperatorConfig() (Config, error) {
	cfg, err := ensureConfig()
	if err != nil {
		return Config{}, err
	}
	return *cfg, nil
}

func ensureConfig() (*Config, error){
	if config == nil {
		f, err := data.Assets.Open("data.json")
		if err != nil {
			return config, errors.Wrap(err,"finding data.json asset")
		}
		defer f.Close()

		d, err := ioutil.ReadAll(f)
		if err != nil {
			return config, errors.Wrap(err, "reading embedded data.json asset")
		}

		cfg := &Config{}
		err = json.Unmarshal(d, cfg)
		if err != nil {
			return config, errors.Wrap(err, "unmarshalling data.json")
		}

		c, ok := os.LookupEnv(EnvCoherenceImage)
		if ok {
			cfg.CoherenceImage = c
		}

		u, ok := os.LookupEnv(EnvUtilsImage)
		if ok {
			cfg.UtilsImage = u
		}

		config = cfg
	}
	return config, nil
}

