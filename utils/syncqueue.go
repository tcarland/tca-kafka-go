/** Synchronized Queue structure for sharing messages across threads.
  * 
  * Author: Timothy C. Arland <tcarland at gmail dot com>
 **/
package utils

import "sync"


type SyncQueue struct {
    queue   []interface{}
    lock     *sync.Mutex
}

// ----------------------------------------------

func NewSyncQueue() *SyncQueue {
    return new(SyncQueue).InitSyncQueue()
}


func (q *SyncQueue) InitSyncQueue() *SyncQueue {
    q.queue = make([]interface{}, 1)
    q.Pop()
    q.lock  = &sync.Mutex{}
    return q
}


func (q *SyncQueue) Pop() interface{} {
    if q.Empty() {
        return nil
    }
    
    rec    := q.queue[0]
    q.queue = q.queue[1:]

    return rec
}


func (q *SyncQueue) Push(val interface{}) {
    q.queue = append(q.queue, val)
}


func (q *SyncQueue) Reset() {
    q.queue = q.queue[:0]
}


func (q *SyncQueue) Lock() {
    q.lock.Lock()
}


func (q *SyncQueue) Unlock() {
    q.lock.Unlock()
}


func (q *SyncQueue) Empty() bool {
    return len(q.queue) == 0
}


func (q *SyncQueue) Size() int {
    return len(q.queue)
}
