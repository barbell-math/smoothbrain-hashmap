goos: linux
goarch: amd64
pkg: github.com/barbell-math/smoothbrain-hashmap
cpu: AMD Ryzen 9 9950X 16-Core Processor            
BenchmarkDifferentGrowthFactors/85%_Full_Growth_Factor-16         	   22273	     54209 ns/op	  108896 B/op	      13 allocs/op
BenchmarkDifferentGrowthFactors/70%_Full_Growth_Factor-16         	   23833	     50034 ns/op	  108896 B/op	      13 allocs/op
BenchmarkDifferentGrowthFactors/65%_Full_Growth_Factor-16         	   24360	     49162 ns/op	  108896 B/op	      13 allocs/op
BenchmarkDifferentGrowthFactors/60%_Full_Growth_Factor-16         	   24860	     48228 ns/op	  108896 B/op	      13 allocs/op
BenchmarkDifferentGrowthFactors/55%_Full_Growth_Factor-16         	   24882	     48103 ns/op	  108896 B/op	      13 allocs/op
BenchmarkDifferentGrowthFactors/50%_Full_Growth_Factor-16         	   25516	     46887 ns/op	  108896 B/op	      13 allocs/op
BenchmarkAgainstMapInsertOnly/1e2_elements/Custom_Map-16          	  125029	      9421 ns/op	   17504 B/op	       9 allocs/op
BenchmarkAgainstMapInsertOnly/1e2_elements/Builtin_Map-16         	  126350	      9324 ns/op	    9832 B/op	      10 allocs/op
BenchmarkAgainstMapInsertOnly/1e3_elements/Custom_Map-16          	   37659	     31695 ns/op	  103520 B/op	      12 allocs/op
BenchmarkAgainstMapInsertOnly/1e3_elements/Builtin_Map-16         	   26481	     45218 ns/op	   79640 B/op	      21 allocs/op
BenchmarkAgainstMapInsertOnly/1e4_elements/Custom_Map-16          	    2730	    455848 ns/op	 1578085 B/op	      16 allocs/op
BenchmarkAgainstMapInsertOnly/1e4_elements/Builtin_Map-16         	    3351	    359233 ns/op	  596858 B/op	      80 allocs/op
BenchmarkAgainstMapInsertOnly/1e6_elements/Custom_Map-16          	      18	  63916064 ns/op	100668533 B/op	      22 allocs/op
BenchmarkAgainstMapInsertOnly/1e6_elements/Builtin_Map-16         	      20	  60647586 ns/op	75537062 B/op	    8201 allocs/op
BenchmarkAgainstMapGetOnly/1e2_elements/Custom_Map-16             	  168049	      6983 ns/op	    5376 B/op	       1 allocs/op
BenchmarkAgainstMapGetOnly/1e2_elements/Builtin_Map-16            	  177687	      6777 ns/op	    5376 B/op	       1 allocs/op
BenchmarkAgainstMapGetOnly/1e3_elements/Custom_Map-16             	   80259	     14859 ns/op	    5376 B/op	       1 allocs/op
BenchmarkAgainstMapGetOnly/1e3_elements/Builtin_Map-16            	  105906	     11367 ns/op	    5376 B/op	       1 allocs/op
BenchmarkAgainstMapGetOnly/1e4_elements/Custom_Map-16             	   13452	     89079 ns/op	    5376 B/op	       1 allocs/op
BenchmarkAgainstMapGetOnly/1e4_elements/Builtin_Map-16            	   17083	     70361 ns/op	    5376 B/op	       1 allocs/op
BenchmarkAgainstMapGetOnly/1e6_elements/Custom_Map-16             	      21	  49088645 ns/op	    5376 B/op	       1 allocs/op
BenchmarkAgainstMapGetOnly/1e6_elements/Builtin_Map-16            	      73	  15909747 ns/op	    5376 B/op	       1 allocs/op
BenchmarkAgainstMapRemoveOnly/1e2_elements/Custom_Map-16          	  175966	      6775 ns/op	    5376 B/op	       1 allocs/op
BenchmarkAgainstMapRemoveOnly/1e2_elements/Builtin_Map-16         	  176710	      6842 ns/op	    5376 B/op	       1 allocs/op
BenchmarkAgainstMapRemoveOnly/1e3_elements/Custom_Map-16          	  114271	     10523 ns/op	    5376 B/op	       1 allocs/op
BenchmarkAgainstMapRemoveOnly/1e3_elements/Builtin_Map-16         	  104463	     11505 ns/op	    5376 B/op	       1 allocs/op
BenchmarkAgainstMapRemoveOnly/1e4_elements/Custom_Map-16          	   15175	     78803 ns/op	    5401 B/op	       1 allocs/op
BenchmarkAgainstMapRemoveOnly/1e4_elements/Builtin_Map-16         	   15784	     75949 ns/op	    5376 B/op	       1 allocs/op
BenchmarkAgainstMapRemoveOnly/1e6_elements/Custom_Map-16          	      28	  46534592 ns/op	    5376 B/op	       1 allocs/op
BenchmarkAgainstMapRemoveOnly/1e6_elements/Builtin_Map-16         	      82	  12697105 ns/op	    5376 B/op	       1 allocs/op
PASS
ok  	github.com/barbell-math/smoothbrain-hashmap	36.234s
