package arcam

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/gorilla/websocket"
)

type client struct {
	ipAddress string
	port      int
	socket    *websocket.Conn
}

func newClient(ipAddress string, port int) client {
	return client{ipAddress: ipAddress, port: port}
}

func (r *client) connect(ctx context.Context) error {
	url := url.URL{Scheme: "ws", Host: fmt.Sprintf("%s:%d", r.ipAddress, r.port)}
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, url.String(), nil)
	if err != nil {
		return err
	}

	r.socket = conn

	return nil
}

func (r *client) disconnect() {
	r.socket.Close()
}

func (r *client) send(request Request) error {
	out := []byte{
		TransmissionStart,
		byte(request.Zone),
		byte(request.Command),
		byte(len(request.Data)),
	}
	out = append(out, request.Data[:]...)
	out = append(out, TransmissionEnd)

	fmt.Printf("TX: [%x]\n", out)

	err := r.socket.WriteMessage(websocket.BinaryMessage, out)
	if err != nil {
		return errors.New(fmt.Sprintf("could not write to the socket: %s", err))
	}

	return nil
}

func (r *client) read(ctx context.Context) (*Response, error) {
	_, response, err := r.socket.ReadMessage()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to read response: %s", err))
	}

	respLen := len(response)

	if respLen == 0 {
		return nil, errors.New("zero len response")
	}

	fmt.Printf("RX: [%x]\n", response)

	dataLen := int(response[4])

	return &Response{
		Zone:        ZoneNumber(response[1]),
		CommandCode: Command(response[2]),
		AnswerCode:  Answer(response[3]),
		Data:        response[5 : dataLen+5],
	}, nil
}
