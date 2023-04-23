package services

import (
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"golearn/api/data"
	models2 "golearn/api/models"
	"log"
	"net/http"
)

func getMovieById(c *gin.Context) {
	title := c.Param("name")

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
	c.IndentedJSON(http.StatusOK, movie)
}
