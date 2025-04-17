package main

import (
	"context"

	sbbs "github.com/barbell-math/smoothbrain-bs"
)

func registerBenchTargets() {
	sbbs.RegisterTarget(
		context.Background(),
		"bench",
		sbbs.CdToRepoRoot(),
		sbbs.Stage(
			"Run go bench",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				arg := "nosimd"
				if len(cmdLineArgs) != 1 {
					sbbs.LogInfo("Defaulting to non-simd unit tests")
					sbbs.LogInfo("Available bench targets: %v", testTargets)
				} else {
					arg = cmdLineArgs[0]
				}

				switch arg {
				case "nosimd":
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
					sbbs.LogErr("An invalid unitTest argument was supplied.")
					sbbs.LogInfo("Usage: ")
					sbbs.LogInfo("\t./bs bench %v", testTargets)
					sbbs.LogQuietInfo("Consider: Re-running with a valid unit test argument")
					return sbbs.StopErr
				}
			},
		),
	)
}
