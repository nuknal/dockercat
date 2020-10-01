package gui

type panel interface {
	name() string
	entries(*Gui)
	setEntries(*Gui)
	updateEntries(*Gui)
	setKeybinding(*Gui)
	focus(*Gui)
	unfocus()
	setFilterWord(string)
}

func (g *Gui) nextPanel() {
	idx := (g.panels.currentPanel + 1) % len(g.panels.panels)
	g.switchPanel(g.panels.panels[idx].name())
}

func (g *Gui) prevPanel() {
	g.panels.currentPanel--

	if g.panels.currentPanel < 0 {
		g.panels.currentPanel = len(g.panels.panels) - 1
	}

	idx := (g.panels.currentPanel) % len(g.panels.panels)
	g.switchPanel(g.panels.panels[idx].name())
}

func (g *Gui) switchPanel(panelName string) {
	for i, panel := range g.panels.panels {
		if panel.name() == panelName {
			g.navigate.update(panelName)
			panel.focus(g)
			g.panels.currentPanel = i
		} else {
			panel.unfocus()
		}
	}
}
