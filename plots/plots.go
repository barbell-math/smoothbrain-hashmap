package main

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/tools/benchmark/parse"
)

type (
	benchResult struct {
		mapType      string
		operation    string
		tags         string
		growthFactor int64
		numElements  int64
		iterations   int
		nsPerOp      float64
		bytesPerOp   uint64
		allocPerOp   uint64
	}

	point struct {
		X float64
		Y float64
	}
)

var (
	allBenchResults = []benchResult{}
	rawDataFiles    = map[string]string{
		"./bs/tmp/builtinBenchmarks.txt": "",
		"./bs/tmp/defaultBenchmarks.txt": "default",
		"./bs/tmp/simd128Benchmarks.txt": "simd128",
		"./bs/tmp/simd256Benchmarks.txt": "simd256",
		"./bs/tmp/simd512Benchmarks.txt": "simd512",
	}
)

func parseAllBenchResults() {
	for f, tags := range rawDataFiles {
		file, err := os.Open(f)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		benchmarks, err := parse.ParseSet(file)
		if err != nil {
			fmt.Println("Error parsing benchmark results:", err)
			return
		}

		for _, benchmark := range benchmarks {
			for _, iterBench := range benchmark {
				iterBenchResult := benchResult{
					iterations: iterBench.N,
					nsPerOp:    iterBench.NsPerOp,
					bytesPerOp: iterBench.AllocedBytesPerOp,
					allocPerOp: iterBench.AllocsPerOp,
				}

				parts := strings.Split(iterBench.Name, "/")
				if len(parts) == 4 {
					iterBenchResult.mapType = "builtin"
					iterBenchResult.operation = parts[2]

					intIdxEnd := strings.Index(parts[3], "_")
					iterBenchResult.numElements, _ = strconv.ParseInt(
						parts[3][0:intIdxEnd], 10, 64,
					)
				} else if len(parts) == 6 {
					iterBenchResult.mapType = "custom"
					iterBenchResult.operation = parts[4]
					iterBenchResult.tags = tags

					intIdxEnd := strings.Index(parts[2], "%")
					iterBenchResult.growthFactor, _ = strconv.ParseInt(
						parts[2][0:intIdxEnd], 10, 64,
					)

					intIdxEnd = strings.Index(parts[5], "_")
					iterBenchResult.numElements, _ = strconv.ParseInt(
						parts[5][0:intIdxEnd], 10, 64,
					)
				}

				allBenchResults = append(allBenchResults, iterBenchResult)
			}
		}
	}
}

func makeNsPerOpLinePlot() {
	f, err := os.Create("./bs/tmp/data.dat")
	if err != nil {
		panic(err)
	}

	builtinPoints := []point{}
	for _, v := range allBenchResults {
		if v.mapType == "builtin" && v.operation == "Put" {
			builtinPoints = append(
				builtinPoints, point{X: float64(v.numElements), Y: v.nsPerOp},
			)
		}
	}
	slices.SortFunc[[]point](builtinPoints, func(a, b point) int {
		return int(a.X - b.X)
	})
	f.WriteString("# Builtin map data block\n")
	f.WriteString("# X Y\n")
	for _, p := range builtinPoints {
		f.WriteString(fmt.Sprintf(" %f %f\n", p.X, p.Y))
	}
	f.WriteString("\n\n")

	tags := slices.Collect(maps.Values(rawDataFiles))
	slices.Sort(tags)
	points := []point{}
	for _, iterTag := range tags {
		if iterTag == "" {
			continue
		}

		// TODO - replace with function call to generate all unique growth factors
		for i := int64(50); i < 95; i += 5 {
			points = []point{}
			for _, v := range allBenchResults {
				if v.mapType == "custom" && v.operation == "Put" && v.tags == iterTag && v.growthFactor == i {
					points = append(
						points, point{X: float64(v.numElements), Y: v.nsPerOp},
					)
				}
			}
			slices.SortFunc[[]point](points, func(a, b point) int {
				return int(a.X - b.X)
			})
			for i, _ := range points {
				points[i].Y -= builtinPoints[i].Y
			}

			f.WriteString(fmt.Sprintf("# %s, %d map data block\n", iterTag, i))
			f.WriteString("# X Y\n")
			for _, p := range points {
				f.WriteString(fmt.Sprintf(" %f %f\n", p.X, p.Y))
			}
			f.WriteString("\n\n")
		}
	}

	f.Close()
}

func main() {
	parseAllBenchResults()
	makeNsPerOpLinePlot()
}
