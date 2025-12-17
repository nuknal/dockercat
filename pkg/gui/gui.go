package gui

import (
	"time"

	"github.com/nuknal/dockercat/pkg/common"
	"github.com/nuknal/dockercat/pkg/docker"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
)

type resources struct {
	containers     []*container
	containerStats []*ContainerStats
	images         []*image
	networks       []*network
	volumes        []*volume
	tasks          []*task
}

type panels struct {
	currentPanel int
	panels       []panel
}

// Gui holds all panels
type Gui struct {
	app             *tview.Application
	panels          panels
	navigate        *navigate
	pages           *tview.Pages
	resources       resources
	stopChans       map[string]chan int
	client          *docker.Docker
	errChan         chan error
	refreshInterval time.Duration

	Log *logrus.Entry
}

func New(client *docker.Docker, refreshInterval time.Duration) *Gui {
	gui := &Gui{
		app:             tview.NewApplication(),
		client:          client,
		stopChans:       make(map[string]chan int),
		errChan:         make(chan error, 10),
		Log:             common.NewLogger(),
		refreshInterval: refreshInterval,
	}

	return gui
}

// Run start app
func (g *Gui) Run() error {
	g.initPanels()
	g.startMonitoring()
	g.errHanlder()
	return g.app.Run()
}

// Stop stop app
func (g *Gui) Stop() error {
	g.app.Stop()
	return nil
}

func (g *Gui) errHanlder() {
	go func() {
		for {
			select {
			case err := <-g.errChan:
				logrus.Error(err)
				g.showError(err.Error())
			}
		}
	}()
}

func (g *Gui) imagePanel() *imagePanel {
	for _, panel := range g.panels.panels {
		if panel.name() == "images" {
			return panel.(*imagePanel)
		}
	}
	return nil
}

func (g *Gui) containerPanel() *containerPanel {
	for _, panel := range g.panels.panels {
		if panel.name() == "containers" {
			return panel.(*containerPanel)
		}
	}
	return nil
}

func (g *Gui) volumePanel() *volumePanel {
	for _, panel := range g.panels.panels {
		if panel.name() == "volumes" {
			return panel.(*volumePanel)
		}
	}
	return nil
}

func (g *Gui) networkPanel() *networkPanel {
	for _, panel := range g.panels.panels {
		if panel.name() == "networks" {
			return panel.(*networkPanel)
		}
	}
	return nil
}

func (g *Gui) infoPanel() *infoPanel {
	for _, panel := range g.panels.panels {
		if panel.name() == "info" {
			return panel.(*infoPanel)
		}
	}
	return nil
}

func (g *Gui) cleanupPanel() *cleanupPanel {
	for _, panel := range g.panels.panels {
		if panel.name() == "cleanup" {
			return panel.(*cleanupPanel)
		}
	}
	return nil
}

func (g *Gui) taskPanel() *tasks {
	for _, panel := range g.panels.panels {
		if panel.name() == "tasks" {
			return panel.(*tasks)
		}
	}
	return nil
}

func (g *Gui) initPanels() {
	containers := newContainerPanel(g)
	images := newImagePanel(g)
	volumes := newVolumePanel(g)
	networks := newNetworkPanel(g)
	cleanup := newCleanupPanel(g)
	navi := newNavigate()
	tasks := newTasks(g)

	infoPanel := newInfoPanel(g)

	g.panels.panels = append(g.panels.panels, containers)
	g.panels.panels = append(g.panels.panels, images)
	g.panels.panels = append(g.panels.panels, volumes)
	g.panels.panels = append(g.panels.panels, networks)
	g.panels.panels = append(g.panels.panels, cleanup)
	g.panels.panels = append(g.panels.panels, infoPanel)
	g.panels.panels = append(g.panels.panels, tasks)
	g.navigate = navi

	left := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(containers, 0, 1, false).
		AddItem(images, 0, 1, false).
		AddItem(volumes, 0, 1, false).
		AddItem(networks, 0, 1, false).
		AddItem(cleanup, 0, 1, false)

	main := tview.NewFlex()
	main.AddItem(left, 0, 1, false).AddItem(infoPanel, 0, 2, false)

	root := tview.NewFlex().SetDirection(tview.FlexRow)
	root.AddItem(main, 0, 4, true).
		AddItem(tasks, 0, 1, true).
		AddItem(navi, 3, 1, false)

	g.pages = tview.NewPages()
	g.pages.AddAndSwitchToPage("main", root, true)

	g.app.SetRoot(g.pages, true)
	g.switchPanel("containers")
}

