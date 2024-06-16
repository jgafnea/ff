package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func main() {
	var (
		directory  string
		ignoreDirs []string
	)

	var rootCmd = &cobra.Command{
		Use:   "ff [pattern]",
		Short: "ff finds files matching a pattern",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			pattern := strings.ToLower(args[0])
			ignoreMap := make(map[string]struct{})
			for _, dir := range ignoreDirs {
				absDir, err := filepath.Abs(dir)
				if err != nil {
					fmt.Printf("Error resolving absolute path for %s: %v\n", dir, err)
					os.Exit(1)
				}
				ignoreMap[absDir] = struct{}{}
			}

			err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				absPath, err := filepath.Abs(path)
				if err != nil {
					return err
				}

				if _, ignored := ignoreMap[filepath.Dir(absPath)]; ignored {
					if info.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}

				if !info.IsDir() && strings.Contains(strings.ToLower(info.Name()), pattern) {
					fmt.Println(path)
				}
				return nil
			})

			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		},
	}

	rootCmd.Flags().StringVarP(&directory, "directory", "d", ".", "Directory to search in")
	rootCmd.Flags().StringSliceVarP(&ignoreDirs, "ignoredirectory", "i", nil, "Directories to ignore")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
