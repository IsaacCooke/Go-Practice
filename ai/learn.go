package ai

import (
	"fmt"
	"math"
	"sort"
)

type Rating struct {
	MovieID int
	Score   float64
}

type User struct {
	ID      int
	Ratings []Rating
	Name    string
}

type Movie struct {
	ID   int
	Name string
}

type Users []User

type Movies []Movie

func (u Users) findUserByID(id int) *User {
	for i := range u {
		if u[i].ID == id {
			return &u[i]
		}
	}
	return nil
}

func (u Users) sharedMovies(user1ID, user2ID int) []int {
	user1 := u.findUserByID(user1ID)
	user2 := u.findUserByID(user2ID)

	movies := make(map[int]bool)
	for _, r := range user1.Ratings {
		movies[r.MovieID] = true
	}
	for _, r := range user2.Ratings {
		if _, ok := movies[r.MovieID]; ok {
			movies[r.MovieID] = false
		}
	}

	shared := make([]int, 0, len(movies))
	for movieID, ok := range movies {
		if !ok {
			shared = append(shared, movieID)
		}
	}

	return shared
}

func findRatingByMovieID(ratings []Rating, movieID int) *Rating {
	for i := range ratings {
		if ratings[i].MovieID == movieID {
			return &ratings[i]
		}
	}
	return nil
}

func (u Users) cosineSimilarity(user1ID, user2ID int) float64 {
	sharedMovies := u.sharedMovies(user1ID, user2ID)
	if len(sharedMovies) == 0 {
		return 0.0
	}

	user1 := u.findUserByID(user1ID)
	user2 := u.findUserByID(user2ID)

	sum1, sum2, sumSquares := 0.0, 0.0, 0.0
	for _, movie := range sharedMovies {
		rating1 := findRatingByMovieID(user1.Ratings, movie)
		rating2 := findRatingByMovieID(user2.Ratings, movie)
		sum1 += rating1.Score
		sum2 += rating2.Score
		sumSquares += rating1.Score * rating1.Score
	}
	mag1 := math.Sqrt(sumSquares)
	mag2 := math.Sqrt(sum2 * sum2)

	dotProduct := 0.0
	for _, movie := range sharedMovies {
		rating1 := findRatingByMovieID(user1.Ratings, movie)
		rating2 := findRatingByMovieID(user2.Ratings, movie)
		dotProduct += rating1.Score * rating2.Score
	}

	return dotProduct / (mag1 * mag2)
}

func getRecommendation(users Users, userID int) []Rating {
	user := users.findUserByID(userID)
	if user == nil {
		return nil
	}

	similarities := make(map[int]float64)
	for _, other := range users {
		if other.ID == userID {
			continue
		}
		similarity := users.cosineSimilarity(userID, other.ID)
		if similarity > 0 {
			similarities[other.ID] = similarity
		}
	}

	recommendations := make(map[int]float64)
	for _, other := range users {
		if other.ID == userID {
			continue
		}
		for _, rating := range other.Ratings {
			//if ok := user.Ratings[rating.MovieID]; ok = nil {
			recommendations[rating.MovieID] += similarities[other.ID] * rating.Score
			//}
		}
	}

	// convert the recommendations map to a slice of Rating structs
	recs := make([]Rating, 0, len(recommendations))
	for movieID, score := range recommendations {
		recs = append(recs, Rating{MovieID: movieID, Score: score})
	}

	// sort the recommendations by score in descending order
	sort.Slice(recs, func(i, j int) bool {
		return recs[i].Score > recs[j].Score
	})

	return recs
}

func Learn() {
	users := Users{
		User{
			ID: 1,
			Ratings: []Rating{
				{
					MovieID: 1,
					Score:   9,
				},
				{
					MovieID: 2,
					Score:   3,
				},
				{
					MovieID: 3,
					Score:   6,
				},
				{
					MovieID: 4,
					Score:   1,
				},
			},
			Name: "User One",
		},
		User{
			ID: 2,
			Ratings: []Rating{
				{
					MovieID: 1,
					Score:   5,
				},
				{
					MovieID: 2,
					Score:   2,
				},
				{
					MovieID: 3,
					Score:   9,
				},
				{
					MovieID: 4,
					Score:   7,
				},
			},
			Name: "User Two",
		},
		User{
			ID: 3,
			Ratings: []Rating{
				{
					MovieID: 1,
					Score:   7,
				},
				{
					MovieID: 2,
					Score:   7,
				},
				{
					MovieID: 3,
					Score:   5,
				},
				{
					MovieID: 4,
					Score:   6,
				},
			},
			Name: "User Three",
		},
		User{
			ID: 4,
			Ratings: []Rating{
				{
					MovieID: 1,
					Score:   8,
				},
				{
					MovieID: 2,
					Score:   9,
				},
				{
					MovieID: 3,
					Score:   2,
				},
			},
			Name: "User Four",
		},
	}

	recommendations := getRecommendation(users, 1)
	fmt.Println(recommendations)
}
