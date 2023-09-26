package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func GetAndRemoveNotes(lesson *string) string {
	var separatedNotes []string
	separatedNotes = (separatedNotes)[:0]
	currentLesson := *lesson
	joinedNotes := ""

	// Regular expression pattern with multiline flag
	re := regexp.MustCompile(`#appengine: (.*)`)
	de := regexp.MustCompile(`(#appengine: .*)|(#appengine:\n)`)

	// Find all matches
	matches := re.FindAllStringSubmatch(currentLesson, -1)

	// Add the titles groups to the slice
	for _, match := range matches {
		separatedNotes = append(separatedNotes, match[1])
	}

	if de.MatchString(currentLesson) {
		joinedNotes = strings.Join(separatedNotes, " ")
	} else {
		joinedNotes = ""
	}

	// Replace lines matching the regular expression with an empty string
	noNotes := de.ReplaceAllString(currentLesson, "")
	*lesson = noNotes
	return joinedNotes

}

func GetAndRemoveCodePath(lesson *string) string {
	currentLesson := *lesson
	re := regexp.MustCompile(`\.play (.*)$`)

	// Find all match
	match := re.FindStringSubmatch(currentLesson)

	// Replace lines matching the regular expression with an empty string
	noCode := re.ReplaceAllString(currentLesson, "")
	*lesson = noCode
	if match != nil {
		return "content/" + match[1]
	}
	return "content/welcome/hello.go"
}

// GetAndRemoveTitle extracts lines starting with '* ' and returns the content after '* ' until the end of the line
func GetAndRemoveTitle(lesson *string) string {
	currentLesson := *lesson
	re := regexp.MustCompile(`(?m)^\*\s(.+)$`)

	// Find all matches
	match := re.FindStringSubmatch(currentLesson)

	// Replace lines matching the regular expression with an empty string
	noTitle := re.ReplaceAllString(currentLesson, "")
	*lesson = noTitle
	return match[1]
}

func ReadSolutionFile(fileName string) string {
	mapping := map[string]string{
		"exercise-stringer.go":                "stringers.go",
		"exercise-web-crawler.go":             "webcrawler.go",
		"exercise-slices.go":                  "slices.go",
		"exercise-rot-reader.go":              "rot13.go",
		"exercise-reader.go":                  "readers.go",
		"exercise-maps.go":                    "maps.go",
		"exercise-loops-and-functions.go":     "loops.go",
		"exercise-images.go":                  "image.go",
		"exercise-equivalent-binary-trees.go": "binarytrees.go",
		"exercise-errors.go":                  "errors.go",
		"exercise-fibonacci-closure.go":       "fib.go",
	}

	solutionFileName, exists := mapping[fileName]
	if !exists {
		return "No solution file found."
	}

	solutionFilePath := filepath.Join("content/solutions", solutionFileName)
	solutionContent, err := os.ReadFile(solutionFilePath)
	re := regexp.MustCompile(`(?s)package main.*`)
	match := re.FindStringSubmatch(string(solutionContent))
	solution := match[0]
	if err != nil {
		log.Printf("Error reading the solution file: %v", err)
		return "Error reading the solution file."
	}

	solution = addFourSpacesToCode(solution)

	return solution
}

func addFourSpacesToCode(code string) string {
	lines := strings.Split(code, "\n")
	for i, line := range lines {
		lines[i] = "    " + line
	}
	return strings.Join(lines, "\n")
}

func GetCode(l string) string {
	filePath := l
	data, err := os.ReadFile(filePath)
	if err != nil {
		return ""
	}
	re := regexp.MustCompile(`(?s)package main.*`)
	match := re.FindStringSubmatch(string(data))
	lessonText := match[0]
	return lessonText
}

// ExtractLessons extracts groups that start with '*' and end with specific patterns
func ExtractLessons(text *string, lessons *[]string, file string) {

	// Regular expression pattern
	re := regexp.MustCompile(`((?s)\*.*?\.play .*?\.go)`)
	de := regexp.MustCompile(`(?s)\* Where to Go from here.*?Visit \[the Go home page\]\(https\:\/\/go.dev\/\) for more\.`)

	// Find all matches
	matches := re.FindAllStringSubmatch(*text, -1)

	// Clear the lessons slice
	*lessons = (*lessons)[:0]

	// Add the lessons groups to the slice
	for _, match := range matches {
		*lessons = append(*lessons, match[1])
	}
	result := de.FindStringSubmatch(*text)
	if result != nil {
		match := result[0]
		*lessons = append(*lessons, match)
	}
}

