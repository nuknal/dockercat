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

func (g *Gui) pullImage() {
	input := tview.NewInputField().SetLabel("Image: ")
	form := tview.NewForm().AddFormItem(input).
		AddButton("Pull", func() {
			ref := input.GetText()
			g.pages.RemovePage("image-pull").ShowPage("main")
			g.addTask("pull image", ref, func(ctx context.Context) error {
				if err := g.client.PullImage(ref); err != nil {
					g.errChan <- err
					return err
				}
				g.imagePanel().updateEntries(g)
				return nil
			})
		}).
		AddButton("Cancel", func() {
			g.pages.RemovePage("image-pull").ShowPage("main")
		})
	form.SetBorder(true).SetTitle("Pull Image").SetTitleAlign(tview.AlignLeft)
	grid := tview.NewGrid().SetColumns(0, 60, 0).SetRows(0, 7, 0).AddItem(form, 1, 1, 1, 1, 0, 0, true)
	g.pages.AddAndSwitchToPage("image-pull", grid, true).ShowPage("main")
}

func (g *Gui) tagImage() {
	image := g.selectedImage()
	if image == nil {
		return
	}
	input := tview.NewInputField().SetLabel("New Tag: ")
	form := tview.NewForm().AddFormItem(input).
		AddButton("Tag", func() {
			target := input.GetText()
			g.pages.RemovePage("image-tag").ShowPage("main")
			g.addTask("tag image", target, func(ctx context.Context) error {
				if err := g.client.TagImage(image.Repo+":"+image.Tag, target); err != nil {
					g.errChan <- err
					return err
				}
				g.imagePanel().updateEntries(g)
				return nil
			})
		}).
		AddButton("Cancel", func() {
			g.pages.RemovePage("image-tag").ShowPage("main")
		})
	form.SetBorder(true).SetTitle("Tag Image").SetTitleAlign(tview.AlignLeft)
	grid := tview.NewGrid().SetColumns(0, 60, 0).SetRows(0, 7, 0).AddItem(form, 1, 1, 1, 1, 0, 0, true)
	g.pages.AddAndSwitchToPage("image-tag", grid, true).ShowPage("main")
}

func (g *Gui) pushImage() {
	image := g.selectedImage()
	if image == nil {
		return
	}
	refInput := tview.NewInputField().SetLabel("Ref: ").SetText(image.Repo + ":" + image.Tag)
	userInput := tview.NewInputField().SetLabel("Username: ")
	passInput := tview.NewInputField().SetLabel("Password: ").SetMaskCharacter('*')
	form := tview.NewForm().AddFormItem(refInput).AddFormItem(userInput).AddFormItem(passInput).
		AddButton("Push", func() {
			ref := refInput.GetText()
			user := userInput.GetText()
			pass := passInput.GetText()
			g.pages.RemovePage("image-push").ShowPage("main")
			g.addTask("push image", ref, func(ctx context.Context) error {
				if err := g.client.PushImage(ref, user, pass); err != nil {
					g.errChan <- err
					return err
				}
				g.imagePanel().updateEntries(g)
				return nil
			})
		}).
		AddButton("Cancel", func() {
			g.pages.RemovePage("image-push").ShowPage("main")
		})
	form.SetBorder(true).SetTitle("Push Image").SetTitleAlign(tview.AlignLeft)
	grid := tview.NewGrid().SetColumns(0, 60, 0).SetRows(0, 10, 0).AddItem(form, 1, 1, 1, 1, 0, 0, true)
	g.pages.AddAndSwitchToPage("image-push", grid, true).ShowPage("main")
}
