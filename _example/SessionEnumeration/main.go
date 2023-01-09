package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"unsafe"

	"github.com/go-ole/go-ole"
	"github.com/moutend/go-wca/pkg/wca"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("error: ")

	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	if err := ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED); err != nil {
		return err
	}

	defer ole.CoUninitialize()

	var mmde *wca.IMMDeviceEnumerator

	if err := wca.CoCreateInstance(wca.CLSID_MMDeviceEnumerator, 0, wca.CLSCTX_ALL, wca.IID_IMMDeviceEnumerator, &mmde); err != nil {
		return err
	}

	defer mmde.Release()

	var deviceCollection *wca.IMMDeviceCollection
	if err := mmde.EnumAudioEndpoints(wca.EAll, wca.DEVICE_STATE_ACTIVE, &deviceCollection); err != nil {
		return err
	}
	defer deviceCollection.Release()

	var deviceCount uint32
	if err := deviceCollection.GetCount(&deviceCount); err != nil {
		return err
	}

	// Enuerate all devices
	for i := uint32(0); i < deviceCount; i++ {
		var device *wca.IMMDevice

		if err := deviceCollection.Item(i, &device); err != nil {
			return err
		}

		var deviceId string
		if err := device.GetId(&deviceId); err != nil {
			return err
		}

		var sessionManager *wca.IAudioSessionManager2
		if err := device.Activate(wca.IID_IAudioSessionManager2, wca.CLSCTX_ALL, nil, &sessionManager); err != nil {
			return err
		}
		defer sessionManager.Release()

		var sessionsEnumerator *wca.IAudioSessionEnumerator
		if err := sessionManager.GetSessionEnumerator(&sessionsEnumerator); err != nil {
			return err
		}
		defer sessionsEnumerator.Release()

		var sessionCount int
		if err := sessionsEnumerator.GetCount(&sessionCount); err != nil {
			return err
		}

		fmt.Printf("Device ID: %s\n", deviceId)

		// Enumerate all sessions of this device
		for j := int(0); j < sessionCount; j++ {
			var session *wca.IAudioSessionControl
			if err := sessionsEnumerator.GetSession(j, &session); err != nil {
				return err
			}
			defer session.Release()

			session2Dispatch, err := session.QueryInterface(wca.IID_IAudioSessionControl2)
			if err != nil {
				return err
			}
			defer session2Dispatch.Release()
			session2 := (*wca.IAudioSessionControl2)(unsafe.Pointer(session2Dispatch))
			defer session2.Release()

			var sessionId string
			if err := session2.GetSessionIdentifier(&sessionId); err != nil {
				return err
			}

			fmt.Printf("\tSession ID: %s\n", sessionId)
		}

		fmt.Println()

		device.Release()

	}

	return nil
}
