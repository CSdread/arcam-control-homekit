package arcam

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/brutella/hap/log"
)

type EventHandler func(ZoneNumber, []byte) error

type Receiver struct {
	model     string
	client    *client
	connected bool

	messages      chan *Response
	eventHandlers map[Command]EventHandler
}

// NewReceiver creates and initializes a new receiver
func NewReceiver(model string, ipAddress string, port int) (Receiver, error) {
	if _, ok := ReceiverModels[model]; !ok {
		return Receiver{}, errors.New(fmt.Sprintf("Invalid receiver model: %s", model))
	}

	netClient := newClient(ipAddress, port)
	return Receiver{
		model:         model,
		client:        &netClient,
		connected:     false,
		eventHandlers: map[Command]EventHandler{},
	}, nil
}

// RegisterEventHandler registers handler methods for responding to commands sent
// from the receiver
func (r *Receiver) RegisterEventHandler(command Command, handler EventHandler) {
	r.eventHandlers[command] = handler
}

// Connect starts the connection to the receiver
func (r *Receiver) Connect(ctx context.Context) error {
	log.Debug.Println("connecting to receiver")
	err := r.client.connect(ctx)
	if err != nil {
		return err
	}

	r.messages = make(chan *Response, 10)

	log.Debug.Println("starting poll go routine")
	go r.startPolling(ctx)

	log.Debug.Println("starting message processor")
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case response, ok := <-r.messages:
				if !ok {
					return
				}
				r.processMessage(ctx, response)
			}
		}
	}()

	time.Sleep(10 * time.Second)

	r.connected = true
	return nil
}

// GetInputs retrieves all possible source inputs for the receiver
func (r *Receiver) GetInputs() []InputSource {
	sources := []InputSource{}
	for k := range InputDisplayNameMap {
		sources = append(sources, k)
	}
	return sources
}

func (r *Receiver) RefreshState(ctx context.Context) {
	log.Debug.Println("trigering state refresh")
	commands := []Command{
		PowerCommand,
		RequestCurrentSource,
		SetRequestVolume,
		RequestMuteStatus,
		RequestDirectModeStatus,
	}
	req := Request{
		Zone: ZoneOne,
		Data: []byte{0xF0},
	}
	for _, command := range commands {
		req.Command = command
		err := r.client.send(req)
		// need to take a second between requests
		time.Sleep(1 * time.Second)
		if err != nil {
			// log failure but continue on
			log.Debug.Printf("ERROR: refresh state request failed for %v", req.Command)
		}
	}
}

func (r *Receiver) PowerOn(ctx context.Context) error {
	req := Request{
		Zone:    ZoneOne,
		Command: SimulateRC5IRCommand,
		Data:    []byte{PowerOn.Data1, PowerOn.Data2},
	}

	return r.client.send(req)
}

func (r *Receiver) PowerOff(ctx context.Context) error {
	req := Request{
		Zone:    ZoneOne,
		Command: SimulateRC5IRCommand,
		Data:    []byte{PowerOff.Data1, PowerOff.Data2},
	}
	return r.client.send(req)
}

func (r *Receiver) EnableDirectMode(ctx context.Context) error {
	data := []byte{DirectModeOn.Data1, DirectModeOn.Data2}

	req := Request{
		Zone:    ZoneOne,
		Command: SimulateRC5IRCommand,
		Data:    data,
	}

	return r.client.send(req)
}

func (r *Receiver) DisableDirectMode(ctx context.Context) error {
	data := []byte{DirectModeOff.Data1, DirectModeOff.Data2}

	req := Request{
		Zone:    ZoneOne,
		Command: SimulateRC5IRCommand,
		Data:    data,
	}

	return r.client.send(req)
}

func (r *Receiver) UnMute(ctx context.Context) error {
	data := []byte{MuteOff.Data1, MuteOff.Data2}

	req := Request{
		Zone:    ZoneOne,
		Command: SimulateRC5IRCommand,
		Data:    data,
	}

	return r.client.send(req)
}

func (r *Receiver) Mute(ctx context.Context) error {
	data := []byte{MuteOn.Data1, MuteOn.Data2}

	req := Request{
		Zone:    ZoneOne,
		Command: SimulateRC5IRCommand,
		Data:    data,
	}

	return r.client.send(req)
}

func (r *Receiver) SetSource(ctx context.Context, source InputSource) error {
	data, ok := InputSourceCommandMap[source]
	if !ok {
		return errors.New("Invalid Input Source")
	}

	req := Request{
		Zone:    ZoneOne,
		Command: SimulateRC5IRCommand,
		Data:    []byte{data.Data1, data.Data2},
	}

	return r.client.send(req)
}

func (r *Receiver) SetVolume(ctx context.Context, newVolume int) error {
	req := Request{
		Zone:    ZoneOne,
		Command: SetRequestVolume,
		Data:    []byte{byte(newVolume)},
	}

	return r.client.send(req)
}

func (r *Receiver) startPolling(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			close(r.messages)
			return
		default:
			message, err := r.client.read(ctx)
			if err != nil {
				log.Debug.Fatalln("error reading from socket")
			}
			r.messages <- message
		}
	}
}

func (r *Receiver) processMessage(ctx context.Context, message *Response) error {
	if handler, ok := r.eventHandlers[message.CommandCode]; ok {
		return handler(message.Zone, message.Data)
	}

	// if we dont have a handler for the Command then we no-op,
	// its a message we dont care about
	return nil
}
