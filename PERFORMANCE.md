# URLVerify Performance Benchmarks

This document provides performance benchmarks for the URLVerify library to help you understand the performance characteristics when processing text with and without URLs.

## Running Benchmarks

To run the benchmarks yourself:

```bash
# Run all benchmarks
go test -bench=. -benchmem

# Run only comparison benchmarks
go test -bench=BenchmarkComparison -benchmem
```

## Performance Summary

### Key Metrics (Apple M1 Pro)

| Scenario | Time per Operation | Memory per Operation | Allocations per Operation |
|----------|-------------------|---------------------|---------------------------|
| **Short text without URLs** | ~3,000 ns | 0 B | 0 |
| **Short text with URLs** | ~3,160 ns | 796 B | 11 |
| **Medium text without URLs** | ~21,800 ns | 0 B | 0 |
| **Medium text with URLs** | ~19,660 ns | 3,730 B | 48 |

### Performance Insights

1. **Minimal Overhead**: Processing text with URLs adds only ~4.8% overhead for short texts
2. **Zero Memory Impact**: Text without URLs uses no memory allocations
3. **Regex Dominates**: The regex matching phase accounts for most of the processing time
4. **Linear Scaling**: Performance scales approximately linearly with text length
5. **Validation Efficiency**: Individual domain validation is very fast (~300-615 ns)

## Detailed Benchmark Results

### Text Processing Benchmarks

- `BenchmarkExtractAll_NoURLs`: Large text without any URLs (~46,000 ns/op)
- `BenchmarkExtractAll_FewURLs`: Text with 4-5 URLs (~18,600 ns/op)  
- `BenchmarkExtractAll_ManyURLs`: Text with 30+ URLs (~34,400 ns/op)
- `BenchmarkExtractAll_Mixed`: Mixed valid/invalid content (~19,400 ns/op)

### Domain Validation Benchmarks

- `BenchmarkValidateDomain_ICANN`: ICANN domains (~425 ns/op)
- `BenchmarkValidateDomain_URL`: Full URLs (~312 ns/op)
- `BenchmarkValidateDomain_IP`: IP addresses (~334 ns/op)
- `BenchmarkValidateDomain_DynamicDNS`: Dynamic DNS (~615 ns/op)
- `BenchmarkValidateDomain_Invalid`: Invalid domains (~401 ns/op)

### Throughput Estimates

Based on the benchmarks, you can expect:

- **1 million short texts without URLs**: ~3 seconds
- **1 million short texts with URLs**: ~3.2 seconds  
- **100,000 medium texts without URLs**: ~2.2 seconds
- **100,000 medium texts with URLs**: ~2.0 seconds

## Memory Usage

- **No URLs**: Zero memory allocations
- **With URLs**: Memory usage scales with the number of URLs found
- **Short text with URLs**: ~800 bytes, 11 allocations
- **Medium text with URLs**: ~3.7 KB, 48 allocations

## Platform Notes

Benchmarks were run on Apple M1 Pro. Performance may vary on different architectures and operating systems. Run the benchmarks on your target platform for accurate measurements.
