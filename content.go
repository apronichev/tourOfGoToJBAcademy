package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const goMod = "module main"
const noteDependecyMissing = "<sub>**Note:** if the dependency in the _import_ section is red, click the dependency, press <span class=\"shortcut\">&shortcut:ShowIntentionActions;</span> and select **Sync dependencies of main**.</sub>\n"

// const taskInfoYaml = "type: %s\nfiles:\n  - name: task.go\n    visible: true\n  - name: go.mod\n    visible: false%s\n"
const contentDir = "./content"

//const outputDir = "./output"

// InsertCode inserts code into Go files
func InsertCode(filePath string, codeToInsert string) {
	// Read the existing file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file %s: %s\n", filePath, err)
		return
	}

	// Insert new code
	newContent := append([]byte(codeToInsert), content...)

	// Write new content back to file
	err = os.WriteFile(filePath, newContent, 0644)
	if err != nil {
		fmt.Printf("Error writing to file %s: %s\n", filePath, err)
		return
	}
}

// AppendCode opens the file for appending. Creates the file if it doesn't exist.
func AppendCode(filePath string, codeToInsert string) {
	if codeToInsert == "  - 'For is Go's \"while\"'" {
		codeToInsert = "  - 'For is Go''s \"while\"'"
	}
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Write lines to the file

	_, err = fmt.Fprintln(file, codeToInsert)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

}

func ProcessingGoFiles(l *Lesson, taskPath string, clean string) {
	root := taskPath
	// File to search for
	taskFile := "task.go"
	goModFile := "go.mod"

	mdFile := "task.md"
	taskInfoYamlFile := "task-info.yaml"

	// Traverse the directory tree
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file is named 'task.go'
		if info.Name() == taskFile {
			InsertCode(path, l.Code)
		}

		// Check if the file is named 'go.mod'
		if info.Name() == goModFile {
			InsertCode(path, goMod)
		}

		// Check if the file is named 'task.md'
		if info.Name() == mdFile {
			if HasLine(l.Description, `.image /tour/static/img/tree.png`) {
				srcPath, _ := FindFile(contentDir, "tree.png")
				replacements := map[string]string{
					".image /tour/static/img/tree.png":                                                 "![Tree Image](tree.png)",
					"\nContinue description on [[javascript:click('.next-page')][next page]].\n\n\n\n": "",
				}
				for src, replace := range replacements {
					l.Description, _ = ReplaceTextInFile(l.Description, src, replace)
				}
				err := copyFile(srcPath, filepath.Dir(path)+"/tree.png")
				if err != nil {
					return err
				}
			}
			InsertCode(path, l.Description)
			if HasLine(l.Code, `"golang\.org/x/tour/`) {
				AppendCode(path, noteDependecyMissing)
			}
			if HasLine(l.Description, `If you omit the loop condition it loops forever`) {
				AppendCode(path, "\n<sub>**Hint:** to stop the infinite loop, press <span class=\"shortcut\">&shortcut:Stop;</span></sub>.")
			}
			if l.Solution != "No solution file found." {
				AppendCode(path, "<div class=\"hint\" title=\"Click to see possible solution\">\n\n"+l.Solution+"\n</div>")
			}
			if l.Notes != "" {
				AppendCode(path, l.Notes)
			}
		}

		// Check if the file is named 'task-info.yaml'
		if info.Name() == taskInfoYamlFile {
			if CheckOutputExercises(clean) {
				taskInfo := ReplaceTypePlaceholder("output", "  - name: test/output.txt\n    visible: false")
				InsertCode(path, taskInfo)
			} else if HasLine(clean, "Equivalent Binary Trees") {
				taskInfo := ReplaceTypePlaceholder("output", "  - name: tree.png\n    visible: false\n  - name: test/output.txt\n    visible: false")
				InsertCode(path, taskInfo)
			} else {
				taskInfo := ReplaceTypePlaceholder("theory", "")
				InsertCode(path, taskInfo)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error scanning directory:", err)
	}
}

func CheckOutputExercises(value string) bool {
	list := []string{
		"Exported names", "Exercise: Loops and Functions", "Exercise: Slices",
		"Exercise: Maps", "Exercise: Fibonacci closure", "Interfaces",
		"Exercise: Stringers", "Exercise: Errors", "Exercise: Readers",
		"Exercise: rot13Reader", "Exercise: Images",
	}

	for _, item := range list {
		if item == value {
			return true
		}
	}
	return false
}

// ReplaceTypePlaceholder replaces the placeholder in the YAML template with the provided taskType.
func ReplaceTypePlaceholder(taskType string, s string) string {
	const taskInfoYamlTemplate = "type: %s\nfiles:\n  - name: task.go\n    visible: true\n  - name: go.mod\n    visible: false\n%s"
	return fmt.Sprintf(taskInfoYamlTemplate, taskType, s)
}

// Copy the file to the root directory
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}
	fmt.Println("File copied.")

	return nil
}

func ReplaceTextInFile(content, oldText, newText string) (string, error) {
	if content == "" {
		return "", fmt.Errorf("content is empty")
	}

	modifiedContent := strings.ReplaceAll(content, oldText, newText)
	return modifiedContent, nil
}

func HasLine(text string, regEx string) bool {
	// Create a regular expression to match the import line
	re := regexp.MustCompile(regEx)
	return re.MatchString(text)
}
