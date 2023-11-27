package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"example.com/webservice/models"
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

	err := models.ConnectDatabase()
	checkErr(err)

	r.Run()

}

func checkErr(err error) {
	if err != nil {
		log.Println("Error:", err)
	}
}

func getPersons(c *gin.Context) {

	persons, err := models.GetPersons(20)
	checkErr(err)

	if persons == nil {
		c.JSON(http.StatusBadRequest, gin.H{"Hata": "Kayıt bulunamadı"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": persons})
	}
}

func getPersonById(c *gin.Context) {
	id := c.Param("id")

	person, err := models.GetPersonById(id)
	checkErr(err)

	if person.FirstName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Hata": "Kayıt bulunamadı"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": person})
	}
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
