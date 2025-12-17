package gui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type navigate struct {
	*tview.TextView
	keybindings map[string]string
}

func newNavigate() *navigate {
	return &navigate{
		TextView: tview.NewTextView().SetTextColor(tcell.ColorYellow),
		keybindings: map[string]string{
			"images":     " Enter: inspect image, d: remove image, s: show parent layers, \n Ctrl+r: refresh images list, Ctrl+d: remove unused and untagged images, p: pull, t: tag, P: push, /: filter",
			"containers": " Enter: inspect container, d: remove container, u: start container, s: stop container, r: restart, p: pause, o: unpause \n Ctrl+r: refresh container list, Ctrl+l: show container logs, Ctrl+s: show container stats, m: select, U/S/D: batch start/stop/remove, /: filter",
			"networks":   " Enter: inspect network, d: remove network, Ctrl+r: refresh network list, Ctrl+d: prune unused networks, /: filter",
			"volumes":    " Enter: inspect volume, d: remove volume, Ctrl+r: refresh volume list, Ctrl+d: prune unused volumes, /: filter",
			"cleanup":    " a: system prune, c: prune containers, i: prune images, n: prune networks, v: prune volumes, Ctrl+r: refresh",
		},
	}
}

func (n *navigate) update(panel string) {
	n.SetText(n.keybindings[panel])
}
