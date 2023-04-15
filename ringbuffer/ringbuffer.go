package ringbuffer

import (
	"errors"
	"fmt"

	"github.com/tmthrgd/go-memset"
)

// Single-Threaded Ring buffer
type RingBuffer struct {
	data              []byte
	len               int
	readPtr, writePtr int
}

func (buffer *RingBuffer) Initialize(len int) {
	buffer.readPtr = 0
	buffer.writePtr = 0
	buffer.len = len
	buffer.data = make([]byte, buffer.len)
}

func (buffer *RingBuffer) SpaceAvailable() int {
	if buffer.readPtr == buffer.writePtr {
		return len(buffer.data)
	} else {
		if buffer.writePtr > buffer.readPtr {
			return len(buffer.data) - (buffer.writePtr - buffer.readPtr)
		} else {
			return len(buffer.data) - (buffer.readPtr - buffer.writePtr)
		}
	}
}

func (buffer *RingBuffer) Size() int {
	if buffer.readPtr == buffer.writePtr {
		return 0
	}
	if buffer.writePtr < buffer.readPtr {
		return buffer.len - buffer.readPtr + buffer.writePtr
	} else {
		return buffer.writePtr - buffer.readPtr
	}
}

/*
Write tries to write data of length len to the buffer only if complete write is possible.
In cases of partial write, it does not write anything and returns error
*/
func (buffer *RingBuffer) Write(data []byte) error {
	if len(data) == 0 {
		return errors.New("no data passed to be written to buffers")
	}
	if len(data) > buffer.SpaceAvailable() {
		return errors.New("not enough space left to write in the buffer")
	}

	if buffer.writePtr > buffer.readPtr &&
		len(data)+buffer.writePtr > buffer.len {
		c1 := buffer.len - buffer.writePtr
		//Have to do 2 copies.
		copy(buffer.data[buffer.writePtr:], data[:c1])
		//c2 := len(data)- c1
		copy(buffer.data, data[c1:])
		buffer.writePtr = len(data) + buffer.writePtr - buffer.len
	} else {
		copy(buffer.data[buffer.writePtr:], data)
		buffer.writePtr += len(data)
	}
	return nil
}

/*
Reads atmost len bytes from buffer
Returns
- Data read from the buffer
- Length of the data read
- Error in case of no data or other errors
*/
func (buffer *RingBuffer) Read(len int) ([]byte, error) {
	bufSize := buffer.Size()
	fmt.Printf("Request to read %d bytes . Buffer Size is %d \n", len, bufSize)

	if bufSize == 0 {
		return nil, errors.New("buffer is empty")
	}
	sizeToRead := len
	if len > bufSize {
		sizeToRead = bufSize
	}
	data := make([]byte, sizeToRead)
	if buffer.readPtr > buffer.writePtr &&
		sizeToRead+buffer.readPtr > buffer.len {
		bytesCopied := copy(data, buffer.data[buffer.readPtr:])
		memset.Memset(buffer.data[buffer.readPtr:], 0)
		copyTillIndex := sizeToRead - bytesCopied
		copy(data[bytesCopied:], buffer.data[:copyTillIndex])
		memset.Memset(buffer.data[:copyTillIndex], 0)
		buffer.readPtr = copyTillIndex
	} else {
		readTill := buffer.readPtr + sizeToRead
		bufSlice := buffer.data[buffer.readPtr:readTill]
		copy(data, bufSlice)
		memset.Memset(bufSlice, 0)
		buffer.readPtr = buffer.readPtr + sizeToRead
	}

	return data, nil

}

func (buffer *RingBuffer) Print() {
	fmt.Printf("Buffer Struct contents %+v \n", buffer)
}
