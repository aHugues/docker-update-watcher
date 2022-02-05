package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/ahugues/docker-update-watcher/config"
	"github.com/ahugues/docker-update-watcher/docker"
	"github.com/ahugues/docker-update-watcher/remotedocker"
	"github.com/sirupsen/logrus"
)

type cliArgs struct {
	ConfigPath string
	AsDaemon   bool
}

func parseArgs() cliArgs {
	confPath := flag.String("config-path", "config.json", "Location of the application JSON config file")
	daemon := flag.Bool("daemon", false, "Run the application as a daemon instead of a one-shot")
	flag.Parse()

	return cliArgs{
		ConfigPath: *confPath,
		AsDaemon:   *daemon,
	}
}

func initLog(logLevel logrus.Level) *logrus.Logger {
	logger := logrus.New()

	logger.SetLevel(logLevel)
	logger.SetOutput(os.Stdout)
	return logger
}

func main() {

	args := parseArgs()

	appConf, err := config.ReadFromFile(args.ConfigPath)
	if err != nil {
		panic(err)
	}

	appLogger := initLog(appConf.Logging.Level)

	appLogger.Info("Application initialized")
	fmt.Printf("%+v", appConf)

	appLogger.Info("Getting initial list of images")
	initialList, err := docker.ReadInitialConfig(appConf.StandaloneConfig.IntialListPath)
	if err != nil {
		appLogger.Fatalf("Failed to get initial list of images: %s", err.Error())
	}
	appLogger.Info(initialList)

	// bearerToken, err := remotedocker.Login(appConf.DockerConfig.Username, appConf.DockerConfig.AuthToken)
	// if err != nil {
	// 	appLogger.Fatalf("Failed to login: %s", err.Error())
	// } else {
	// 	appLogger.Infof("Got token %q", bearerToken)
	// }

	appLogger.Info("Comparing remote images")
	for _, img := range *initialList {
		if _, err := remotedocker.GetRemote(context.Background(), http.DefaultClient, img.Namespace, img.Name); err != nil {
			appLogger.Errorf("Failed to get remote %s/%s: %s", img.Namespace, img.Name, err.Error())
		} else {
			appLogger.Infof("Got images for tag %s/%s", img.Namespace, img.Name)

		}
	}
}
