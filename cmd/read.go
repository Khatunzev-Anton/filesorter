package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/Khatunzev-Anton/filesorter/internal/mainconfig"
	"github.com/Khatunzev-Anton/filesorter/internal/models"
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
		fmt.Println("config", cfgFile)
		var err error
		config, err := mainconfig.NewMainConfig(cfgFile)
		if err != nil {
			return err
		}

		rd := config.GetRecordDescriptor()
		inputfile, err := os.Open(readfilename)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer inputfile.Close()

		var line uint32 = 0

		for {
			b := make([]byte, rd.Size())
			_, err := inputfile.ReadAt(b, int64(readoffset+line)*int64(rd.Size()))
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				return fmt.Errorf("failed to read record: %w", err)
			}
			rec, err := models.FromBytes(b, rd)
			if err != nil {
				return fmt.Errorf("failed to deserialize record: %w", err)
			}
			str, err := rec.String()
			if err != nil {
				return fmt.Errorf("failed to serialize record: %w", err)
			}
			fmt.Println(str)
			line++
			if line >= readcount {
				break
			}
		}

		return nil
	},
}
