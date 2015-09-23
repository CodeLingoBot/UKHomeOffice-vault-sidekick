/*
Copyright 2015 Home Office All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"fmt"
	"net/url"
	"time"
)

// config ... the command line configuration
type config struct {
	// the url for th vault server
	vaultURL string
	// a file containing the authenticate options
	vaultAuthFile string
	// the authentication options
	vaultAuthOptions map[string]string
	// the place to write the resources
	outputDir string
	// switch on dry run
	dryRun bool
	// skip tls verify
	tlsVerify bool
	// the resource items to retrieve
	resources *VaultResources
	// the interval for producing statistics
	statsInterval time.Duration
}

var (
	options config
)

func init() {
	// step: setup some defaults
	options.resources = new(VaultResources)
	options.vaultAuthOptions = map[string]string{VaultAuth: "token"}

	flag.StringVar(&options.vaultURL, "vault", getEnv("VAULT_ADDR", "https://127.0.0.1:8200"), "the url the vault service is running behind (VAULT_ADDR if available)")
	flag.StringVar(&options.vaultAuthFile, "auth", "", "a configuration file in a json or yaml containing authentication arguments")
	flag.StringVar(&options.outputDir, "output", getEnv("VAULT_OUTPUT", "/etc/secrets"), "the full path to write the protected resources (VAULT_OUTPUT if available)")
	flag.BoolVar(&options.dryRun, "dryrun", false, "perform a dry run, printing the content to screen")
	flag.BoolVar(&options.tlsVerify, "tls-skip-verify", false, "whether to check and verify the vault service certificate")
	flag.DurationVar(&options.statsInterval, "stats", time.Duration(5)*time.Minute, "the interval to produce statistics on the accessed resources")
	flag.Var(options.resources, "cn", "a resource to retrieve and monitor from vault (e.g. pki:name:cert.name, secret:db_password, aws:s3_backup)")
}

// parseOptions ... validate the command line options and validates them
func parseOptions() error {
	flag.Parse()

	return validateOptions(&options)
}

// validateOptions ... parses and validates the command line options
func validateOptions(cfg *config) error {
	// step: validate the vault url
	_, err := url.Parse(cfg.vaultURL)
	if err != nil {
		return fmt.Errorf("invalid vault url: '%s' specified", cfg.vaultURL)
	}

	// step: read in the token if required
	if cfg.vaultAuthFile != "" {
		if exists, _ := fileExists(cfg.vaultAuthFile); !exists {
			return fmt.Errorf("the token file: %s does not exists, please check", cfg.vaultAuthFile)
		}

		if options.vaultAuthOptions, err = readConfigFile(options.vaultAuthFile); err != nil {
			return fmt.Errorf("unable to read in authentication options from: %s, error: %s", cfg.vaultAuthFile, err)
		}
	}

	return nil
}
