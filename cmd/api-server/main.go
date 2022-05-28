package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"

	"github.com/crypto-pe/api/config"
	"github.com/crypto-pe/api/server"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "api",
	Short: "Cryptope API Server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return run()
	},
}

var (
	configFile string
	logFile    string

	cpuprofile string
	memprofile string

	cfg = &config.Config{}
)

func init() {
	cobra.OnInitialize(initCmd)
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "path to config file")
	rootCmd.PersistentFlags().StringVar(&logFile, "log", "", "path to log file")

	rootCmd.PersistentFlags().StringVar(&cpuprofile, "cpuprofile", "", "path for cpuprofile output")
	rootCmd.PersistentFlags().StringVar(&memprofile, "memprofile", "", "path for memprofile output")
}

func initCmd() {
	must(initConfig())
	// initLogging()
	must(initCPUProfile())
	must(initMemProfile())
}

func initConfig() error {
	err := config.NewFromFile(configFile, os.Getenv("CONFIG"), cfg)
	if err != nil {
		return fmt.Errorf("Error initializing config: %s", err)
	}

	return nil
}

// func initLogging() {
// 	if logFile != "" {
// 		cfg.Logging.LogFile = logFile
// 	}
// }

func initCPUProfile() error {
	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			return fmt.Errorf("Could not create CPU profile: %s", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			return fmt.Errorf("Could not start CPU profile: %s", err)
		}
	}

	return nil
}

func initMemProfile() error {
	if memprofile != "" {
		f, err := os.Create(memprofile)
		if err != nil {
			return fmt.Errorf("Could not create memory profile: %s", err)
		}
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			return fmt.Errorf("Could not write memory profile: %s", err)
		}
	}

	return nil
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	s, err := server.NewServer(cfg)
	if err != nil {
		return err
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		// Shutdown the server
		s.Stop()
	}()

	// Run the server
	err = s.Run()
	if err != nil {
		s.Fatal("server runtime error: %v", err)
	}

	s.End()

	return nil
}
