package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var readfilename string
var readcount uint32
var readoffset uint32

func init() {
	readCmd.Flags().StringVarP(&readfilename, "file", "f", "", "File name")
	readCmd.MarkFlagRequired("file")
	readCmd.Flags().Uint32Var(&readcount, "count", 10000, "Count of records to read")
	readCmd.MarkFlagRequired("count")
	readCmd.Flags().Uint32Var(&readoffset, "offset", 0, "Read offset")
	readCmd.MarkFlagRequired("offset")

	rootCmd.AddCommand(readCmd)
}

var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Read binary file",
	Long:  `This command reads binary name-value and outputs the result into console or in the output file`,
	RunE: func(cmd *cobra.Command, args []string) error {
		inputfile, err := os.Open(readfilename)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer inputfile.Close()

		var i uint32
		recordsize := application.RecordSize()
		for i = 0; i < readcount; i++ {
			rec, err := application.ReadSerializableRecord(inputfile, int64(readoffset+i)*recordsize)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				return err
			}
			str, err := rec.String()
			if err != nil {
				return fmt.Errorf("failed to serialize record: %w", err)
			}
			fmt.Println(str)
		}

		return nil
	},
}
