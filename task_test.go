package main

import (
	"github.com/thedevsaddam/task/taskmanager"
)

func Example_showTask() {
	showTask(taskmanager.Task{
		Id:          1,
		UID:         "213e9bb0-79e8-4647-8902-8421271e1809",
		Description: "Watch Pirates of the Caribbean: Dead Men Tell No Tales",
		Tag:         "low",
		Created:     "Fri, 07/21/17, 12:13PM",
		Updated:     "Fri, 07/21/17, 12:15PM",
		Completed:   "Fri, 07/21/17, 24:10PM",
	})
	//output:
	//
	//Task Details view
	//--------------------------------
	//ID: 1
	//UID: 213e9bb0-79e8-4647-8902-8421271e1809
	//Description: Watch Pirates of the Caribbean: Dead Men Tell No Tales
	//Tag: low
	//Created: Fri, 07/21/17, 12:13PM
	//Updated: Fri, 07/21/17, 12:15PM
	//
}
