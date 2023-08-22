package utils

import (
    "testing"
)


func TestList_PushBack(t *testing.T) {
    t.Parallel()
    testCases := []struct {
        name      string
        msg       string
        msgcnt    int
    }{
        {"Add 1st message", "this is a test message", 1},
        {"Add 2nd message", "this is another test", 2},
    }

    l := SyncList{}
    l.InitSyncList()

    l.Lock()
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            l.PushBack(tc.msg)
            sz := l.Size()
            if sz != tc.msgcnt {
                t.Errorf("Expecting %v items but got: %v", tc.msgcnt, sz)
            }
        })
    }
    l.Unlock()
}

func TestList_PopBack(t *testing.T) {
    t.Parallel()

    testCases := []struct {
        name    string
        excnt   int
    }{
        {"Pop last pos", 3},
        {"Pob back again", 2},
    }

    l := NewSyncList()

    l.PushBack(1)
    l.PushBack(2)
    l.PushBack(3)
    l.PushBack(4)

    l.Lock()
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            v := l.PopBack()
            if l.Size() != tc.excnt || v != (l.Size() + 1) {
                t.Errorf("Invalid resulting list size: %v, expected %v, val=%v", l.Size(), tc.excnt, v) 
            }
        })
    }
    l.Unlock()
}

func TestList_PopFront(t *testing.T) {
    t.Parallel()

    testCases := []struct {
        name     string
        excnt    int
    }{
        {name: "Pop 1st message", excnt: 2},
        {name: "Pop 2nd message", excnt: 1},
    }

    l := SyncList{}
    l.InitSyncList()
    l.PushBack("First Message")
    l.PushBack("Second Message")
    l.PushBack("Third Message")

    l.Lock()
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            l.PopFront()
            if l.Size() != tc.excnt {
                t.Errorf("Invalid number of items in the list: %v", l.Size())
            }
        })
    }
    l.Unlock()
}


func TestList_Reset(t *testing.T) {
    t.Parallel()

    testCases := []struct {
        name  string
        excnt int
    }{
        {name: "Reset Queue to zero", excnt: 0},
    }

    l := NewSyncList()
    l.PushBack("First message")
    l.PushBack("Second message")
    l.PushBack("Third message")
    l.PushBack("Fourth message")
    l.PushBack("Fifth message")

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            l.Clear()
            if l.Size() != tc.excnt {
                t.Errorf("Invalid queue count, not 0, size= %v", l.Size())
            }
        })
    }
}
