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

func main() {
	registerMtncTargets()
	registerReadmeTargets()
	registerBuildTargets()
	registerUnitTestTargets()
	registerBenchTargets()
	registerDataCollectionAndGenerationTargets()
	registerCITargets()
	sbbs.Main("bs")
}
