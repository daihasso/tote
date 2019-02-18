package tote

import (
    "fmt"
    "io/ioutil"
    "os"
    "strings"
    "reflect"

    logging "github.com/daihasso/slogging"
    "gopkg.in/yaml.v2"
    "github.com/pkg/errors"
)

func readConfigFile(path string) ([]byte, error) {
    configYAMLData, err := ioutil.ReadFile(path)
    if os.IsNotExist(err) {
        return nil, errors.New(
            fmt.Sprintf("Config file '%s' doesn't exist.", path),
        )
    }
    if err != nil {
        return nil, err
    }

    return configYAMLData, nil
}

func unmarshalConfigBytes(
    configFileData []byte, config interface{},
) error {
    err := yaml.Unmarshal(configFileData, config)
    if err != nil {
        return errors.Wrap(err, "Error while unmarshaling file to yaml.")
    }

    return nil
}

func readEmbeddedConfig(
    configFileData []byte, configKey string, embeddedConfigStruct interface{},
) error {
    var rawMapItems yaml.MapSlice
    err := yaml.Unmarshal(configFileData, &rawMapItems)
    if err != nil {
        return err
    }

    for _, mapItem := range rawMapItems {
        stringValue, ok := mapItem.Key.(string)
        if !ok {
            continue
        }
        if strings.ToLower(stringValue) == strings.ToLower(configKey) {
            marshaledValue, err := yaml.Marshal(mapItem.Value)
            if err != nil {
                return err
            }
            err = yaml.Unmarshal(
                marshaledValue, embeddedConfigStruct,
            )

            return err
        }
    }

    return nil
}

// ReadConfig reads in yaml data at provided paths into the provided interface.
func ReadConfig(config interface{}, allOptions ...Option) error {
    if reflect.TypeOf(config).Kind() != reflect.Ptr {
        return errors.Errorf(
            "ReadConfig requires a pointer to your config struct, %[1]T " +
                "was provided but a *%[1]T is required",
            config,
        )
    }
    opts := newOptions(allOptions)
    paths := opts.configPath
    envPrefix := opts.envVarPrefix
    envVar := fmt.Sprintf("%s_CONFIG_FILE", envPrefix)
    if newPath, exists := os.LookupEnv(envVar); exists {
        paths = append([]string{newPath}, paths...)
    }

    found := false
    for _, path := range paths {
        configReader, err := opts.pathReader.Read(path)
        if err != nil {
            logging.Warn("Failed to load config data frome file.").With(
                "error", err,
            ).And(
                "path", path,
            ).Send()
            continue
        }
        configBytes, err := ioutil.ReadAll(configReader)
        if err != nil {
            return errors.Wrap(
                err, "Error while reading data returned from path",
            )
        }

        err = unmarshalConfigBytes(configBytes, config)
        if err != nil {
            return errors.Wrap(
                err, "Error while reading data from config file",
            )
        }

        for key, embedded := range opts.embeddedConfigs {
            if reflect.TypeOf(embedded).Kind() != reflect.Ptr {
                return errors.Errorf(
                    "ReadConfig requires a pointer to your embedded " +
                        "config(s), %[1]T was provided but a *%[1]T is " +
                        "required",
                    embedded,
                )
            }
            err = readEmbeddedConfig(configBytes, key, embedded)
            if err != nil {
                return errors.Wrap(
                    err, "Error while reading embedded config",
                )
            }
        }

        found = true
        break
    }

    if !found {
        return errors.Errorf(
            "Couldn't load any of the provided paths %v", paths,
        )
    }

    err := readConfigFromEnvironment(config, envPrefix)
    if err != nil {
        return errors.Wrap(
            err, "Error while reading config values from environment",
        )
    }

    for key, embedded := range opts.embeddedConfigs {
        err := readConfigFromEnvironment(
            embedded, envPrefix, strings.ToUpper(key),
        )
        if err != nil {
            return errors.Wrap(
                err,
                "Error while reading embedded config values from environment",
            )
        }
    }

    return nil
}
