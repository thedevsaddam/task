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
	"time"
)

type (
	Task struct {
		Id          int    `json:"id"`
		UID         string `json:"uid"`
		Description string `json:"description"`
		Tag         string `json:"tag"`
		Created     string `json:"created"`
		Updated     string `json:"updated"`
		Completed   string `json:"completed"`
	}

	Tasks []Task
)

const (
	DB_FILE     = ".task.json"
	TIME_LAYOUT = "Mon, 01/02/06, 03:04PM"
)

func New() Tasks {
	return readDBFile()
}

//create a new task
func (t *Tasks) Add(description, tag string) Task {
	_t := Task{Id: t.GetNextId(), UID: uid(), Description: description, Tag: tag, Created: time.Now().Format(TIME_LAYOUT), Completed: ""}
	*t = append(*t, _t)
	writeDBFile(*t)
	return _t
}

//get all tasks
func (t Tasks) GetAllTasks() Tasks {
	sort.Sort(t)
	return t
}

//get completed tasks
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

//get pending tasks
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

//get a task
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

//update a task by id
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

//update a task tag by id
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

//mark a task as completed by id
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

//mark a task as pending by id
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

//remove a task by id
//ref https://stackoverflow.com/questions/18566499/how-to-remove-an-item-from-a-slice-by-calling-a-method-on-the-slice
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

//get total task
func (t Tasks) TotalTask() int {
	return len(t)
}

//get total completed task
func (t Tasks) CompletedTask() int {
	completed_task := 0
	for _, i := range t {
		if i.Completed != "" {
			completed_task++
		}
	}
	return completed_task
}

//get total pending task
func (t Tasks) PendingTask() int {
	return len(t) - t.CompletedTask()
}

//get last inserted id
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

//get next id
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
func (t Tasks) FlushDB() error {
	removeDBFileIfExist()
	createDBFileIfNotExist()
	return nil
}

//implement the sort interface
func (t Tasks) Len() int {
	return len(t)
}

func (t Tasks) Less(i, j int) bool {
	return t[i].Id > t[j].Id
}

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
