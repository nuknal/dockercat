package gui

import (
	"context"

	"github.com/nuknal/dockercat/pkg/common"
)

func (g *Gui) runTask() {
	g.Log.Info("start running task")
LOOP:
	for {
		select {
		case task := <-g.taskPanel().tasks:
			go func() {
				if err := task.Func(task.Ctx); err != nil {
					task.Status = err.Error()
				} else {
					task.Status = success
				}
				g.updateTask()
			}()
		case <-g.stopChans["task"]:
			g.Log.Info("stop monitoring task")
			break LOOP
		}
	}
}

func (g *Gui) addTask(taskName, target string, f func(ctx context.Context) error) {
	ctx, cancel := context.WithCancel(context.Background())

	task := &task{
		Name:      taskName,
		Status:    executing,
		Target:    target,
		CreatedAt: common.DateNow(),
		Func:      f,
		Ctx:       ctx,
		Cancel:    cancel,
	}

	g.resources.tasks = append(g.resources.tasks, task)
	go g.updateTask()
	g.taskPanel().tasks <- task
}

func (g *Gui) cancelTask() {
	taskPanel := g.taskPanel()
	row, _ := taskPanel.GetSelection()

	task := g.resources.tasks[row-1]
	if task.Status == executing {
		task.Cancel()
		task.Status = cancel
		g.updateTask()
	}
}

func (g *Gui) updateTask() {
	g.app.QueueUpdateDraw(func() {
		g.taskPanel().setEntries(g)
	})
}
