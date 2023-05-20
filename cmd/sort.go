package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var sortfilename string
var sortby string

func init() {
	sortCmd.Flags().StringVarP(&sortfilename, "file", "f", "", "File name")
	sortCmd.MarkFlagRequired("file")
	sortCmd.Flags().StringVarP(&sortby, "sortby", "s", "", "Sort by")
	sortCmd.MarkFlagRequired("sortby")

	rootCmd.AddCommand(sortCmd)
}

var sortCmd = &cobra.Command{
	Use:   "sort",
	Short: "Sort binary file",
	Long:  `This command sorts binary file`,
	RunE: func(cmd *cobra.Command, args []string) error {
		f, err := os.OpenFile(sortfilename, os.O_RDWR, os.ModeExclusive)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer f.Close()
		fi, err := f.Stat()
		if err != nil {
			return fmt.Errorf("failed to get file statistics: %w", err)
		}
		left, size := 0, fi.Size()

		err = application.Sort(f, sortby, int64(left), size)

		return err
	},
}
