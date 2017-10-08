package taskmanager

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type (
	// Task describes a task object
	Task struct {
		Id          int    `json:"id"`
		UID         string `json:"uid"`
		Description string `json:"description"`
		Tag         string `json:"tag"`
		Created     string `json:"created"`
		Updated     string `json:"updated"`
		RemindAt    string `json:"remind_at"`
		Completed   string `json:"completed"`
	}

	// Tasks represents a list of Task object
	Tasks []Task
)

const (
	// DB_FILE is the default storage file path
	DB_FILE = ".task.json"
	// TIME_LAYOUT default time layout for task application
	TIME_LAYOUT = "Mon, 01/02/06, 03:04PM"
)

var mutex sync.Mutex

// New return a Task list instance
func New() Tasks {
	return readDBFile()
}

//Add create a new task
func (t *Tasks) Add(description, tag string, remind string) Task {
	_t := Task{Id: t.GetNextId(), UID: uid(), Description: description, Tag: tag, Created: time.Now().Format(TIME_LAYOUT), RemindAt: remind, Completed: ""}
	*t = append(*t, _t)
	writeDBFile(*t)
	return _t
}

//GetAllTasks fetch all tasks
func (t Tasks) GetAllTasks() Tasks {
	sort.Sort(t)
	return t
}

//GetCompletedTasks fetch all completed tasks
func (t Tasks) GetCompletedTasks() Tasks {
	var completedTasks Tasks
	for _, item := range t {
		if item.Completed != "" {
			completedTasks = append(completedTasks, item)
		}
	}
	sort.Sort(completedTasks)
	return completedTasks
}

//GetPendingTasks fetch all pending tasks
func (t Tasks) GetPendingTasks() Tasks {
	var pendingTasks Tasks
	for _, item := range t {
		if item.Completed == "" {
			pendingTasks = append(pendingTasks, item)
		}
	}
	sort.Sort(pendingTasks)
	return pendingTasks
}

//GetReminderTasks fetch all the reminder tasks
func (t Tasks) GetReminderTasks() Tasks {
	var reminderList Tasks
	for _, item := range t {
		if item.RemindAt != "" && item.Completed == "" {
			reminderList = append(reminderList, item) //only uncompleted reminder
		}
	}
	return reminderList
}

//GetTask fetch a single task
func (t Tasks) GetTask(id int) (Task, error) {
	if err := t.isValidId(id); err != nil {
		return Task{}, err
	}
	i, err := t.getIndexIdNo(id)
	if err != nil {
		return Task{}, err
	}
	return t[i], nil
}

//UpdateTask update a task by id
func (t *Tasks) UpdateTask(id int, description string) (string, error) {
	if err := t.isValidId(id); err != nil {
		return fmt.Sprintf("Unable to update %s", description), err
	}
	i, err := t.getIndexIdNo(id)
	if err != nil {
		return "", err
	}
	oldDescription := (*t)[i].Description
	(*t)[i].Description = description
	(*t)[i].Updated = time.Now().Format(TIME_LAYOUT)
	writeDBFile(*t)
	return fmt.Sprintf("Task Updated: %s --> %s", oldDescription, description), nil
}

//UpdateTaskTag update a task's tag by id
func (t *Tasks) UpdateTaskTag(id int, tag string) (string, error) {
	if err := t.isValidId(id); err != nil {
		return fmt.Sprintf("Unable to update %s", tag), err
	}
	i, err := t.getIndexIdNo(id)
	if err != nil {
		return "", nil
	}
	oldTag := (*t)[i].Tag
	(*t)[i].Tag = tag
	(*t)[i].Updated = time.Now().Format(TIME_LAYOUT)
	writeDBFile(*t)
	return fmt.Sprintf("Task Updated: %s --> %s", oldTag, tag), nil
}

//MarkAsCompleteTask mark a task as completed by id
func (t *Tasks) MarkAsCompleteTask(id int) (Task, error) {
	if err := t.isValidId(id); err != nil {
		return Task{}, err
	}
	i, err := t.getIndexIdNo(id)
	if err != nil {
		return Task{}, err
	}
	(*t)[i].Completed = time.Now().Format(TIME_LAYOUT)
	writeDBFile(*t)
	return (*t)[i], nil
}

