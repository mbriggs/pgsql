package pgsql_test

import (
	"testing"

	"github.com/jackc/pgsql"
	"github.com/stretchr/testify/assert"
)

func TestPlaceholder(t *testing.T) {
	assert.Equal(t, "$42", pgsql.Placeholder(42).String())
}

func TestArgs(t *testing.T) {
	args := &pgsql.Args{}
	assert.Len(t, args.Values(), 0)

	placeholder := args.Use(42)
	assert.Equal(t, pgsql.Placeholder(1), placeholder)
	assert.Equal(t, []interface{}{42}, args.Values())

	placeholder = args.Use(7)
	assert.Equal(t, pgsql.Placeholder(2), placeholder)
	assert.Equal(t, []interface{}{42, 7}, args.Values())

	placeholder = args.Use(42)
	assert.Equal(t, pgsql.Placeholder(1), placeholder)
	assert.Equal(t, []interface{}{42, 7}, args.Values())

	assert.Equal(t, "array[$3, $1, $2]", args.Format("array[?, ?, ?]", 1, 42, 7))
}

func BenchmarkArgs_1_Uses_1_Values(b *testing.B) {
	benchmarkArgs(b, 1, 1)
}

func BenchmarkArgs_5_Uses_5_Values(b *testing.B) {
	benchmarkArgs(b, 5, 5)
}

func BenchmarkArgs_5_Uses_4_Values(b *testing.B) {
	benchmarkArgs(b, 5, 4)
}

func BenchmarkArgs_10_Uses_10_Values(b *testing.B) {
	benchmarkArgs(b, 10, 10)
}

func BenchmarkArgs_10_Uses_9_Values(b *testing.B) {
	benchmarkArgs(b, 10, 9)
}

func BenchmarkArgs_10_Uses_5_Values(b *testing.B) {
	benchmarkArgs(b, 10, 5)
}

func BenchmarkArgs_30_Uses_30_Values(b *testing.B) {
	benchmarkArgs(b, 30, 30)
}

func BenchmarkArgs_30_Uses_27_Values(b *testing.B) {
	benchmarkArgs(b, 30, 27)
}

func BenchmarkArgs_30_Uses_15_Values(b *testing.B) {
	benchmarkArgs(b, 10, 5)
}

func BenchmarkArgs_100_Uses_100_Values(b *testing.B) {
	benchmarkArgs(b, 100, 100)
}

func BenchmarkArgs_100_Uses_97_Values(b *testing.B) {
	benchmarkArgs(b, 100, 97)
}

func BenchmarkArgs_100_Uses_50_Values(b *testing.B) {
	benchmarkArgs(b, 100, 50)
}

func BenchmarkArgs_1000_Uses_1000_Values(b *testing.B) {
	benchmarkArgs(b, 1000, 1000)
}

func BenchmarkArgs_1000_Uses_900_Values(b *testing.B) {
	benchmarkArgs(b, 1000, 900)
}

func BenchmarkArgs_1000_Uses_500_Values(b *testing.B) {
	benchmarkArgs(b, 1000, 500)
}

func BenchmarkArgs_50000_Uses_50000_Values(b *testing.B) {
	benchmarkArgs(b, 50000, 50000)
}

func BenchmarkArgs_50000_Uses_45000_Values(b *testing.B) {
	benchmarkArgs(b, 50000, 45000)
}

func BenchmarkArgs_50000_Uses_25000_Values(b *testing.B) {
	benchmarkArgs(b, 50000, 25000)
}

func benchmarkArgs(b *testing.B, UseCount int, valueCount int) {
	// This factors out the allocation from boxing an int in an interface.
	preboxedValues := make([]interface{}, valueCount)
	for i := 0; i < valueCount; i++ {
		preboxedValues[i] = i
	}

	for i := 0; i < b.N; i++ {
		args := &pgsql.Args{}
		for j := 0; j < UseCount; j++ {
			args.Use(preboxedValues[j%valueCount])
		}

		if len(args.Values()) != valueCount {
			b.Fatalf("expected len(args.Values()) to be %d, got %d", valueCount, len(args.Values()))
		}
	}
}
