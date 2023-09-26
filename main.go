package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// Lesson is a struct that holds information about each lesson
type Lesson struct {
	id          int
	Description string
	Notes       string
	FilePath    string
	Title       string
	Code        string
	Solution    string
	yamlTitle   string
}

type ArticleFile struct {
	id         int
	name       string
	shortdescr string
	path       string
	fulltext   string
	yamlName   string
}

func main() {
	var (
		sections    []string
		id          int
		title       string
		description string
		codePath    string
		note        string
		code        string
		solution    string
	)
	var courseInfoYamlText = "type: marketplace\ntitle: Tour of Go\nlanguage: English\nsummary: |-\n  This course copies the official Tour of Go that gives you a comprehensive tour of the Go programming language. The\n  tour covers the most important features of the language. The original Tour of Go is available at \n  (https://go.dev/tour/list)[https://go.dev/tour/list]. Please note that the \"Welcome!\" chapter from the original tour,\n  which focuses on controls and specifics for the web version, has not been included in this course.\nprogramming_language: Go\nmode: Study\ncontent:\n"
	const contentDir = "./content"
	const outputDir = "./output"

	// Delete and create the 'output' directory
	DeleteContents(outputDir)
	CreateDir(outputDir)

	var articles = FindArticleFiles(contentDir)

	// In 'output', create 'course-info.yaml' file
	courseInfoYamlPath := filepath.Join(outputDir, "course-info.yaml")
	CreateFile(courseInfoYamlPath)
	InsertCode(courseInfoYamlPath, courseInfoYamlText)
	length := 0

	for n, article := range articles {
		articleName := filepath.Base(article)
		a := ArticleFile{
			id:         n + 1,
			name:       "",
			yamlName:   "",
			shortdescr: "",
			path:       "",
			fulltext:   "",
		}
		a.path = article
		data, err := os.ReadFile(a.path)
		if err != nil {
			return
		}
		a.fulltext = ReplacePatterns(string(data))
		a.fulltext = TransformLinks(a.fulltext)
		ExtractLessons(&a.fulltext, &sections, articleName)
		a.name = ExtractArticleName(a.fulltext)
		lessonName := "  - " + "'" + a.name + "'"
		lessonDirPath := outputDir + "/" + a.name
		lessonInfoYamlPath := CreateLessonStructure(a.name, outputDir)

		AppendCode(courseInfoYamlPath, lessonName)

		length += len(sections)

		for i, section := range sections {
			id = i + 1
			lesson := Lesson{
				id:          id,
				Title:       title,
				yamlTitle:   "",
				Description: description,
				Notes:       note,
				FilePath:    codePath,
				Code:        code,
				Solution:    solution,
			}
			titleClean := GetAndRemoveTitle(&section)
			lesson.Title = "# " + titleClean + "\n"
			lesson.FilePath = GetAndRemoveCodePath(&section)
			lesson.Notes = GetAndRemoveNotes(&section)
			lesson.Code = GetCode(lesson.FilePath)
			lesson.Solution = ReadSolutionFile(filepath.Base(lesson.FilePath))
			lesson.Description = lesson.Title + section

			taskNumber := "  - " + "'" + titleClean + "'"
			taskPath := lessonDirPath + "/" + titleClean

			CreateTaskStructure(titleClean, lessonDirPath)

			AppendCode(lessonInfoYamlPath, taskNumber)

			ProcessingGoFiles(&lesson, taskPath, titleClean)
		}
	}

	fmt.Printf("Task directories created: %d \n", countDirectories(outputDir)-6)
	fmt.Printf("Sections from ARTICLE files processed: %d", length)
}

func countDirectories(path string) int {
	var count int
	err := filepath.Walk(path, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			count++
		}
		return nil
	})
	if err != nil {
		return 0
	}

	// Subtract 1 to exclude the root directory from the count
	return count - 1
}
