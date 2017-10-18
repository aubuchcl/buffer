package capbuff

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

type Buffer struct {
	*bytes.Buffer
	cap   int
	mutex sync.Mutex
}

// object that has a bytes.Buffer pointer, an int, and a sync.Mutex
//pointer to the bytes buffer is so you can change the size
//mutex is something that handles goroutines

var ErrorBufferFull = fmt.Errorf("Buffer is full")

func NewBuffer(cap int) (b *Buffer) {
	b = &Buffer{
		//this syntax actually makes a new buffer that returns the *Buffer
		//not sure but I think the 0 is the init value and the cap is max value
		Buffer: bytes.NewBuffer(make([]byte, 0, cap)),
		cap:    cap,
	}
	// return the new buffer*
	return b
}

//call if you need to see how much room is left in the buffer
func (b *Buffer) unused() int {
	return b.cap - b.Buffer.Len()
}

//
func (b *Buffer) Write(p []byte) (n int, err error) {
	//lock the mutex or block until its open
	b.mutex.Lock()
	//see write func below
	n, err = b.write(p)
	// Unlock the locked mutex
	b.mutex.Unlock()

	return n, err
}

//exported way of accessing the cap value for a Buffer
func (b *Buffer) Cap() (n int) {
	return b.cap
}

//change the value of cap by passing the memory address of the buffer
func (b *Buffer) Grow(cap int) {
	b.mutex.Lock()
	// grow the associate int
	b.cap = b.cap + cap
	// grow the actual bytes.Buffer capacity
	b.Buffer.Grow(cap)
	b.mutex.Unlock()
	//do you lock and unlock these so they arnt used until the updates are finished?
}

func (b *Buffer) write(p []byte) (n int, err error) {
	//if there isnt enough room use Write to expand b by the appropriate amount
	if len(p) > b.unused() {
		//append p to bytes buffer
		n, err = b.Buffer.Write(p[0:b.unused()])

		if err != nil {
			return n, err
		} else {
			return n, ErrorBufferFull
		}
	}
	//else just do it
	return b.Buffer.Write(p)
}

func (b *Buffer) ReadFrom(r io.Reader) (n int64, err error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	//create a byte slice that is the size of the unused bytes in b
	//and points at teh firs b.unused() values

	alloc := make([]byte, b.unused(), b.unused())
	//if there is no error reading the data from r
	//push it into alloc and then return
	//n will represent how many bytes into r was read
	//this will error if there isnt enough room in alloc

	if n, err := r.Read(alloc); err != nil {
		// returns bytes read and what the error was
		return int64(n), err
	}
	//if there was no error in the first if statement
	//try to write b into alloc

	if n, err := b.write(alloc); err != nil {
		return int64(n), err
	} else {
		// this only returns if both the above statements dont return errors

		return int64(n), nil
	}
	//return int64(n), nil wouldnt work here although suggested because n would be out of scope here?
}

//convert a string to bytes and write it to the bytes.buffer
//appends it to the end.
func (b *Buffer) WriteString(s string) (n int, err error) {
	b.mutex.Lock()
	n, err = b.write([]byte(s))
	b.mutex.Unlock()

	return n, err
}
