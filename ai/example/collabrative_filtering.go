package example

import (
	"fmt"
	"math"
	"sort"
)

type Ratings map[int]map[int]float64

func cosineSimilarity(ratings Ratings, user1 int, user2 int) float64 {
	sharedMovies := []int{}
	for movie := range ratings[user1] {
		if _, ok := ratings[user2][movie]; ok {
			sharedMovies = append(sharedMovies, movie)
		}
	}

	if len(sharedMovies) == 0 {
		return 0.0
	}

	sum1, sum2, sumSquares := 0.0, 0.0, 0.0
	for _, movie := range sharedMovies {
		rating1 := ratings[user1][movie]
		rating2 := ratings[user2][movie]
		sum1 += rating1
		sum2 += rating2
		sumSquares += rating1 * rating1
	}
	mag1 := math.Sqrt(sumSquares)
	mag2 := math.Sqrt(sum2 * sum2)

	dotProduct := 0.0
	for _, movie := range sharedMovies {
		rating1 := ratings[user1][movie]
		rating2 := ratings[user2][movie]
		dotProduct += rating1 * rating2
	}

	return dotProduct / (mag1 * mag2)
}

func getRecommendations(ratings Ratings, userID int, numRecs int) []int {
	similarities := map[int]float64{}
	for otherID := range ratings {
		if otherID == userID {
			continue
		}
		similarities[otherID] = cosineSimilarity(ratings, userID, otherID)
	}

	recommendations := map[int]float64{}
	for movieID := range ratings[1] {
		ratingSum, simSum := 0.0, 0.0
		for otherID, similarity := range similarities {
			if rating, ok := ratings[otherID][movieID]; ok {
				ratingSum += similarity * rating
				simSum += similarity
			}
		}
		if simSum > 0 {
			recommendations[movieID] = ratingSum / simSum
		}
	}

	sortedRecs := make([]int, 0, len(recommendations))
	for movieID := range recommendations {
		sortedRecs = append(sortedRecs, movieID)
		sort.Slice(sortedRecs, func(i, j int) bool {
			return recommendations[sortedRecs[i]] > recommendations[sortedRecs[j]]
		})
	}

	if numRecs < len(sortedRecs) {
		return sortedRecs[:numRecs]
	}
	return sortedRecs
}

func main() {
	// Create a sample ratings dataset
	ratings := Ratings{
		1: {1: 5, 2: 3, 3: 4, 4: 4},
		2: {1: 3, 2: 1, 3: 2, 4: 3, 5: 3},
		3: {1: 4, 2: 3, 3: 4, 5: 5},
		4: {2: 4, 3: 5, 4: 4, 5: 3},
		5: {1: 3, 3: 4, 4: 2, 5: 1},
	}

	// Get recommendations for user 1
	recs := getRecommendations(ratings, 1, 3)

	// Print the recommended movie IDs
	fmt.Println(recs)
}
