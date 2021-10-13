package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string
var db *sql.DB
var home string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "convo",
	Version: "0.1.3",
	Short:   "convo: A CLI for your iMessage conversations",
	// DisableFlagsInUseLine: true,
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

	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $HOME/.convo.yaml)")

	var err error
	home, err = os.UserHomeDir()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
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

	// Read in environment variables that match the format CONVO_[A-Z]+ format.
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	viper.ReadInConfig()
	// err := viper.ReadInConfig()
	// if err == nil {
	// 	fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	// }
}
