package Common

import "sync"

var Users = make(map[string]User)

var Complaints = make(map[string]Complaint)

var Mu sync.Mutex
