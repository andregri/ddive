package registry

import (
	"context"
	"io"
	"log"
	"os"

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

	if _, err := os.Stat("redis_latest.tar"); err != nil {
		// if file does not exists
		manifestBytes, err := copy.Image(context.Background(), policyCtx, destRef, srcRef, nil)
		if err != nil {
			log.Printf("failed to copy image from %s to %s\n", srcImageName, destImageName)
			return nil, err
		}
		log.Println(string(manifestBytes))
	}

	file, err := os.Open("redis_latest.tar")
	if err != nil {
		log.Println("Error while reading tar")
		return nil, err
	}

	readCloser := io.NopCloser(file)

	img, err := docker.NewImageArchive(readCloser)
	if err != nil {
		return nil, err
	}
	return img.ToImage()
}

func (r *engineResolver) Build(args []string) (*image.Image, error) {
	panic("Not available")
}
