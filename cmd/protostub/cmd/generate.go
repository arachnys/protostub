package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

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
		st := time.Now()
		defer func(st time.Time) {
			fmt.Printf("\nTime taken: %s\n", time.Since(st))
		}(st)

		protos, err := cmd.Flags().GetStringSlice("proto")
		if err != nil {
			log.Fatalln(err)
		}

		mypy := cmd.Flag("mypy").Value.String()

		totalProtos := len(protos)

		if totalProtos == 0 {
			fmt.Println("You must provide a protobuf file with -p. See help.")
			return
		} else if totalProtos == 1 {
			generateFile(protos[0], mypy)
			return
		}

		var wg sync.WaitGroup
		wg.Add(totalProtos)

		for _, proto := range protos {
			go func(proto string) {
				defer wg.Done()
				generateFile(proto, "")
			}(proto)
		}

		wg.Wait()
	},
}

func generateFile(proto string, mypy string) {
	if mypy == "" {
		mypy = strings.Replace(proto, ".proto", "_pb2.pyi", -1)
	}

	mf, err := os.Create(mypy)

	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		if err := mf.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	if *verbose {
		fmt.Println(fmt.Sprintf("Generating %s", mypy))
	}

	pf, err := os.Open(proto)

	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		if err := pf.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	// first parse the protobuf
	p := protostub.New(pf)

	if err := p.Parse(); err != nil {
		log.Fatalln(err)
	}

	if err := gen.Gen(mf, p, true); err != nil {
		log.Fatalln(err)
	}
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	generateCmd.Flags().StringSliceP("proto", "p", []string{}, "Specify the protobuf file to read from")
	generateCmd.Flags().StringP("mypy", "m", "", "Specify the output file to write the MyPy stub to")
	verbose = generateCmd.Flags().BoolP("verbose", "v", false, "Enable logging")
}