//MarkAsPendingTask mark a task as pending by id
func (t *Tasks) MarkAsPendingTask(id int) (Task, error) {
	if err := t.isValidId(id); err != nil {
		return Task{}, err
	}
	i, err := t.getIndexIdNo(id)
	if err != nil {
		return Task{}, err
	}
	(*t)[i].Completed = ""
	writeDBFile(*t)
	return (*t)[i], nil
}

//RemoveTask delete a task by id
func (t *Tasks) RemoveTask(id int) error {
	if err := t.isValidId(id); err != nil {
		return err
	}
	i, err := t.getIndexIdNo(id)
	if err != nil {
		return err
	}
	*t = append((*t)[:i], (*t)[i+1:]...)
	writeDBFile(*t)
	return nil
}

//TotalTask return total task count
func (t Tasks) TotalTask() int {
	return len(t)
}

//CompletedTask return total completed task count
func (t Tasks) CompletedTask() int {
	completed_task := 0
	for _, i := range t {
		if i.Completed != "" {
			completed_task++
		}
	}
	return completed_task
}

//PendingTask return total pending task count
func (t Tasks) PendingTask() int {
	return len(t) - t.CompletedTask()
}

//GetLastId return last inserted id
func (t Tasks) GetLastId() int {
	if t.TotalTask() <= 0 {
		return 0
	}
	maxId := t[0].Id
	for _, item := range t {
		if item.Id >= maxId {
			maxId = item.Id
		}
	}
	return maxId
}

//GetNextId return next id
func (t Tasks) GetNextId() int {
	return t.GetLastId() + 1
}

//check if id valid
func (t Tasks) isValidId(id int) error {
	if id < 0 || id == 0 {
		return errors.New("Negative id not accepted!")
	}
	index, err := t.getIndexIdNo(id)
	if err != nil {
		return err
	}
	if index > len(t) {
		return errors.New("Id " + strconv.Itoa(id) + " not exist!")
	}
	return nil
}

// get indexIdNo from id
func (t Tasks) getIndexIdNo(id int) (int, error) {
	for i, task := range t {
		if task.Id == id {
			return i, nil
		}
	}
	return 0, errors.New("Invalid Id!")
}

//flush database
func (t *Tasks) FlushDB() error {
	*t = Tasks{}
	removeDBFileIfExist()
	createDBFileIfNotExist()
	return nil
}

//implement the sort interface
// Len return total length of task list
func (t Tasks) Len() int {
	return len(t)
}

// Less order the task as descending order
func (t Tasks) Less(i, j int) bool {
	return t[i].Id > t[j].Id
}

// Swap tasks
func (t Tasks) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

//===========================helpers
//generate a uid
func uid() string {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return ""
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

//get file path
func dbFile() string {
	env := os.Getenv("TASK_DB_FILE_PATH")
	if env != "" {
		if strings.HasSuffix(env, ".json") {
			return env
		} else {
			return filepath.Join(filepath.Clean(env), DB_FILE)
		}
	}

	usr, err := user.Current()
	if err != nil {
		panic(err)

	}
	return filepath.Join(usr.HomeDir, DB_FILE)
}

//load database
func readDBFile() Tasks {
	//load the json to task
	mutex.Lock()
	defer mutex.Unlock()
	file, e := ioutil.ReadFile(dbFile())
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	var tasks Tasks
	json.Unmarshal(file, &tasks)
	return tasks
}

//write to json
func writeDBFile(tasks Tasks) {
	mutex.Lock()
	defer mutex.Unlock()
	removeDBFileIfExist()
	taskJson, _ := json.Marshal(tasks)
	e := ioutil.WriteFile(dbFile(), taskJson, 0644)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
}

//create a db file if not exist
func createDBFileIfNotExist() {
	if _, err := os.Stat(dbFile()); os.IsNotExist(err) {
		os.Create(dbFile())
	}
}

//delete a db file if exist
func removeDBFileIfExist() {
	if _, err := os.Stat(dbFile()); !os.IsNotExist(err) {
		os.Remove(dbFile())
	}
}

func init() {
	//create .task.json file if not exist
	createDBFileIfNotExist()
}
