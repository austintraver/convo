package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var with string

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search REGEXP",
	Short: "find conversations matching a regular expression",
	Run:   handleSearch,
	Args:  cobra.ExactArgs(1),
}

// handleSearch finds messages that match the provided pattern,
// and prints each match to the command line
func handleSearch(cmd *cobra.Command, args []string) {
	pattern := args[0]
	query := `
	SELECT
	  chat.chat_identifier AS identifier,
	  DATETIME(message.date / 1000000000 + STRFTIME("%s", "2001-01-01"), "unixepoch", "localtime") AS moment,
	  text
	FROM
	  chat
	  JOIN chat_message_join ON chat.rowid = chat_message_join.chat_id
	  JOIN message ON chat_message_join.message_id = message.rowid
	WHERE
	  TRUE
	  AND text LIKE ?
	  AND identifier = ?
	ORDER BY
	  moment DESC
		`
	// rows, err := db.Query(query, limit)
	rows, err := db.Query(query, pattern, with)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()
	fmt.Printf("%s\t%s\n", "identifier", "messages")
	for rows.Next() {
		var identifer string
		var moment string
		var text string
		err = rows.Scan(&identifer, &moment, &text)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("%s\t%s\t%s\n", identifer, moment, text)
	}
	err = rows.Err()
	if err != nil {
		log.Fatalln(err)
	}
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringVarP(&with, "with", "w", with, "only search conversations with a specific person")
}
