package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zetsub0u/void_archives/http"
	"github.com/zetsub0u/void_archives/loaders"
	"github.com/zetsub0u/void_archives/store"
	"os"
	"os/signal"
	"strings"
	"time"
)

var flagHelper viperPFlagHelper

var (
	addressFlag       string
	portFlag          int
	inMemFlag         bool
	loaderTimerFlag   time.Duration
	mainChannelIDFlag string
)

func printSettings() {
	log.Infof("Bind Address: %s", addressFlag)
	log.Infof("Bind Port: %d", portFlag)
}

func start() {
	printSettings()

	// wire signal handlers
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// hold the binary version
	versionObj := http.Version{
		Version: version,
		Branch:  branch,
		Commit:  commit,
	}

	s := store.NewInMem()
	if err := store.LoadDummyData(s); err != nil {
		log.Fatalf("failed loading dummy data: %v", err)
	}

	// create the youtube loader
	ytLoader, err := loaders.NewYoutube(mainChannelIDFlag)
	if err != nil {
		log.Fatalf("failed initializing youtube loader: %v", err)
	}

	// setup the loader loop
	loader := loaders.NewRunner(s, loaderTimerFlag, ytLoader)
	go loader.Start()

	// http api server
	log.Info("cmd: initializing http api")
	apiConfig := http.ServerConfig{Address: addressFlag, Port: portFlag, Version: versionObj}
	apiServer := http.NewServer(&apiConfig).WithMetrics().WithArchive(s)
	apiServer.Setup()

	log.Info("cmd: starting http api...")
	go apiServer.Start()

	<-quit
	log.Info("cmd: warm shutdown initiated...")

	// on double signal exit immediately
	go func() {
		<-quit
		log.Info("cmd: cold shutdown requested, bye.")
		os.Exit(1)
	}()

	loader.Stop()
	apiServer.Stop()

	log.Info("cmd: exiting...")
	os.Exit(0)
}

func init() {
	// set default prefix for environment variable overrides, nested variables use _ instead of . for sublevels
	viper.SetEnvPrefix("va")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// Base FLags
	startCmd.Flags().StringVarP(&addressFlag, "bind-addr", "b", "localhost", "bind address for the web server")
	startCmd.Flags().IntVarP(&portFlag, "bind-port", "p", 8080, "bind port for the web server")

	// Store flags
	startCmd.Flags().BoolVarP(&inMemFlag, "in-mem", "m", false, "use in memory storage, for testing")
	startCmd.Flags().DurationVarP(&loaderTimerFlag, "loader-timer", "t", 10*time.Minute, "how often to run the loaders")
	startCmd.Flags().StringVarP(&mainChannelIDFlag, "channel-id", "c", "UCLINibmUEYa5PhuBVpyd-7w", "youtube channel id to find subscriptions to poll")

	// Add commands
	RootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start the void archives server",
	Run: func(cmd *cobra.Command, args []string) {
		start()
	},
}
