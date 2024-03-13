package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	ignore "github.com/codeskyblue/dockerignore"
	"github.com/spf13/cobra"
)

var ignorePatterns []string
var needToCopyPackages = false

func main() {

	var projects []string
	var docker bool

	rootCmd := &cobra.Command{
		Use:   "prune",
		Short: "Prune Go packages",
		Run: func(cmd *cobra.Command, args []string) {
			runPrune(projects, docker)
		},
	}

	rootCmd.Flags().StringSliceVarP(&projects, "project", "p", []string{}, "project names (required)")
	rootCmd.MarkFlagRequired("project")
	rootCmd.Flags().BoolVarP(&docker, "docker", "d", true, "use Docker")

	dockerignoreFile, err := ignore.ReadIgnoreFile(".dockerignore")
	if err != nil {
		log.Fatal(err)
	}
	ignorePatterns = dockerignoreFile

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error executing command:", err)
		os.Exit(1)
	}
}

func runPrune(projects []string, docker bool) {
	err := os.RemoveAll("out")
	if err != nil {
		fmt.Println("Error removing out folder:", err)
		return
	}

	cmdGowork := exec.Command("go", "work", "use", "-r", ".")
	cmdGowork.Stdout = os.Stdout
	cmdGowork.Stderr = os.Stderr
	err = cmdGowork.Run()
	if err != nil {
		fmt.Println("Error running 'go work use -r .':", err)
		return
	}

	for _, project := range projects {
		pruneProject(project, docker)
	}

	if needToCopyPackages {
		copyAllGoModFiles("packages/", filepath.Join("out", "json", "packages"))
		copyAllGoPackages("packages/", filepath.Join("out", "full", "packages"))
	}

	fmt.Println("Prune completed successfully!")
}

func pruneProject(project string, docker bool) {
	// Run only if current project directory is a Go project
	_, err := os.Stat(filepath.Join("apps", project, "go.mod"))
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Skipping project:", project, "as it is not a Go project")
			return
		}

		log.Panicln("Error checking if project is a Go project:", err)
	}

	needToCopyPackages = true

	// Step 2: Copy go.work out/json and out/full folders
	copyFolder("go.work", "./out/json/go.work")
	copyFolder("go.work", "./out/full/go.work")
	copyFolder("go.work", "./out/go.work")

	// Step 3: Copy go.work.sum file to out/json folder
	copyFile("go.work.sum", "./out/json/go.work.sum")
	copyFile("go.work.sum", "./out/go.work.sum")

	// Step 4: Go to apps/$project folder
	err = os.Chdir(filepath.Join("apps", project))
	if err != nil {
		fmt.Println("Error changing directory to 'apps/"+project+"':", err)
		return
	}

	// Step 5: Copy go.mod file to out/json/apps/$project folder
	copyFile("go.mod", filepath.Join("..", "..", "out", "json", "apps", project, "go.mod"))

	// Step 6: Copy $project folder to out/full/apps/$project
	copyFolder(".", filepath.Join("..", "..", "out", "full", "apps", project))

	// Step 7: Go to packages folder
	err = os.Chdir(filepath.Join("..", ".."))
	if err != nil {
		fmt.Println("Error changing directory to 'packages':", err)
		return
	}
}

func copyFolder(src, dest string) {
	fmt.Println("Copying folder:", src, "to", dest)
	if shouldIgnore(src, ignorePatterns) {
		return
	}

	createDirIfNotExist(filepath.Dir(dest))

	cmd := exec.Command("cp", "-r", src, dest)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error copying folder:", err)
	}
}

func copyFile(src, dest string) {
	fmt.Println("Copying file:", src, "to", dest)
	if shouldIgnore(src, ignorePatterns) {
		fmt.Println("Ignoring file:", src)
		return
	}

	srcFile, err := os.Open(src)
	if err != nil {
		fmt.Println("Error opening source file:", err)
		return
	}
	defer srcFile.Close()

	createDirIfNotExist(filepath.Dir(dest))

	destFile, err := os.Create(dest)
	if err != nil {
		fmt.Println("Error creating destination file:", err)
		return
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		fmt.Println("Error copying file:", err)
		return
	}

	fmt.Println("File copied successfully!")
}

func copyAllGoModFiles(src, dest string) {
	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".mod" {
			dest := filepath.Join(filepath.Dir(dest), path)
			fmt.Println("Copying go.mod file:", path, "to", dest)
			copyFile(path, dest)
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error copying go.mod files:", err)
	}
}

func copyAllGoPackages(src, dest string) {
	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".mod" {
			path := filepath.Dir(path)
			dest := filepath.Join(filepath.Dir(dest), path)
			fmt.Println("Copying go package:", path, "to", dest)
			copyFolder(path, dest)
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error copying go packages:", err)
	}
}

func shouldIgnore(file string, patterns []string) bool {
	isSkip, err := ignore.Matches(file, patterns)
	if err != nil {
		log.Fatal(err)
	}

	return isSkip
}

// function that recursively creates directories if it does not exist
func createDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}
