package view

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func CreateView() *ui.Grid {
	imagesWidget := widgets.NewParagraph()
	imagesWidget.Text = "Text"
	imagesWidget.Title = "Images"

	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	grid.Set(
		ui.NewRow(1.0/2,
			ui.NewCol(1.0/2, imagesWidget),
			ui.NewCol(1.0/2, imagesWidget),
		),
		ui.NewRow(1.0/2,
			ui.NewCol(1.0/2, imagesWidget),
			ui.NewCol(1.0/2, imagesWidget),
		),
	)
	return grid
}
func RefreshView() {

}
