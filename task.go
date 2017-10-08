package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/0xAX/notificator"
	"github.com/ProtonMail/go-autostart"
	"github.com/fatih/color"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
	"github.com/olekukonko/tablewriter"
	"github.com/segmentio/go-prompt"
	"github.com/thedevsaddam/task/taskmanager"
)

const usage = `Usage:
	Name:
		Terminal Task
	Description:
		Your favorite terminal task manager and reminder!
	Version:
		1.0.0
	$ task
		Show all tasks
	$ task p
		Show all pending tasks
	$ task a Watch Games of thrones
		Add a new task [Watch Games of thrones] to list
	$ task remind Meeting with John tomorrow at 10:30pm
		This will send you a desktop notification
	$ task del
		Remove latest task from list
	$ task rm ID
		Remove task of ID from list
	$ task s ID
		Show detail view task of ID
	$ task c ID
		Mark task of ID as completed
	$ task m ID Pirates of the Caribbean
		Modify a task
	$ task p ID
		Mark task of ID as pending
	$ task flush
		Flush the database!
	$ task service-start
		Run task as service if you are using reminder
	$ task service-stop
		Unregister Task from service!
`

const (
	COMPLETED_MARK   = "\u2713"
	PENDING_MARK     = "\u2613"
	DATE_TIME_LAYOUT = "2006-01-02 15:04"
	REFRESH_RATE     = 40
)

var (
	//task manager instance
	tm = taskmanager.New()
	//notifier
	notify *notificator.Notificator
	//run a service
	service = autostart.App{
		Name:        "thedevsaddam_terminal_task",
		DisplayName: "Task",
		Exec:        []string{"/usr/local/bin/task", "listen-reminder-queue"},
	}
)

func main() {

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
		flag.PrintDefaults()
	}
	flag.Parse()
	cmd, args, argsLen := flag.Arg(0), flag.Args(), len(flag.Args())

	switch {
	case cmd == "" || cmd == "l" || cmd == "ls" && argsLen == 1:
		showTasksInTable(tm.GetAllTasks())
	case cmd == "a" || cmd == "add" && argsLen >= 1:
		if len(args[1:]) <= 0 {
			warningText(" Task description can not be empty \n")
			return
		}
		tm.Add(strings.Join(args[1:], " "), "", "")
		successText(" Added to list: " + strings.Join(args[1:], " ") + " ")
	case cmd == "reminder" || cmd == "remind" || cmd == "remind-me" && argsLen >= 1:
		if len(args[1:]) <= 0 {
			warningText(" Task/Reminder description can not be empty \n")
			return
		}
		reminder := strings.Join(args[1:], " ")
		action, action_when := parseReminder(reminder)
		tm.Add(action, "", action_when)
		successText(" Reminder Added: " + action + " ")
	case cmd == "p" || cmd == "pending" && argsLen == 1:
		showTasksInTable(tm.GetPendingTasks())
	case cmd == "del" || cmd == "delete" && argsLen == 1:
		p := prompt.Choose("Do you want to delete latest task?", []string{"yes", "no"})
		if p == 1 {
			warningText(" Task delete aboarted! ")
			return
		}
		err := tm.RemoveTask(tm.GetLastId())
		if err != nil {
			errorText(err.Error())
			return
		}
		successText(" Removed latest task ")
	case cmd == "r" || cmd == "rm" && argsLen == 2:
		id, _ := strconv.Atoi(flag.Arg(1))
		p := prompt.Choose("Do you want to delete task of id "+flag.Arg(1)+" ?", []string{"yes", "no"})
		if p == 1 {
			warningText(" Task delete aboarted! ")
			return
		}
		err := tm.RemoveTask(id)
		if err != nil {
			errorText(err.Error())
			return
		}
		successText(" Task " + strconv.Itoa(id) + " removed! ")
	case cmd == "e" || cmd == "m" || cmd == "u" && argsLen >= 2:
		id, _ := strconv.Atoi(flag.Arg(1))
		ok, _ := tm.UpdateTask(id, strings.Join(args[2:], " "))
		successText(ok)
	case cmd == "c" || cmd == "d" || cmd == "done" && argsLen >= 2:
		id, _ := strconv.Atoi(flag.Arg(1))
		task, err := tm.MarkAsCompleteTask(id)
		if err != nil {
			errorText(err.Error())
			return
		}
		successText(" " + COMPLETED_MARK + " " + task.Description)
	case cmd == "i" || cmd == "p" || cmd == "pending" && argsLen >= 2:
		id, _ := strconv.Atoi(flag.Arg(1))
		task, err := tm.MarkAsPendingTask(id)
		if err != nil {
			errorText(err.Error())
			return
		}
		successText(" " + pendingMark() + " " + task.Description)
	case cmd == "s" && argsLen == 2:
		id, _ := strconv.Atoi(flag.Arg(1))
		task, err := tm.GetTask(id)
		if err != nil {
			errorText(err.Error())
			return
		}
		showTask(task)
	case cmd == "flush":
		p := prompt.Choose("Do you want to delete all tasks?", []string{"yes", "no"})
		if p == 1 {
			warningText(" Flush aborted! ")
			return
		}
		err := tm.FlushDB()
		if err != nil {
			errorText(err.Error())
			return
		}
		successText(" Database flushed successfully! ")
	case cmd == "service-start" && argsLen == 1:
		serviceStart()
	case cmd == "service-force-start" && argsLen == 1:
		serviceForceStart()
	case cmd == "service-stop" && argsLen == 1:
		serviceStop()
	case cmd == "listen-reminder-queue" && argsLen == 1:
		listenReminderQueue()
	case cmd == "h" || cmd == "v":
		fmt.Fprint(os.Stderr, usage)
	default:
		errorText(" [No command found by " + cmd + "] ")
		fmt.Fprint(os.Stderr, "\n"+usage)
	}

}

