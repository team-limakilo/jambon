package jambon

import (
	"fmt"
	"strings"

	"github.com/b1naryth1ef/jambon/tacview"
	"github.com/urfave/cli/v2"
)

// CommandRecord handles recording TacView files from a real-time server
var CommandRecord = cli.Command{
	Name:        "record",
	Description: "record a tacview acmi file from a real time server",
	Action:      commandRecord,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "server",
			Usage:    "connection string for the TacView realtime server",
			Required: true,
		},
		&cli.PathFlag{
			Name:     "output",
			Usage:    "path to the output ACMI file",
			Required: true,
		},
		&cli.StringFlag{
			Name:  "username",
			Usage: "username to use when connecting to the realtime server",
			Value: "jambon-record",
		},
		&cli.StringFlag{
			Name:  "password",
			Usage: "password to use when connecting to the realtime server",
			Value: "",
		},
	},
}

func commandRecord(ctx *cli.Context) error {
	serverStr := ctx.String("server")
	if strings.Index(serverStr, ":") == -1 {
		serverStr = fmt.Sprintf("%s:42674", serverStr)
	}

	reader, err := tacview.NewRealTimeReader(serverStr, ctx.String("username"), ctx.String("password"))
	if err != nil {
		return err
	}

	outputFile, err := OpenWritableTacView(ctx.Path("output"))
	if err != nil {
		return err
	}

	writer, err := tacview.NewWriter(outputFile, &reader.Header)
	if err != nil {
		return err
	}
	defer writer.Close()

	data := make(chan *tacview.TimeFrame, 1)
	go reader.ProcessTimeFrames(1, data)

	for {
		frame, ok := <-data
		if !ok {
			break
		}

		err = writer.WriteTimeFrame(frame)
		if err != nil {
			return err
		}
	}

	return nil
}
