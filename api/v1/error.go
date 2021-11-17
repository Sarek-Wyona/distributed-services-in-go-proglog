package log_v1

import (
	"fmt"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

//ErrOffsetOutOfRange is custom error that the server will send back to the client when the client tries to consume
//an offset that’s outside of the log
type ErrOffsetOutOfRange struct {
	Offset uint64
}

func (e ErrOffsetOutOfRange) GRPCStatus()  *status.Status {
	st := status.New(404, fmt.Sprintf("offset out of range: %d", e.Offset))
	msg := fmt.Sprintf("The requested offset is outside the log's range: %d", e.Offset)

	d := &errdetails.LocalizedMessage{Locale: "en-US", Message: msg}
	std, err := st.WithDetails(d)
	if err != nil{
		return st
	}
	return std

}

func (e ErrOffsetOutOfRange) Error() string  {
	return e.GRPCStatus().Err().Error()

}