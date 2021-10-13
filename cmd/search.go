package cmd

import (
	"fmt"
	"log"

	"github.com/logrusorgru/aurora/v3"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
)

var with string

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search PATTERN",
	Short: "Find conversations matching a SQL formatted pattern",
	Run:   handleSearch,
	Args:  cobra.ExactArgs(1),
}

// handleSearch finds messages that match the provided pattern,
// and prints each match to the command line
func handleSearch(cmd *cobra.Command, args []string) {
	pattern := args[0]
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
			AS identifier,
		text
	FROM
		chat
		JOIN chat_message_join
			ON chat.rowid = chat_message_join.chat_id
		JOIN message
			ON chat_message_join.message_id = message.rowid
	WHERE
		TRUE
		AND text LIKE ?
		-- filter out empty messages
		AND text IS NOT NULL
		AND trim(text, ' ') <> ''
		AND text <> '￼'
		-- filter out tapbacks
		AND NOT text LIKE 'Loved “%”'
		AND NOT text LIKE 'Liked “%”'
		AND NOT text LIKE 'Laughed at “%”'
		AND NOT text LIKE 'Disliked “%”'
		AND NOT text LIKE 'Emphasized “%”'
	 	AND NOT text LIKE '%an image'
		-- filter out Fitness app
		AND NOT text LIKE '$(kIMTranscriptPluginBreadcrumb%'
		`
	if with != "" {
		query += "AND identifier = ?"
	}
	// rows, err := db.Query(query, limit)
	rows, err := db.Query(query, pattern, with)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()
	for rows.Next() {
		var moment string
		var way string
		var identifier string
		var content string
		err = rows.Scan(&moment, &way, &identifier, &content)
		if err != nil {
			log.Fatalln(err)
		}
		if isatty.IsTerminal(uintptr(unix.Stdout)) {
			fmt.Printf("[%s] %s (%s): %s\n", aurora.Blue(moment), aurora.Bold(way), aurora.Yellow(identifier), content)
		} else {
			fmt.Printf("[%s] %s (%s): %s\n", moment, way, identifier, content)
		}
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
