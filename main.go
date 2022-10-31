package main

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v2"

	l "github.com/sirupsen/logrus"
)

// This is the type that the YAML detection rule is unmarshalled into. This
// has been designed with pattern matching via regular expressions to implement
// generic searches that find instances of usages of a library and the secrets
// associated with using said library.
type ApplicationRule struct {
	Example       string `yaml:"example"`
	Language      string `yaml:"language"`
	Library       string `yaml:"library"`
	RuleName      string `yaml:"ruleName"`
	SecretPattern string `yaml:"secretPattern"`
}

// loadRules attempts to read and unmarshal the detection rules from a given directory path
func loadRules(path string) []ApplicationRule {
	rules := []ApplicationRule{}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		l.Error("yamlFile.Get err: ", err)
	}

	for _, f := range files {
		yamlFile, err := ioutil.ReadFile(path + "/" + f.Name())
		if err != nil {
			l.Error("yamlFile.Get err: ", err)
		}
		rule := new(ApplicationRule)
		err = yaml.Unmarshal(yamlFile, rule)
		rules = append(rules, *rule)
	}

	return rules
}

// findSource finds source code for a given directory and file extension
func findSource(dir, ext string) []string {
	var filePaths []string
	filepath.WalkDir(dir, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			filePaths = append(filePaths, s)
		}
		return nil
	})

	return filePaths
}

// scanSource checks the provided source code with the provided detection rules
// and displays matches to the terminal.
func scanSource(rules []ApplicationRule, source []string) {
	for r := range rules {
		for s := range source {
			// Build the regular expression. Panics if the pattern does not compile.
			re := regexp.MustCompile(rules[r].SecretPattern)
			re.SubexpNames()
			// Open the source code file
			body, err := ioutil.ReadFile(source[s])
			if err != nil {
				l.Error("unable to read file: %v", err)
				continue
			}
			// Convert the byte stream to a string
			file := string(body[:])
			// Search the source with the regex from the rule
			result := re.FindAllStringSubmatch(file, -1)
			for k := range result {
				if k == 0 {
					l.Print("Secrets detected!")
					l.Print("Rule: ", rules[r].RuleName)
					l.Print("Language: ", rules[r].Language)
					l.Print("File: ", source[s])
					l.Print("String containing secrets: ", result[k][k])
				}
				break
			}
		}
	}
}

// setupLogging configures logrus to output full timestamp messages
func setupLogging() {
	customFormatter := new(l.TextFormatter)
	l.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true
}

// Entrypoint to the detection program
func main() {
	setupLogging()
	l.Info("Detector started")
	dir, err := os.Getwd()
	if err != nil {
		// Bail out if we can't find where we're working
		l.Fatal(err)
	}

	// Get the detection rules
	rulesDir := dir + "/rules"
	rules := loadRules(rulesDir)
	if len(rules) == 0 {
		l.Fatal("No detection rules found")
	}

	// Filter the rules for each language
	pythonRules := []ApplicationRule{}
	goRules := []ApplicationRule{}
	for f := range rules {
		switch rules[f].Language {
		case "Python":
			pythonRules = append(pythonRules, rules[f])
		case "Go":
			goRules = append(goRules, rules[f])
		default:
			l.Warn("Unknown rule type %s", rules[f].Language)
		}
	}

	// Load up the source code to scan using the detection rules
	pythonSources := findSource(dir+"/fixtures", ".py")
	goSources := findSource(dir+"/fixtures", ".go")

	// Begin the detection process of the rules and the sources
	if len(pythonRules) > 0 && len(pythonSources) > 0 {
		scanSource(pythonRules, pythonSources)
	}
	if len(goRules) > 0 && len(goSources) > 0 {
		scanSource(goRules, goSources)
	}
}
