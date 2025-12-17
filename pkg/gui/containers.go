package gui

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/gdamore/tcell"
	"github.com/jesseduffield/asciigraph"
	"github.com/mcuadros/go-lookup"
	"github.com/nuknal/dockercat/pkg/common"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
)

type container struct {
	ID      string
	Name    string
	State   string
	Status  string
	LogType string
}

var statsConfig *StatsConfig

type containerStat struct {
	ID              string
	StatHistory     []RecordedStats
	ContainerMutex  sync.Mutex
	MonitoringStats bool
	StatsConfig     *StatsConfig
	stopChan        chan int
}

// implement the panel interface
type containerPanel struct {
	*tview.Table
	filterWord     string
	log            *logrus.Entry
	containersStat map[string]*containerStat
	selected       map[string]bool
}

func newContainerPanel(g *Gui) *containerPanel {
	containers := &containerPanel{
		Table:          tview.NewTable().SetSelectable(true, false).Select(0, 0).SetFixed(1, 1),
		containersStat: make(map[string]*containerStat),
		selected:       make(map[string]bool),
	}

	statsConfig = &StatsConfig{
		MaxDuration: time.Minute * 5,
		Graphs: []GraphConfig{
			{
				Caption:  "CPU (%)",
				StatPath: "DerivedStats.CPUPercentage",
				Color:    "cyan",
			},
			{
				Caption:  "Memory (%)",
				StatPath: "DerivedStats.MemoryPercentage",
				Color:    "green",
			},
		},
	}

	containers.SetTitle(" Containers ").SetTitleAlign(tview.AlignLeft)
	containers.SetBorder(true)
	containers.SetBorderColor(CurrentTheme.Border)
	containers.SetTitleColor(CurrentTheme.Title)
	containers.setEntries(g)
	containers.setKeybinding(g)
	return containers
}

func (c *containerPanel) name() string {
	return "containers"
}

func (c *containerPanel) setKeybinding(g *Gui) {
	c.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		g.setGlobalKeybinding(event)
		switch event.Key() {
		case tcell.KeyEnter:
			g.inspectContainer()
		// case tcell.KeyCtrlE:
		// 	g.attachContainerForm()
		case tcell.KeyCtrlL:
			g.tailContainerLog()
		case tcell.KeyCtrlR:
			c.setEntries(g)
		case tcell.KeyCtrlS:
			g.containerStats()
		}

		switch event.Rune() {
		case 'd':
			g.removeContainer()
		case 'u':
			g.startContainer()
		case 's':
			g.stopContainer()
		case 'm':
			c.toggleSelectCurrent(g)
		case 'U':
			g.batchStartContainers()
		case 'S':
			g.batchStopContainers()
		case 'D':
			g.batchRemoveContainers()
		case 'r':
			g.restartContainer()
		case 'p':
			g.pauseContainer()
		case 'o':
			g.unpauseContainer()
		}

		return event
	})
}

func (c *containerPanel) entries(g *Gui) {
	containers, err := g.client.Containers(types.ContainerListOptions{All: true})
	if err != nil {
		return
	}

	g.resources.containers = make([]*container, 0)

	for _, con := range containers {
		if strings.Index(con.Names[0][1:], c.filterWord) == -1 {
			continue
		}

		// TODO if a container has been removed, we should stop the monitor
		if cstat, ok := c.containersStat[con.ID]; ok {
			if !cstat.MonitoringStats {
				go cstat.statMonitor(g)
			}
		} else {
			c.containersStat[con.ID] = &containerStat{
				ID:          con.ID,
				StatsConfig: statsConfig,
				stopChan:    make(chan int),
			}
			go c.containersStat[con.ID].statMonitor(g)
		}

		g.resources.containers = append(g.resources.containers, &container{
			ID:     con.ID,
			Name:   con.Names[0][1:],
			Status: con.Status,
			State:  con.State,
		})
	}

	present := make(map[string]struct{}, len(g.resources.containers))
	for _, con := range g.resources.containers {
		present[con.ID] = struct{}{}
	}
	for id, cstat := range c.containersStat {
		if _, ok := present[id]; !ok {
			select {
			case cstat.stopChan <- 1:
			default:
			}
			delete(c.containersStat, id)
		}
	}
}

