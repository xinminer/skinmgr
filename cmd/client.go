package cmd

import (
	"context"
	"github.com/gogf/gf/v2/os/gtimer"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/urfave/cli/v2"
	"net"
	"skinmgr/internal/balance"
	"skinmgr/internal/exec"
	"skinmgr/internal/file"
	"skinmgr/internal/log"
	"time"
)

func init() {
	registerCommand(clientCmd)
}

var clientCmd = &cli.Command{
	Name:    "client",
	Aliases: []string{"cli"},
	Usage:   "managing copy business",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "consul-addr",
			Usage: "consul ip",
			Value: "0.0.0.0",
		},
		&cli.IntFlag{
			Name:  "consul-port",
			Usage: "consul port",
			Value: 8500,
		},
		&cli.StringFlag{
			Name:  "path",
			Usage: "scan file path",
		},
		&cli.StringFlag{
			Name:  "suffix",
			Usage: "scan file suffix",
		},
		&cli.IntFlag{
			Name:  "interval",
			Usage: "mv file interval time, default 25s",
			Value: 25,
		},
		&cli.StringFlag{
			Name:  "preferred",
			Usage: "preferred call consul service",
		},
		&cli.IntFlag{
			Name:  "parallel",
			Usage: "plot copy count",
			Value: 5,
		},
	},
	Action: func(ctx *cli.Context) error {
		consulAddr := ctx.String("consul-addr")
		consulPort := ctx.Int("consul-port")
		path := ctx.String("path")
		suffix := ctx.String("suffix")
		interval := ctx.Int("interval")
		preferred := ctx.String("preferred")
		parallel := ctx.Int("parallel")

		consulAddress := &net.TCPAddr{
			IP:   net.ParseIP(consulAddr),
			Port: consulPort,
		}

		gtimer.AddSingleton(ctx.Context, time.Second, func(ctx context.Context) {
			cps, err := exec.PlotCopyCount("chia_plot_copy")
			if err != nil {
				log.Log.Errorf("plot copy count error: %v", err)
				return
			}

			if cps >= parallel {
				time.Sleep(time.Duration(interval) * time.Second)
				return
			}

			f, err := file.ScanFile(path, suffix)
			if err != nil {
				log.Log.Errorf("scan file error: %v", err)
				return
			}
			svr, err := balance.Random(consulAddress.String(), "storage-manage-service", preferred)
			if err != nil {
				log.Log.Errorf("discovery error: %v", err)
				return
			}
			svr = gstr.Split(svr, ":")[0]
			if err = exec.PlotCopy(f, svr); err != nil {
				log.Log.Errorf("plot copy error: %v", err)
				return
			}

			time.Sleep(time.Duration(interval) * time.Second)
		})

		select {}

	},
}
