package log

import (
	api "github.com/Sarek-Wyona/proglog/api/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"os"
	"testing"
)

//testAppendRead(*testing.T, *log.Log) tests that we can successfully append to and read from the log.
//When we append a record to the log, the log returns the offset it associated that record with. So,
//when we ask the log for the record at that offset, we expect to get the same record that we appended.
func testAppendRead(t *testing.T, log *Log) {
	appendValue := &api.Record{Value: []byte("hello world")}
	off, err := log.Append(appendValue)
	require.NoError(t, err)
	require.Equal(t, uint64(0), off)

	read, err := log.Read(off)
	require.NoError(t, err)
	require.Equal(t, appendValue.Value, read.Value)
}

//testOutOfRangeErr(*testing.T, *log.Log) tests that the log returns an error when we try to read an offset
//that’s outside of the range of offsets the log has stored.
func testOutOfRangeErr(t *testing.T, log *Log) {
	read, err := log.Read(1)
	require.Nil(t, read)
	apiErr := err.(api.ErrOffsetOutOfRange)
	require.Error(t, apiErr)
}

//testInitExisting(*testing.T, *log.Log) tests that when we create a log, the log bootstraps itself from the data
//stored by prior log instances. We append three records to the original log before closing it. Then we create a
//new log configured with the same directory as the old log. Finally, we confirm that the new log set itself up
//from the data stored by the original log.
func testInitExisting(t *testing.T, log *Log) {
	appendValue := &api.Record{Value: []byte("hello world")}
	for i := 0; i < 3; i++ {
		_, err := log.Append(appendValue)
		require.NoError(t, err)
	}
	require.NoError(t, log.Close())

	off, err := log.LowestOffset()
	require.NoError(t, err)
	require.Equal(t, uint64(0), off)

	off, err = log.HighestOffset()
	require.NoError(t, err)
	require.Equal(t, uint64(2), off)

	n, err := NewLog(log.Dir, log.Config)
	require.NoError(t, err)

	off, err = n.LowestOffset()
	require.NoError(t, err)
	require.Equal(t, uint64(0), off)
	off, err = n.HighestOffset()
	require.NoError(t, err)
	require.Equal(t, uint64(2), off)
}

//testReader(*testing.T, *log.Log) tests that we can read the full, raw log as it’s stored on disk so that we can
//snapshot and restore the logs in Finite-State Machine
func testReader(t *testing.T, log *Log) {
	appendValue := &api.Record{Value: []byte("hello world"),}
	off, err := log.Append(appendValue)
	require.NoError(t, err)
	require.Equal(t, uint64(0), off)

	reader := log.Reader()
	b, err := ioutil.ReadAll(reader)
	require.NoError(t, err)

	read := &api.Record{}
	err = proto.Unmarshal(b[lenWidth:], read)
	require.NoError(t, err)
	require.Equal(t, appendValue.Value, read.Value)

}

//testTruncate(*testing.T, *log.Log) tests that we can truncate the log and remove old segments that we don’t
//need any more.
func testTruncate(t *testing.T, log *Log)  {
	appendValue := &api.Record{Value: []byte("hello world")}

	for i :=0; i <3; i++ {
		_, err := log.Append(appendValue)
		require.NoError(t, err)
	}

	err := log.Truncate(1)
	require.NoError(t, err)

	_, err = log.Read(0)
	require.Error(t, err)

	_, err = log.Read(2)
	require.NoError(t, err)

}

//TestLog(*testing.T) defines a table of tests to test the log. Table is used to write the log tests so we
//don’t have to repeat the code that creates a new log for every test case.
func TestLog(t *testing.T) {
	for scenario, fn := range map[string]func(
		t *testing.T, log *Log,
	){
		"append and read a record succeeds": testAppendRead,
		"offset our of range error":         testOutOfRangeErr,
		"init with existing segments":       testInitExisting,
		"reader":                            testReader,
		"turncate":                          testTruncate,
	} {
		t.Run(scenario, func(t *testing.T) {
			dir, err := ioutil.TempDir("", "store-test")
			require.NoError(t, err)
			defer os.RemoveAll(dir)

			c := Config{}
			c.Segment.MaxStoreBytes = 32
			log, err := NewLog(dir, c)
			require.NoError(t, err)

			fn(t, log)
		})
	}

}
