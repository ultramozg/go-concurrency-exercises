//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"sync"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
	mut       sync.Mutex
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	done := make(chan struct{})
	var timeLast int64
	if !u.IsPremium {
		timeLast = 10
	} else {
		u.mut.Lock()
		timeLast = 10 - u.TimeUsed
		u.mut.Unlock()
	}

	go func() {
		start := time.Now()
		process()
		complete := time.Since(start)
		if u.IsPremium {
			u.mut.Lock()
			u.TimeUsed += int64(complete)
			u.mut.Unlock()
		}
		done <- struct{}{}
	}()

	select {
	case <-time.After((time.Duration)(timeLast) * time.Second):
		return false
	case <-done:
		return true
	}
	return true
}

func main() {
	RunMockServer()
}
