## PBRelay Producers

A collection of producers for [pbrelay](https://github.com/psychobummer/pbrelay)

## Installation

`go get github.com/psychobummer/pbrelay-producer`

## Useage

```
$ go build
$ ./pbrelay-producer

$ A tool for sending arbitrary data to a PsychoBummer(t)(r)(tm) relayserver

Usage:
  pbrelay-producer [command]

Available Commands:
  help        Help about any command
  midi        Stream MIDI data to a relay server

Flags:
  -h, --help   help for pbrelay-producer

Use "pbrelay-producer [command] --help" for more information about a command.
```

## Producers

* [midi](https://github.com/psychobummer/pbrelay-producer/blob/master/midiproducer/midiproducer.go) use [pbmidi](https://github.com/psychobummer/pbmidi/) to capture midi events from a local device and stream them to any number of remote endpoints.
