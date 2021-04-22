package cmd

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"io"
	"os"

	"github.com/psychobummer/pbrelay-producer/midiproducer"
	pbk "github.com/psychobummer/pbrelay/rpc/keystore"
	pbr "github.com/psychobummer/pbrelay/rpc/relay"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var midiCmd = &cobra.Command{
	Use:   "midi",
	Short: "Stream MIDI data to a relay server",
	Run:   doMidiCmd,
}

func init() {
	midiCmd.Flags().StringP("connect", "c", "", "ip:port of the server to connect (e.g: 1.2.3.4:9999) (required)")
	midiCmd.Flags().BoolP("list", "l", false, "print a list of available MIDI devices to monitor and exit")
	midiCmd.Flags().IntP("device", "d", 0, "specify which MIDI device we should read from. Use -l to list known devices (default: 0)")
	rootCmd.AddCommand(midiCmd)
}

func doMidiCmd(cmd *cobra.Command, args []string) {
	printList, _ := cmd.Flags().GetBool("list")
	if printList {
		availableInputs, err := midiproducer.Inputs()
		if err != nil {
			log.Error().Msg(err.Error())
			os.Exit(1)
		}
		for i, name := range availableInputs {
			fmt.Printf("%d %s\n", i, name)
		}
		os.Exit(0)
	}

	midiDevice, _ := cmd.Flags().GetInt("device")
	midiProducer, err := midiproducer.New(midiDevice)
	if err != nil {
		log.Error().Msgf("could not open midi device %d: %s", midiDevice, err.Error())
	}

	go func() {
		if err := midiProducer.Start(); err != nil {
			log.Error().Msgf("couldnt start midi read stream: %s", err.Error())
		}
	}()

	// TODO: extract connection creation stuff; it's boilerplate
	// dial server
	hostAddr, _ := cmd.Flags().GetString("connect")
	conn, err := grpc.Dial(hostAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal().Msgf("couldnt dial server: %v", err)
	}
	defer conn.Close()

	// create a signing keypair for this session
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	keystoreClient := pbk.NewKeystoreServiceClient(conn)
	keyResp, err := keystoreClient.CreateKey(context.Background(), &pbk.CreateKeyRequest{
		PublicKey: pubKey,
	})
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	log.Info().Msgf("YOUR PRODUCER ID: %s", keyResp.GetId())

	relayClient := pbr.NewRelayServiceClient(conn)
	stream, err := relayClient.CreateStream(context.Background())

	// publish data from the midi stream
	for payload := range midiProducer.Stream() {
		msg := &pbr.StreamMessage{
			Id:        keyResp.GetId(),
			Data:      payload,
			Signature: ed25519.Sign(privKey, payload),
		}
		if err := stream.Send(msg); err != nil {
			if err == io.EOF {
				log.Debug().Msgf("stream closed")
				break
			}
			_, grpcErr := stream.CloseAndRecv()
			e, _ := status.FromError(grpcErr)
			log.Error().Msgf("%v", e)
			break
		}
	}
	midiProducer.Stop()
}