func (g *Gui) selectedContainer() *container {
	row, _ := g.containerPanel().GetSelection()
	if len(g.resources.containers) == 0 {
		return nil
	}
	if row-1 < 0 {
		return nil
	}

	return g.resources.containers[row-1]
}

func (g *Gui) selectedImage() *image {
	row, _ := g.imagePanel().GetSelection()
	if len(g.resources.images) == 0 {
		return nil
	}
	if row-1 < 0 {
		return nil
	}

	return g.resources.images[row-1]
}

func (g *Gui) selectedVolume() *volume {
	row, _ := g.volumePanel().GetSelection()
	if len(g.resources.volumes) == 0 {
		return nil
	}
	if row-1 < 0 {
		return nil
	}

	return g.resources.volumes[row-1]
}

func (g *Gui) selectedNetwork() *network {
	row, _ := g.networkPanel().GetSelection()
	if len(g.resources.networks) == 0 {
		return nil
	}
	if row-1 < 0 {
		return nil
	}

	return g.resources.networks[row-1]
}

func (g *Gui) startMonitoring() {
	g.stopChans["task"] = make(chan int, 1)
	g.stopChans["info"] = make(chan int, 1)
	g.stopChans["image"] = make(chan int, 1)
	g.stopChans["volume"] = make(chan int, 1)
	g.stopChans["network"] = make(chan int, 1)
	g.stopChans["container"] = make(chan int, 1)
	go g.runTask()
	go g.imagePanel().monitoringImages(g)
	go g.networkPanel().monitoringNetworks(g)
	go g.volumePanel().monitoringVolumes(g)
	go g.containerPanel().monitoringContainers(g)
}

func (g *Gui) confirm(message, fromPanel string, doneFunc func()) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"Confirm", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			g.closeAndSwitchPage("confirm", fromPanel)
			if buttonLabel == "Confirm" {
				doneFunc()
			}
		})

	confirmDialog := tview.NewGrid().
		SetColumns(0, 80, 0).
		SetRows(0, 29, 0).
		AddItem(modal, 1, 1, 1, 1, 0, 0, true)
	g.pages.AddAndSwitchToPage("confirm", confirmDialog, true).ShowPage("main")
}

func (g *Gui) closeAndSwitchPage(removePage, switchPanel string) {
	g.pages.RemovePage(removePage).ShowPage("main")
	g.switchPanel(switchPanel)
}

func (g *Gui) showError(message string) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"Close"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			g.pages.RemovePage("error").ShowPage("main")
		})
	dialog := tview.NewGrid().
		SetColumns(0, 80, 0).
		SetRows(0, 10, 0).
		AddItem(modal, 1, 1, 1, 1, 0, 0, true)
	g.pages.AddAndSwitchToPage("error", dialog, true).ShowPage("main")
}

func (g *Gui) filter() {
	input := tview.NewInputField().SetLabel("Filter: ")
	form := tview.NewForm().AddFormItem(input).
		AddButton("Apply", func() {
			word := input.GetText()
			if g.panels.currentPanel >= 0 && g.panels.currentPanel < len(g.panels.panels) {
				p := g.panels.panels[g.panels.currentPanel]
				p.setFilterWord(word)
				p.setEntries(g)
			}
			g.pages.RemovePage("filter").ShowPage("main")
		}).
		AddButton("Cancel", func() {
			g.pages.RemovePage("filter").ShowPage("main")
		})
	form.SetBorder(true).SetTitle("Filter").SetTitleAlign(tview.AlignLeft)
	grid := tview.NewGrid().SetColumns(0, 60, 0).SetRows(0, 7, 0).AddItem(form, 1, 1, 1, 1, 0, 0, true)
	g.pages.AddAndSwitchToPage("filter", grid, true).ShowPage("main")
}
