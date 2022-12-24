package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/NodyHub/hashref/pkg/hashref"
	"github.com/NodyHub/hashref/pkg/util"
	"github.com/alecthomas/kong"
)

var CLI struct {
	Config    string `short:"c" optional:"" type:"path" help:"Path to hashref config (default: ~/.hashref), can be overwritten by environment."`
	Details   bool   `short:"d" optional:"" help:"Show details to hash. The fiel \"success\" will be added to the json, even if not stored on the server."`
	Generate  bool   `short:"g" optional:"" help:"Generate client configuration"`
	Meta      string `short:"m" optional:"" type:"path" help:"Read metadata from JSON file, comma separated file list, existing keys are overwritten. Empty values are removed from metadata."`
	Remove    bool   `short:"r" optional:"" help:"Remove hash from db"`
	Set       bool   `short:"s" optional:"" help:"Set metadata for input. (Extends/update existing)"`
	Self      bool   `optional:"" help:"Set/get metadata to yourself"`
	Output    string `short:"o" optional:"" help:"Specify output" type:"path"`
	Publisher string `short:"p" optional:"" help:"Limit request to data from publisher"`
	Verbose   bool   `short:"v" optional:"" help:"Show verbose output"`
	Yes       bool   `short:"y" optional:"" help:"Always confirm"`

	Input []string `arg:"" name:"input" optional:"" help:"Files, strings, hashes"`
}

