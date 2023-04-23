package example

import (
	"fmt"
	"math"
	"sort"
)

// Movie Define a struct to represent a movie
type Movie struct {
	Name     string
	Features []float64
}

// Features Nine Features
type Features struct {
	Comedy      float64
	Horror      float64
	Action      float64
	SciFi       float64
	Documentary float64
	Historical  float64
	Drama       float64
	Superhero   float64
	Fantasy     float64
}

// Define a function to calculate the Euclidean distance between two movies
func euclideanDistance(a, b Movie) float64 {
	sum := 0.0
	for i := range a.Features {
		sum += math.Pow(a.Features[i]-b.Features[i], 2)
	}
	return math.Sqrt(sum)
}

// Define a function to find the k nearest neighbors for a given movie
func findKNearestNeighbors(movies []Movie, m Movie, k int) []Movie {
	distances := make(map[float64][]Movie)
	for i := range movies {
		distance := euclideanDistance(movies[i], m)
		distances[distance] = append(distances[distance], movies[i])
	}
	var keys []float64
	for k := range distances {
		keys = append(keys, k)
	}
	sort.Float64s(keys)
	var neighbors []Movie
	for _, key := range keys {
		if len(neighbors) >= k {
			break
		}
		neighbors = append(neighbors, distances[key]...)
	}
	return neighbors
}

// Define a function to recommend movies based on a given movie
func recommendMovies(movies []Movie, m Movie, k int) []Movie {
	neighbors := findKNearestNeighbors(movies, m, k)
	recommendations := make(map[string]Movie)
	for i := range neighbors {
		for j := range neighbors[i].Features {
			if neighbors[i].Features[j] > 0 && m.Features[j] == 0 {
				recommendations[neighbors[i].Name] = neighbors[i]
			}
		}
	}
	var recommendedMovies []Movie
	for _, r := range recommendations {
		recommendedMovies = append(recommendedMovies, r)
	}
	return recommendedMovies
}

func Learn() {
	// Define some sample movies
	movies := []Movie{
		{"The Shawshank Redemption", []float64{9.3, 0, 0, 0, 0, 0, 0, 0, 0}},
		{"The Godfather", []float64{9.2, 0, 0, 0, 0, 0, 0, 0, 0}},
		{"The Dark Knight", []float64{9.0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{"The Godfather: Part II", []float64{9.0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{"12 Angry Men", []float64{8.9, 0, 0, 0, 9, 0, 0, 9, 0}},
		{"Schindler's List", []float64{8.9, 0, 0, 0, 0, 0, 0, 0, 0}},
		{"The Lord of the Rings: The Return of the King", []float64{8.9, 0, 0, 0, 0, 0, 0, 0, 0}},
		{"Pulp Fiction", []float64{8.9, 0, 9, 0, 0, 0, 0, 0, 9}},
	}

	// Define a sample movie for which to recommend similar movies
	m := Movie{"The Dark Knight Rises", []float64{0, 0, 9, 0, 9, 0, 8.4, 9, 9}}

	// Find the movies most similar to the sample movie
	k := 3
	neighbors := findKNearestNeighbors(movies, m, k)

	// Print the names of the most similar movies
	fmt.Printf("The %d most similar movies to %q are:\n", k, m.Name)
	for i := range neighbors {
		fmt.Printf("%d. %s\n", i+1, neighbors[i].Name)
	}

	// Recommend movies based on the sample movie
	recommendedMovies := recommendMovies(movies, m, k)

	// Print the names of the recommended movies
	fmt.Printf("\nRecommended movies based on %q:\n", m.Name)
	for i := range recommendedMovies {
		fmt.Printf("%d. %s\n", i+1, recommendedMovies[i].Name)
	}
}
