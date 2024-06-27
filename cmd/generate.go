package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var generateoutput string
var generatecount uint32

func init() {
	generateCmd.Flags().StringVarP(&generateoutput, "output", "o", "", "Output file name")
	generateCmd.MarkFlagRequired("output")
	generateCmd.Flags().Uint32Var(&generatecount, "count", 10000, "Count of records to generate")
	generateCmd.MarkFlagRequired("count")

	rootCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate test file",
	Long:  `This command generates binary name-value file-`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := application.GenerateFile(generateoutput, int(generatecount))
		if err != nil {
			return err
		}

		fmt.Println("FILE GENERATED:", generateoutput)
		return nil
	},
}
