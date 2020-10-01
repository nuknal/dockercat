package gui

import (
	"context"
	"fmt"
	"strings"

	"github.com/nuknal/dockercat/pkg/common"
	"github.com/rivo/tview"
)

func (g *Gui) inspectImage() {
	image := g.selectedImage()

	inspect, err := g.client.InspectImage(image.ID)
	if err != nil {
		g.Log.Errorf("cannot inspect image %s", err)
		return
	}

	infoPanel := g.infoPanel()
	from := infoPanel.key
	to := "image-" + image.ID + "-detail"
	infoPanel.setKey(to)

	g.inspect(common.StructToJSON(inspect), from, to)
}

func (g *Gui) removeUnusedImage() {
	g.confirm("Do you want to remove all unused and untagged images?", "images", func() {
		g.addTask("remove unused images", "unused images", func(ctx context.Context) error {
			if err := g.client.RemoveDanglingImages(); err != nil {
				g.errChan <- err
				return err
			}
			g.imagePanel().updateEntries(g)
			return nil
		})
	})

}

func (g *Gui) removeImage() {
	image := g.selectedImage()
	g.confirm("Do you want to remove the images", "images", func() {
		g.addTask("remove image", image.Repo, func(ctx context.Context) error {
			if err := g.client.RemoveImage(image.ID); err != nil {
				g.errChan <- err
				return err
			}
			g.imagePanel().updateEntries(g)
			return nil
		})
	})

}

func (g *Gui) historyImage() {
	image := g.selectedImage()
	hist, err := g.client.HistoryImage(image.ID)
	if err != nil {
		g.errChan <- err
		return
	}

	var text string
	strs := make([]string, len(hist))
	for i, h := range hist {
		createdAt := common.ParseDateToString(h.Created)
		size := common.ParseSizeToString(h.Size)
		strs[len(hist)-1-i] = fmt.Sprintf("%s\n%s\n%s", h.CreatedBy, createdAt, size)
	}
	text = strings.Join(strs, "\n-------------------------\n")

	infoPanel := g.infoPanel()
	fromKey := infoPanel.key
	toKey := "images-" + image.ID + "-history"
	infoPanel.setKey(toKey)
	tv := tview.NewTextView()
	tv.SetTitle("Image Layers").SetTitleAlign(tview.AlignLeft)
	tv.SetBorder(true)
	tv.SetText(text)
	infoPanel.switchItemTextView(tv, fromKey, toKey)

}
