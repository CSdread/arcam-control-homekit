package main

import (
	"arcam-controller/arcam"
	"context"
	"net/http"

	"github.com/brutella/hap/characteristic"
	"github.com/brutella/hap/log"
)

func requestCurrentSourceCallback(receiver *HomekitReceiver) func(arcam.ZoneNumber, []byte) error {
	return func(zone arcam.ZoneNumber, data []byte) error {
		return receiver.SetSource(int(data[0]))
	}
}

func powerCommandCallback(receiver *HomekitReceiver) func(arcam.ZoneNumber, []byte) error {
	return func(zone arcam.ZoneNumber, data []byte) error {
		active := characteristic.ActiveInactive
		if arcam.PowerStatus(data[0]) == arcam.PowerStatusActive {
			active = characteristic.ActiveActive
		}
		return receiver.SetPowerState(active)
	}
}

func muteCommandCallback(receiver *HomekitReceiver) func(arcam.ZoneNumber, []byte) error {
	return func(zone arcam.ZoneNumber, data []byte) error {
		isMuted := arcam.MuteState(data[0]) == arcam.MuteStateMuted
		receiver.SetMute(isMuted)
		return nil
	}
}

func directModeCommandCallback(receiver *HomekitReceiver) func(arcam.ZoneNumber, []byte) error {
	return func(zone arcam.ZoneNumber, data []byte) error {
		enable := arcam.DirectModeStatus(data[0]) == arcam.DirectModeActive
		return receiver.SetDirectMode(enable)
	}
}

func volumeCommandCallback(receiver *HomekitReceiver) func(arcam.ZoneNumber, []byte) error {
	return func(zone arcam.ZoneNumber, data []byte) error {
		return receiver.SetVolume(int(data[0]))
	}
}

func directModeStatusCallback(ctx context.Context, arcamClient arcam.Receiver) boolCallback {
	return func(newVal, oldVal bool, r *http.Request) {
		var err error
		if newVal {
			err = arcamClient.EnableDirectMode(ctx)
		} else {
			err = arcamClient.DisableDirectMode(ctx)
		}

		if err != nil {
			log.Info.Fatalln("")
		}
	}
}

func muteCallback(ctx context.Context, arcamClient arcam.Receiver) boolCallback {
	return func(newVal, oldVal bool, r *http.Request) {
		var err error
		if newVal {
			err = arcamClient.Mute(ctx)
		} else {
			err = arcamClient.UnMute(ctx)
		}

		if err != nil {
			log.Info.Fatalln("")
		}
	}
}

func volumeCallback(ctx context.Context, arcamClient arcam.Receiver) intCallback {
	return func(newVal, oldVal int, r *http.Request) {
		err := arcamClient.SetVolume(ctx, newVal)
		if err != nil {
			log.Info.Fatalln("")
		}
	}
}

func powerCallback(ctx context.Context, arcamClient arcam.Receiver) intCallback {
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

func sourceCallback(ctx context.Context, arcamClient arcam.Receiver) intCallback {
	return func(newVal, oldVal int, r *http.Request) {
		// TODO: should have some logic here that if we are going with the analog input (record player) we should
		// do stereo analog pass through, maybe make it an option to the plugin if this is what you want to do
		arcamClient.SetSource(ctx, arcam.InputSource(newVal))
	}
}
