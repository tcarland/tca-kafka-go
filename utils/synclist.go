/** Synchronized list structure for storing a kafka message list
  *
  * Author: Timothy C. Arland <tcarland at gmail dot com>
 **/
package utils

import "sync"

type SyncList struct {
    list  []interface{}
    lock    sync.Mutex
}


func NewSyncList() *SyncList {
    return new(SyncList).InitSyncList()
}


func (l *SyncList) InitSyncList() *SyncList {
    l.list = make([]interface{}, 0)
    return l
}


func (l *SyncList) PushBack(val interface{}) {
    l.list = append(l.list, val)
}


func (l *SyncList) Front() interface{} {
    if l.Empty() {
        return nil
    }

    return l.list[0]
}


func (l *SyncList) Back() interface{} {
    if l.Empty() {
        return nil
    }

    return l.list[l.Size() - 1]
}


func (l *SyncList) PopFront() interface{} {
    if l.Empty() {
        return nil
    }
    rec    := l.list[0]
    l.list  = l.list[1:]
    return rec
}


func (l *SyncList) PopBack() interface{} {
    if l.Empty() {
        return nil
    }
    last   := l.Size() - 1
    rec    := l.list[last]
    l.list  = l.list[:last]

    return rec
}


func (l *SyncList) At(pos int) interface{} {
    if l.Empty() || pos >= l.Size() {
        return nil
    }
    return l.list[pos]
}


func (l *SyncList) Empty() bool {
    return len(l.list) == 0
}


func (l *SyncList) Size() int {
    return len(l.list)
}


func (l *SyncList) Clear() {
    l.list = l.list[:0]
}


func (l *SyncList) Lock() {
    l.lock.Lock()
}


func (l *SyncList) Unlock() {
    l.lock.Unlock()
}
