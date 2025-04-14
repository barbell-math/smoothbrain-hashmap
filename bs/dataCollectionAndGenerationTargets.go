package main

import (
	"context"
	"os"

	sbbs "github.com/barbell-math/smoothbrain-bs"
)

func registerDataCollectionAndGenerationTargets() {
	sbbs.RegisterTarget(
		context.Background(),
		"collectProfiles",
		sbbs.Stage(
			"Default profile",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				if err := sbbs.RunStdout(
					ctxt, "go", "test", "./",
				); err != nil {
					return err
				}
				return os.Rename("./bs/tmp/testProf.prof", "./bs/tmp/defaultProfile.prof")
			},
		),
		sbbs.Stage(
			"simd128 profile",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				if err := sbbs.RunStdout(
					ctxt, "go", "test", "-tags=sbmap_simd128", "./",
				); err != nil {
					return err
				}
				return os.Rename("./bs/tmp/testProf.prof", "./bs/tmp/simd128Profile.prof")
			},
		),
		sbbs.Stage(
			"simd256 profile",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				if err := sbbs.RunStdout(
					ctxt, "go", "test", "-tags=sbmap_simd256", "./",
				); err != nil {
					return err
				}
				return os.Rename("./bs/tmp/testProf.prof", "./bs/tmp/simd256Profile.prof")
			},
		),
		sbbs.Stage(
			"simd512 profile",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				if err := sbbs.RunStdout(
					ctxt, "go", "test", "-tags=sbmap_simd512", "./",
				); err != nil {
					return err
				}
				return os.Rename("./bs/tmp/testProf.prof", "./bs/tmp/simd512Profile.prof")
			},
		),
	)

	sbbs.RegisterTarget(
		context.Background(),
		"collectBenchmarks",
		sbbs.Stage(
			"Run builtin bench",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				builtinResults, err := os.Create("./bs/tmp/builtinBenchmarks.txt")
				if err != nil {
					return err
				}
				defer builtinResults.Close()
				return sbbs.Run(
					ctxt, builtinResults, "go", "test",
					"-bench=BuiltinMap", "-benchmem",
					"./",
				)
			},
		),
		sbbs.Stage(
			"Run default bench",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				defaultResults, err := os.Create("./bs/tmp/defaultBenchmarks.txt")
				if err != nil {
					return err
				}
				defer defaultResults.Close()
				return sbbs.Run(
					ctxt, defaultResults, "go", "test",
					"-bench=CustomMap", "-benchmem",
					"./",
				)
			},
		),
		sbbs.Stage(
			"Run simd128 bench",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				simd128Results, err := os.Create("./bs/tmp/simd128Benchmarks.txt")
				if err != nil {
					return err
				}
				defer simd128Results.Close()
				return sbbs.Run(
					ctxt, simd128Results, "go", "test",
					"-tags=sbmap_simd128", "-bench=CustomMap", "-benchmem",
					"./",
				)
			},
		),
		sbbs.Stage(
			"Run simd256 bench",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				simd256Results, err := os.Create("./bs/tmp/simd256Benchmarks.txt")
				if err != nil {
					return err
				}
				defer simd256Results.Close()
				return sbbs.Run(
					ctxt, simd256Results, "go", "test",
					"-tags=sbmap_simd256", "-bench=CustomMap", "-benchmem",
					"./",
				)
			},
		),
		sbbs.Stage(
			"Run simd512 bench",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				simd512Results, err := os.Create("./bs/tmp/simd512Benchmarks.txt")
				if err != nil {
					return err
				}
				defer simd512Results.Close()
				return sbbs.Run(
					ctxt, simd512Results, "go", "test",
					"-tags=sbmap_simd512", "-bench=CustomMap", "-benchmem",
					"./",
				)
			},
		),
	)

	sbbs.RegisterTarget(
		context.Background(),
		"generatePlots",
		sbbs.Stage(
			"Run plot generator",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				return sbbs.RunStdout(ctxt, "./bs/bin/plots")
			},
		),
	)

	sbbs.RegisterTarget(
		context.Background(),
		"collectBenchmarksAndGeneratePlots",
		sbbs.TargetAsStage("collectBenchmarks"),
		sbbs.TargetAsStage("buildPlots"),
		sbbs.TargetAsStage("generatePlots"),
	)
}
