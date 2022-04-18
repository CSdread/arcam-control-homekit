package arcam

import (
	"context"
	"errors"
	"fmt"
)

type Receiver struct {
	model  string
	client *client
}

func NewReceiver(model string, ipAddress string, port int) (Receiver, error) {
	if _, ok := ReceiverModels[model]; !ok {
		return Receiver{}, errors.New(fmt.Sprintf("Invalid receiver model: %s", model))
	}

	netClient := newClient(ipAddress, port)
	return Receiver{
		model:  model,
		client: &netClient,
	}, nil
}

func (r *Receiver) Connect(ctx context.Context) error {
	fmt.Println("connecting to receiver")
	return r.client.connect(ctx)
}

func (r *Receiver) IsOn(ctx context.Context) (bool, error) {
	req := Request{
		Zone:    ZoneOne,
		Command: PowerCommand,
		Data:    []byte{0xF0},
	}

	err := r.client.send(req)
	if err != nil {
		return false, err
	}

	resp, err := r.client.read(ctx)
	if err != nil {
		return false, err
	}

	if resp.AnswerCode != StatusUpdate {
		return false, errors.New("")
	}

	isOn := int(resp.Data[0]) == 0x01

	return isOn, nil
}

func (r *Receiver) PowerOn(ctx context.Context) error {
	req := Request{
		Zone:    ZoneOne,
		Command: SimulateRC5IRCommand,
		Data:    []byte{PowerOn.Data1, PowerOn.Data2},
	}

	err := r.client.send(req)
	if err != nil {
		return err
	}

	resp, err := r.client.read(ctx)
	if err != nil {
		return err
	}

	if resp.AnswerCode != StatusUpdate {
		return errors.New("")
	}

	return nil
}

func (r *Receiver) PowerOff(ctx context.Context) error {
	req := Request{
		Zone:    ZoneOne,
		Command: SimulateRC5IRCommand,
		Data:    []byte{PowerOff.Data1, PowerOff.Data2},
	}
	err := r.client.send(req)
	if err != nil {
		return err
	}

	resp, err := r.client.read(ctx)
	if err != nil {
		fmt.Println(err)
	}

	if resp.AnswerCode != StatusUpdate {
		return errors.New("")
	}

	return nil
}

/*

func (r *Receiver) GetVolume(ctx context.Context) (int, error) {
	req := Request{
		Zone:    ZoneOne,
		Command: SetRequestVolume,
		Data:    []byte{0xF0},
	}

	resp, err := r.client.send(ctx, req)
	if err != nil {
		return -1, err
	}

	if resp.AnswerCode != StatusUpdate {
		return -1, errors.New(fmt.Sprintf("Failed to retrieve current volume: %x", resp.AnswerCode))
	}

	if len(resp.Data) != 1 {
		return -1, errors.New(fmt.Sprintf("Invalid response data: %s", resp.Data))
	}

	return int(resp.Data[0]), nil
}

func (r *Receiver) SetVolume(ctx context.Context, newVolume int) error {
	req := Request{
		Zone:    ZoneOne,
		Command: SetRequestVolume,
		Data:    []byte{byte(newVolume)},
	}

	resp, err := r.client.send(ctx, req)
	if err != nil {
		return err
	}

	if resp.AnswerCode != StatusUpdate {
		return errors.New(fmt.Sprintf("Invalid answer code: %x", resp.AnswerCode))
	}

	if int(resp.Data[0]) != newVolume {
		return errors.New("New volume failed to set")
	}

	return nil
}

func (r *Receiver) VolumeUp(ctx context.Context) error {
	currentVol, err := r.GetVolume(ctx)
	if err != nil {
		return err
	}

	if currentVol == 99 {
		return nil
	}

	return r.SetVolume(ctx, currentVol+1)
}
func (r *Receiver) VolumeDown(ctx context.Context) error {
	currentVol, err := r.GetVolume(ctx)
	if err != nil {
		return errors.New(fmt.Sprintf("could not fetch current volume: %s", err))
	}

	if currentVol == 0 {
		return nil
	}

	return r.SetVolume(ctx, currentVol-1)
}
func (r *Receiver) IsMuted(ctx context.Context) (bool, error) {
	req := Request{
		Zone:    ZoneOne,
		Command: RequestMuteStatus,
		Data:    []byte{0xF0},
	}

	resp, err := r.client.send(ctx, req)
	if err != nil {
		return false, err
	}

	if resp.AnswerCode != StatusUpdate {
		return false, errors.New("")
	}

	return resp.Data[0] == 0x00, nil
}

func (r *Receiver) ToggleMute(ctx context.Context) error {
	muted, err := r.IsMuted(ctx)
	if err != nil {
		return err
	}

	data := []byte{MuteOn.Data1, MuteOn.Data2}
	if muted {
		data = []byte{MuteOff.Data1, MuteOff.Data2}
	}

	req := Request{
		Zone:    ZoneOne,
		Command: SimulateRC5IRCommand,
		Data:    data,
	}

	resp, err := r.client.send(ctx, req)
	if err != nil {
		return err
	}

	if resp.AnswerCode != StatusUpdate {
		return errors.New("")
	}

	return nil
}

type NowPlaying struct {
	TrackTitle string
	Artist     string
	Album      string
}

func (r *Receiver) GetNowPlaying(ctx context.Context) (*NowPlaying, error) {
	dataCodes := []byte{0xF0, 0xF1, 0xF2}

	values := []string{}
	for code := range dataCodes {
		req := Request{
			Zone:    ZoneOne,
			Command: NowPlayingInfo,
			Data:    []byte{byte(code)},
		}

		resp, err := r.client.send(ctx, req)
		if err != nil {
			return nil, err
		}

		if resp.AnswerCode != StatusUpdate {
			return nil, errors.New("")
		}

		data := string(resp.Data)
		values = append(values, data)
	}

	result := NowPlaying{
		TrackTitle: values[0],
		Artist:     values[1],
		Album:      values[2],
	}

	return &result, nil
}

func (r *Receiver) GetSource(ctx context.Context) (InputSource, error) {
	req := Request{
		Zone:    ZoneOne,
		Command: RequestCurrentSource,
		Data:    []byte{0xF0},
	}

	resp, err := r.client.send(ctx, req)
	if err != nil {
		return 0x00, err
	}

	if resp.AnswerCode != StatusUpdate {
		return 0x00, errors.New("")
	}

	input := InputSource(resp.Data[0])
	return input, nil
}

func (r *Receiver) SetSource(ctx context.Context, source InputSource) error {
	rc5CommandDataLookup := map[InputSource]AVRC5CommandCode{
		InputCD:       CD,
		InputBD:       BD,
		InputAV:       AV,
		InputSAT:      Sat,
		InputPVR:      PVR,
		InputUHD:      UHD,
		InputAUX:      Aux,
		InputDISPLAY:  Display,
		InputTUNERFM:  FM,
		InputTUNERDAB: DAB,
		InputNET:      Net,
		InputSTB:      STB,
		InputGAME:     Game,
		InputBT:       BT,
	}

	data, ok := rc5CommandDataLookup[source]
	if !ok {
		return errors.New("Invalid Input Source")
	}

	req := Request{
		Zone:    ZoneOne,
		Command: SimulateRC5IRCommand,
		Data:    []byte{data.Data1, data.Data2},
	}

	resp, err := r.client.send(ctx, req)
	if err != nil {
		return err
	}

	if resp.AnswerCode != StatusUpdate {
		return errors.New("")
	}

	return nil
}
*/
