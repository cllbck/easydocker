package main

import (
	"easydocker/docker"
	"easydocker/view"
	ui "github.com/gizak/termui/v3"
	"log"
	"time"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	docker.InitDocker()

	newView := view.CreateView()
	dockerInfo := docker.GetDockerInfo()
	view.RefreshView(dockerInfo)
	ui.Render(newView)

	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second).C
	view.RefreshView(dockerInfo)
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			}
		case <-ticker:
			dockerInfo := docker.GetDockerInfo()
			view.RefreshView(dockerInfo)
			ui.Render(newView)
		}
	}
}
