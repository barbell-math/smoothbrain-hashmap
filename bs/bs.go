package main

import (
	sbbs "github.com/barbell-math/smoothbrain-bs"
)

var (
	testTargets       = []string{"nosimd", "128", "256", "512"}
	testTargetsAndAll = append(testTargets, "all")
)

func main() {
	sbbs.RegisterBsBuildTarget()
	sbbs.RegisterUpdateDepsTarget()
	sbbs.RegisterGoMarkDocTargets()
	sbbs.RegisterCommonGoCmdTargets(sbbs.GoTargets{
		GenericFmtTarget: true,
	})
	sbbs.RegisterMergegateTarget(sbbs.MergegateTargets{
		CheckDepsUpdated:     true,
		CheckReadmeGomarkdoc: true,
		CheckFmt:             true,
		CheckUnitTests:       true,
	})

	registerUnitTestTargets()
	registerBenchTargets()
	registerDataCollectionTargets()
	registerPlotTargets()

	// TODO - is this really necessary?
	// sbbs.RegisterTarget(
	// 	context.Background(),
	// 	"collectBenchmarksAndGeneratePlots",
	// 	sbbs.TargetAsStage("collectBenchmarks"),
	// 	sbbs.TargetAsStage("generatePlots"),
	// )

	sbbs.Main("bs")
}
