package main

import (
	"bytes"
	"context"
	"strings"

	sbbs "github.com/barbell-math/smoothbrain-bs"
)

func registerMtncTargets() {
	// Register a target that updates all dependences. Dependencies that are in
	// the `barbell-math` repo will always be pinned at latest and all other
	// dependencies will be updated to the latest version.
	sbbs.RegisterTarget(
		context.Background(),
		"updateDeps",
		sbbs.Stage(
			"barbell math package cmds",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				var packages bytes.Buffer
				if err := sbbs.Run(
					ctxt, &packages, "go", "list", "-m", "-u", "all",
				); err != nil {
					return err
				}

				lines := strings.Split(packages.String(), "\n")
				// First line is the current package, skip it
				for i := 1; i < len(lines); i++ {
					iterPackage := strings.SplitN(lines[i], " ", 2)
					if !strings.Contains(iterPackage[0], "barbell-math") {
						continue
					}

					if err := sbbs.RunStdout(
						ctxt, "go", "get", iterPackage[0]+"@latest",
					); err != nil {
						return err
					}
				}
				return nil
			},
		),
		sbbs.Stage(
			"Non barbell math package cmds",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				if err := sbbs.RunStdout(ctxt, "go", "get", "-u", "./..."); err != nil {
					return err
				}
				if err := sbbs.RunStdout(ctxt, "go", "mod", "tidy"); err != nil {
					return err
				}

				return nil
			},
		),
	)
	sbbs.RegisterTarget(
		context.Background(),
		"fmt",
		sbbs.Stage(
			"Run go fmt",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				return sbbs.RunStdout(ctxt, "go", "fmt", "./...")
			},
		),
	)
}
