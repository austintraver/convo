package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// count filters out entries that do not have a minimum number of
// messages that have already been sent. If the value of count is negative,
// (which it is by default), then the value of this filter is disabled.
var count = -1

// listCommand represents the list command
var listCommand = &cobra.Command{
	Use:   "list",
	Short: "show all conversations",
	Run:   handleList,
	Args:  cobra.NoArgs,
}

// handleList finds all conversations that are currently stored
// in the Messages app, and prints out an overview. The overview displays
// who the conversation is with, along with a count of the number of messages
// that have been exchanged.
func handleList(cmd *cobra.Command, args []string) {
	query := `
	SELECT
    chat.chat_identifier AS identifier,
    COUNT(chat.chat_identifier) AS messages
	FROM
		chat
		JOIN chat_message_join ON chat.rowid = chat_message_join.chat_id
		JOIN message ON chat_message_join.message_id = message.rowid
	GROUP BY
		chat.chat_identifier
	HAVING messages > ?
	ORDER BY
		messages DESC;
	`
	rows, err := db.Query(query, count)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()
	fmt.Printf("%s\t%s\n", "identifier", "messages")
	for rows.Next() {
		var id string
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("%s\t%s\n", id, name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatalln(err)
	}
}

func init() {
	rootCmd.AddCommand(listCommand)

	listCommand.Flags().IntVarP(
		&count,
		"count",
		"c",
		count,
		"only show conversations with more than this many messages",
	)
}
