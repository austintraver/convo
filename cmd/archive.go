package cmd

import (
	"fmt"
	"github.com/mattn/go-isatty"
	"golang.org/x/sys/unix"
	"log"

	"github.com/logrusorgru/aurora/v3"
	"github.com/spf13/cobra"
)

// archiveCmd represents the search command
var archiveCmd = &cobra.Command{
	Use:   "archive ID",
	Short: "archive the text contents of a conversation to standard output",
	Run:   handleArchive,
	Args:  cobra.ExactArgs(1),
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
		-- filter out tapbacks
		AND NOT CONTENT LIKE 'Loved “%”'
		AND NOT CONTENT LIKE 'Liked “%”'
		AND NOT CONTENT LIKE 'Laughed at “%”'
		AND NOT CONTENT LIKE 'Disliked “%”'
		AND NOT CONTENT LIKE 'Emphasized “%”'
		AND NOT CONTENT LIKE '%an image'
		-- filter out Fitness app
		AND NOT text LIKE '$(kIMTranscriptPluginBreadcrumb%'
	ORDER BY moment ASC
	`
	// rows, err := db.Query(query, limit)
	rows, err := db.Query(query, identifier)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()
	for rows.Next() {
		var moment string
		var way string
		var id string
		var content string
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
