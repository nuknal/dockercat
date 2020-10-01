package gui

import (
	"context"
	"fmt"

	"github.com/nuknal/dockercat/pkg/common"
	"github.com/rivo/tview"
)

func (g *Gui) inspectContainer() {
	container := g.selectedContainer()

	if container == nil {
		fmt.Println("container doesn't exist")
	}

	inspect, err := g.client.InspectContainer(container.ID)
	if err != nil {
		return
	}

	infoPanel := g.infoPanel()
	from := infoPanel.key
	to := "container-" + container.ID + "-detail"
	infoPanel.setKey(to)

	g.inspect(common.StructToJSON(inspect), from, to)
}

func (g *Gui) tailContainerLog() {
	container := g.selectedContainer()
	if container == nil {
		return
	}

	infoPanel := g.infoPanel()
	fromKey := infoPanel.key
	toKey := "container-" + container.ID + "-logs"
	infoPanel.setKey(toKey)
	text := tview.NewTextView()
	text.SetTitle("Logs").SetTitleAlign(tview.AlignLeft)
	text.SetBorder(true)
	infoPanel.switchItemTextView(text, fromKey, toKey)

	go infoPanel.containerLogs(g, container.ID)
}

func (g *Gui) containerStats() {
	container := g.selectedContainer()
	if container == nil {
		return
	}

	infoPanel := g.infoPanel()
	fromKey := infoPanel.key
	toKey := "container-" + container.ID + "-stats"
	infoPanel.setKey(toKey)

	text := tview.NewTextView()
	text.SetTitle("Stats").SetTitleAlign(tview.AlignLeft)
	text.SetBorder(true)
	infoPanel.switchItemTextView(text, fromKey, toKey)
	cstat, ok := g.containerPanel().containersStat[container.ID]
	if ok {
		_, _, w, _ := infoPanel.GetRect()
		content, _ := cstat.RenderStats(w)
		infoPanel.itemTextView.SetText(content)
	}

	go infoPanel.containerStats(g, container)
}

func (g *Gui) startContainer() {
	container := g.selectedContainer()
	g.addTask("start container", container.Name, func(ctx context.Context) error {
		if err := g.client.StartContainer(container.ID); err != nil {
			g.Log.Errorf("cannot start container %s", err)
			return err
		}

		g.containerPanel().updateEntries(g)
		return nil
	})
}

func (g *Gui) stopContainer() {
	container := g.selectedContainer()
	g.addTask("stop container", container.Name, func(ctx context.Context) error {
		if err := g.client.StopContainer(container.ID); err != nil {
			g.Log.Errorf("cannot stop container %s", err)
			return err
		}

		g.containerPanel().updateEntries(g)
		return nil
	})
}

func (g *Gui) removeContainer() {
	container := g.selectedContainer()
	g.addTask("remove container", container.Name, func(ctx context.Context) error {
		if err := g.client.RemoveContainer(container.ID); err != nil {
			g.Log.Errorf("cannot remove container %s", err)
			return err
		}
		g.containerPanel().updateEntries(g)
		return nil
	})

}
