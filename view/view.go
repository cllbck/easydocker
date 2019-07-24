package view

import (
	"easydocker/docker"
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var (
	imagesWidget          *widgets.Paragraph
	containersWidget      *Table
	containersStatsWidget *Table
)

func CreateView() *ui.Grid {
	imagesWidget = widgets.NewParagraph()
	imagesWidget.Title = "Images"

	containersWidget = NewCustomTable()
	containersWidget.Title = "Containers"
	containersWidget.RowSeparator = false

	containersStatsWidget = NewCustomTable()
	containersStatsWidget.Title = "Stats"
	containersStatsWidget.RowSeparator = false

	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	grid.Set(
		ui.NewRow(1.0/2,
			ui.NewCol(1.0, containersStatsWidget),
		),
		ui.NewRow(1.0/2,
			ui.NewCol(3.0/10, imagesWidget),
			ui.NewCol(7.0/10, containersWidget),
		),
	)
	return grid
}
func RefreshView(info docker.DockerInfo) {
	imagesWidget.Text = ""
	var totalSize int64
	for _, image := range info.Images {
		imagesWidget.Text += fmt.Sprintf("%s - %s\n", image.RepoTags[0], convert(float64(image.Size)))
		totalSize += image.Size
	}
	imagesWidget.Title = fmt.Sprintf("Images (total size: %s)", convert(float64(totalSize)))

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
	containersStatsWidget.Rows = [][]string{
		{"CONTAINER ID", "NAME", "CPU %", "MEM USAGE / LIMIT", "MEM %", "NET I/O", "BLOCK I/O", "PIDS"},
	}
	for _, stat := range info.Stats {
		blkRead, blkWrite := calculateBlockIO(stat.Stats)
		rx, tx := calculateNetwork(stat.Stats)
		containersStatsWidget.Rows = append(
			containersStatsWidget.Rows,
			[]string{
				stat.Container.ID[:12],
				stat.Container.Names[0],
				fmt.Sprintf("%.2f", calculateCpuPercent(stat.Stats)),
				fmt.Sprintf("%s / %s", convert(calculateMem(stat.Stats)), convert(calculateMemoryLimit(stat.Stats))),
				fmt.Sprintf("%.2f", calculateMemoryPercentage(stat.Stats)),
				fmt.Sprintf("%s / %s", convert(rx), convert(tx)),
				fmt.Sprintf("%s / %s", convert(blkRead), convert(blkWrite)),
				fmt.Sprintf("%d", stat.Stats.PidsStats.Current),
			})
	}
}
