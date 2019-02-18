// +build !nos3

package tote

import (
    "github.com/aws/aws-sdk-go/service/s3/s3iface"
    "github.com/daihasso/peechee"
)

// WithS3Client adds an S3Client to the PathReader so that it can correctly
// grab data for S3 paths.
func WithS3Client(s3Client s3iface.S3API) Option {
    return func(opts *options) {
        opts.pathReader.AddOption(peechee.WithS3(s3Client))
    }
}
