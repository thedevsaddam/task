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

### **License**
The **task** is a open-source software licensed under the [MIT License](LICENSE.md).