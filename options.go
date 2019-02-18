package tote

import (
	"strings"

    "github.com/daihasso/peechee"
)

type options struct {
    configPath []string
    embeddedConfigs map[string]interface{}
	envVarPrefix string
    pathReader *peechee.PathReader
}

// Option is an option for reading a config.
type Option func(*options)

// AddPaths adds the provided paths to the search paths which will be checked
// in the order provided.
func AddPaths(paths ...string) Option {
    return func(opts *options) {
        opts.configPath = append(opts.configPath, paths...)
    }
}

// AddEmbedded adds a struct (in) that's embedded under some key.
func AddEmbedded(key string, in interface{}) Option {
    return func(opts *options) {
        opts.embeddedConfigs[key] = in
    }
}

// OverrideEnvVarPrefix overrides the defaultEnvironmentPrefix used in finding
// environment variables that override config values.
func OverrideEnvVarPrefix(prefix string) Option {
	return func(opts *options) {
		opts.envVarPrefix = strings.ToUpper(prefix)
	}
}

// WithPathReader lets you override the internal PathReader with one you've
// defined seperately if needed.
func WithPathReader(pathReader *peechee.PathReader) Option {
    return func(opts *options) {
        opts.pathReader = pathReader
    }
}

func newOptions(allOptions []Option) *options {
    opts := &options{
		configPath: make([]string, 0),
		embeddedConfigs: make(map[string]interface{}, len(allOptions)),
		envVarPrefix: defaultEnvironmentPrefix,
        pathReader: peechee.NewPathReader(peechee.WithFilesystem()),
	}
    for _, opt := range allOptions {
        opt(opts)
    }

	return opts
}
