package main

import (
	"log"
	"os"
	"strconv"

	"github.com/urfave/cli/v2"
	p "github.com/whxaxes/go-proxy/proxy"
)

func main() {
	app := &cli.App{
		Name:  "go-proxy",
		Usage: "Create a proxy server",

		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "port",
				Usage:   "The port of proxy server",
				Value:   3000,
				Aliases: []string{"p"},
			},
		},

		Action: func(c *cli.Context) error {
			proxy := p.TCPProxy{
				Dest: c.Args().Get(0),
				// Through: func(b []byte, flush tp.ThroughFlush, tc *tp.TCPConn) error {
				// 	return flush(b)
				// },
			}

			proxy.Listen(strconv.Itoa(c.Int("port")))
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
