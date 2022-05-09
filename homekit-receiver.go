package main

import (
	"context"
	"net/http"

	"github.com/brutella/hap/accessory"
	"github.com/brutella/hap/characteristic"
	"github.com/brutella/hap/service"
)

type Source struct {
	name        byte
	displayName string
}

type HomekitReceiver struct {
	model   string
	sources []Source

	Bridge *accessory.Bridge
	Tv     *accessory.Television
}

type intCallback func(int, int, *http.Request)
type boolCallback func(oldVal, newVal bool, r *http.Request)

func NewHomekitReceiver(ctx context.Context, model string, sources []Source) *HomekitReceiver {
	receiver := HomekitReceiver{
		model:   model,
		sources: sources,
	}
	receiver.init(ctx)

	return &receiver
}

func (r *HomekitReceiver) init(ctx context.Context) {
	r.Bridge = accessory.NewBridge(accessory.Info{
		Name:         "ARCAM Receiver Bridge",
		Firmware:     Version,
		Model:        "ARCAM Receiver Bridge",
		Manufacturer: "CSdread",
	})

	tv := accessory.NewTelevision(accessory.Info{
		Name:         "ARCAM Receiver",
		SerialNumber: Version,
		Model:        r.model,
		Manufacturer: "Arcam",
	})
	tv.Television.ConfiguredName.SetValue("Arcam Receiver")
	tv.Television.SleepDiscoveryMode.SetValue(characteristic.SleepDiscoveryModeAlwaysDiscoverable)
	r.Tv = tv

	for _, source := range r.sources {
		inputSource := createSource(source)
		r.Tv.Television.AddS(inputSource.S)
		r.Tv.A.AddS(inputSource.S)
	}
}

func (r *HomekitReceiver) SetSource(source int) error {
	return r.Tv.Television.ActiveIdentifier.SetValue(source)
}

func (r *HomekitReceiver) SetPowerState(state int) error {
	return r.Tv.Television.Active.SetValue(state)
}

func (r *HomekitReceiver) RegisterPowerCallback(callback intCallback) {
	r.Tv.Television.Active.OnValueUpdate(callback)
}

func (r *HomekitReceiver) RegisterSourceCallback(callback intCallback) {
	r.Tv.Television.ActiveIdentifier.OnValueUpdate(callback)
}

func createSource(source Source) *service.InputSource {
	inputSource := service.NewInputSource()

	inputId := int(source.name)

	id := characteristic.NewIdentifier()
	id.SetValue(inputId)
	inputSource.AddC(id.C)

	inputSource.ConfiguredName.SetValue(source.displayName)
	inputSource.IsConfigured.SetValue(characteristic.IsConfiguredConfigured)
	inputSource.InputSourceType.SetValue(characteristic.InputSourceTypeHdmi)
	inputSource.CurrentVisibilityState.SetValue(characteristic.CurrentVisibilityStateShown)
	inputSource.ConfiguredName.Permissions = []string{characteristic.PermissionRead}

	return inputSource
}
