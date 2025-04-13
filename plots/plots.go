package main

import (
	"fmt"
	"os"

	"golang.org/x/tools/benchmark/parse"
)

func main() {
	file, err := os.Open("./bs/defaultBenchmarks.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Parse the benchmark results
	benchmarks, err := parse.ParseSet(file)
	if err != nil {
		fmt.Println("Error parsing benchmark results:", err)
		return
	}

	// Iterate over the parsed benchmarks
	for _, benchmark := range benchmarks {
		for _, iterBench := range benchmark {
			fmt.Printf("Benchmark: %s\n", iterBench.Name)
			fmt.Printf("  Iterations: %d\n", iterBench.N)
			fmt.Printf("  ns/op: %.2f\n", float64(iterBench.NsPerOp))
		}
	}
}
