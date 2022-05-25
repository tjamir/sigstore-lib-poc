package main

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/sigstore/cosign/cmd/cosign/cli/fulcio"
	"github.com/sigstore/cosign/pkg/cosign"
	rekor "github.com/sigstore/rekor/pkg/generated/client"
)

func main() {
	args := os.Args

	if len(args) == 1 {
		fmt.Errorf("At least one image is required")
	}
	imageName := args[1]
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ref, err := name.ParseReference(imageName)
	if err != nil {
		panic("Error parsing image reference")
	}

	co := &cosign.CheckOpts{}

	// Set the rekor client
	rekorURL := url.URL{
		Scheme: rekor.DefaultSchemes[0],
		Host:   rekor.DefaultHost,
		Path:   rekor.DefaultBasePath,
	}
	co.RekorClient = rekor.NewHTTPClientWithConfig(nil, rekor.DefaultTransportConfig().WithBasePath(rekorURL.Path).WithHost(rekorURL.Host))

	co.RootCerts = fulcio.GetRoots()

	cosign.VerifyImageSignatures(ctx, ref, co)
}
