# go-swarmed
> Load docker swarm secrets at runtime in golang.

# About

To load secrets at runtime, use the following code:

```
import "github.com/blaskovicz/go-swarmed"

func main() {
  err := swarmed.LoadSecrets()
  if err != nil {
    panic(err)
  }
}
```

`swarmed.LoadSecrets()` translates all files in `/var/secrets` to a corresponding env var in the process.
For example, `/var/secrets/db_password` (with contents password) would be translated to
env variable `DB_PASSWORD` with value password.

Alternatively, you can call `swarmed.LoadSecretsWithOpts`:

```
import "github.com/blaskovicz/go-swarmed"

func main() {
  err := swarmed.LoadSecretsWithOpts("app-name", true)
  if err != nil {
    panic(err)
  }
}
```

`swarmed.LoadSecretsWithOpts()` is the same as `swarmed.LoadSecrets()`, but applies additional data transforms
 to env variables. For example, if you define the secret `app-name_redis_url_v1`, if you set
 prefix to app-name (any casing), the env var will be exposed as `REDIS_URL_V1`. Furthermore,
 if you enable version suffix stripping, the `_v1` in this example wil also be removed, exposing the variable
 as `REDIS_URL` (useful for versioned deploys).

To disable logging of when an evironment variable is overridden, set `SWARMED_LOG_OVERRIDES=false` in your env.


This is a feature most useful to [Docker swarm](https://docs.docker.com/engine/swarm/secrets/)
users.
