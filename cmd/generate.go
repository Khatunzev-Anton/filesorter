package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Khatunzev-Anton/filesorter/internal/mainconfig"
	"github.com/Khatunzev-Anton/filesorter/internal/repositories"
	"github.com/Khatunzev-Anton/filesorter/internal/services/fgenerator"
	"github.com/Khatunzev-Anton/filesorter/internal/services/recordgenerator"
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
	Long:  `This command generates binary name-value file`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("config", cfgFile)
		var err error
		config, err := mainconfig.NewMainConfig(cfgFile)
		if err != nil {
			return err
		}
		ex, err := os.Executable()
		if err != nil {
			return fmt.Errorf("failed to get current executable directory: %w", err)
		}
		namesrepo, err := repositories.NewNameRepository(fmt.Sprintf("%[1]s/internal/data/names.txt", filepath.Dir(ex))) //???
		if err != nil {
			return fmt.Errorf("failed to initialize namesrepo: %w", err)
		}

		g, err := recordgenerator.NewNameSalaryGenerator(config.GetRecordDescriptor(), namesrepo, 100000, 250000)
		if err != nil {
			return fmt.Errorf("failed to initialize generator: %w", err)
		}

		fgen, err := fgenerator.NewFGenerator(g)
		if err != nil {
			return fmt.Errorf("failed to initialize file generator: %w", err)
		}

		err = fgen.GenerateFile(generateoutput, int(generatecount))
		if err != nil {
			return fmt.Errorf("failed to generate file: %w", err)
		}

		fmt.Println("FILE GENERATED:", generateoutput)
		return nil
	},
}
