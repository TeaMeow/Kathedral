package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Kathedral"
	app.Version = "1.0.0"
	app.Usage = "Launch the Kathrdral daemon."
	app.Action = func(c *cli.Context) {
		bot(c)
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "token",
			Usage: "The telegram bot token.",
		},
		cli.StringFlag{
			Name:  "addr",
			Usage: "The url of the file server, will be used at the link buttons.",
			Value: "example.com",
		},
		cli.StringFlag{
			Name:  "port",
			Usage: "The port that Kathedral runs the Golang native file server.",
			Value: "8888",
		},
		cli.BoolFlag{
			Name:  "with-port",
			Usage: "Show the link with the port?",
		},
	}

	app.Run(os.Args)
}
