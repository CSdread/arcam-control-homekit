package main

import (
	"arcam-controller/arcam"
	"context"
	"net/http"

	"github.com/brutella/hap/accessory"
	"github.com/brutella/hap/characteristic"
	"github.com/brutella/hap/log"
	"github.com/brutella/hap/service"
)

type intCallback func(int, int, *http.Request)
type boolCallback func(oldVal, newVal bool, r *http.Request)

func attachInputs(ctx context.Context, arcamClient *arcam.Receiver, tv *accessory.Television) {
	for _, input := range arcamClient.GetAllInputs() {
		inputSource := service.NewInputSource()

		displayName := arcam.InputDisplayNameMap[input]
		inputId := int(input)

		id := characteristic.NewIdentifier()
		id.SetValue(inputId)
		inputSource.AddC(id.C)

		inputSource.ConfiguredName.SetValue(displayName)
		inputSource.IsConfigured.SetValue(characteristic.IsConfiguredConfigured)
		inputSource.InputSourceType.SetValue(characteristic.InputSourceTypeHdmi)
		inputSource.CurrentVisibilityState.SetValue(characteristic.CurrentVisibilityStateShown)
		inputSource.ConfiguredName.Permissions = []string{characteristic.PermissionRead}

		tv.Television.AddS(inputSource.S)
		tv.A.AddS(inputSource.S)
	}
}

func powerOnCallback(ctx context.Context, arcamClient *arcam.Receiver) intCallback {
	return func(newVal, oldVal int, r *http.Request) {
		if newVal == characteristic.ActiveActive {
			err := arcamClient.PowerOn(ctx)
			if err != nil {
				log.Info.Fatalln("")
			}
			return
		}

		err := arcamClient.PowerOff(ctx)
		if err != nil {
			log.Info.Fatalln("")
		}
	}
}

func inputSelectionCallback(ctx context.Context, arcamClient *arcam.Receiver) intCallback {
	return func(newVal, oldVal int, r *http.Request) {
		// TODO: should have some logic here that if we are going with the analog input (record player) we should
		// do stereo analog pass through, maybe make it an option to the plugin if this is what you want to do
		arcamClient.SetSource(ctx, arcam.InputSource(newVal))
	}
}

func speakerVolumeCallback(ctx context.Context, arcamClient *arcam.Receiver) intCallback {
	return func(newVal, oldVal int, r *http.Request) {
		arcamClient.SetVolume(ctx, newVal)
	}
}

func muteCallback(ctx context.Context, arcamClient *arcam.Receiver) boolCallback {
	return func(oldVal, newVal bool, r *http.Request) {
		// TODO; this logic is wrong
		arcamClient.ToggleMute(ctx)
	}
}

func brightnessVolCallback(ctx context.Context, arcamClient *arcam.Receiver) intCallback {
	return func(oldVal, newVal int, r *http.Request) {
		arcamClient.SetVolume(ctx, newVal)
	}
}
