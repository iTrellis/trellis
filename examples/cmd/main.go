//

// tbuilder template

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-trellis/common/builder"
	"github.com/go-trellis/common/logger"
	"github.com/go-trellis/config"
	"github.com/go-trellis/trellis/configure"
	"github.com/go-trellis/trellis/service"
	"github.com/go-trellis/trellis/version"

	_ "github.com/go-trellis/trellis/examples/services"
	_ "github.com/go-trellis/trellis/service/api"

	"github.com/spf13/cobra"
)

var cfgFile string

func main() {
	Execute()
}

var rootCmd = &cobra.Command{
	Use:     "sample",
	Short:   "sample project",
	Version: "v1",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	runCmd.Flags().StringVar(&cfgFile, "config", ".sample.yaml", "config file (default is .sample.yaml)")

	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(runCmd)
}

// infoCmd represents the builder command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "show project info",
	Long:  "show the project of build info",
	Run:   func(*cobra.Command, []string) { fmt.Println(version.BuildInfo()) },
}

// runCmd represents the builder command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run project",
	Long:  "run sample project",
	Run: func(*cobra.Command, []string) {
		builder.Show()

		c := &configure.Config{}

		err := config.NewSuffixReader().Read(cfgFile, c)
		if err != nil {
			panic(err)
		}

		log := logger.NewLogger()
		defer log.ClearSubscribers()

		chanWriter, err := logger.ChanWriter(log,
			logger.ChanWiterLevel(c.Project.Logger.Level),
			logger.ChanWiterSeparator(c.Project.Logger.Separator),
			logger.ChanWiterBuffer(c.Project.Logger.ChanBuffer),
		)
		if err != nil {
			panic(err)
		}
		log.Subscriber(chanWriter)

		_, err = log.Subscriber(chanWriter)
		if err != nil {
			panic(err)
		}

		defer chanWriter.Stop()

		err = service.Run(c.Project, log)
		if err != nil {
			time.Sleep(time.Second)
			return
		}

		service.BlockStop()
	},
}
