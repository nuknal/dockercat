package gui

import (
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type cleanupPanel struct {
	*tview.Table
}

func newCleanupPanel(g *Gui) *cleanupPanel {
	p := &cleanupPanel{
		Table: tview.NewTable().SetSelectable(true, false).Select(0, 0).SetFixed(1, 1),
	}
	p.SetTitle(" Cleanup ").SetTitleAlign(tview.AlignLeft)
	p.SetBorder(true)
	p.SetBorderColor(CurrentTheme.Border)
	p.SetTitleColor(CurrentTheme.Title)
	p.setEntries(g)
	p.setKeybinding(g)
	return p
}

func (c *cleanupPanel) name() string {
	return "cleanup"
}

func (c *cleanupPanel) setKeybinding(g *Gui) {
	c.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		g.setGlobalKeybinding(event)
		switch event.Key() {
		case tcell.KeyCtrlR:
			c.setEntries(g)
		}

		switch event.Rune() {
		case 'a':
			g.systemPrune()
		case 'c':
			g.pruneContainersAction()
		case 'i':
			g.pruneImagesAction()
		case 'n':
			g.pruneNetworks()
		case 'v':
			g.pruneVolumes()
		}

		return event
	})
}

func (c *cleanupPanel) entries(g *Gui) {
	containers, _ := g.client.Containers(types.ContainerListOptions{All: true})
	stopped := 0
	for _, con := range containers {
		if con.State == "exited" {
			stopped++
		}
	}

	imgs, _ := g.client.Images(types.ImageListOptions{
		Filters: filters.NewArgs(filters.Arg("dangling", "true")),
	})
	dangling := len(imgs)
	var danglingSize int64
	for _, i := range imgs {
		danglingSize += i.Size
	}

	nets, _ := g.client.Networks(types.NetworkListOptions{})
	unusedN := 0
	for _, n := range nets {
		netDetail, _ := g.client.InspectNetwork(n.ID)
		if len(netDetail.Containers) == 0 {
			unusedN++
		}
	}

	du := g.getDiskUsage()

	c.Table.Clear()
	c.Table.SetSelectedStyle(CurrentTheme.SelectedFg, CurrentTheme.SelectedBg, 0)

	headers := []string{"  Action", "Targets", "Count", "Size"}
	for i, h := range headers {
		c.Table.SetCell(0, i, &tview.TableCell{
			Text:            h,
			NotSelectable:   true,
			Align:           tview.AlignLeft,
			Color:           CurrentTheme.Header,
			BackgroundColor: CurrentTheme.Bg,
			Attributes:      tcell.AttrBold,
		})
	}

	rows := [][]string{
		{"  c: prune containers", "stopped containers", fmt.Sprintf("%d", stopped), ""},
		{"  i: prune images", "dangling images", fmt.Sprintf("%d", dangling), du.imagesSize},
		{"  n: prune networks", "unused networks", fmt.Sprintf("%d", unusedN), ""},
		{"  v: prune volumes", "unused volumes", "", du.volumesSize},
		{"  a: system prune", "all unused (safe)", "", ""},
	}

	for i, r := range rows {
		for j, col := range r {
			expansion := 1
			if j == 0 {
				expansion = 3
			} else if j == 1 {
				expansion = 2
			}
			c.Table.SetCell(i+1, j, tview.NewTableCell(col).
				SetTextColor(CurrentTheme.CleanupItems).
				SetMaxWidth(1).
				SetExpansion(expansion))
		}
	}
}

func (c *cleanupPanel) setEntries(g *Gui) {
	c.entries(g)
}

func (c *cleanupPanel) updateEntries(g *Gui) {
	g.app.QueueUpdateDraw(func() {
		c.setEntries(g)
	})
}

func (c *cleanupPanel) focus(g *Gui) {
	c.SetSelectable(true, false)
	g.app.SetFocus(c)
}

func (c *cleanupPanel) unfocus() {
	c.SetSelectable(false, false)
}

func (c *cleanupPanel) setFilterWord(word string) {}
