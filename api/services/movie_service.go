package services

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"golearn/api/data"
	models2 "golearn/api/models"
	"log"
)

var getAllMovies = &graphql.Field{
	Type: graphql.NewList(models2.MovieType),
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		configuration := data.ParseConfiguration()
		driver, err := configuration.NewDriver()
		if err != nil {
			log.Fatal(err)
		}
		defer data.UnsafeClose(driver)

		session := driver.NewSession(neo4j.SessionConfig{
			AccessMode:   neo4j.AccessModeRead,
			DatabaseName: configuration.Database,
		})
		defer data.UnsafeClose(session)

		movies, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			records, err := tx.Run(
				`MATCH (movie:Movie)
                  OPTIONAL MATCH (movie)<-[r]-(person:Person)
                  WITH movie.title as title,
                         collect({name:person.name,
                         job:head(split(toLower(type(r)),'_')),
                         role:r.roles}) as cast 
                  UNWIND cast as c 
                  RETURN title, c.name as name, c.job as job, c.role as role`,
				map[string]interface{}{})
			if err != nil {
				return nil, err
			}
			var result []models2.Movie
			currentMovie := models2.Movie{}
			for records.Next() {
				record := records.Record()
				title, _ := record.Get("title")
				name, _ := record.Get("name")
				job, _ := record.Get("job")
				role, _ := record.Get("role")
				if title.(string) != currentMovie.Title {
					if currentMovie.Title != "" {
						result = append(result, currentMovie)
					}
					currentMovie = models2.Movie{Title: title.(string)}
				}
				switch role.(type) {
				case []interface{}:
					currentMovie.Cast = append(currentMovie.Cast, models2.Person{Name: name.(string), Job: job.(string), Role: data.ToStringSlice(role.([]interface{}))})
				default: // handle nulls or unexpected stuff
					currentMovie.Cast = append(currentMovie.Cast, models2.Person{Name: name.(string), Job: job.(string)})
				}
			}
			if currentMovie.Title != "" {
				result = append(result, currentMovie)
			}
			return result, nil
		})
		if movies == nil {
			return nil, fmt.Errorf("no movies found")
		}
		return movies, nil
	},
}

var getMovieByTitle = &graphql.Field{
	Type: models2.MovieType,
	Args: graphql.FieldConfigArgument{
		"title": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		title := params.Args["title"].(string)

		configuration := data.ParseConfiguration()
		driver, err := configuration.NewDriver()
		if err != nil {
			log.Fatal(err)
		}
		defer data.UnsafeClose(driver)

		session := driver.NewSession(neo4j.SessionConfig{
			AccessMode:   neo4j.AccessModeRead,
			DatabaseName: configuration.Database,
		})
		defer data.UnsafeClose(session)

		movie, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			records, err := tx.Run(
				`MATCH (movie:Movie {title:$title})
				  OPTIONAL MATCH (movie)<-[r]-(person:Person)
				  WITH movie.title as title,
						 collect({name:person.name,
						 job:head(split(toLower(type(r)),'_')),
						 role:r.roles}) as cast 
				  LIMIT 1
				  UNWIND cast as c 
				  RETURN title, c.name as name, c.job as job, c.role as role`,
				map[string]interface{}{"title": title})
			if err != nil {
				return nil, err
			}
			var result models2.Movie
			for records.Next() {
				record := records.Record()
				title, _ := record.Get("title")
				result.Title = title.(string)
				name, _ := record.Get("name")
				job, _ := record.Get("job")
				role, _ := record.Get("role")
				switch role.(type) {
				case []interface{}:
					result.Cast = append(result.Cast, models2.Person{Name: name.(string), Job: job.(string), Role: data.ToStringSlice(role.([]interface{}))})
				default: // handle nulls or unexpected stuff
					result.Cast = append(result.Cast, models2.Person{Name: name.(string), Job: job.(string)})
				}
			}
			return result, nil
		})
		if movie == nil {
			return nil, fmt.Errorf("movie not found")
		}
		return movie, nil
	},
}

