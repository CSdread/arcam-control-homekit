package arcam

import (
	"context"
	"errors"
	"fmt"
	"net"
)

type client struct {
	ipAddress string
	port      int
}

func newClient(ipAddress string, port int) client {
	return client{ipAddress: ipAddress, port: port}
}

func (r *client) send(ctx context.Context, request Request) (*Response, error) {
	var d net.Dialer

	conn, err := d.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", r.ipAddress, r.port))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("could not connect to reciever: %s", err))
	}

	defer conn.Close()

	out := []byte{
		TransmissionStart,
		byte(request.Zone),
		byte(request.Command),
		byte(len(request.Data)),
	}
	out = append(out, request.Data[:]...)
	out = append(out, TransmissionEnd)

	_, err = conn.Write(out)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("could not write to the socket: %s", err))
	}

	var response []byte
	respLen, err := conn.Read(response)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to read response: %s", err))
	}

	if respLen == 0 {
		return nil, errors.New("zero len response")
	}

	dataLen := int(response[4])

	return &Response{
		Zone:        ZoneNumber(response[1]),
		CommandCode: Command(response[2]),
		AnswerCode:  Answer(response[3]),
		Data:        response[5 : dataLen-2],
	}, nil
}
