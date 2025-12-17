package gui

import (
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/gdamore/tcell"
	"github.com/nuknal/dockercat/pkg/common"
	"github.com/rivo/tview"
)

type image struct {
	ID      string
	Repo    string
	Tag     string
	Created string
	Size    string
}

type imagePanel struct {
	*tview.Table
	filterWord string
}

func newImagePanel(g *Gui) *imagePanel {
	images := &imagePanel{
		Table: tview.NewTable().SetSelectable(true, false).Select(0, 0).SetFixed(1, 1),
	}

	images.SetTitle(" Images ").SetTitleAlign(tview.AlignLeft)
	images.SetBorder(true)
	images.SetBorderColor(CurrentTheme.Border)
	images.SetTitleColor(CurrentTheme.Title)
	images.setEntries(g)
	images.setKeybinding(g)
	return images
}

func (i *imagePanel) name() string {
	return "images"
}

func (i *imagePanel) setKeybinding(g *Gui) {
	i.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		g.setGlobalKeybinding(event)
		switch event.Key() {
		case tcell.KeyEnter:
			g.inspectImage()
		case tcell.KeyCtrlD:
			g.removeUnusedImage()
		case tcell.KeyCtrlR:
			i.setEntries(g)
		}

		switch event.Rune() {
		case 's':
			g.historyImage()
		case 'd':
			g.removeImage()
		case 'p':
			g.pullImage()
		case 't':
			g.tagImage()
		case 'P':
			g.pushImage()
		}

		return event
	})
}

func (i *imagePanel) entries(g *Gui) {
	images, err := g.client.Images(types.ImageListOptions{})
	if err != nil {
		return
	}

	g.resources.images = make([]*image, 0)

	for _, imgInfo := range images {
		for _, repoTag := range imgInfo.RepoTags {
			repo, tag := parseRepoTag(repoTag)
			if strings.Index(repo, i.filterWord) == -1 {
				continue
			}

			g.resources.images = append(g.resources.images, &image{
				ID:      imgInfo.ID[7:19],
				Repo:    repo,
				Tag:     tag,
				Created: common.ParseDateToString(imgInfo.Created),
				Size:    common.ParseSizeToString(imgInfo.Size),
			})
		}
	}
}

func (i *imagePanel) setEntries(g *Gui) {
	i.entries(g)
	table := i.Clear()

	table.SetSelectedStyle(CurrentTheme.SelectedFg, CurrentTheme.SelectedBg, 0)

	headers := []string{
		"  Repo",
		"Tag",
		"Size",
	}

	for i, header := range headers {
		table.SetCell(0, i, &tview.TableCell{
			Text:            header,
			NotSelectable:   true,
			Align:           tview.AlignLeft,
			Color:           CurrentTheme.Header,
			BackgroundColor: CurrentTheme.Bg,
			Attributes:      tcell.AttrBold,
		})
	}

	for i, image := range g.resources.images {
		table.SetCell(i+1, 0, tview.NewTableCell("  "+image.Repo).
			SetTextColor(CurrentTheme.Images).
			SetMaxWidth(1).
			SetExpansion(4))

		table.SetCell(i+1, 1, tview.NewTableCell(image.Tag).
			SetTextColor(CurrentTheme.Images).
			SetMaxWidth(1).
			SetExpansion(2))
		table.SetCell(i+1, 2, tview.NewTableCell(image.Size).
			SetTextColor(CurrentTheme.Images).
			SetMaxWidth(1).
			SetExpansion(1))
	}
}

func (i *imagePanel) updateEntries(g *Gui) {
	g.app.QueueUpdateDraw(func() {
		i.setEntries(g)
	})
}

func (i *imagePanel) focus(g *Gui) {
	i.SetSelectable(true, false)
	g.app.SetFocus(i)
}

func (i *imagePanel) unfocus() {
	i.SetSelectable(false, false)
}

func (i *imagePanel) setFilterWord(word string) {
	i.filterWord = word
}

func (i *imagePanel) monitoringImages(g *Gui) {
	ticker := time.NewTicker(g.refreshInterval)

LOOP:
	for {
		select {
		case <-ticker.C:
			i.updateEntries(g)
		case <-g.stopChans["image"]:
			ticker.Stop()
			break LOOP
		}
	}
}

func parseRepoTag(repoTag string) (string, string) {
	tmp := strings.Split(repoTag, ":")
	tag := tmp[len(tmp)-1]
	repo := strings.Join(tmp[0:len(tmp)-1], ":")
	return repo, tag
}
