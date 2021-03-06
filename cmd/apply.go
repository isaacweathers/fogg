package cmd

import (
	"os"

	"github.com/chanzuckerberg/fogg/apply"
	"github.com/chanzuckerberg/fogg/templates"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func init() {
	applyCmd.Flags().StringP("config", "c", "fogg.json", "Use this to override the fogg config file.")
	applyCmd.Flags().BoolP("verbose", "v", false, "use this to turn on verbose output")
	rootCmd.AddCommand(applyCmd)
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply model defined in fogg.json to the current tree.",
	Long:  "This command will take the model defined in fogg.json, build a plan and generate the appropriate files from templates.",
	Run: func(cmd *cobra.Command, args []string) {
		logLevel := log.InfoLevel
		if debug { // debug overrides quiet
			logLevel = log.DebugLevel
		} else if quiet {
			logLevel = log.FatalLevel
		}
		log.SetLevel(logLevel)

		var e error
		// Set up fs
		pwd, e := os.Getwd()
		if e != nil {
			log.Panic(e)
		}
		fs := afero.NewBasePathFs(afero.NewOsFs(), pwd)

		// handle flags
		verbose, e := cmd.Flags().GetBool("verbose")
		if e != nil {
			log.Panic(e)
		}
		configFile, e := cmd.Flags().GetString("config")
		if e != nil {
			log.Panic(e)
		}

		// check that we are at root of initialized git repo
		openGitOrExit(pwd)

		config, err := readAndValidateConfig(fs, configFile, verbose)

		exitOnConfigErrors(err)

		// apply
		e = apply.Apply(fs, config, templates.Templates)
		if e != nil {
			log.Panic(e)
		}
	},
}
