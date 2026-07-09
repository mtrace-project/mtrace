package cmd

import (
	"fmt"
	"io/fs"
	"log/slog"
	"path/filepath"
	"strings"
)

func getTestPaths(args []string) ([]string, error) {
	baseDir := Config.Directory
	if baseDir == "" {
		baseDir = "." // current directory as default
	}

	var testPaths []string

	if len(args) > 0 {
		for _, arg := range args {
			if strings.HasSuffix(arg, "/") {
				return nil, fmt.Errorf("invalid argument: '%s'. Only file names are allowed, not paths", arg)
			}

			testPath := arg
			if !strings.HasSuffix(arg, ".mt.yaml") {
				if filepath.Ext(arg) != "" {
					return nil, fmt.Errorf("invalid argument: '%s'. Only .mt.yaml files are allowed", arg)
				}
				testPath = fmt.Sprintf("%s.mt.yaml", arg)
			}

			fullPath := filepath.Join(baseDir, testPath)
			testPaths = append(testPaths, fullPath)
		}
	} else {
		baseDirStr := baseDir
		if baseDirStr == "." {
			baseDirStr = "the current directory"
		}

		slog.Info("No files specified, executing tests", "baseDir", baseDirStr)

		err := filepath.WalkDir(baseDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("error accessing path '%s': %w", path, err)
			}

			if !d.IsDir() && strings.HasSuffix(d.Name(), ".mt.yaml") {
				testPaths = append(testPaths, path)
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("error while scanning the directory: %w", err)
		}
	}

	return testPaths, nil
}
