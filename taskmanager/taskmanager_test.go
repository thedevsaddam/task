package taskmanager

import (
	"os/user"
	"path/filepath"
	"testing"
)

var tasksList = []struct {
	description string
	uuid        string
	tag         string
}{
	{
		description: "Go to store", tag: "low",
	},
	{
		description: "Learn golang testing", tag: "high",
	},
	{
		description: "Watch Pirates of the carribean", tag: "medium",
	},
}

var tm Tasks

func TestMain(m *testing.M) {
	createDBFileIfNotExist()
	tm = New()
	m.Run()
	removeDBFileIfExist()
}

func TestTasks_Add(t *testing.T) {
	for _, t := range tasksList {
		tm.Add(t.description, t.tag)
	}
	if 3 != len(tasksList) {
		t.Error("Task count does not matched!")
	}
}

func TestTasks_GetAllTasks(t *testing.T) {
	tasks := tm.GetAllTasks()
	if len(tasks) != 3 {
		t.Error("Failed to get all tasks!")
	}
}

func TestTasks_GetTask(t *testing.T) {
	taskId := 3
	task, err := tm.GetTask(taskId)
	t.Log(task)
	if err != nil {
		t.Error("Failed to get task by id")
	}
	if task.Id != taskId {
		t.Error("Task id does not match in get task by id")
	}
}

func TestTasks_UpdateTask(t *testing.T) {
	task, err := tm.UpdateTask(1, "Go to USA")
	t.Log(task)
	if err != nil {
		t.Error("Unable to update task")
	}
}

func TestTasks_UpdateTaskTag(t *testing.T) {
	tag, err := tm.UpdateTaskTag(1, "important")
	t.Log(tag)
	if err != nil {
		t.Error("Unable to Update task tag")
	}
}

func TestTasks_MarkAsCompleteTask(t *testing.T) {
	taskId := 2
	task, err := tm.MarkAsCompleteTask(taskId)
	t.Log(task)
	if err != nil {
		t.Error("Unable to mark task as complete")
	}
	if task.Id != taskId {
		t.Error("Task id does not match in mark as complete")
	}
}

func TestTasks_GetCompletedTasks(t *testing.T) {
	tasks := tm.GetCompletedTasks()
	t.Log(tasks)
	if len(tasks) != 1 {
		t.Error("Failed to match number of completed tasks")
	}
	if tasks[0].Id != 2 {
		t.Error("Failed to match the completed task id")
	}
	for _, _task := range tasks {
		if len(_task.Completed) == 0 {
			t.Error("Failed match completed tasks status")
		}
	}
}

func TestTasks_MarkAsPendingTask(t *testing.T) {
	taskId := 1
	task, err := tm.MarkAsPendingTask(taskId)
	t.Log(task)
	if err != nil {
		t.Error("Unable to mark task as pending")
	}
	if task.Id != taskId {
		t.Error("Task id does not match in mark as pending")
	}
}

func TestTasks_GetPendingTasks(t *testing.T) {
	tasks := tm.GetPendingTasks()
	if len(tasks) != 2 {
		t.Error("Failed to match the pending tasks total count!")
	}

}

func TestTasks_GetLastId(t *testing.T) {
	if tm.GetLastId() != 3 {
		t.Error("Last inserted id does not correct!")
	}
}

func TestTasks_GetNextId(t *testing.T) {
	if tm.GetNextId() != 4 {
		t.Error("Next id does not correct!")
	}
}

func TestTasks_TotalTask(t *testing.T) {
	if tm.TotalTask() != 3 {
		t.Error("Failed to count total task!")
	}
}

func TestTasks_Len(t *testing.T) {
	if tm.Len() != 3 {
		t.Error("Failed to count let of tasks")
	}
}

func TestTasks_RemoveTask(t *testing.T) {
	err := tm.RemoveTask(3)
	if err != nil {
		t.Error("Failed to remove task")
	}
	if tm.TotalTask() != 2 {
		t.Error("Task did not remove properly!")
	}
}

func TestTasks_FlushDB(t *testing.T) {
	err := tm.FlushDB()
	if err != nil {
		t.Error("Failed to flush database!")
	}
}

func TestTasks_dbFile(t *testing.T) {
	usr, _ := user.Current()
	if dbFile() != filepath.Join(filepath.Clean(usr.HomeDir), ".task.json") {
		t.Error("Task file path incorrect!")
	}
}

func BenchmarkTasks_Add(b *testing.B) {
	for _, task := range tasksList {
		tm.Add(task.description, task.tag)
	}
}

func BenchmarkTasks_GetAllTasks(b *testing.B) {
	tm.GetAllTasks()
}

func BenchmarkTasks_GetCompletedTasks(t *testing.B) {
	tm.GetCompletedTasks()
}

func BenchmarkTasks_GetPendingTasks(b *testing.B) {
	tm.GetPendingTasks()
}

func BenchmarkTasks_GetTask(b *testing.B) {
	tm.GetTask(2)
}

func BenchmarkTasks_MarkAsCompleteTask(b *testing.B) {
	tm.MarkAsCompleteTask(1)
}

func BenchmarkTasks_MarkAsPendingTask(b *testing.B) {
	tm.MarkAsPendingTask(3)
}
