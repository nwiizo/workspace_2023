package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "wiki-toc-generator [wikiRepoURL]",
		Short: "wiki-toc-generator is a tool to clone a GitLab Wiki repository and generate a table of contents",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			wikiRepoURL := args[0]
			cloneDir, err := ioutil.TempDir("", "cloned_wiki")
			if err != nil {
				fmt.Printf("Error creating temporary directory: %v\n", err)
				return
			}
			defer os.RemoveAll(cloneDir)

			err = cloneWiki(wikiRepoURL, cloneDir)
			if err != nil {
				fmt.Printf("Error cloning wiki: %v\n", err)
				return
			}

			toc, err := generateTOC(cloneDir)
			if err != nil {
				fmt.Printf("Error generating TOC: %v\n", err)
				return
			}

			fmt.Println("Table of Contents:")
			fmt.Println(toc)
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func cloneWiki(repoURL, destDir string) error {
	cmd := exec.Command("git", "clone", repoURL, destDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func generateTOC(dir string) (string, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		if filepath.Ext(relPath) != ".md" {
			return nil
		}

		paths = append(paths, relPath)
		return nil
	})
	if err != nil {
		return "", err
	}

	sort.Strings(paths)

	var toc strings.Builder
	for _, relPath := range paths {
		depth := strings.Count(relPath, string(os.PathSeparator))
		indent := ""
		if depth > 0 {
			indent = strings.Repeat("  ", depth-1)
		}
		title := strings.TrimSuffix(filepath.Base(relPath), ".md")
		relPathWithoutExt := strings.TrimSuffix(relPath, ".md")
		toc.WriteString(fmt.Sprintf("%s* [[%s|%s]]\n", indent, title, relPathWithoutExt))
	}

	return toc.String(), nil
}
