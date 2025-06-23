package main

import "sync"

var users = make(map[string]User)
var complaints = make(map[string]Complaint)
var mu sync.Mutex