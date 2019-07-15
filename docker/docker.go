package docker

import (
	docker "github.com/fsouza/go-dockerclient"
	"log"
)

type DockerInfo struct {
	Images     []docker.APIImages
	Containers []docker.APIContainers
}

var (
	client *docker.Client
	err    error
)

func InitDocker() {
	client, err = docker.NewClientFromEnv()
	if err != nil {
		log.Fatalf("failed create docker client: %v", err)
	}

}

func GetDockerInfo() DockerInfo {
	images, err := client.ListImages(docker.ListImagesOptions{All: true})
	if err != nil {
		log.Fatalf("failed received docker images: %v", err)
	}
	containers, err := client.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		log.Fatalf("failed received docker containers: %v", err)
	}
	return DockerInfo{
		images,
		containers,
	}
}
