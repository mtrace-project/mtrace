package parser

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v3"
)

func ParseTests(filePaths []string) ([]*TestDTO, error) {
	if len(filePaths) == 0 {
		return nil, fmt.Errorf("no .mt.yaml files found to execute")
	}

	var tests []*TestDTO

	for _, filePath := range filePaths {
		testData, err := parseTest(filePath)
		if err != nil {
			return nil, fmt.Errorf("error parsing test file '%s': %w", filePath, err)
		}
		testData.FilePath = filePath
		tests = append(tests, testData)
	}

	return tests, nil
}

func parseTest(filePath string) (*TestDTO, error) {
	file, err := os.Open(filePath) // nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("error opening file '%s': %w", filePath, err)
	}

	defer file.Close() // nolint:errcheck

	var testData *TestDTO

	decoder := yaml.NewDecoder(file)

	if err := decoder.Decode(&testData); err != nil {
		return nil, fmt.Errorf("error parsing YAML file '%s': %w", filePath, err)
	}

	if err := testData.Validate(); err != nil {
		return nil, fmt.Errorf("validation error in file '%s': %w", filePath, err)
	}

	return testData, nil
}
