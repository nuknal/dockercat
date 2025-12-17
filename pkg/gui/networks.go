package gui

import (
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/gdamore/tcell"
	"github.com/nuknal/dockercat/pkg/common"
	"github.com/rivo/tview"
)

type network struct {
	ID         string
	Name       string
	Driver     string
	Scope      string
	Containers string
}

type networkPanel struct {
	*tview.Table
	filterWord string
}

func newNetworkPanel(g *Gui) *networkPanel {
	networks := &networkPanel{
		Table: tview.NewTable().SetSelectable(true, false),
	}

	networks.SetTitle("networks").SetTitleAlign(tview.AlignLeft)
	networks.SetBorder(true)
	networks.setEntries(g)
	networks.setKeybinding(g)
	return networks
}

func (n *networkPanel) name() string {
	return "networks"
}

func (n *networkPanel) setKeybinding(g *Gui) {
	n.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		g.setGlobalKeybinding(event)
		switch event.Key() {
		case tcell.KeyEnter:
			g.inspectNetwork()
		case tcell.KeyCtrlR:
			n.setEntries(g)
		case tcell.KeyCtrlD:
			g.pruneNetworks()
		}

		switch event.Rune() {
		case 'd':
			g.removeNetwork()
		}

		return event
	})
}

func (n *networkPanel) entries(g *Gui) {
	networks, err := g.client.Networks(types.NetworkListOptions{})
	if err != nil {
		return
	}

	keys := make([]string, 0, len(networks))
	tmpMap := make(map[string]*network)

	for _, net := range networks {
		if strings.Index(net.Name, n.filterWord) == -1 {
			continue
		}

		var containers string

		net, err := g.client.InspectNetwork(net.ID)
		if err != nil {
			continue
		}

		for _, endpoint := range net.Containers {
			containers += fmt.Sprintf("%s ", endpoint.Name)
		}

		tmpMap[net.ID[:12]] = &network{
			ID:         net.ID,
			Name:       net.Name,
			Driver:     net.Driver,
			Scope:      net.Scope,
			Containers: containers,
		}

		keys = append(keys, net.ID[:12])

	}

	g.resources.networks = make([]*network, 0)

	for _, key := range common.SortKeys(keys) {
		g.resources.networks = append(g.resources.networks, tmpMap[key])
	}
}

func (n *networkPanel) setEntries(g *Gui) {
	n.entries(g)
	table := n.Clear()

	headers := []string{
		"Name",
		"Containers",
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

	for i, network := range g.resources.networks {
		table.SetCell(i+1, 0, tview.NewTableCell(network.Name).
			SetTextColor(tcell.ColorLightSkyBlue).
			SetMaxWidth(1).
			SetExpansion(1))
		table.SetCell(i+1, 1, tview.NewTableCell(network.Containers).
			SetTextColor(tcell.ColorLightSkyBlue).
			SetMaxWidth(1).
			SetExpansion(1))
	}
}

func (n *networkPanel) focus(g *Gui) {
	n.SetSelectable(true, false)
	g.app.SetFocus(n)
}

func (n *networkPanel) unfocus() {
	n.SetSelectable(false, false)
}

func (n *networkPanel) updateEntries(g *Gui) {
	g.app.QueueUpdateDraw(func() {
		n.setEntries(g)
	})
}

func (n *networkPanel) setFilterWord(word string) {
	n.filterWord = word
}

func (n *networkPanel) monitoringNetworks(g *Gui) {
	ticker := time.NewTicker(g.refreshInterval)

LOOP:
	for {
		select {
		case <-ticker.C:
			n.updateEntries(g)
		case <-g.stopChans["network"]:
			ticker.Stop()
			break LOOP
		}
	}
}
