package main

import (
	"context"

	sbbs "github.com/barbell-math/smoothbrain-bs"
)

func registerReadmeTargets() {
	sbbs.RegisterTarget(
		context.Background(),
		"updateReadme",
		sbbs.Stage(
			"Run gomarkdoc",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				err := sbbs.RunStdout(
					ctxt, "gomarkdoc", "--embed", "--output", "README.md", ".",
				)
				if err != nil {
					sbbs.LogQuietInfo("Consider running build system with installGoMarkDoc target if gomarkdoc is not installed")
				}
				return err
			},
		),
	)

	sbbs.RegisterTarget(
		context.Background(),
		"installGoMarkDoc",
		sbbs.Stage(
			"Install gomarkdoc",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				return sbbs.RunStdout(
					ctxt, "go",
					"install", "github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest",
				)
			},
		),
	)
}
