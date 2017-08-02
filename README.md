# Task
Terminal tasks todo tool for geek

![Task screenshot](https://raw.githubusercontent.com/thedevsaddam/task/master/screenshot.png)

**[Dwnloads Binaries from here](https://github.com/thedevsaddam/task_binaries)**

Mac/Linux download the binary
```bash
$ cp task /usr/local/bin/task
$ sudo chmod +x /usr/local/bin/task
```
For windows download the binary and set environment variable so that you can access the binary from terminal

### Custom File Path
If you are interested to sync the task in dropbox/google drive, you can set a custom path. To set a custom path
 open your `.bashrc` or `.bash_profile` and add this line `export TASK_DB_FILE_PATH=Your file path`
 
 Example File path
 ```bash
export TASK_DB_FILE_PATH=/home/thedevsaddam/Dropbox  # default file name will be .task.json
export TASK_DB_FILE_PATH=/home/thedevsaddam/Dropbox/mytasks.json
```

### Usage
* List all the tasks
    ```bash
    $ task
    ```
* Add a new task to list
    ```bash
    $ task a Pirates of the Caribbean: Dead Men Tell No Tales
    ```
* Add a reminder task to list
    ```bash
    $ task reminder Meeting with Jane next wednesday at 2:30pm
    ```
* List all pending tasks
    ```bash
    $ task p
    ```
* Show a task details
    ```bash
    $ task s ID
    ```
* Mark a task as completed
    ```bash
    $ task c ID
    ```
* Mark a task as pending
    ```bash
    $ task p ID
    ```
* Modify a task task
    ```bash
    $ task m ID Watch Game of Thrones
    ```    
* Delete latest task
    ```bash
    $ task del
    ```
* Remove a specific task by id
    ```bash
    $ task r ID
    ```
* Flush/Delete all the tasks
    ```bash
    $ task flush
    ```
* To start the program as service (Note: Must use as service if you are using **reminder**)
    ```bash
    $ task srvice-start # Start a service
    $ task srvice-force-start # Forcefully start service
    $ task service-stop #stop service
    ```
    
##### Examples of reminder
```bash
$ task remind Take a cup of coffee in 30min
$ task remind Watch game of thrones season 7 today 8:30pm
$ task remind Watch despicable me 3 next friday at 3pm
$ task remind Bug fix of the docker and send PR next thursday
```

### Build yourself

Go to your $GOPATH/src and get the package
```bash
$ go get github.com/thedevsaddam/task
```

Install dependency management tool go [govendor](https://github.com/kardianos/govendor)
```bash
$ go get -u github.com/kardianos/govendor
```

To install dependencies go to project root and `$ cd vendor`
```bash
$ govendor sync
```

In unix system use
```bash
$ sh build.sh
```

### Roadmap
- [x] Add reminder
- [x] Desktop notifier
- [ ] Recurring reminder
- [ ] Reminder voice
- [ ] Code review

### Some awesome packages are used to make this awesome task :)
* [Notifier](github.com/0xAX/notificator)
* [Auto start/service](github.com/ProtonMail/go-autostart)
* [Color](github.com/fatih/color)
* [Natural date parser](github.com/olebedev/when)
* [Table writter](github.com/olekukonko/tablewriter)
* [Go prompt](github.com/segmentio/go-prompt)
* [Task manager](github.com/thedevsaddam/task/taskmanager)

### Contribution
There are some tasks that need to be done. I have tried to make a minimal setup, need more code refactoring, review, bug fixing and adding features.
If you are interested to make this application better please send pull requests.

### **License**
The **task** is a open-source software licensed under the [MIT License](LICENSE.md).