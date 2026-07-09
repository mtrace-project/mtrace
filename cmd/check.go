/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/moby/moby/client"
	"github.com/spf13/cobra"
	"gitlab.m31.com/m31/academy/devops/cloud-trace-testing/mtrace/export"
	idgenerator "gitlab.m31.com/m31/academy/devops/cloud-trace-testing/mtrace/idGenerator"
	"gitlab.m31.com/m31/academy/devops/cloud-trace-testing/mtrace/parser"
	"gitlab.m31.com/m31/academy/devops/cloud-trace-testing/mtrace/test"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check the syntax of a test case file",
	Long:  `Check the syntax of a test case file by validating its structure and content.`,
	RunE:  CheckTestCase,
}

func init() {
	rootCmd.AddCommand(checkCmd)
}

func CheckTestCase(cmd *cobra.Command, args []string) error {
	testPaths, err := getTestPaths(args)
	if err != nil {
		return err
	}

	dtos, err := parser.ParseTests(testPaths)
	if err != nil {
		return err
	}

	var ctx context.Context
	if cmd != nil {
		ctx = cmd.Context()
	} else {
		ctx = context.Background()
	}

	traceAdapter, err := Config.NewTraceAdapterFromConfig(ctx)
	if err != nil {
		return fmt.Errorf("error creating trace adapter: %w", err)
	}

	idGenerator := &idgenerator.IdGeneratorV1{}

	cli, err := client.New(
		client.FromEnv,
	)
	if err != nil {
		return fmt.Errorf("error creating Docker client: %w", err)
	}

	var suitesList []*test.TestSuite

	for _, dto := range dtos {
		dtoPath := filepath.Dir(dto.FilePath)
		_, err := test.NewTraceTest(dto, idGenerator, cli, traceAdapter, test.TraceTestOptions{}, dtoPath, ctx)
		var res *test.TestResult
		if err != nil {
			slog.Info("Test case file is invalid", "testName", dto.Name, "error", err)

			res = &test.TestResult{
				Passed: false,
				Args:   []any{"message", "Test case file is invalid", "error", err.Error()},
			}
		} else {
			slog.Info("Test case file is valid", "testName", dto.Name)
			res = &test.TestResult{
				Passed: true,
			}
		}
		suitesList = append(suitesList, test.NewTestSuite(
			dto.Name,
			[]*test.TestResult{res},
		))
	}

	if !Config.Quiet {
		err = export.DisplayTestsSummary(suitesList, "VALID", "INVALID")
		if err != nil {
			slog.Error("error displaying the tests summary", "error", err)
		}
	}

	return nil
}
