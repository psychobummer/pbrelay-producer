package midiproducer

import (
	"encoding/json"

	"github.com/psychobummer/pbmidi"
	"github.com/rs/zerolog/log"
)

type MidiProducer struct {
	stream pbmidi.MidiStream
}

func New(deviceNumber int) (*MidiProducer, error) {
	device, err := pbmidi.New(0)
	if err != nil {
		return nil, err
	}
	mp := MidiProducer{
		stream: device,
	}

	return &mp, nil
}

func (m *MidiProducer) Stream() <-chan []byte {
	streamData := make(chan []byte)
	go func() {
		for midiMsg := range m.stream.Stream() {
			data, err := json.Marshal(midiMsg)
			if err != nil {
				log.Error().Msgf("malformed midi message: %v", err)
				continue
			}
			streamData <- data
		}
	}()
	return streamData
}

func (m *MidiProducer) Start() error {
	return m.stream.Start()
}

func (m *MidiProducer) Stop() {
	m.stream.Stop()
}

func Inputs() ([]string, error) {
	return pbmidi.Inputs()
}
