# Tote

## Description
A configuration fetching framework to take some of the pain out of reading
config files in go.

## Usage
By default tote reads from a file defined by the environment variable
`TOTE_CONFIG_FILE`.

In this example let's say you've exported `TOTE_CONFIG_FILE=./config.yaml` and
the contents of `./config.yaml` is as follows:

``` yaml
database:
    port: 1234
appname: my-test-app
```
\* Refer to the [go-yaml/yaml](https://godoc.org/gopkg.in/yaml.v2) for more
information about the mapping between your struct and the yaml file.

``` go
package main

import (
    "fmt"

    "github.com/daihasso/tote"
)

type MyConfig struct {
    Database struct {
        Port int
    }
    AppName string
}

func main() {
    myConfig := MyConfig{}
    tote.ReadConfig(&myConfig)

    fmt.Println("Database port:", myConfig.Database.Port)
    fmt.Println("App name:", myConfig.Name)
}
```

``` shell
Database port: 1234
App name: my-test-app
```

## Defining Multiple Search Paths
You can define any number of potential search paths (with mixed backends) via
the `AddPaths` option like so:

``` go
tote.ReadConfig(
    &myConfig, tote.AddPaths("config.yaml", "s3://my-bucket/config.yaml"),
)
```
The provided paths will be searched in order with the first found taking
precedence. (If the environment variable `TOTE_CONFIG_FILE` is specified it
will always take precedence over the other defined paths)

## S3 Compatibility
### Usage
By default tote is only configured to check on disk for files but it is built
to support S3 as well. If you'd like to search for a config in S3 as well simply
use the `WithS3Client` option like so:

``` go
    tote.ReadConfig(
        myConfig, tote.WithS3Client(myS3Client),
    )
```

### Disabling S3
If you're not using S3 and are concerned by the added bulk of the s3 libraries
(I don't blame you, it's YUUUGE) just build your binary with the build flag
`nos3` like so:

``` shell
go build . -tags nos3
```
