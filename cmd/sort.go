package cmd

import (
	"fmt"

	"github.com/Khatunzev-Anton/filesorter/internal/mainconfig"
	"github.com/Khatunzev-Anton/filesorter/internal/services/fsorter"
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
		fmt.Println("config", cfgFile)
		var err error
		config, err := mainconfig.NewMainConfig(cfgFile)
		if err != nil {
			return err
		}

		rd := config.GetRecordDescriptor()

		sorter, err := fsorter.NewFSorterQuick(rd)
		if err != nil {
			return err
		}

		err = sorter.Sort(sortfilename, sortby)

		return err
	},
}
