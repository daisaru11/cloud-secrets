# cloud-secrets

A tool to get secrets located in cloud providers.
cloud-secrets traverses environment variables and replaces the variables with decoded one if its values start with `cloud-secrets://` prefix.

## Usage

### CLI

```
$ export FOO=cloud-secrets://aws-parameter-store/my-secrets/foo
$ cloud-secrets exec sh -c 'echo $FOO'
FOO_SECRET_VALUE
```

### library

```

import (
  "log"
  "os"

  "github.com/daisaru11/cloud-secrets/env"
)

func main() {
  log.Println(os.Getenv("FOO")) # => cloud-secrets://aws-parameter-store/my-secrets/foo

  env.ReplaceEnvironmentVariables()

  log.Println(os.Getenv("FOO")) # => FOO_SECRET_VALUE
}

```
