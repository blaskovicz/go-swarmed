package swarmed

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const secretPath = "/run/secrets"

/* LoadSecrets translates all files in /var/secrets to a corresponding env var in the process.
   For example, /var/secrets/db_password (with contents password) would be translated to
	 env variable DB_PASSWORD with value password.

	 This is a feature most useful to [Docker swarm](https://docs.docker.com/engine/swarm/secrets/)
	 users.
*/
func LoadSecrets() error {
	return LoadSecretsWithOpts("", false)
}

/* LoadSecretsWithOpts is the same as LoadSecrets but applies additional data transforms
   to env variables. For example, if you define the secret app-name_redis_url_v1, if you set
	 prefix to app-name (any casing), the env var will be exposed as REDIS_URL_V1. Furthermore,
	 if you enable version suffix, the _v1 in the example wil also be stripped, exposing the variable
	 as REDIS_URL (useful for version deploys).
*/
var delim = []string{"_", "-"}

func LoadSecretsWithOpts(prefix string, removeVersionSuffix bool) error {
	logOverrides := os.Getenv("SWARMED_LOG_OVERRIDES")

	if prefix != "" {
		prefix = strings.ToUpper(prefix)
	}

	f, err := os.Open(secretPath)
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	files, err := f.Readdirnames(0)
	if err != nil {
		return err
	}

	for _, f := range files {
		f = path.Base(f)
		b, err := ioutil.ReadFile(path.Join(secretPath, f))
		if err != nil {
			return err
		}

		val := string(b)
		key := strings.ToUpper(f)

		// app-prefix_foo_v1 -> foo_v1
		if prefix != "" {
			for _, c := range delim {
				if strings.HasPrefix(key, prefix+c) {
					key = strings.TrimPrefix(key, prefix+c)
					break
				}
			}
		}

		// foo_v1 -> foo
		if removeVersionSuffix {
			for _, c := range delim {
				parts := strings.Split(key, c)
				partc := len(parts)
				if partc > 0 && strings.HasPrefix(parts[partc-1], "V") {
					key = strings.Join(parts[:partc-1], c)
					break
				}
			}
		}

		if prev := os.Getenv(key); prev != "" && logOverrides != "false" {
			fmt.Printf("[swarmed] %s env value overridden.\n", key)
		}

		err = os.Setenv(key, val)
		if err != nil {
			return err
		}
	}
	return nil
}
