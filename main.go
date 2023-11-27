package main

import (
	"log"
	"net/http"
	"strconv"

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

	var json models.Person

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Hata": err.Error()})
		return
	}

	success, err := models.AddPerson(json)

	if success {
		c.JSON(http.StatusOK, gin.H{"MSG": "BAŞARILI !!! PERSON EKLENDİ"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"HATA": err})
	}
}

func updatePerson(c *gin.Context) {

	var json models.Person

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"HATA": err.Error()})
		return
	}

	personId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"HATA": "GEÇERSİZ ID !"})
	}

	success, err := models.UpdatePerson(json, personId)

	if success {
		c.JSON(http.StatusOK, gin.H{"MSG": "BAŞARILI !!! BİLGİLER DEĞİŞTİRİLDİ"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"HATA": err})
	}
}

func deletePerson(c *gin.Context) {

	personId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"HATA": "GEÇERSİZ ID !"})
	}

	success, err := models.DeletePerson(personId)

	if success {
		c.JSON(http.StatusOK, gin.H{"MSG": "BAŞARILI !!! BİLGİLER SİLİNDİ"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"HATA": err})
	}
}

func options(c *gin.Context) {
	secenekler := "200 OK\n" +
		"METOTLAR: GET,POST,PUT,DELETE,OPTIONS\n" +
		"HOST: http://localhost:8080\n"

	c.String(200, secenekler)
}
