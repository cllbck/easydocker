package view

import (
	"easydocker/docker"
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var (
	imagesWidget     *widgets.Paragraph
	containersWidget *widgets.Paragraph
)

func CreateView() *ui.Grid {
	imagesWidget = widgets.NewParagraph()
	imagesWidget.Title = "Images"

	containersWidget = widgets.NewParagraph()
	containersWidget.Title = "Containers"

	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	grid.Set(
		ui.NewRow(1.0/2,
			ui.NewCol(1.0/2, containersWidget),
			ui.NewCol(1.0/2, imagesWidget),
		),
		//ui.NewRow(1.0/2,
		//	ui.NewCol(1.0/2, imagesWidget),
		//	ui.NewCol(1.0/2, imagesWidget),
		//),
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

	containersWidget.Text = ""
	for _, container := range info.Containers {
		containersWidget.Text += fmt.Sprintf("%s - %s\n", container.ID[:12], container.Names[0])
	}
}
