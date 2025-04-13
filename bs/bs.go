package main

import (
	"bytes"
	"context"
	"os"
	"strings"

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
				return os.Rename("./bs/testProf.prof", "./bs/defaultProfile.prof")
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
				return os.Rename("./bs/testProf.prof", "./bs/simd128Profile.prof")
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
				return os.Rename("./bs/testProf.prof", "./bs/simd256Profile.prof")
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
				return os.Rename("./bs/testProf.prof", "./bs/simd512Profile.prof")
			},
		),
	)

	sbbs.RegisterTarget(
		context.Background(),
		"collectBenchmarks",
		sbbs.Stage(
			"Run builtin bench",
			func(ctxt context.Context, cmdLineArgs ...string) error {
				builtinResults, err := os.Create("./bs/builtinBenchmarks.txt")
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
				defaultResults, err := os.Create("./bs/defaultBenchmarks.txt")
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
				simd128Results, err := os.Create("./bs/simd128Benchmarks.txt")
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
				simd256Results, err := os.Create("./bs/simd256Benchmarks.txt")
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
				simd512Results, err := os.Create("./bs/simd512Benchmarks.txt")
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
				return sbbs.RunStdout(ctxt, "./bs/plots")
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

	sbbs.Main("bs")
}
