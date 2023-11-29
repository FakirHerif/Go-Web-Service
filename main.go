package main

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"

	"example.com/webservice/models"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "example.com/webservice/docs"
)

// @title Web Service API
// @version 1.0
// @description This is a sample CRUD API for managing persons.
// @host localhost:8080
// @BasePath /
// @contact.name   Ali
// @contact.url    https://github.com/FakirHerif/Go-Web-Service
// @contact.email  alibasdemir@gmail.com
func main() {

	r := gin.Default()

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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

func handleRequest(f func(*gin.Context), c *gin.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	f(c)
}

func getPersons(c *gin.Context) {

	var wg sync.WaitGroup
	wg.Add(1)

	go handleRequest(func(c *gin.Context) {
		persons, err := models.GetPersons(20)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Hata": "Veritabanından kişiler alınamadı"})
			return
		}

		if persons == nil {
			c.JSON(http.StatusBadRequest, gin.H{"Hata": "Kayıt bulunamadı"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": persons})
	}, c, &wg)

	wg.Wait()
}

func getPersonById(c *gin.Context) {

	var wg sync.WaitGroup
	wg.Add(1)

	go handleRequest(func(c *gin.Context) {
		id := c.Param("id")
		person, err := models.GetPersonById(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"HATA": "Veritabanında kişi aranırken bir hata oluştu"})
			return
		}

		if person.FirstName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"Hata": "Kayıt bulunamadı"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": person})
	}, c, &wg)

	wg.Wait()
}

func addPerson(c *gin.Context) {

	var wg sync.WaitGroup
	wg.Add(1)

	go handleRequest(func(c *gin.Context) {
		var json models.Person

		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Hata": err.Error()})
			return
		}

		if json.FirstName == "" || json.LastName == "" || json.Email == "" || json.IpAddress == "" {
			c.JSON(http.StatusBadRequest, gin.H{"Hata": "Geçersiz giriş verisi"})
			return
		}

		success, err := models.AddPerson(json)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Hata": "Kişi eklenirken bir hata oluştu"})
			return
		}

		if success {
			c.JSON(http.StatusOK, gin.H{"MSG": "BAŞARILI !!! PERSON EKLENDİ"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"HATA": "Kişi eklenemedi"})
		}
	}, c, &wg)

	wg.Wait()
}

func updatePerson(c *gin.Context) {

	var wg sync.WaitGroup
	wg.Add(1)

	go handleRequest(func(c *gin.Context) {
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

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Hata": "Kişi güncellenirken bir hata oluştu"})
			return
		}

		if success {
			c.JSON(http.StatusOK, gin.H{"MSG": "BAŞARILI !!! BİLGİLER DEĞİŞTİRİLDİ"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"HATA": "Kişi bilgileri güncellenemedi"})
		}
	}, c, &wg)

	wg.Wait()
}

func deletePerson(c *gin.Context) {

	var wg sync.WaitGroup
	wg.Add(1)

	go handleRequest(func(c *gin.Context) {
		personId, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"HATA": "GEÇERSİZ ID !"})
		}

		success, err := models.DeletePerson(personId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Hata": "Kişi silinirken bir hata oluştu"})
			return
		}

		if success {
			c.JSON(http.StatusOK, gin.H{"MSG": "BAŞARILI !!! BİLGİLER SİLİNDİ"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"HATA": "Bilgiler silinemedi"})
		}
	}, c, &wg)

	wg.Wait()
}

// @Summary Get available options
// @Description Get available options for the API
// @Tags persons
// @Produce plain
// @Success 200 {string} string "Available options for the API"
// @Router /api/v1/person [options]
func options(c *gin.Context) {

	var wg sync.WaitGroup
	wg.Add(1)

	go handleRequest(func(c *gin.Context) {
		secenekler := "200 OK\n" +
			"METOTLAR: GET,POST,PUT,DELETE,OPTIONS\n" +
			"HOST: http://localhost:8080\n"

		c.String(200, secenekler)
	}, c, &wg)

	wg.Wait()
}
