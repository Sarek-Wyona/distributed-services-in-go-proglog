package log

import (
	"bufio"
	"encoding/binary"
	"os"
	"sync"
)

// Defines the encoding to persist record sizes and index entries
var enc = binary.BigEndian

// Number of bytes used to store the records length
const lenWidth = 8


//The store struct is a simple wrapper around a file with two APIs to append and read bytes to and from the file.
type store struct {
	*os.File
	mu   sync.Mutex
	buf  *bufio.Writer // Append to the store
	size uint64
}


//The newStore(*os.File) function creates a store for the given file. The function calls os.Stat(name string) to get the
//file’s current size, in case we’re re-creating the store from a file that has existing data, which would happen if,
//for example, our service had restarted.
func newStore(f *os.File) (*store, error) {

	// Get the files current size
	fi, err := os.Stat(f.Name())
	if err != nil {
		return nil, err
	}
	size := uint64(fi.Size())

	// Return the file
	return &store{
		File: f,
		size: size,
		buf:  bufio.NewWriter(f),
	}, nil
}

// Append persists the given bytes to the store. We write the length of the record so that, when we read the
//record, we know how many bytes to read. We write to the buffered writer instead of directly to the file to reduce the
//number of system calls and improve performance. If a user wrote a lot of small records, this would help a lot. Then
//we return the number of bytes written, which similar Go APIs conventionally do, and the position where the store
//holds the record in its file. The segment will use this position when it creates an associated index entry for this
//record.
func (s *store) Append(p []byte) (n uint64, pos uint64, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Get the size of the file
	pos = s.size

	//Write the binary representation of data
	if err := binary.Write(s.buf, enc, uint64(len(p))); err != nil {
		return 0, 0, err
	}

	//Write the content onto a buffer
	w, err := s.buf.Write(p)
	if err != nil {
		return 0, 0, err
	}

	w += lenWidth       //New length
	s.size += uint64(w) //Update size
	return uint64(w), pos, nil
}

//Read returns the record stored at the given position. First it flushes the writer buffer, in case
//we’re about to try to read a record that the buffer has’t flushed to disk yet. We find out how many bytes we
//have to read to get the whole record, and then we fetch and return the record. The compiler allocates byte
//slices that don’t escape the functions they’re declared in on the stack. A value escapes when it lives beyond
//the lifetime of the function call—if you return the value, for example.
func (s *store) Read(pos uint64) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Flush out the buffer before we begin
	if err := s.buf.Flush(); err != nil {
		return nil, err
	}

	// ReadAt the offset provided by pos, for the length of size
	size := make([]byte, lenWidth)
	if _, err := s.File.ReadAt(size, int64(pos)); err != nil {
		return nil, err
	}

	// TODO why are we reading it twice here?
	// Read at offset+8
	b := make([]byte, enc.Uint64(size))
	if _, err := s.File.ReadAt(b, int64(pos+lenWidth)); err != nil {
		return nil, err
	}
	return b, nil
}

//ReadAt reads len(p) bytes into p beginning at the off offset in the store’s file.
//It implements io.ReaderAt on the store type.
func (s *store) ReadAt(p []byte, off int64) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.buf.Flush(); err != nil {
		return 0, err
	}
	return s.File.ReadAt(p, off)
}

// Close persists any buffered data before closing the file.
func (s *store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.buf.Flush()
	if err != nil {
		return err
	}
	return s.File.Close()
}
