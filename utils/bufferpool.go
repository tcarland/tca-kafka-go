/** A BufferPool channel of byte buffers for use between goroutines 
  *
  * Author: Timothy C. Arland <tcarland at gmail dot com>
 **/
package utils

import (
    "bytes"
)

type BufferPool struct {
    c chan *bytes.Buffer
}


func NewBufferPool(size int) (pool *BufferPool) {
    return &BufferPool{ c: make(chan *bytes.Buffer, size) }
}

func (pool *BufferPool) Get() (buffer *bytes.Buffer) {
    select {
    case buffer = <-pool.c:
    default:
        buffer = bytes.NewBuffer([]byte{})
    }
    return
}

func (pool *BufferPool) Put(buffer *bytes.Buffer) {
    buffer.Reset()
    select {
    case pool.c <- buffer:
    default:
    }
}

func (pool *BufferPool) Size() int {
    return len(pool.c)
}