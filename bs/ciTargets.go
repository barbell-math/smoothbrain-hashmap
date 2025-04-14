package main

import (
	"bytes"
	"context"

	sbbs "github.com/barbell-math/smoothbrain-bs"
)

// Generate a stage that runs `git diff` and returns an error if there are any
// differences. This is mainly used by ci targets.
func gitDiffStage(errMessage string, targetToRun string) sbbs.StageFunc {
	return sbbs.Stage(
		"Run Diff",
		func(ctxt context.Context, cmdLineArgs ...string) error {
			var buf bytes.Buffer
			if err := sbbs.Run(ctxt, &buf, "git", "diff"); err != nil {
				return err
			}
			if buf.Len() > 0 {
				sbbs.LogErr(errMessage)
				sbbs.LogQuietInfo(buf.String())
				sbbs.LogErr(
					"Run build system with %s and push any changes",
					targetToRun,
				)
				return sbbs.StopErr
			}
			return nil
		},
	)
}

func registerCITargets() {
	// Registers a target that will update all deps and run a diff to make sure
	// that the commited code is using all of the correct dependencies.
	sbbs.RegisterTarget(
		context.Background(),
		"ciCheckDeps",
		sbbs.TargetAsStage("updateDeps"),
		gitDiffStage("Out of date packages were detected", "updateDeps"),
	)
	// Registers a target that will install gomarkdoc, update the readme, and
	// run a diff to make sure that the commited readme is up to date.
	sbbs.RegisterTarget(
		context.Background(),
		"ciCheckReadme",
		sbbs.TargetAsStage("installGoMarkDoc"),
		sbbs.TargetAsStage("updateReadme"),
		gitDiffStage("Readme is out of date", "updateReadme"),
	)
	// Registers a target that will run go fmt and then run a diff to make sure
	// that the commited code is properly formated.
	sbbs.RegisterTarget(
		context.Background(),
		"ciCheckFmt",
		sbbs.TargetAsStage("fmt"),
		gitDiffStage("Fix formatting to get a passing run!", "fmt"),
	)
	// Registers a target that will run all mergegate checks. This includes:
	//	- checking that the code is formatted
	//	- checking that the readme is up to date
	//	- checking that the dependencies are up to date
	//	- checking that all unit tests pass
	sbbs.RegisterTarget(
		context.Background(),
		"mergegate",
		sbbs.TargetAsStage("ciCheckFmt"),
		sbbs.TargetAsStage("ciCheckReadme"),
		sbbs.TargetAsStage("ciCheckDeps"),
		sbbs.TargetAsStage("unitTests"),
	)

}
