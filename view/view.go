package view

import (
	"easydocker/docker"
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"log"
)

var (
	imagesWidget     *widgets.Paragraph
	containersWidget *Table
)

func CreateView() *ui.Grid {
	imagesWidget = widgets.NewParagraph()
	imagesWidget.Title = "Images"

	containersWidget = NewCustomTable()
	containersWidget.Title = "Containers"
	log.Print(containersWidget.Block.BorderLeft)
	containersWidget.RowSeparator = false

	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	grid.Set(
		ui.NewRow(1.0/2,
			ui.NewCol(1.0, containersWidget),
		),
		ui.NewRow(1.0/2,
			ui.NewCol(1.0/2, imagesWidget),
			//ui.NewCol(1.0/2, imagesWidget),
		),
	)
	return grid
}
func RefreshView(info docker.DockerInfo) {
	imagesWidget.Text = ""
	var totalSize int64
	for _, image := range info.Images {
		imagesWidget.Text += fmt.Sprintf("%s - %s\n", image.RepoTags[0], convert(image.Size))
		totalSize += image.Size
	}
	imagesWidget.Title = fmt.Sprintf("Images (total size: %s)", convert(totalSize))

	containersWidget.Rows = [][]string{
		{"CONTAINER ID", "IMAGE", "COMMAND", "CREATED", "STATUS", "PORTS", "NAMES"},
	}
	for _, container := range info.Containers {
		portsInfo := ""
		if len(container.Ports) > 0 {
			portsInfo = fmt.Sprintf("%d:%d - %s", container.Ports[0].PublicPort, container.Ports[0].PrivatePort, container.Ports[0].Type)
		}
		containersWidget.Rows = append(
			containersWidget.Rows,
			[]string{
				container.ID[:12],
				container.Image,
				container.Command,
				timeElapsed(container.Created, false),
				container.Status,
				portsInfo,
				container.Names[0],
			})
	}
}
