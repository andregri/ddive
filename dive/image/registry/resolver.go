package registry

import (
	"archive/tar"
	"bytes"
	"context"
	"io"
	"log"

	"github.com/andregri/ddive/dive/image"
	"github.com/andregri/ddive/dive/image/docker"
	"github.com/containers/image/v5/copy"
	"github.com/containers/image/v5/signature"
	"github.com/containers/image/v5/transports/alltransports"
)

type engineResolver struct{}

func NewResolverFromRegistry() *engineResolver {
	return &engineResolver{}
}

func (r *engineResolver) Fetch(id string) (*image.Image, error) {
	srcImageName := "docker://quay.io/quay/redis:latest"
	srcRef, err := alltransports.ParseImageName(srcImageName)
	if err != nil {
		log.Printf("Invalid source name %s: %v\n", srcImageName, err)
		return nil, err
	}

	destImageName := "docker-archive:./redis_latest.tar"
	destRef, err := alltransports.ParseImageName(destImageName)
	if err != nil {
		log.Printf("Invalid destination name %s: %v\n", destImageName, err)
		return nil, err
	}

	policy := &signature.Policy{Default: []signature.PolicyRequirement{signature.NewPRInsecureAcceptAnything()}}
	policyCtx, err := signature.NewPolicyContext(policy)
	if err != nil {
		log.Println("invalid policy")
		return nil, err
	}

	manifestBytes, err := copy.Image(context.Background(), policyCtx, destRef, srcRef, nil)
	if err != nil {
		log.Printf("failed to copy image from %s to %s\n", srcImageName, destImageName)
		return nil, err
	}

	log.Println(string(manifestBytes))

	var buf bytes.Buffer
	reader := tar.NewReader(&buf)
	readCloser := io.NopCloser(reader)

	img, err := docker.NewImageArchive(readCloser)
	if err != nil {
		return nil, err
	}
	return img.ToImage()
}

func (r *engineResolver) Build(args []string) (*image.Image, error) {
	panic("Not available")
}
