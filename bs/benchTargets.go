package main

import (
	"context"

	sbbs "github.com/barbell-math/smoothbrain-bs"
)

func registerBenchTargets() {
	sbbs.RegisterTarget(
		context.Background(),
		"bench",
		sbbs.Stage(
			"Run go bench",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				if len(cmdLineArgs) != 1 {
					sbbs.LogErr("The bench build target requires one argument.")
					sbbs.LogInfo("Usage: ")
					sbbs.LogInfo("\t./bs bench [default | 128 | 256 | 512]")
					sbbs.LogQuietInfo("Consider: Re-running with a valid unit test argument")
					return sbbs.StopErr
				}

				switch cmdLineArgs[0] {
				case "default":
					return sbbs.RunStdout(
						ctxt, "go", "test",
						"-bench=CustomMap", "-benchmem",
						"./",
					)
				case "128":
					return sbbs.RunStdout(
						ctxt, "go", "test",
						"-tags=sbmap_simd128", "-bench=CustomMap", "-benchmem",
						"./",
					)
				case "256":
					return sbbs.RunStdout(
						ctxt, "go", "test",
						"-tags=sbmap_simd256", "-bench=CustomMap", "-benchmem",
						"./",
					)
				case "512":
					return sbbs.RunStdout(
						ctxt, "go", "test",
						"-tags=sbmap_simd512", "-bench=CustomMap", "-benchmem",
						"./",
					)
				default:
					sbbs.LogErr("An invalid bench argument was supplied.")
					sbbs.LogInfo("Usage: ")
					sbbs.LogInfo("\t./bs bench [default | 128 | 256 | 512]")
					sbbs.LogQuietInfo("Consider: Re-running with a valid unit test argument")
					return sbbs.StopErr
				}
			},
		),
	)
}
