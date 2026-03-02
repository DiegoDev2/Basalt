package main

import (
	"fmt"
	"os"

	"github.com/DiegoDev2/basalt/internal/config"
	"github.com/DiegoDev2/basalt/internal/generator"
	"github.com/DiegoDev2/basalt/internal/parser"
	"github.com/DiegoDev2/basalt/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

var version = "0.1.0"

func main() {
	var rootCmd = &cobra.Command{
		Use:   "basalt",
		Short: "Basalt is a beautiful backend code generator",
		Long:  `Basalt reads .bs DSL files and generates production-ready backends. Built with Bubble Tea.`,
	}

	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Interactive project wizard",
		RunE: func(cmd *cobra.Command, args []string) error {
			results, confirmed, err := tui.RunInitWizard()
			if err != nil {
				return err
			}

			if !confirmed {
				fmt.Println("Project initialization cancelled.")
				return nil
			}

			cfg := config.BasaltConfig{
				Project:   results["project"],
				Database:  results["db"],
				Framework: results["framework"],
				Auth:      results["auth"],
				Language:  results["lang"],
			}

			if err := config.WriteConfig(cfg); err != nil {
				return err
			}

			if err := config.WriteStarterDSL(cfg); err != nil {
				return err
			}

			fmt.Println("\n ✨ Project " + cfg.Project + " initialized!")
			fmt.Println(" 📄 Created basalt.config.json and main.bs")
			fmt.Println("\n Next steps:")
			fmt.Println("   1. Edit main.bs to define your schema")
			fmt.Println("   2. Run 'basalt generate' to build your backend")
			return nil
		},
	}

	var generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate from main.bs",
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath, _ := cmd.Flags().GetString("file")
			if filePath == "" {
				filePath = "main.bs"
			}
			outDir, _ := cmd.Flags().GetString("out")
			if outDir == "" {
				outDir = "generated"
			}

			data, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}

			l := parser.NewLexer(string(data))
			p := parser.NewParser(l)
			file, err := p.ParseFile()
			if err != nil {
				fmt.Printf(" ✗ Parse error  %v\n", err)
				return nil
			}

			prog := tui.NewProgressModel()
			tp := tea.NewProgram(prog)

			g := generator.NewGenerator(file, outDir)
			g.OnProgress = func(msg string) {
				tp.Send(tui.ProgressMsg(msg))
			}

			go func() {
				err := g.Generate()
				tp.Send(tui.DoneMsg(err))
			}()

			_, err = tp.Run()
			return err
		},
	}

	generateCmd.Flags().StringP("file", "f", "main.bs", "DSL file to generate from")
	generateCmd.Flags().StringP("out", "o", "generated", "Output directory")

	var validateCmd = &cobra.Command{
		Use:   "validate",
		Short: "Parse and report errors, no generation",
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath, _ := cmd.Flags().GetString("file")
			if filePath == "" {
				filePath = "main.bs"
			}

			data, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}

			l := parser.NewLexer(string(data))
			p := parser.NewParser(l)
			_, err = p.ParseFile()
			if err != nil {
				fmt.Printf(" ✗ Parse error  %v\n", err)
				return nil
			}

			fmt.Println(" ✓ DSL validated successfully")
			return nil
		},
	}

	validateCmd.Flags().StringP("file", "f", "main.bs", "DSL file to validate")

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Basalt",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Basalt v%s\n", version)
		},
	}

	var devCmd = &cobra.Command{
		Use:   "dev",
		Short: "Watch mode, regenerate on save",
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath, _ := cmd.Flags().GetString("file")
			if filePath == "" {
				filePath = "main.bs"
			}

			watcher, err := fsnotify.NewWatcher()
			if err != nil {
				return err
			}
			defer watcher.Close()

			done := make(chan bool)
			go func() {
				for {
					select {
					case event, ok := <-watcher.Events:
						if !ok {
							return
						}
						if event.Op&fsnotify.Write == fsnotify.Write {
							fmt.Println(" 🔄 Change detected, regenerating...")
							// Re-run generate logic (simplified here)
							data, _ := os.ReadFile(filePath)
							l := parser.NewLexer(string(data))
							p := parser.NewParser(l)
							file, err := p.ParseFile()
							if err == nil {
								g := generator.NewGenerator(file, "generated")
								g.Generate()
								fmt.Println(" ✓ Regenerated")
							} else {
								fmt.Printf(" ✗ Parse error: %v\n", err)
							}
						}
					case err, ok := <-watcher.Errors:
						if !ok {
							return
						}
						fmt.Println("error:", err)
					}
				}
			}()

			err = watcher.Add(filePath)
			if err != nil {
				return err
			}
			fmt.Printf(" 👀 Watching %s for changes...\n", filePath)
			<-done
			return nil
		},
	}

	devCmd.Flags().StringP("file", "f", "main.bs", "DSL file to watch")

	rootCmd.AddCommand(initCmd, generateCmd, devCmd, validateCmd, versionCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
