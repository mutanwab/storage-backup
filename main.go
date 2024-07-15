package main

import (
	"fmt"
	"github.com/mutanwab/storage-backup/pkg/app"
	"github.com/mutanwab/storage-backup/pkg/config"
	"github.com/mutanwab/storage-backup/pkg/server"

	"github.com/rancher/wrangler/pkg/signals"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func cmdNotFound(c *cli.Context, command string) {
	panic(fmt.Errorf("unrecognized command: %s", command))
}

func onUsageError(c *cli.Context, err error, isSubcommand bool) error {
	panic(fmt.Errorf("usage error, please check your command"))
}

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	//a := cli.NewApp()
	//a.Usage = "storage backup"
	//
	//a.Before = func(c *cli.Context) error {
	//	if c.GlobalBool("debug") {
	//		logrus.SetLevel(logrus.DebugLevel)
	//	}
	//	if c.GlobalBool("trace") {
	//		logrus.SetLevel(logrus.TraceLevel)
	//	}
	//	if c.GlobalBool("log-json") {
	//		logrus.SetFormatter(&logrus.JSONFormatter{})
	//	}
	//	return nil
	//}
	//
	//a.Flags = []cli.Flag{
	//	cli.BoolFlag{
	//		Name:   "debug, d",
	//		Usage:  "enable debug logging level",
	//		EnvVar: "RANCHER_DEBUG",
	//	},
	//	cli.BoolFlag{
	//		Name:   "trace, t",
	//		Usage:  "enable trace logging level",
	//		EnvVar: "RANCHER_TRACE",
	//	},
	//	cli.BoolFlag{
	//		Name:   "log-json, j",
	//		Usage:  "log in json format",
	//		EnvVar: "RANCHER_LOG_JSON",
	//	},
	//}
	//a.Commands = []cli.Command{
	//	app.BackupCmd(),
	//	app.RestoreCmd(),
	//}
	//a.CommandNotFound = cmdNotFound
	//a.OnUsageError = onUsageError

	//if err := a.Run(os.Args); err != nil {
	//	logrus.Fatalf("Critical error: %v", err)
	//}

	var options config.Options
	flags := []cli.Flag{
		cli.IntFlag{
			Name:        "threadiness",
			EnvVar:      "THREADINESS",
			Usage:       "Specify controller threads",
			Value:       10,
			Destination: &options.Threadiness,
		},
		cli.IntFlag{
			Name:        "http-port",
			EnvVar:      "HARVESTER_SERVER_HTTP_PORT",
			Usage:       "HTTP listen port",
			Value:       8080,
			Destination: &options.HTTPListenPort,
		},
		cli.IntFlag{
			Name:        "https-port",
			EnvVar:      "HARVESTER_SERVER_HTTPS_PORT",
			Usage:       "HTTPS listen port",
			Value:       8443,
			Destination: &options.HTTPSListenPort,
		},
		cli.StringFlag{
			Name:        "namespace",
			EnvVar:      "NAMESPACE",
			Destination: &options.Namespace,
			Usage:       "The default namespace to store management resources",
			Required:    true,
		},
		cli.BoolFlag{
			Name:        "hci-mode",
			EnvVar:      "HCI_MODE",
			Usage:       "Enable HCI mode. Additional controllers are registered in HCI mode",
			Destination: &options.HCIMode,
		},
		cli.BoolFlag{
			Name:        "rancher-embedded",
			EnvVar:      "RANCHER_EMBEDDED",
			Usage:       "Specify whether the Harvester is running with embedded Rancher mode, default to false",
			Destination: &options.RancherEmbedded,
		},
		cli.StringFlag{
			Name:        "rancher-server-url",
			EnvVar:      "RANCHER_SERVER_URL",
			Usage:       "Specify the URL to connect to the Rancher server",
			Destination: &options.RancherURL,
			Hidden:      true,
		},
	}

	appDaemon := app.NewApp("storage-backup server", "", flags, func(commonOptions *config.CommonOptions) error {
		return run(commonOptions, options)
	})
	appDaemon.Run()
}

func run(commonOptions *config.CommonOptions, options config.Options) error {
	logrus.Info("Starting controller")
	ctx := signals.SetupSignalContext()

	kubeConfig, err := server.GetConfig(commonOptions.KubeConfig)
	if err != nil {
		return fmt.Errorf("failed to find kubeconfig: %v", err)
	}

	harv, err := server.New(ctx, kubeConfig, options)
	if err != nil {
		return fmt.Errorf("failed to create harvester server: %v", err)
	}
	return harv.ListenAndServe(nil, options)
}
