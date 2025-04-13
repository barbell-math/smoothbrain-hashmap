package main

import (
	"context"

	sbbs "github.com/barbell-math/smoothbrain-bs"
)

func registerUnitTestTargets() {
	sbbs.RegisterTarget(
		context.Background(),
		"unitTests",
		sbbs.Stage(
			"Run go test",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				if len(cmdLineArgs) != 1 {
					sbbs.LogErr("The unitTests build target requires one argument.")
					sbbs.LogInfo("Usage: ")
					sbbs.LogInfo("\t./bs unitTests [default | 128 | 256 | 512 | all]")
					sbbs.LogQuietInfo("Consider: Re-running with a valid unit test argument")
					return sbbs.StopErr
				}

				switch cmdLineArgs[0] {
				case "default":
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
					sbbs.LogInfo("\t./bs unitTests [default | 128 | 256 | 512 | all]")
					sbbs.LogQuietInfo("Consider: Re-running with a valid unit test argument")
					return sbbs.StopErr
				}
			},
		),
	)

	sbbs.RegisterTarget(
		context.Background(),
		"unitTestExe",
		sbbs.Stage(
			"Run go test -c",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				if len(cmdLineArgs) != 1 {
					sbbs.LogErr("The unitTestExe build target requires one argument.")
					sbbs.LogInfo("Usage: ")
					sbbs.LogInfo("\t./bs unitTestExe [default | 128 | 256 | 512]")
					sbbs.LogQuietInfo("Consider: Re-running with a valid unit test argument")
					return sbbs.StopErr
				}

				switch cmdLineArgs[0] {
				case "default":
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
					sbbs.LogErr("An invalid unitTestExe argument was supplied.")
					sbbs.LogInfo("Usage: ")
					sbbs.LogInfo("\t./bs unitTestExe [default | 128 | 256 | 512]")
					sbbs.LogQuietInfo("Consider: Re-running with a valid unit test argument")
					return sbbs.StopErr
				}
			},
		),
	)
}
