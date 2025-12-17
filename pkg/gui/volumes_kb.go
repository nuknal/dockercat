package gui

import (
	"context"

	"github.com/nuknal/dockercat/pkg/common"
)

func (g *Gui) inspectVolume() {
	volume := g.selectedVolume()

	inspect, err := g.client.InspectVolume(volume.Name)
	if err != nil {
		g.Log.Errorf("cannot inspect volume %s", err)
		return
	}

	infoPanel := g.infoPanel()
	from := infoPanel.key
	to := "volume-" + volume.Name + "-detail"
	infoPanel.setKey(to)

	g.inspect(common.StructToJSON(inspect), from, to)
}

func (g *Gui) removeVolume() {
	volume := g.selectedVolume()
	g.confirm("Do you want to remove volume", "volumes", func() {
		g.addTask("remove volume", volume.Name, func(ctx context.Context) error {
			if err := g.client.RemoveVolume(volume.Name); err != nil {
				g.errChan <- err
				return err
			}
			g.volumePanel().updateEntries(g)
			return nil
		})
	})
}

func (g *Gui) pruneVolumes() {
	g.confirm("Do you want to remove unused volumes", "volumes", func() {
		g.addTask("prune volumes", "unused volumes", func(ctx context.Context) error {
			if err := g.client.PruneVolumes(); err != nil {
				g.errChan <- err
				return err
			}
			g.volumePanel().updateEntries(g)
			return nil
		})
	})
}
