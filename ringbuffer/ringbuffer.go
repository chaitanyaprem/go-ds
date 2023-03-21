package ringbuffer

import "errors"

// Single-Threaded Ring buffer
type RingBuffer struct {
	data       []byte
	len        int
	head, tail int
}

func (buffer *RingBuffer) Initialize(len int) {
	buffer.head = 0
	buffer.tail = 0
	buffer.len = len
	buffer.data = make([]byte, buffer.len)
}

func (buffer *RingBuffer) SpaceAvailable() int {
	if buffer.head == buffer.tail {
		return len(buffer.data)
	} else {
		if buffer.tail > buffer.head {
			return len(buffer.data) - (buffer.tail - buffer.head)
		} else {
			return len(buffer.data) - (buffer.head - buffer.tail)
		}
	}
}

func (buffer *RingBuffer) Size() int {
	if buffer.head == buffer.tail {
		return 0
	}
	if buffer.tail > buffer.head {
		return buffer.len - buffer.tail + buffer.head
	} else {
		return buffer.tail - buffer.head
	}
}

/*
Write tries to write data of length len to the buffer only if complete write is possible.
In cases of partial write, it does not write anything and returns error
*/
func (buffer *RingBuffer) write(data []byte) error {
	if len(data) == 0 {
		return errors.New("No data passed to be written to buffers")
	}
	if len(data) > buffer.SpaceAvailable() {
		return errors.New("not enough space left to write in the buffer")
	}
	if buffer.head == buffer.tail {
		copy(buffer.data, data)
		buffer.tail = len(data)
	} else {
		if buffer.tail > buffer.head {
			if len(data)+buffer.tail > buffer.len {
				c1 := buffer.len - buffer.tail
				//Have to do 2 copies.
				copy(buffer.data[buffer.tail:], data[:c1])
				//c2 := len(data)- c1
				copy(buffer.data, data[c1:])
				buffer.tail = len(data) + buffer.tail - buffer.len
			} else {
				copy(buffer.data[buffer.tail:], data)
				buffer.tail += len(data)
			}
		} else {
			copy(buffer.data[buffer.tail:], data)
			buffer.tail += len(data)
		}
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
func (buffer *RingBuffer) read(len int) ([]byte, error) {
	bufSize := buffer.Size()
	if bufSize == 0 {
		return nil, errors.New("buffer is empty")
	}
	sizeToRead := len
	if len > bufSize {
		sizeToRead = bufSize
	}
	data := make([]byte, sizeToRead)
	if buffer.head > buffer.tail {

	} else {

	}
	//copy()

	return data, nil
}

func main() {

	var buffer RingBuffer
	buffer.Initialize(100)

	data := []byte{102, 97, 108, 99, 111, 110}

	buffer.write(data)

}
