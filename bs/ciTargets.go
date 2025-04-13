package main

import (
	"context"

	sbbs "github.com/barbell-math/smoothbrain-bs"
)

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
