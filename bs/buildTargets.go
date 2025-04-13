package main

import (
	"context"

	sbbs "github.com/barbell-math/smoothbrain-bs"
)

func registerBuildTargets() {
	sbbs.RegisterTarget(
		context.Background(),
		"buildBs",
		sbbs.Stage(
			"Run go build",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				return sbbs.RunStdout(ctxt, "go", "build", "-o", "./bs", "./bs")
			},
		),
	)

	sbbs.RegisterTarget(
		context.Background(),
		"buildPlots",
		sbbs.Stage(
			"Run go build",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				return sbbs.RunStdout(ctxt, "go", "build", "-o", "./bs", "./plots")
			},
		),
	)
}
