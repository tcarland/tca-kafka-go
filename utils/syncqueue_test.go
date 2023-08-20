package utils

import (
    "testing"
)


func TestQueue_Push(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name      string
		msg       string
		msgcnt    int
	}{
		{"Add 1st message", "this is a test message", 1},
		{"Add 2nd message", "this is another test", 2},
	}

	q := SyncQueue{}
    q.InitSyncQueue()

    q.Lock()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q.Push(tc.msg)

			sz := q.Size()
			if sz != tc.msgcnt {
				t.Errorf("Expecting %v items but got: %v", tc.msgcnt, sz)
			}
		})
	}
    q.Unlock()
}


func TestQueue_Pop(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		excnt    int
	}{
		{name: "Pop 1st message", excnt: 2},
		{name: "Pop 2nd message", excnt: 1},
	}

	q := SyncQueue{}
    q.InitSyncQueue()
    q.Push("First Message")
    q.Push("Second Message")
    q.Push("Third Message")

    q.Lock()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q.Pop()
			if q.Size() != tc.excnt {
				t.Errorf("Invalid number of items in the queue: %v", q.Size())
			}
		})
	}
    q.Unlock()
}


func TestQueue_Reset(t *testing.T) {
    t.Parallel()

    testCases := []struct {
        name  string
        excnt int
    }{
        {name: "Reset Queue to zero", excnt: 0},
    }

    q := NewSyncQueue()
    q.Push("First message")
    q.Push("Second message")
    q.Push("Third message")
    q.Push("Fourth message")
    q.Push("Fifth message")

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            q.Reset()
            if q.Size() != tc.excnt {
                t.Errorf("Invalid queue count, not 0, size= %v", q.Size())
            }
        })
    }
}