//show tasks list in table
func showTasksInTable(tasks taskmanager.Tasks) {
	fmt.Fprintln(os.Stdout, "")
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Description", COMPLETED_MARK + "/" + pendingMark(), "Created"})
	table.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: false})
	table.SetFooter([]string{"", "Total: " + strconv.Itoa(tm.TotalTask()), "", "Pending: " + strconv.Itoa(tm.PendingTask())})
	table.SetCenterSeparator("|")
	table.SetRowLine(true)
	for _, task := range tasks {
		//set completed icon
		status := PENDING_MARK
		if task.Completed != "" {
			status = COMPLETED_MARK
		} else {
			status = pendingMark()
		}
		table.Append([]string{
			strconv.Itoa(task.Id),
			task.Description,
			status,
			task.Created,
		})
	}
	table.Render()
	fmt.Fprintln(os.Stdout, "")
}

//show a single tasks
func showTask(task taskmanager.Task) {
	fmt.Fprintln(os.Stdout, "")
	printText("Task Details view")
	printText("--------------------------------")
	printText("ID: " + strconv.Itoa(task.Id))
	printText("UID: " + task.UID)
	printText("Description: " + task.Description)
	printText("Tag: " + task.Tag)
	printText("Created: " + task.Created)
	printText("Updated: " + task.Updated)
	fmt.Fprintln(os.Stdout, "")
}

func printText(str string) {
	fmt.Fprintf(os.Stdout, str+"\n")
}

func printBoldText(str string) {
	if runtime.GOOS == "windows" {
		fmt.Fprintf(os.Stdout, str+"\n")
	} else {
		bold := color.New(color.Bold).FprintlnFunc()
		bold(os.Stdout, str)
	}
}

func successText(str string) {
	if runtime.GOOS == "windows" {
		fmt.Fprintf(color.Output, color.GreenString(str))
	} else {
		success := color.New(color.Bold, color.BgGreen, color.FgWhite).FprintlnFunc()
		success(os.Stdout, str)
	}
}

func warningText(str string) {
	if runtime.GOOS == "windows" {
		fmt.Fprintf(color.Output, color.YellowString(str))
	} else {
		warning := color.New(color.Bold, color.BgYellow, color.FgBlack).FprintlnFunc()
		warning(os.Stdout, str)
	}
}

func errorText(str string) {
	if runtime.GOOS == "windows" {
		fmt.Fprintf(color.Output, color.RedString(str))
	} else {
		error_ := color.New(color.Bold, color.BgRed, color.FgWhite).FprintlnFunc()
		error_(os.Stdout, str)
	}
}

func pendingMark() string {
	pending := PENDING_MARK
	if runtime.GOOS == "windows" {
		pending = "x"
	}
	return pending
}

//parse reminder
func parseReminder(reminder string) (string, string) {
	defer func() {
		if r := recover(); r != nil {
			errorText(" Your reminder does not contain any date time reference! ")
			os.Exit(1)
		}
	}()
	w := when.New(nil)
	w.Add(en.All...)
	w.Add(common.All...)
	r, _ := w.Parse(reminder, time.Now())
	action := strings.Replace(reminder, reminder[r.Index:r.Index+len(r.Text)], "", -1)
	actionTime := r.Time.Format(DATE_TIME_LAYOUT)
	return action, actionTime
}

//listen for reminder queue
func listenReminderQueue() {
	for {
		rm := taskmanager.New()
		reminderList := rm.GetReminderTasks()
		now := time.Now().Format(DATE_TIME_LAYOUT)
		for _, r := range reminderList {
			if r.RemindAt == now {
				desktopNotifier("Task Reminder!", r.Description)
				rm.MarkAsCompleteTask(r.Id)
			}
		}
		time.Sleep(time.Second * REFRESH_RATE)
	}
}

//send desktop notification
func desktopNotifier(title, body string) {
	notify = notificator.New(notificator.Options{
		DefaultIcon: "default.png",
		AppName:     "Terminal Task",
	})
	notify.Push(title, body, "default.png", notificator.UR_NORMAL)
}

//enable auto start
func serviceStart() {
	if service.IsEnabled() {
		warningText("Task is already enabled as service!")
	} else {
		if err := service.Enable(); err != nil {
			errorText(err.Error())
		}
		successText("Task has been registered as service!")
	}
}

//disable auto start
func serviceStop() {
	if service.IsEnabled() {
		if err := service.Disable(); err != nil {
			errorText(err.Error())
		}
		successText("Task has been removed from service!")
	} else {
		warningText("Task was not registered as service!")
	}
}

//force start
func serviceForceStart() {
	serviceStop()
	serviceStart()
}
