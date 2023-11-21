package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	v1 := r.Group("/api/v1")

	{
		v1.GET("person", getPersons)
		v1.GET("person/:id", getPersonById)
		v1.POST("person", addPerson)
		v1.PUT("person/:id", updatePerson)
		v1.DELETE("person/:id", deletePerson)
		v1.OPTIONS("person", options)
	}

	r.Run()

}

func getPersons(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "Persons Çağırıldı (GET)"})
}

func getPersonById(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"msg": "ID'si " + id + " Olan Person Çağırıldı (GET)"})
}

func addPerson(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "Person Eklendi (POST)"})
}

func updatePerson(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "Person Güncellendi (PUT)"})
}

func deletePerson(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"msg": "ID'si " + id + " Olan Person Silindi (DELETE)"})
}

func options(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "Options Çağırıldı (OPTIONS)"})
}
