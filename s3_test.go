package tote

import (
	"bytes"
    "testing"
	"io/ioutil"
	"net/http"

    logging "github.com/daihasso/slogging"
    "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/aws/aws-sdk-go/service/s3/s3iface"
	gm "github.com/onsi/gomega"
)

var testBucket = "foo-config"
var fakeRequestId = "fake"
var reqDone bool

type testWriter struct {
	t *testing.T
}

func (tw testWriter) Write(p []byte) (n int, err error) {
	tw.t.Log(string(p))
	return len(p), nil
}

type fakeS3 struct {
	s3iface.S3API
}

func (self fakeS3) GetObjectWithContext(
	aws.Context, *s3.GetObjectInput, ...request.Option,
) (*s3.GetObjectOutput, error) {
	if reqDone {
		return nil, awserr.NewRequestFailure(
			nil, http.StatusRequestedRangeNotSatisfiable, fakeRequestId,
		)
	}
	reqDone = true
	output := new(s3.GetObjectOutput)
	output.Body = ioutil.NopCloser(bytes.NewReader(fakeYamlBytes))
	return output, nil
}

func TestReadConfigBasicS3(t *testing.T) {
	g := gm.NewGomegaWithT(t)

    err := logging.GetRootLogger().SetLogLevel(logging.DEBUG)
    g.Expect(err).ToNot(gm.HaveOccurred())
    logging.GetRootLogger().SetWriters(testWriter{t})

    fakeConfig := fakeYamlConfigStruct{}

	s3c := fakeS3{}

	reqDone = false
	err = ReadConfig(
		&fakeConfig,
        AddPaths("s3://" + testBucket + "/test/config.yaml"),
        WithS3Client(s3c),
	)
    g.Expect(err).To(gm.BeNil())

	t.Logf("fakeConfig: %#+v", fakeConfig)
	g.Expect(fakeConfig.Test.Foo).To(gm.Equal(1))
	g.Expect(fakeConfig.Test.Bar).To(gm.Equal("baz"))
}
