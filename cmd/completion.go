package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion { bash | zsh | fish | powershell }",
	Short: "Generate shell completions",
	Args: cobra.ExactArgs(1),
	ValidArgs: []string{
		"bash\t generate autocompletions for bash",
		"zsh\t generate autocompletions for zsh",
		"fish\t generate autocompletions for fish",
		"powershell\t generate autocompletions for powershell",
	},
	Run: func(cmd *cobra.Command, args []string) {
		shell := args[0]
		switch shell {
		case "bash":
			cmd.Root().GenBashCompletionV2(os.Stdout, true)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		default:
			fmt.Fprintf(os.Stderr, "%s - invalid argument '%s'\n\n", cmd.Name(), shell)
			cmd.Usage()
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
