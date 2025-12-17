package gui

import (
	"context"
)

func (g *Gui) systemPrune() {
	g.confirm("Do you want to prune all unused data", "cleanup", func() {
		g.addTask("system prune", "all unused", func(ctx context.Context) error {
			if err := g.client.PruneAll(); err != nil {
				g.errChan <- err
				return err
			}
			g.containerPanel().updateEntries(g)
			g.imagePanel().updateEntries(g)
			g.networkPanel().updateEntries(g)
			g.volumePanel().updateEntries(g)
			g.infoPanel().getDockerInfo(g)
			g.cleanupPanel().updateEntries(g)
			return nil
		})
	})
}

func (g *Gui) pruneContainersAction() {
	g.confirm("Do you want to prune stopped containers", "cleanup", func() {
		g.addTask("prune containers", "stopped containers", func(ctx context.Context) error {
			if err := g.client.PruneContainers(); err != nil {
				g.errChan <- err
				return err
			}
			g.containerPanel().updateEntries(g)
			g.cleanupPanel().updateEntries(g)
			return nil
		})
	})
}

func (g *Gui) pruneImagesAction() {
	g.confirm("Do you want to prune dangling images", "cleanup", func() {
		g.addTask("prune images", "dangling images", func(ctx context.Context) error {
			if err := g.client.PruneImages(); err != nil {
				g.errChan <- err
				return err
			}
			g.imagePanel().updateEntries(g)
			g.cleanupPanel().updateEntries(g)
			return nil
		})
	})
}
