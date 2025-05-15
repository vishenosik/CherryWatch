package collections

import (
	"fmt"
	"math/rand"
	"slices"
	"testing"

	"github.com/pkg/profile"
	"github.com/stretchr/testify/require"

	"github.com/brianvoe/gofakeit/v6"
)

func TestFilter(t *testing.T) {

	numbers := intSlice(100_000_000)

	filtered := Filter(Iter(numbers), func(i int) bool {
		return i%2 == 0
	})

	s := slices.Collect(filtered)

	t.Log(len(s))

}

func TestFilter1(t *testing.T) {

	numbers := intSlice(100_000_000)

	filtered, cnt := FilterCount(Iter(numbers), func(i int) bool {
		return i%2 == 0
	})

	out := make([]int, 0, cnt)
	for i := range filtered {
		out = append(out, i)
	}

	t.Log(len(out))

}

func TestFilter2(t *testing.T) {

	numbers := intSlice(100_000_000)

	var out []int
	for _, n := range numbers {
		if n%2 == 0 {
			out = append(out, n)
		}
	}

	t.Log(len(out))

}

func intSlice(n int) []int {
	out := make([]int, n)
	for i := 0; i < n; i++ {
		out[i] = i
	}
	return out
}

func randomIntSlice(n int) []int {
	slice := make([]int, n)
	for i := 0; i < n; i++ {
		slice[i] = rand.Intn(n / 2) // This will ensure some duplicates
	}
	return slice
}

func randomString(n int) []string {
	slice := make([]string, n)
	for i := 0; i < n; i++ {
		slice[i] = gofakeit.Email()
	}
	return slice
}

func randomStringSlice(n int) []string {
	random := randomString(n)
	slice := make([]string, n)
	for i := 0; i < n; i++ {
		slice[i] = gofakeit.RandomString(random)
	}
	return slice
}

func BenchmarkUnique(b *testing.B) {
	// Start memory profiling
	pp := profile.Start(
		profile.MemProfile, profile.ProfilePath("./pprof.out"),
		profile.MemProfileRate(1), profile.NoShutdownHook,
	)
	defer pp.Stop()

	sizes := []int{
		100_000,
		1_000_000,
	}
	for _, size := range sizes {
		b.Run(fmt.Sprintf("unique_int_size_%d", size), func(b *testing.B) {
			slice := randomIntSlice(size)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				Unique(slice)
			}
		})
		b.Run(fmt.Sprintf("unique_string_size_%d", size), func(b *testing.B) {
			slice := randomStringSlice(size)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				Unique(slice)
			}
		})
	}
}

func TestUnique(t *testing.T) {

	// Start memory profiling
	pp := profile.Start(
		profile.MemProfile, profile.ProfilePath("./pprof.out"),
		profile.MemProfileRate(1), profile.NoShutdownHook,
	)
	defer pp.Stop()

	slice := randomIntSlice(1000000)
	result := Unique(slice)

	require.NotEqual(t, 0, len(result))
}

func TestUniqueIntMemProfile(t *testing.T) {
	// Skip this test unless explicitly run with -test.run=TestUniqueIntMemProfile
	if testing.Short() {
		t.Skip("Skipping memory profile test in short mode")
	}

	// Start memory profiling
	pp := profile.Start(
		profile.MemProfile, profile.ProfilePath("./pprof.out"),
		profile.MemProfileRate(1), profile.NoShutdownHook,
	)
	defer pp.Stop()

	slice := randomIntSlice(1_000_000)

	result := Unique(slice)

	// Use the result to prevent compiler optimization
	if len(result) == 0 {
		t.Fatal("Unexpected empty result")
	}
}
