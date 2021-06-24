package cmd

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string
var db *sql.DB
var home string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "convo",
	Short: "A brief description of your application",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	pathname := fmt.Sprintf("%s/Library/Messages/chat.db", home)

	var err error
	db, err = sql.Open("sqlite3", pathname)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $HOME/.dm.yaml)")

	var err error
	home, err = os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if configFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(configFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".convo" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".convo")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
