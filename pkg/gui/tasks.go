package gui

import (
	"context"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

var (
	success   = "Success"
	executing = "Executing"
	cancel    = "canceled"
)

type task struct {
	Name      string
	CreatedAt string
	Status    string
	Target    string
	Func      func(ctx context.Context) error
	Ctx       context.Context
	Cancel    context.CancelFunc
}

type tasks struct {
	*tview.Table
	tasks chan *task
}

func newTasks(g *Gui) *tasks {
	tasks := &tasks{
		Table: tview.NewTable().SetSelectable(true, false).Select(0, 0).SetFixed(1, 1),
		tasks: make(chan *task, 10),
	}

	tasks.SetTitle(" Tasks ").SetTitleAlign(tview.AlignLeft)
	tasks.SetBorder(true)
	tasks.SetBorderColor(CurrentTheme.Border)
	tasks.SetTitleColor(CurrentTheme.Title)
	tasks.setEntries(g)
	tasks.setKeybinding(g)
	return tasks
}

func (t *tasks) name() string {
	return "tasks"
}

func (t *tasks) setKeybinding(g *Gui) {
	t.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		g.setGlobalKeybinding(event)

		switch event.Key() {
		}

		switch event.Rune() {
		case 'c':
			g.cancelTask()
		}

		return event
	})
}

func (t *tasks) entries(g *Gui) {
	// do nothing
}

func (t *tasks) setEntries(g *Gui) {
	t.entries(g)
	if len(g.resources.tasks) == 0 {
		return
	}
	table := t.Clear()
	table.SetSelectedStyle(CurrentTheme.SelectedFg, CurrentTheme.SelectedBg, 0)
	headers := []string{
		"  Name",
		"Target",
		"Status",
		"Created",
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

	for i, task := range g.resources.tasks {
		row := len(g.resources.tasks) - i
		table.SetCell(row, 0, tview.NewTableCell("  "+task.Name).
			SetTextColor(CurrentTheme.Tasks).
			SetMaxWidth(1).
			SetExpansion(2))

		table.SetCell(row, 1, tview.NewTableCell(task.Target).
			SetTextColor(CurrentTheme.Tasks).
			SetMaxWidth(1).
			SetExpansion(2))

		table.SetCell(row, 2, tview.NewTableCell(task.Status).
			SetTextColor(CurrentTheme.Tasks).
			SetMaxWidth(1).
			SetExpansion(1))

		table.SetCell(row, 3, tview.NewTableCell(task.CreatedAt).
			SetTextColor(CurrentTheme.Tasks).
			SetMaxWidth(1).
			SetExpansion(2))

	}
}

func (t *tasks) focus(g *Gui) {
	t.SetSelectable(true, false)
	g.app.SetFocus(t)
}

func (t *tasks) unfocus() {
	t.SetSelectable(false, false)
}

func (t *tasks) setFilterWord(word string) {
	// do nothings
}

func (t *tasks) updateEntries(g *Gui) {
	// do nothings
}
