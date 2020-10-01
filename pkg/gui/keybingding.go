package gui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func (g *Gui) setGlobalKeybinding(event *tcell.EventKey) {
	switch event.Rune() {
	case 'h':
		g.prevPanel()
	case 'l':
		g.nextPanel()
	case 'q':
		g.Stop()
		// case '/':
		// 	g.filter()
	}

	switch event.Key() {
	case tcell.KeyTab:
		g.nextPanel()
	case tcell.KeyBacktab:
		g.prevPanel()
	case tcell.KeyRight:
		g.nextPanel()
	case tcell.KeyLeft:
		g.prevPanel()
	}
}

func (g *Gui) inspect(data, from, to string) {
	text := tview.NewTextView()
	text.SetTitle("Detail").SetTitleAlign(tview.AlignLeft)
	text.SetBorder(true)
	text.SetText(data)

	text.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc || event.Rune() == 'q' {
			// g.pages.RemovePage("inspect").ShowPage("usage")
		}
		return event
	})

	g.infoPanel().switchItemTextView(text, from, to)
}
