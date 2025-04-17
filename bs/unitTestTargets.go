package main

import (
	"context"

	sbbs "github.com/barbell-math/smoothbrain-bs"
)

func registerUnitTestTargets() {
	sbbs.RegisterTarget(
		context.Background(),
		"test",
		sbbs.CdToRepoRoot(),
		sbbs.Stage(
			"Run go test",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				arg := "nosimd"
				if len(cmdLineArgs) != 1 {
					sbbs.LogInfo("Defaulting to non-simd unit tests")
					sbbs.LogInfo("Available test targets: %v", testTargetsAndAll)
				} else {
					arg = cmdLineArgs[0]
				}

				switch arg {
				case "nosimd":
					return sbbs.RunStdout(ctxt, "go", "test", "-v", "./...")
				case "128":
					return sbbs.RunStdout(ctxt, "go", "test", "-tags=sbmap_simd128", "-v", "./...")
				case "256":
					return sbbs.RunStdout(ctxt, "go", "test", "-tags=sbmap_simd256", "-v", "./...")
				case "512":
					return sbbs.RunStdout(ctxt, "go", "test", "-tags=sbmap_simd512", "-v", "./...")
				case "all":
					err := sbbs.RunStdout(ctxt, "go", "test", "-v", "./...")
					if err != nil {
						return err
					}
					err = sbbs.RunStdout(ctxt, "go", "test", "-tags=sbmap_simd128", "-v", "./...")
					if err != nil {
						return err
					}
					err = sbbs.RunStdout(ctxt, "go", "test", "-tags=sbmap_simd256", "-v", "./...")
					if err != nil {
						return err
					}
					return sbbs.RunStdout(ctxt, "go", "test", "-tags=sbmap_simd512", "-v", "./...")
				default:
					sbbs.LogErr("An invalid unitTest argument was supplied.")
					sbbs.LogInfo("Usage: ")
					sbbs.LogInfo("\t./bs test %v", testTargetsAndAll)
					sbbs.LogQuietInfo("Consider: Re-running with a valid unit test argument")
					return sbbs.StopErr
				}
			},
		),
	)

	sbbs.RegisterTarget(
		context.Background(),
		"testexe",
		sbbs.CdToRepoRoot(),
		sbbs.Stage(
			"Run go test -c",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				arg := "nosimd"
				if len(cmdLineArgs) != 1 {
					sbbs.LogInfo("Defaulting to non-simd unit tests")
					sbbs.LogInfo("Available test targets: %v", testTargets)
				} else {
					arg = cmdLineArgs[0]
				}

				switch arg {
				case "nosimd":
					return sbbs.RunStdout(
						ctxt, "go", "test",
						"-gcflags", "-N", "-ldflags=-compressdwarf=false",
						"-c", "./...",
					)
				case "128":
					return sbbs.RunStdout(
						ctxt, "go", "test",
						"-tags=sbmap_simd128",
						"-gcflags", "-N", "-ldflags=-compressdwarf=false",
						"-c", "./...",
					)
				case "256":
					return sbbs.RunStdout(
						ctxt, "go", "test",
						"-tags=sbmap_simd256",
						"-gcflags", "-N", "-ldflags=-compressdwarf=false",
						"-c", "./...",
					)
				case "512":
					return sbbs.RunStdout(
						ctxt, "go", "test",
						"-tags=sbmap_simd512",
						"-gcflags", "-N", "-ldflags=-compressdwarf=false",
						"-c", "./...",
					)
				default:
					sbbs.LogErr("An invalid unitTest argument was supplied.")
					sbbs.LogInfo("Usage: ")
					sbbs.LogInfo("\t./bs testexe %v", testTargets)
					sbbs.LogQuietInfo("Consider: Re-running with a valid unit test argument")
					return sbbs.StopErr
				}
			},
		),
	)
}
