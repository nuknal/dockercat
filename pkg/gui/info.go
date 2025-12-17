package gui

import (
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

const (
	about = `------------------------------------
[cyan::b]Dockercat[-::-]  Terminal UI for Docker

[green::b]Quick Start[-::-]
  [yellow::b]h/l[-::-] or [yellow::b]←/→[-::-], [yellow::b]Tab[-::-]   Switch panels
  [yellow::b]/[-::-]                       Filter list in current panel
  [yellow::b]Enter[-::-]                   Inspect item (detail on right)
  [yellow::b]q[-::-]                       Quit

[green::b]Containers[-::-]
  [yellow::b]m[-::-] Select   [yellow::b]U/S/D[-::-] Batch start/stop/remove
  [yellow::b]u/s/r[-::-] Start/Stop/Restart   [yellow::b]p/o[-::-] Pause/Unpause
  [yellow::b]Ctrl+l[-::-] Logs   [yellow::b]Ctrl+s[-::-] Stats

[green::b]Images[-::-]
  [yellow::b]d[-::-] Remove   [yellow::b]p[-::-] Pull   [yellow::b]t[-::-] Tag   [yellow::b]P[-::-] Push

[green::b]Volumes / Networks[-::-]
  [yellow::b]Enter[-::-] Inspect   [yellow::b]d[-::-] Remove   [yellow::b]Ctrl+d[-::-] Prune (where supported)

[green::b]Cleanup[-::-]
  [yellow::b]a[-::-] System prune   [yellow::b]c/i/n/v[-::-] Prune by type

[green::b]Detail View[-::-]
  [yellow::b]Ctrl+j / Ctrl+k[-::-] Scroll detail
  [gray]Tip: Bottom bar shows panel-specific shortcuts.[-]
------------------------------------
`
)

// the panel for display detail information
type infoPanel struct {
	*tview.Flex                         //
	key             string              // a key likes "container-id-logs"
	stopReadStreams map[string]chan int // container logs stream and stat stream
	itemTextView    *tview.TextView
}

func newInfoPanel(g *Gui) *infoPanel {
	info := &infoPanel{
		Flex: tview.NewFlex(),
	}

	stopChan := make(chan int)
	info.stopReadStreams = make(map[string]chan int)
	info.stopReadStreams["container-logs"] = stopChan
	info.stopReadStreams["container-stats"] = stopChan

	info.setEntries(g)
	info.setKeybinding(g)

	info.key = "system-docker-info"
	info.itemTextView = tview.NewTextView().SetText(about)
	info.itemTextView.SetBorder(true)
	info.itemTextView.SetDynamicColors(true)
	info.itemTextView.SetBorderColor(CurrentTheme.Border)
	info.itemTextView.SetTextColor(CurrentTheme.Fg)
	info.AddItem(info.itemTextView, 0, 1, true)

	info.getDockerInfo(g)

	return info
}

func (i *infoPanel) getDockerInfo(g *Gui) {
	go func() {
		dockerInfo, _ := g.getDockerInfo()
		du := g.getDiskUsage()

		var dInfo string
		if dockerInfo != nil {
			dInfo = fmt.Sprintf("[green]Docker\n  Host: [%s][%s] \n  Endpoint: [%s]\n  %s Mem | %d Containers | %d Images %s | volumes %s[white]\n\n",
				dockerInfo.HostName,
				dockerInfo.ServerVersion,
				dockerInfo.Endpoint,
				dockerInfo.MemTotal,
				dockerInfo.Containers,
				dockerInfo.Images,
				du.imagesSize,
				du.volumesSize,
			)
		}

		g.app.QueueUpdateDraw(func() {
			// TODO mutex read key
			if i.key == "system-docker-info" {
				i.itemTextView.SetText(dInfo + about)
			}
		})
	}()
}

func (i *infoPanel) name() string {
	return "info"
}

func (i *infoPanel) setKey(key string) {
	i.key = key
}

func (i *infoPanel) entries(*Gui) {}

func (i *infoPanel) setEntries(g *Gui) {}

func (i *infoPanel) updateEntries(g *Gui) {}

func (i *infoPanel) setKeybinding(g *Gui) {
	i.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		g.setGlobalKeybinding(event)
		// switch event.Rune() {
		// case 'q':
		// }

		return event
	})
}

func (i *infoPanel) focus(g *Gui) {
	g.app.SetFocus(i)
}

func (i *infoPanel) unfocus() {

}

func (i *infoPanel) setFilterWord(string) {

}

func (i *infoPanel) scrollDown() {
	row, col := i.itemTextView.GetScrollOffset()
	i.itemTextView.ScrollTo(row+1, col)
}

func (i *infoPanel) scrollUp() {
	row, col := i.itemTextView.GetScrollOffset()
	if row > 0 {
		i.itemTextView.ScrollTo(row-1, col)
	}
}

func (i *infoPanel) switchItemTextView(v *tview.TextView, fromKey, toKey string) {
	// stop pre streams
	from := strings.Split(fromKey, "-")
	if len(from) > 0 {
		switch from[len(from)-1] {
		case "logs":
			i.stopReadStreams["container-logs"] <- 1
		case "stats":
			i.stopReadStreams["container-stats"] <- 1

		}
	}
	i.itemTextView = v
	i.Clear().AddItem(v, 0, 1, true)
}

func (i *infoPanel) readContainerLogs(g *Gui, id string, stop <-chan int) (<-chan string, error) {
	logs := make(chan string, 10)

	go func() {
		reader, err := g.client.TailContainerLogStream(id, "5")
		if err != nil {
			return
		}

		hdr := make([]byte, 8)
	LOOP:
		for {
			select {
			case <-stop:
				reader.Close()
				break
			default:
				_, err := reader.Read(hdr)
				if err != nil {
					break LOOP
				}
				count := binary.BigEndian.Uint32(hdr[4:])
				dat := make([]byte, count)
				_, err = reader.Read(dat)
				logs <- string(dat)
			}
		}
	}()

	return logs, nil
}

func (i *infoPanel) containerLogs(g *Gui, id string) error {
	stopChan := make(chan int)
	logs, err := i.readContainerLogs(g, id, stopChan)
	if err != nil {
		return err
	}

LOOP:
	for {
		select {
		case clog := <-logs:
			g.app.QueueUpdateDraw(func() {
				i.itemTextView.Write([]byte(clog))
			})
		case <-i.stopReadStreams["container-logs"]:
			stopChan <- 1
			break LOOP
		}
	}

	return nil
}

func (i *infoPanel) containerStats(g *Gui, container *container) {
	ticker := time.NewTicker(g.refreshInterval)

LOOP:
	for {
		select {
		case <-ticker.C:
			cstat, ok := g.containerPanel().containersStat[container.ID]
			if ok {
				_, _, w, _ := i.GetRect()
				content, _ := cstat.RenderStats(w)
				g.app.QueueUpdateDraw(func() {
					i.itemTextView.SetText(content)
				})
			} else {
				g.Log.Error(container.ID, " has not been monitored")
			}
		case <-i.stopReadStreams["container-stats"]:
			ticker.Stop()
			break LOOP
		}
	}
}
