/*
Copyright © 2020 Henry Huang <hhh@rutcode.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

// tbuilder template

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-trellis/trellis/configure"
	"github.com/go-trellis/trellis/service"
	"github.com/go-trellis/trellis/version"

	"github.com/go-trellis/common/builder"
	"github.com/go-trellis/common/logger"
	"github.com/go-trellis/config"
	"github.com/spf13/cobra"

	_ "github.com/go-trellis/trellis/examples/services"
	_ "github.com/go-trellis/trellis/service/api"
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

		r, err := service.Run(c.Project, log)
		if err != nil {
			time.Sleep(time.Second)
			return
		}
		defer r.Stop()

		service.BlockStop()
	},
}
