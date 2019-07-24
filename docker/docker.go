package docker

import (
	docker "github.com/fsouza/go-dockerclient"
	"log"
)

type DockerInfo struct {
	Images     []docker.APIImages
	Containers []docker.APIContainers
	Stats      []Stats
}

type Stats struct {
	Container docker.APIContainers
	Stats     docker.Stats
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

	var stats []Stats

	containersForStats, err := client.ListContainers(docker.ListContainersOptions{All: false})
	if err != nil {
		log.Fatalf("failed received docker containers: %v", err)
	}

	for _, container := range containersForStats {
		stats = append(stats, Stats{
			container,
			*getContainerStat(container.ID),
		})
	}

	return DockerInfo{
		images,
		containers,
		stats,
	}
}

func getContainerStat(id string) *docker.Stats {
	stats := make(chan *docker.Stats)

	go func() {
		if err := client.Stats(docker.StatsOptions{ID: id, Stats: stats, Stream: false}); err != nil {
			log.Fatalf("failed received  container id-%s stats: %v", id, err)
		}
	}()

	stat := <-stats
	return stat
}
