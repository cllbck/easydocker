package docker

import (
	docker "github.com/fsouza/go-dockerclient"
	"log"
)

type DockerInfo struct {
	images []docker.APIImages
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
	imgs, err := client.ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		log.Fatalf("failed received docker images: %v", err)
	}
	return DockerInfo{
		imgs,
	}
}