func (c *containerPanel) render(g *Gui) {
	table := c.Clear()

	table.SetSelectedStyle(CurrentTheme.SelectedFg, CurrentTheme.SelectedBg, 0)

	headers := []string{
		"  Name",
		"Status",
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

	for i, container := range g.resources.containers {
		mark := "  "
		if c.selected[container.ID] {
			mark = "● "
		}

		table.SetCell(i+1, 0, tview.NewTableCell(mark+container.Name).
			SetTextColor(CurrentTheme.ListItem).
			SetMaxWidth(1).
			SetExpansion(3))

		table.SetCell(i+1, 1, tview.NewTableCell(container.Status).
			SetTextColor(GetStatusColor(container.State)).
			SetMaxWidth(1).
			SetExpansion(2))
	}
}

func (c *containerPanel) setEntries(g *Gui) {
	c.entries(g)
	c.render(g)
}

func (c *containerPanel) focus(g *Gui) {
	c.SetSelectable(true, false)
	g.app.SetFocus(c)
}

func (c *containerPanel) unfocus() {
	c.SetSelectable(false, false)
}

func (c *containerPanel) updateEntries(g *Gui) {
	g.app.QueueUpdateDraw(func() {
		c.setEntries(g)
	})
}

func (c *containerPanel) setFilterWord(word string) {
	c.filterWord = word
}

func (c *containerPanel) toggleSelectCurrent(g *Gui) {
	con := g.selectedContainer()
	if con == nil {
		return
	}
	c.selected[con.ID] = !c.selected[con.ID]
	g.app.QueueUpdateDraw(func() {
		c.render(g)
	})
}

func (c *containerPanel) monitoringContainers(g *Gui) {
	ticker := time.NewTicker(g.refreshInterval)

LOOP:
	for {
		select {
		case <-ticker.C:
			c.updateEntries(g)
		case <-g.stopChans["container"]:
			ticker.Stop()
			break LOOP
		}
	}
}

func (c *containerStat) statMonitor(g *Gui) {
	g.Log.Info("start stats monitor ", c.ID)
	stream, err := g.client.ContainerStatsStream(c.ID)
	if err != nil {
		return
	}

	c.MonitoringStats = true
	defer stream.Close()

	scanner := bufio.NewScanner(stream)
LOOP:
	for scanner.Scan() {
		select {
		case <-c.stopChan:
			break LOOP
		default:
			data := scanner.Bytes()
			var stats ContainerStats
			json.Unmarshal(data, &stats)

			recordedStats := RecordedStats{
				ClientStats: stats,
				DerivedStats: DerivedStats{
					CPUPercentage:    stats.CalculateContainerCPUPercentage(),
					MemoryPercentage: stats.CalculateContainerMemoryUsage(),
				},
				RecordedAt: time.Now(),
			}

			c.ContainerMutex.Lock()
			c.StatHistory = append(c.StatHistory, recordedStats)
			c.eraseOldHistory()
			c.ContainerMutex.Unlock()
		}
	}

	c.MonitoringStats = false
}

func (c *containerStat) eraseOldHistory() {
	for i, stat := range c.StatHistory {
		if time.Since(stat.RecordedAt) < 5*time.Minute {
			c.StatHistory = c.StatHistory[i:]
			return
		}
	}
}

// RenderStats returns a string containing the rendered stats of the container
func (c *containerStat) RenderStats(viewWidth int) (string, error) {
	history := c.StatHistory
	if len(history) == 0 {
		return "", nil
	}
	currentStats := history[len(history)-1]

	graphSpecs := c.StatsConfig.Graphs
	graphs := make([]string, len(graphSpecs))
	for i, spec := range graphSpecs {
		graph, err := c.PlotGraph(spec, viewWidth-10)
		if err != nil {
			return "", err
		}
		graphs[i] = graph
	}

	pidsCount := fmt.Sprintf("PIDs: %d", currentStats.ClientStats.PidsStats.Current)
	dataReceived := fmt.Sprintf("Traffic received: %s", common.FormatDecimalBytes(currentStats.ClientStats.Networks.Eth0.RxBytes))
	dataSent := fmt.Sprintf("Traffic sent: %s", common.FormatDecimalBytes(currentStats.ClientStats.Networks.Eth0.TxBytes))

	originalJSON, err := json.MarshalIndent(currentStats, "", "  ")
	if err != nil {
		return "", err
	}

	contents := fmt.Sprintf("\n\n%s\n\n%s\n\n%s\n%s\n\n%s",
		strings.Join(graphs, "\n\n"),
		pidsCount,
		dataReceived,
		dataSent,
		string(originalJSON),
	)

	return contents, nil
}

// PlotGraph returns the plotted graph based on the graph spec and the stat history
func (c *containerStat) PlotGraph(spec GraphConfig, width int) (string, error) {
	data := make([]float64, len(c.StatHistory))

	max := spec.Max
	min := spec.Min
	for i, stats := range c.StatHistory {
		value, err := lookup.LookupString(stats, spec.StatPath)
		if err != nil {
			return "Could not find key: " + spec.StatPath, nil
		}
		floatValue, err := common.GetFloat(value.Interface())
		if err != nil {
			return "", err
		}
		if spec.MinType == "" {
			if i == 0 {
				min = floatValue
			} else if floatValue < min {
				min = floatValue
			}
		}

		if spec.MaxType == "" {
			if i == 0 {
				max = floatValue
			} else if floatValue > max {
				max = floatValue
			}
		}

		data[i] = floatValue
	}

	height := 10
	if spec.Height > 0 {
		height = spec.Height
	}

	return asciigraph.Plot(
		data,
		asciigraph.Height(height),
		asciigraph.Width(width),
		asciigraph.Min(min),
		asciigraph.Max(max),
		asciigraph.Caption(fmt.Sprintf("%s: %0.2f (%v)", spec.Caption, data[len(data)-1], time.Since(c.StatHistory[0].RecordedAt).Round(time.Second))),
	), nil
}
