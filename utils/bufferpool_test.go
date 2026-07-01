package utils

import (
    "bytes"
    "testing"
)


func TestBufferPool_Put(t *testing.T) {
    t.Parallel()
    testCases := []struct {
        name    string
        msg     string
        bufcnt  int
    }{
        {"Put 1st buffer", "this is a test message", 1},
        {"Put 2nd buffer", "this is another test", 2},
    }

    pool := NewBufferPool(4)

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            buf := bytes.NewBufferString(tc.msg)
            pool.Put(buf)

            sz := pool.Size()
            if sz != tc.bufcnt {
                t.Errorf("Expecting %v buffers but got: %v", tc.bufcnt, sz)
            }
        })
    }
}


func TestBufferPool_Get(t *testing.T) {
    t.Parallel()

    testCases := []struct {
        name    string
        excnt   int
    }{
        {name: "Get 1st buffer", excnt: 2},
        {name: "Get 2nd buffer", excnt: 1},
    }

    pool := NewBufferPool(4)
    pool.Put(bytes.NewBufferString("First Buffer"))
    pool.Put(bytes.NewBufferString("Second Buffer"))
    pool.Put(bytes.NewBufferString("Third Buffer"))

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            buf := pool.Get()
            if buf == nil {
                t.Errorf("Get returned a nil buffer")
            }
            if pool.Size() != tc.excnt {
                t.Errorf("Invalid number of buffers in the pool: %v", pool.Size())
            }
        })
    }
}


func TestBufferPool_GetEmpty(t *testing.T) {
    t.Parallel()

    testCases := []struct {
        name   string
        excnt  int
    }{
        {name: "Get from empty pool allocates a new buffer", excnt: 0},
    }

    pool := NewBufferPool(4)

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            buf := pool.Get()
            if buf == nil {
                t.Errorf("Get returned a nil buffer")
            }
            if buf.Len() != 0 {
                t.Errorf("Expected an empty buffer but got length: %v", buf.Len())
            }
            if pool.Size() != tc.excnt {
                t.Errorf("Invalid pool count, expected %v, size= %v", tc.excnt, pool.Size())
            }
        })
    }
}


func TestBufferPool_PutResets(t *testing.T) {
    t.Parallel()

    testCases := []struct {
        name   string
        exlen  int
    }{
        {name: "Put resets the buffer contents", exlen: 0},
    }

    pool := NewBufferPool(4)

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            buf := pool.Get()
            buf.WriteString("some data to be discarded on Put")
            pool.Put(buf)

            reused := pool.Get()
            if reused.Len() != tc.exlen {
                t.Errorf("Expected reset buffer of length %v but got: %v", tc.exlen, reused.Len())
            }
        })
    }
}


func TestBufferPool_PutFull(t *testing.T) {
    t.Parallel()

    testCases := []struct {
        name   string
        excnt  int
    }{
        {name: "Put on a full pool discards the buffer", excnt: 2},
    }

    pool := NewBufferPool(2)
    pool.Put(bytes.NewBufferString("First Buffer"))
    pool.Put(bytes.NewBufferString("Second Buffer"))

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            pool.Put(bytes.NewBufferString("Overflow Buffer"))
            if pool.Size() != tc.excnt {
                t.Errorf("Invalid pool count, expected %v, size= %v", tc.excnt, pool.Size())
            }
        })
    }
}
