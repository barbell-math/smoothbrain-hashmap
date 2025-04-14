package main

import (
	sbbs "github.com/barbell-math/smoothbrain-bs"
)

func main() {
	registerMtncTargets()
	registerReadmeTargets()
	registerBuildTargets()
	registerUnitTestTargets()
	registerBenchTargets()
	registerDataCollectionAndGenerationTargets()
	registerCITargets()
	sbbs.Main("bs")
}