func ExtractArticleName(text string) string {
	// Split the text by newline character
	lines := strings.Split(text, "\n")

	// Initialize the article name as an empty string
	articleName := ""

	// Iterate through the lines
	for _, line := range lines {
		// Trim leading whitespace and dots from the line
		trimmedLine := strings.TrimRight(line, ".")

		// If the trimmed line is not empty, set it as the article name and break
		if trimmedLine != "" {
			articleName = trimmedLine
			break
		}
	}

	return articleName
}

func TransformLinks(text string) string {
	// Regular expression to find links in the format [[https://...][text]]
	re := regexp.MustCompile(`\[\[https://([^\]]+)\]\[([^\]]+)\]\]`)
	httpsLinksReplaced := re.ReplaceAllString(text, "[$2](https://$1)")
	tourLinksReplaced := replaceTourLinksWithMap(httpsLinksReplaced)
	goDevLinksReplaced := replacePathLinksToMarkdown(tourLinksReplaced)
	return goDevLinksReplaced
}

func replaceTourLinksWithMap(text string) string {
	// Define a map for link replacements
	linkMap := map[string]string{
		"[[/tour/flowcontrol/8][earlier exercise]]": `[earlier exercise](course://%20Flow%20control%20statements%3A%20for%2C%20if%2C%20else%2C%20switch%20and%20defer/Exercise%3A%20Loops%20and%20Functions/taks.go)`,
		"[[/tour/moretypes/18][picture generator]]": `[picture generator](course://%20Methods%20and%20interfaces/Exercise%3A%20Images/task.go)`,
	}

	// Replace links
	for oldLink, newLink := range linkMap {
		text = strings.ReplaceAll(text, oldLink, newLink)
	}
	return text
}

func replacePathLinksToMarkdown(text string) string {
	// Define the regular expression for matching [[/path...][...]] format
	re := regexp.MustCompile(`\[\[\/([^]]+)\]\[([^]]+)\]\]`)

	// Replace the old-style links with Markdown links
	return re.ReplaceAllStringFunc(text, func(match string) string {
		submatches := re.FindStringSubmatch(match)
		newLink := fmt.Sprintf("[%s](https://go.dev/%s)", submatches[2], submatches[1])
		return newLink
	})
}

func ReplacePatterns(text string) string {
	// List of replacements
	replacements := map[string]string{
		"`(`)`":                            "`( )`",
		"`{`}`":                            "`{ }`",
		"`func`Sqrt`":                      "`func Sqrt`",
		"`switch`true`":                    "`switch true`",
		"_interface_type_":                 "_interface type_",
		"`package`rand`":                   "`package rand`",
		"_type_switch_":                    "_type switch_",
		"`IPAddr{1,`2,`3,`4}`":             "`IPAddr{1, 2, 3, 4}`",
		"cannot`Sqrt`negative`number:`-2":  "cannot Sqrt negative number: -2",
		"`image.Rect(0,`0,`w,`h)`":         "`image.Rect(0, 0, w, h)`",
		"`color.RGBA{v,`v,`255,`255}`":     "`color.RGBA{v, v, 255, 255}`",
		"`for`i`:=`range`c`":               "`for i := range c`",
		"*Another*note:*":                  "*Another note:*",
		"`Same(tree.New(1),`tree.New(1))`": "`Same(tree.New(1), tree.New(1))`",
		"`Same(tree.New(1),`tree.New(2))`": "`Same(tree.New(1), tree.New(2))`",
		"_zero_value_":                     "_zero value_",
		"`var`=`":                          "`var =`",
		"#appengine: You can get started by\n#appengine: [[/doc/install/][installing Go]].\n\n#appengine: Once you have Go installed, the": "You can get started by [[/doc/install/][installing Go]].\n\n",
		"The\n[[/doc/][Go Documentation]] is a great place to\n#appengine: continue.\nstart.":                                              "The [[/doc/][Go Documentation]] is a great place to start.",
		"Visit [[/][the Go home page]] for more.": "Visit [the Go home page](https://go.dev/) for more.",
	}

	for old, changed := range replacements {
		text = strings.ReplaceAll(text, old, changed)
	}

	return text
}
