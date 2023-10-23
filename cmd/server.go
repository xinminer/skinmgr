package cmd

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/urfave/cli/v2"
	"net"
	"skinmgr/internal/balance"
	"skinmgr/internal/log"
)

func init() {
	registerCommand(serverCmd)
}

var serverCmd = &cli.Command{
	Name:    "server",
	Aliases: []string{"srv"},
	Usage:   "start a server that receives files and listens on a specified port.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "target-addr",
			Usage: "target on a specified ip",
			Value: "0.0.0.0",
		},
		&cli.IntFlag{
			Name:  "target-port",
			Usage: "target on a specified port",
			Value: 9999,
		},
		&cli.StringFlag{
			Name:  "listen-addr",
			Usage: "listen on a specified ip",
			Value: "0.0.0.0",
		},
		&cli.IntFlag{
			Name:  "listen-port",
			Usage: "listen n a specified port",
			Value: 8888,
		},
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
			Name:  "preferred",
			Usage: "preferred call consul service",
		},
	},
	Action: func(ctx *cli.Context) error {
		targetAddr := ctx.String("target-addr")
		targetPort := ctx.Int("target-port")
		listenAddr := ctx.String("listen-addr")
		listenPort := ctx.Int("listen-port")
		consulAddr := ctx.String("consul-addr")
		consulPort := ctx.Int("consul-port")
		preferred := ctx.String("preferred")

		targetAddress := &net.TCPAddr{
			IP:   net.ParseIP(targetAddr),
			Port: targetPort,
		}

		consulAddress := &net.TCPAddr{
			IP:   net.ParseIP(consulAddr),
			Port: consulPort,
		}

		dis := balance.DiscoveryConfig{
			ID:      guid.S(),
			Name:    "storage-manage-service",
			Tags:    []string{},
			Port:    listenPort,
			Address: listenAddr,
			Meta: map[string]string{
				"preferred": preferred,
				"target":    targetAddress.String(),
			},
		}

		if err := balance.RegisterService(consulAddress.String(), dis); err != nil {
			log.Log.Errorf("register consul error: %v", err)
		}

		svr := g.Server()
		svr.SetPort(listenPort)
		svr.BindHandler("/", func(r *ghttp.Request) {
			r.Response.Write("storage-manage！")
		})
		svr.Run()
		return nil
	},
}
