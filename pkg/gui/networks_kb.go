package gui

import (
	"context"

	"github.com/nuknal/dockercat/pkg/common"
)

func (g *Gui) inspectNetwork() {
	network := g.selectedNetwork()

	inspect, err := g.client.InspectNetwork(network.ID)
	if err != nil {
		g.Log.Errorf("cannot inspect network %s", err)
		return
	}

	infoPanel := g.infoPanel()
	from := infoPanel.key
	to := "network-" + network.ID + "-detail"
	infoPanel.setKey(to)

	g.inspect(common.StructToJSON(inspect), from, to)
}

func (g *Gui) removeNetwork() {
	network := g.selectedNetwork()
	g.confirm("Do you want to remove network", "networks", func() {
		g.addTask("remove network", network.Name, func(ctx context.Context) error {
			if err := g.client.RemoveNetwork(network.ID); err != nil {
				g.errChan <- err
				return err
			}
			g.networkPanel().updateEntries(g)
			return nil
		})
	})
}
