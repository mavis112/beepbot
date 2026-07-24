//go:build windows && cgo

package audio

import (
	"encoding/binary"
	"fmt"
	"log"
	"strings"
	"sync"
	"unsafe"

	"github.com/gen2brain/malgo"
	"github.com/gopxl/beep/v2"
)

type MalgoPlayer struct {
	context *malgo.AllocatedContext
	device  *malgo.Device
	mtx     sync.Mutex
	mixer   beep.Mixer
	buffer  [][2]float64
}

func NewMalgoPlayer(sampleRate int, deviceName string) (*MalgoPlayer, error) {
	ctx, err := malgo.InitContext([]malgo.Backend{malgo.BackendWasapi}, malgo.ContextConfig{}, nil)
	if err != nil {
		return nil, fmt.Errorf("context initialization failed: %w", err)
	}
	player := &MalgoPlayer{context: ctx}
	deviceConfig := malgo.DefaultDeviceConfig(malgo.Playback)
	deviceConfig.Playback.Format = malgo.FormatS16
	deviceConfig.Playback.Channels = 2
	deviceConfig.SampleRate = uint32(sampleRate)

	var deviceId *malgo.DeviceID
	devices, err := ctx.Devices(malgo.Playback)
	if err != nil || deviceName == "" || strings.ToLower(deviceName) == "list" {
		deviceId = nil
		if deviceName == "" {
			log.Println("Routing audio to system default output device")
		}
	} else {
		for i := range devices {
			if strings.Contains(strings.ToLower(devices[i].Name()), strings.ToLower(deviceName)) {
				deviceId = &devices[i].ID
				log.Printf("Audio successfully routed to: %s", devices[i].Name())
				break
			}
		}
		if deviceId == nil {
			log.Printf("Audio device '%s' not found", deviceName)
			showDevicesList(devices)
		}
	}
	if strings.ToLower(deviceName) == "list" {
		showDevicesList(devices)
	}
	var deviceIdPtr unsafe.Pointer

	if deviceId != nil {
		deviceIdPtr = deviceId.Pointer()
	}

	deviceConfig.Playback.DeviceID = deviceIdPtr
	var callback malgo.DeviceCallbacks
	callback.Data = func(pOutputSample, pInputSamples []byte, framecount uint32) {
		if len(player.buffer) != int(framecount) {
			player.mtx.Lock()
			player.buffer = make([][2]float64, framecount)
			player.mtx.Unlock()
		}
		player.mtx.Lock()
		n, _ := player.mixer.Stream(player.buffer)
		player.mtx.Unlock()
		for i := 0; i < n; i++ {
			sample := player.buffer[i]
			left := int16(sample[0] * 32767.0)
			right := int16(sample[1] * 32767.0)
			offset := i * 4
			binary.LittleEndian.PutUint16(pOutputSample[offset:], uint16(left))
			binary.LittleEndian.PutUint16(pOutputSample[offset+2:], uint16(right))
		}
	}
	device, err := malgo.InitDevice(ctx.Context, deviceConfig, callback)
	if err != nil {
		_ = ctx.Uninit()
		ctx.Free()
		return nil, fmt.Errorf("player is failed to init: %w", err)
	}
	if err := device.Start(); err != nil {
		device.Uninit()
		_ = ctx.Uninit()
		ctx.Free()
		return nil, fmt.Errorf("player is failed to init: %w", err)
	}
	player.device = device
	return player, nil
}

func (m *MalgoPlayer) Play(streamer beep.Streamer) {
	m.mtx.Lock()
	m.mixer.Add(streamer)
	m.mtx.Unlock()
}

func (m *MalgoPlayer) Stop() {
	m.mtx.Lock()
	m.mixer.Clear()
	m.mtx.Unlock()
}

func (m *MalgoPlayer) Close() {
	if m.device != nil {
		m.device.Uninit()
	}
	if m.context != nil {
		_ = m.context.Uninit()
		m.context.Free()
	}
}

func showDevicesList(devices []malgo.DeviceInfo) {
	log.Println("Routing audio to system default output device")
	log.Println("Available system audio output devices:")
	for _, device := range devices {
		log.Println(device.Name())
	}
}
