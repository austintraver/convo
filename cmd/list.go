package cmd

import (
	"fmt"
	"log"

	"github.com/logrusorgru/aurora/v3"
	"github.com/mattn/go-isatty"
	"golang.org/x/sys/unix"

	"github.com/spf13/cobra"
)

// count filters out entries that do not have a minimum number of
// messages that have already been sent. If the value of count is negative,
// (which it is by default), then the value of this filter is disabled.
var count = -1

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show all conversations",
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
    chat.chat_identifier AS id,
    count(chat.chat_identifier) AS messages
	FROM
		chat
		JOIN chat_message_join ON chat."ROWID" = chat_message_join.chat_id
		JOIN message ON chat_message_join.message_id = message."ROWID"
	WHERE TRUE
	-- filter out message reactions
	AND text IS NOT NULL
	AND associated_message_type == 0
	-- filter out empty messages
	AND trim(text, ' ') <> ''
	AND text <> '￼'
	GROUP BY
		chat.chat_identifier
	HAVING messages > ?
	ORDER BY
		messages DESC, id DESC;
	`
	rows, err := db.Query(query, count)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		var messages string
		err = rows.Scan(&id, &messages)
		if err != nil {
			log.Fatalln(err)
		}
		if isatty.IsTerminal(uintptr(unix.Stdout)) {
			fmt.Printf("%s\t%s\n", aurora.Yellow(id), aurora.Blue(messages))
		} else {
			fmt.Printf("%s\t%s\n", id, messages)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatalln(err)
	}
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().IntVarP(
		&count,
		"count",
		"c",
		count,
		"only show conversations with more than this many messages",
	)
}
