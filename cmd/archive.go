package cmd

import (
	"fmt"
	"github.com/mattn/go-isatty"
	"golang.org/x/sys/unix"
	"log"
	"strings"

	"github.com/logrusorgru/aurora/v3"
	"github.com/spf13/cobra"
)

// archiveCmd represents the search command
var archiveCmd = &cobra.Command{
	Use:   "archive ID",
	Short: "Archive the text contents of a conversation to standard output",
	Run:   handleArchive,
	Args:  cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		query := `
		SELECT
		DISTINCT
		chat.chat_identifier as id
		FROM
		chat
		`
		rows, err := db.Query(query)
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		var id string
		var output []string
		for rows.Next() {
			err = rows.Scan(&id)
			if strings.Contains(id, toComplete) {
				output = append(output, id)
			}
		}
		return output, cobra.ShellCompDirectiveNoFileComp
	},
}

// handleArchive finds messages that match the provided pattern,
// and prints each match to the command line
func handleArchive(cmd *cobra.Command, args []string) {
	identifier := args[0]
	query := `
	SELECT
		DATETIME(message.date / 1000000000 + STRFTIME("%s", "2001-01-01"), "unixepoch", "localtime")
			AS moment,
		CASE
			WHEN is_from_me = 0
				THEN "<-"
			WHEN is_from_me = 1
				THEN "->"
		END
			AS way,
		chat.chat_identifier
			AS id,
		text AS content
	FROM
		chat
		JOIN chat_message_join
			ON chat.rowid = chat_message_join.chat_id
		JOIN message
			ON chat_message_join.message_id = message.rowid
	WHERE
		id = ?
	  	-- filter out empty messages
		AND content IS NOT NULL
		AND trim(content, ' ') <> ''
		AND CONTENT <> '￼'
		-- filter out message reactions
		AND text IS NOT NULL
		AND associated_message_type == 0
		-- filter out empty messages
		AND trim(text, ' ') <> ''
		AND text <> '￼'
		-- filter out Fitness app
	ORDER BY moment ASC
	`
	// rows, err := db.Query(query, limit)
	rows, err := db.Query(query, identifier)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()
	var moment, way, id, content string
	for rows.Next() {
		err = rows.Scan(&moment, &way, &id, &content)
		if err != nil {
			log.Fatalln(err)
		}
		if isatty.IsTerminal(uintptr(unix.Stdout)) {
			fmt.Printf("(%s) @ [%s] %s %s\n", aurora.Yellow(id), aurora.Bold(moment), aurora.Bold(way), content)
		} else {
			fmt.Printf("(%s) @ [%s] %s %s\n", id, moment, way, content)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatalln(err)
	}
}

func init() {
	rootCmd.AddCommand(archiveCmd)
}
