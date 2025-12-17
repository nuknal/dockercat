package gui

import (
	"strings"
	"time"

	"github.com/gdamore/tcell"
	"github.com/nuknal/dockercat/pkg/common"
	"github.com/rivo/tview"
)

var replacer = strings.NewReplacer("T", " ", "Z", "")

type volume struct {
	Name       string
	MountPoint string
	Driver     string
	Created    string
}

type volumePanel struct {
	*tview.Table
	filterWord string
}

func newVolumePanel(g *Gui) *volumePanel {
	volumes := &volumePanel{
		Table: tview.NewTable().SetSelectable(true, false).Select(0, 0),
	}

	volumes.SetTitle("volumes").SetTitleAlign(tview.AlignLeft)
	volumes.SetBorder(true)
	volumes.setEntries(g)
	volumes.setKeybinding(g)
	return volumes
}

func (v *volumePanel) name() string {
	return "volumes"
}

func (v *volumePanel) setKeybinding(g *Gui) {
	v.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		g.setGlobalKeybinding(event)
		switch event.Key() {
		case tcell.KeyEnter:
			g.inspectVolume()
		case tcell.KeyCtrlR:
			v.setEntries(g)
		case tcell.KeyCtrlD:
			g.pruneVolumes()
		}

		switch event.Rune() {
		case 'd':
			g.removeVolume()
		}

		return event
	})
}

func (v *volumePanel) entries(g *Gui) {
	volumes, err := g.client.Volumes()
	if err != nil {
		return
	}

	keys := make([]string, 0, len(volumes))
	tmpMap := make(map[string]*volume)

	for _, vo := range volumes {
		if strings.Index(vo.Name, v.filterWord) == -1 {
			continue
		}

		tmpMap[vo.Name] = &volume{
			Name:       vo.Name,
			MountPoint: vo.Mountpoint,
			Driver:     vo.Driver,
			Created:    replacer.Replace(vo.CreatedAt),
		}

		keys = append(keys, vo.Name)
	}

	g.resources.volumes = make([]*volume, 0)
	for _, key := range common.SortKeys(keys) {
		g.resources.volumes = append(g.resources.volumes, tmpMap[key])
	}
}

func (v *volumePanel) setEntries(g *Gui) {
	v.entries(g)
	table := v.Clear()

	headers := []string{
		"Name",
		"MountPoint",
	}

	for i, header := range headers {
		table.SetCell(0, i, &tview.TableCell{
			Text:            header,
			NotSelectable:   true,
			Align:           tview.AlignLeft,
			Color:           tcell.ColorWhite,
			BackgroundColor: tcell.ColorDefault,
			Attributes:      tcell.AttrBold,
		})
	}

	for i, network := range g.resources.volumes {
		table.SetCell(i+1, 0, tview.NewTableCell(network.Name).
			SetTextColor(tcell.ColorLightPink).
			SetMaxWidth(1).
			SetExpansion(1))
		table.SetCell(i+1, 1, tview.NewTableCell(network.MountPoint).
			SetTextColor(tcell.ColorLightPink).
			SetMaxWidth(1).
			SetExpansion(1))
	}
}

func (v *volumePanel) focus(g *Gui) {
	v.SetSelectable(true, false)
	g.app.SetFocus(v)
}

func (v *volumePanel) unfocus() {
	v.SetSelectable(false, false)
}

func (v *volumePanel) updateEntries(g *Gui) {
	g.app.QueueUpdateDraw(func() {
		v.setEntries(g)
	})
}

func (v *volumePanel) setFilterWord(word string) {
	v.filterWord = word
}

func (v *volumePanel) monitoringVolumes(g *Gui) {
	ticker := time.NewTicker(g.refreshInterval)

LOOP:
	for {
		select {
		case <-ticker.C:
			v.updateEntries(g)
		case <-g.stopChans["volume"]:
			ticker.Stop()
			break LOOP
		}
	}
}
