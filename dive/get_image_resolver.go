package dive

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/andregri/ddive/dive/image"
	"github.com/andregri/ddive/dive/image/docker"
	"github.com/andregri/ddive/dive/image/podman"
	"github.com/andregri/ddive/dive/image/registry"
)

const (
	SourceUnknown ImageSource = iota
	SourceDockerEngine
	SourcePodmanEngine
	SourceDockerArchive
	SourceRegistry
)

type ImageSource int

var ImageSources = []string{
	SourceDockerEngine.String(),
	SourcePodmanEngine.String(),
	SourceDockerArchive.String(),
	SourceRegistry.String(),
}

func (r ImageSource) String() string {
	return [...]string{"unknown", "docker", "podman", "docker-archive", "registry"}[r]
}

func ParseImageSource(r string) ImageSource {
	switch r {
	case SourceDockerEngine.String():
		return SourceDockerEngine
	case SourcePodmanEngine.String():
		return SourcePodmanEngine
	case SourceDockerArchive.String():
		return SourceDockerArchive
	case "docker-tar":
		return SourceDockerArchive
	case SourceRegistry.String():
		return SourceRegistry
	default:
		return SourceUnknown
	}
}

func DeriveImageSource(image string) (ImageSource, string) {
	u, err := url.Parse(image)
	if err != nil {
		return SourceUnknown, ""
	}

	imageSource := strings.TrimPrefix(image, u.Scheme+"://")

	switch u.Scheme {
	case SourceDockerEngine.String():
		return SourceDockerEngine, imageSource
	case SourcePodmanEngine.String():
		return SourcePodmanEngine, imageSource
	case SourceDockerArchive.String():
		return SourceDockerArchive, imageSource
	case "docker-tar":
		return SourceDockerArchive, imageSource
	case SourceRegistry.String():
		return SourceRegistry, imageSource
	}
	return SourceUnknown, ""
}

func GetImageResolver(r ImageSource) (image.Resolver, error) {
	switch r {
	case SourceDockerEngine:
		return docker.NewResolverFromEngine(), nil
	case SourcePodmanEngine:
		return podman.NewResolverFromEngine(), nil
	case SourceDockerArchive:
		return docker.NewResolverFromArchive(), nil
	case SourceRegistry:
		return registry.NewResolverFromRegistry(), nil
	}

	return nil, fmt.Errorf("unable to determine image resolver")
}