var moviesWithinThreeRelations = &graphql.Field{
	Type: graphql.NewList(models2.MovieType),
	Args: graphql.FieldConfigArgument{
		"title": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		title := params.Args["title"].(string)

		configuration := data.ParseConfiguration()
		driver, err := configuration.NewDriver()
		if err != nil {
			log.Fatal(err)
		}
		defer data.UnsafeClose(driver)

		session := driver.NewSession(neo4j.SessionConfig{
			AccessMode:   neo4j.AccessModeRead,
			DatabaseName: configuration.Database,
		})
		defer data.UnsafeClose(session)

		movies, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			records, err := tx.Run(
				`MATCH (initial:Movie {title:$title})-[*1..3]-(movies:Movie)
                  OPTIONAL MATCH (movies)<-[r]-(person:Person)
                  WITH movies.title as title,
                         collect({name:person.name,
                         job:head(split(toLower(type(r)),'_')),
                         role:r.roles}) as cast 
                  UNWIND cast as c 
                  RETURN title, c.name as name, c.job as job, c.role as role`,
				map[string]interface{}{"title": title})
			if err != nil {
				return nil, err
			}
			var result []models2.Movie
			currentMovie := models2.Movie{}
			for records.Next() {
				record := records.Record()
				title, _ := record.Get("title")
				name, _ := record.Get("name")
				job, _ := record.Get("job")
				role, _ := record.Get("role")
				if title.(string) != currentMovie.Title {
					if currentMovie.Title != "" {
						result = append(result, currentMovie)
					}
					currentMovie = models2.Movie{Title: title.(string)}
				}
				switch role.(type) {
				case []interface{}:
					currentMovie.Cast = append(currentMovie.Cast, models2.Person{Name: name.(string), Job: job.(string), Role: data.ToStringSlice(role.([]interface{}))})
				default: // handle nulls or unexpected stuff
					currentMovie.Cast = append(currentMovie.Cast, models2.Person{Name: name.(string), Job: job.(string)})
				}
			}
			if currentMovie.Title != "" {
				result = append(result, currentMovie)
			}
			return result, nil
		})
		if movies == nil {
			return nil, fmt.Errorf("no movies found")
		}
		return movies, nil
	},
}

var moviesByDirector = &graphql.Field{
	Type: graphql.NewList(models2.MovieType),
	Args: graphql.FieldConfigArgument{
		"name": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		name := params.Args["name"].(string)

		configuration := data.ParseConfiguration()
		driver, err := configuration.NewDriver()
		if err != nil {
			log.Fatal(err)
		}
		defer data.UnsafeClose(driver)

		session := driver.NewSession(neo4j.SessionConfig{
			AccessMode:   neo4j.AccessModeRead,
			DatabaseName: configuration.Database,
		})
		defer data.UnsafeClose(session)

		movies, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			records, err := tx.Run(
				`MATCH (p:Person {name: $name})-[:DIRECTED]->(movies:Movie)
				OPTIONAL MATCH (movies)<-[r]-(person:Person)
				WITH movies.title as title,
					collect({name:person.name,
					job:head(split(toLower(type(r)),'_')),
					role:r.roles}) as cast 
				UNWIND cast as c 
				RETURN title, c.name as name, c.job as job, c.role as role`,
				map[string]interface{}{"name": name})
			if err != nil {
				return nil, err
			}
			var result []models2.Movie
			currentMovie := models2.Movie{}
			for records.Next() {
				record := records.Record()
				title, _ := record.Get("title")
				name, _ := record.Get("name")
				job, _ := record.Get("job")
				role, _ := record.Get("role")
				if title.(string) != currentMovie.Title {
					if currentMovie.Title != "" {
						result = append(result, currentMovie)
					}
					currentMovie = models2.Movie{Title: title.(string)}
				}
				switch role.(type) {
				case []interface{}:
					currentMovie.Cast = append(currentMovie.Cast, models2.Person{Name: name.(string), Job: job.(string), Role: data.ToStringSlice(role.([]interface{}))})
				default: // handle nulls or unexpected stuff
					currentMovie.Cast = append(currentMovie.Cast, models2.Person{Name: name.(string), Job: job.(string)})
				}
			}
			if currentMovie.Title != "" {
				result = append(result, currentMovie)
			}
			return result, nil
		})
		if movies == nil {
			return nil, fmt.Errorf("no movies found")
		}
		return movies, nil
	},
}
