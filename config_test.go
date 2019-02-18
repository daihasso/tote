package tote

import (
    "testing"
	"io/ioutil"
	"os"

    "github.com/daihasso/peechee"
	gm "github.com/onsi/gomega"
)

type fakeYamlConfigStruct struct {
    Test struct{
        Foo int
        Bar string
    }
}

type fakeEmbeddedConfig struct {
	Name string
	Age int
}

var fakeYamlBytes = []byte(`
test:
  foo: 1
  bar: baz

embedded:
   name: Joe
   age: 27
`,
)[1:]

func TestReadConfigFile(t *testing.T) {
	g := gm.NewGomegaWithT(t)

	tempFile, err := ioutil.TempFile("", "config.yaml")
	defer os.Remove(tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile(tempFile.Name(), fakeYamlBytes, 0644)
	if err != nil {
		t.Fatal(err)
	}

	configBytes, err := readConfigFile(tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("configBytes:\n%s", string(configBytes))
	t.Logf("fakeYamlBytes:\n%s", string(fakeYamlBytes))
	g.Expect(string(configBytes)).To(gm.Equal(string(fakeYamlBytes)))
}

func TestReadConfigData(t *testing.T) {
	g := gm.NewGomegaWithT(t)

    fakeConfig := fakeYamlConfigStruct{}
    err := unmarshalConfigBytes(fakeYamlBytes, &fakeConfig)
    if err != nil {
        t.Fatal(err)
    }

	t.Logf("fakeConfig: %#+v", fakeConfig)
	g.Expect(fakeConfig.Test.Foo).To(gm.Equal(1))
	g.Expect(fakeConfig.Test.Bar).To(gm.Equal("baz"))
}

func TestReadEmbeddedConfig(t *testing.T) {
	g := gm.NewGomegaWithT(t)

	embeddedConfig := fakeEmbeddedConfig{}
	err := readEmbeddedConfig(fakeYamlBytes, "embedded", &embeddedConfig)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("embeddedConfig: %#+v", embeddedConfig)
	g.Expect(embeddedConfig.Name).To(gm.Equal("Joe"))
	g.Expect(embeddedConfig.Age).To(gm.Equal(27))
}

func TestReadConfigBasicFile(t *testing.T) {
	g := gm.NewGomegaWithT(t)

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

    pathReader := peechee.NewPathReader(peechee.WithFilesystem())

	err = ReadConfig(
        &fakeConfig, AddPaths(tempFile.Name()), WithPathReader(pathReader),
    )
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("fakeConfig: %#+v", fakeConfig)
	g.Expect(fakeConfig.Test.Foo).To(gm.Equal(1))
	g.Expect(fakeConfig.Test.Bar).To(gm.Equal("baz"))
}

func TestReadConfigWithEmbedded(t *testing.T) {
	g := gm.NewGomegaWithT(t)

    fakeConfig := fakeYamlConfigStruct{}
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
		&fakeConfig,
		AddPaths(tempFile.Name()),
		AddEmbedded("embedded", &embeddedConfig),
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("fakeConfig: %#+v", fakeConfig)
	t.Logf("embeddedConfig: %#+v", embeddedConfig)
	g.Expect(fakeConfig.Test.Foo).To(gm.Equal(1))
	g.Expect(fakeConfig.Test.Bar).To(gm.Equal("baz"))
	g.Expect(embeddedConfig.Name).To(gm.Equal("Joe"))
	g.Expect(embeddedConfig.Age).To(gm.Equal(27))
}

func TestReadConfigFailureNotPointer(t *testing.T) {
	g := gm.NewGomegaWithT(t)

    fakeConfig := fakeYamlConfigStruct{}

	err := ReadConfig(fakeConfig)

    t.Log(err)
	t.Logf("fakeConfig: %#+v", fakeConfig)
    g.Expect(err.Error()).Should(gm.MatchRegexp(
        `^ReadConfig requires a pointer to your config struct, .*`,
    ))
}

func TestReadEmbeddedConfigFailureNotPointer(t *testing.T) {
	g := gm.NewGomegaWithT(t)

    fakeConfig := fakeYamlConfigStruct{}
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
		&fakeConfig,
		AddPaths(tempFile.Name()),
		AddEmbedded("embedded", embeddedConfig),
	)

	t.Logf("fakeConfig: %#+v", fakeConfig)
	t.Logf("embeddedConfig: %#+v", embeddedConfig)
    g.Expect(err.Error()).Should(gm.MatchRegexp(
        `^ReadConfig requires a pointer to your embedded config\(s\), .*`,
    ))
}
