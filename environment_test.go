package tote

import (
    "testing"
	"io/ioutil"
	"os"

	gm "github.com/onsi/gomega"
)

func TestReadConfigFileWithEnvironment(t *testing.T) {
	g := gm.NewGomegaWithT(t)

    lookupEnv = func(key string) (string, bool) {
        if key == "TOTE_TEST_FOO" {
            return "15", true
        }
        return "", false
    }
    defer func() { lookupEnv = os.LookupEnv }()

    fakeConfig := fakeYamlConfigStruct{}

	tempFile, err := ioutil.TempFile("", "config.yaml")
	defer os.Remove(tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile(tempFile.Name(), fakeYamlBytes, 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = ReadConfig(&fakeConfig, AddPaths(tempFile.Name()))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("fakeConfig: %#+v", fakeConfig)
	g.Expect(fakeConfig.Test.Foo).To(gm.Equal(15))
	g.Expect(fakeConfig.Test.Bar).To(gm.Equal("baz"))
}

func TestReadEmbeddedConfigFileWithEnvironment(t *testing.T) {
	g := gm.NewGomegaWithT(t)

    lookupEnv = func(key string) (string, bool) {
        if key == "TOTE_EMBEDDED_NAME" {
            return "Steve", true
        }
        return "", false
    }
    defer func() { lookupEnv = os.LookupEnv }()

	embeddedConfig := fakeEmbeddedConfig{}

	tempFile, err := ioutil.TempFile("", "config.yaml")
	defer os.Remove(tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile(tempFile.Name(), fakeYamlBytes, 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = ReadConfig(
        &struct{}{},
        AddPaths(tempFile.Name()),
		AddEmbedded("embedded", &embeddedConfig),
    )
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("embeddedConfig: %#+v", embeddedConfig)
	g.Expect(embeddedConfig.Name).To(gm.Equal("Steve"))
	g.Expect(embeddedConfig.Age).To(gm.Equal(27))
}

func TestReadConfigFileWithEnvironmentPrefixOverride(t *testing.T) {
	g := gm.NewGomegaWithT(t)

    testPrefix := "SECRET"
    lookupEnv = func(key string) (string, bool) {
        if key == (testPrefix + "_TEST_FOO") {
            return "29", true
        }

        return "", false
    }
    defer func() { lookupEnv = os.LookupEnv }()

	fakeConfig := fakeYamlConfigStruct{}

	tempFile, err := ioutil.TempFile("", "config.yaml")
	defer os.Remove(tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile(tempFile.Name(), fakeYamlBytes, 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = ReadConfig(
        &fakeConfig,
        AddPaths(tempFile.Name()),
        OverrideEnvVarPrefix(testPrefix),
    )
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("fakeConfig: %#+v", fakeConfig)
	g.Expect(fakeConfig.Test.Foo).To(gm.Equal(29))
	g.Expect(fakeConfig.Test.Bar).To(gm.Equal("baz"))
}