func main() {
	_ = kong.Parse(&CLI)
	// Check for verbose output
	if CLI.Verbose {
		log.SetOutput(os.Stderr)
	} else {
		log.SetOutput(ioutil.Discard)
	}
	log.Printf("Flags: %+v\n", CLI)

	// Adjust output
	output := os.Stderr
	if CLI.Output == "-" {
		output = os.Stdout
	} else if CLI.Output != "" {
		// Check if file needs to be overwritten
		if _, err := os.Stat(CLI.Output); err == nil {
			// File exist, do we have --yes flag or ask?
			if !CLI.Yes && !util.YesOrNoQuestion(fmt.Sprintf("Overwrite existing file %v?", CLI.Output)) {
				log.Println("Aborted")
				os.Exit(-1)
			} else {
				log.Printf("Overwrite %v\n", CLI.Output)
			}
		}
		// Open file for output
		out, err := os.OpenFile(CLI.Output, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			log.Fatal(err)
		}
		output = out
	}

	// Generate hashref config
	if CLI.Generate {
		cfg := hashref.NewConfig()
		cfg.LoadEnvValues()
		b, err := json.MarshalIndent(cfg, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(output, "%v\n", string(b))
		os.Exit(0)
	}

	// Load local cfg file for client
	cfg := hashref.LoadConfig(CLI.Config)
	cfg.LoadEnvValues()
	hc := hashref.NewClient(cfg)

	// handle management of our own data
	if CLI.Self {
		if CLI.Set {
			// Collect all the metadata
			meta := map[string]interface{}{
				"user": cfg.Publisher,
				"hash": hashref.CalculateHash([]byte(cfg.Publisher)),
				"type": hashref.Lookup[hashref.Publisher],
			}

			// Extend with metadata from config
			for k, v := range cfg.DefaultMeta {
				meta[k] = v
			}

			// Extend with metadata files
			for k, v := range util.LoadMultipleJsonFiles(CLI.Meta) {
				meta[k] = v
			}

			// Remove empty fields
			for k, v := range meta {
				if len(fmt.Sprintf("%v", v)) == 0 {
					delete(meta, k)
				}
			}

			// Perform request
			if success := hc.SetSelf(meta); success {

				// Detailed output or status?
				if CLI.Details {

					// Pretty print details
					if pretty, err := util.GetPrettyJsonFromMap(meta); err != nil {
						log.Printf("%v\n", err)
					} else {
						fmt.Fprintf(output, "%v\n", pretty)
					}

				} else {
					// Just print the success
					fmt.Fprintf(output, "Self-metadata set :)\n")
				}

			} else {
				fmt.Fprintf(output, "Error setting self-metadata :(\n")
			}

		} else {
			meta := hc.GetSelf()
			if pretty, err := util.GetPrettyJsonFromMap(meta); err != nil {
				log.Printf("%v\n", err)
			} else {
				fmt.Fprintf(output, "%v\n", pretty)
			}
		}
		os.Exit(0)
	}

	// Handle hash removal
	if CLI.Remove {
		for _, input := range CLI.Input {
			log.Printf("Process input %v\n", input)
			_, calculatedHash := hashref.GetHashTypeAndValue(input)
			if success := hc.RemoveHash(CLI.Yes, input, calculatedHash); success {
				fmt.Fprintf(output, "%v removed :)\n", input)
			} else {
				fmt.Fprintf(output, "%v not removed :(\n", input)
			}
		}
	}

	// iterate over input
	if len(CLI.Input) > 0 {

		// track status overall
		successAll := true

		// track processed files
		allKeys := make(map[string]bool)

		// iterate over input
		for _, input := range CLI.Input {

			// check if already processed
			if _, isProcessed := allKeys[input]; !isProcessed {

				// Start input processing
				log.Printf("Process input %v\n", input)
				inputType, calculatedHash := hashref.GetHashTypeAndValue(input)

				// Ignore existing data and overwrite
				if CLI.Set {

					// Get based on the type metadata
					meta := hc.CollectLocalMetadata(inputType, input, calculatedHash)

					// Extend with metadata from config
					for k, v := range cfg.DefaultMeta {
						meta[k] = v
					}

					// Extend with metadata from provided files
					for k, v := range util.LoadMultipleJsonFiles(CLI.Meta) {
						meta[k] = v
					}

					// Remove empty fields
					for k, v := range meta {
						if len(fmt.Sprintf("%v", v)) == 0 {
							delete(meta, k)
						}
					}

					// finalize
					if success := hc.SetRemoteData(inputType, input, calculatedHash, meta); success {

						// Detailed output or status?
						if CLI.Details {

							// Pretty print details
							if pretty, err := util.GetPrettyJsonFromMap(meta); err != nil {
								log.Printf("%v\n", err)
							} else {
								fmt.Fprintf(output, "%v\n", pretty)
							}

						} else {
							// Just print the success
							fmt.Fprintf(output, "%v metadata set :)\n", input)
						}

					} else {
						fmt.Fprintf(output, "%v metadata not set :(\n", input)
						successAll = false
					}

				} else {
					// Check if for dedicated publisher before fetch remote data
					var success bool
					var meta map[string]interface{}

					if len(CLI.Publisher) > 0 {
						success, meta = hc.GetRemoteDataFromPublisher(inputType, input, calculatedHash, CLI.Publisher)
					} else {
						success, meta = hc.GetRemoteData(inputType, input, calculatedHash)
					}

					// Fetch data from all publisher
					if success {

						// Detailed output or status?
						if CLI.Details {

							// Pretty print details
							if pretty, err := util.GetPrettyJsonFromMap(meta); err != nil {
								log.Printf("%v\n", err)
							} else {
								fmt.Fprintf(output, "%v\n", pretty)
							}

						} else {
							// Just print the success
							fmt.Fprintf(output, "%v found :)\n", input)
						}

					} else {

						// Detailed output for sad state?
						if CLI.Details {

							// Pretty print details
							meta["input"] = input
							if pretty, err := util.GetPrettyJsonFromMap(meta); err != nil {
								log.Printf("%v\n", err)
							} else {
								fmt.Fprintf(output, "%v\n", pretty)
							}

						} else {
							// Just print the sad state
							fmt.Fprintf(output, "%v not found :(\n", input)
						}

						// Remember non-success
						successAll = false
					}
				}

				// Note that input is processed
				allKeys[input] = true

			} else {
				// Input already processed in privous step
				log.Printf("Skip %v, already processed!\n", input)
			}

		}
		if !successAll {
			os.Exit(-1)
		}
	}

}
