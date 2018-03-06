package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/arachnys/protostub"
	"github.com/arachnys/protostub/gen"
)

var verbose *bool

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a stub from a given proto file",
	Run: func(cmd *cobra.Command, args []string) {
		proto := cmd.Flag("proto").Value.String()

		if len(proto) == 0 {
			fmt.Println("You must provide a protobuf file with -p. See help.")
			return
		}

		mypy := cmd.Flag("mypy").Value.String()

		if len(mypy) == 0 {
			mypy = strings.Replace(proto, ".proto", "_pb2.pyi", -1)
		}

		mf, err := os.Create(mypy)

		if err != nil {
			panic(err)
		}

		defer func() {
			err := mf.Close()

			if err != nil {
				panic(err)
			}
		}()

		if *verbose {
			fmt.Println(fmt.Sprintf("Generating %s", mypy))
		}

		pf, err := os.Open(proto)

		if err != nil {
			panic(err)
		}

		defer func() {
			err := pf.Close()

			if err != nil {
				panic(err)
			}
		}()

		// first parse the protobuf
		p := protostub.New(pf)
		err = p.Parse()

		if err != nil {
			panic(err)
		}

		err = gen.Gen(mf, p)

		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	generateCmd.Flags().StringP("proto", "p", "", "Specify the protobuf file to read from")
	generateCmd.Flags().StringP("mypy", "m", "", "Specify the output file to write the MyPy stub to")
	verbose = generateCmd.Flags().BoolP("verbose", "v", false, "Enable logging")
}
