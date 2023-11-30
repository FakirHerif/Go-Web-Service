package main

import (
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"example.com/webservice/models"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "example.com/webservice/docs"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_request_duration_seconds",
		Help: "Duration of HTTP requests in seconds",
	}, []string{"handler", "method"})

	// Sayaç metriği: Toplam başarılı ve başarısız işlem sayısı için
	crudOperations = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "crud_operations_total",
		Help: "Total number of CRUD operations",
	}, []string{"operation", "status"})
)

func init() {
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(crudOperations)
}

// @title Web Service API
// @version 1.0
// @description This is a sample CRUD API for managing persons. Uses Prometheus for monitoring HTTP request durations and CRUD operations.
// @host localhost:8080
// @BasePath /
// @contact.name   Ali
// @contact.url    https://github.com/FakirHerif/Go-Web-Service
// @contact.email  alibasdemir@gmail.com
func main() {

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()

		requestDuration.WithLabelValues(c.FullPath(), c.Request.Method).Observe(duration)
	})

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

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

	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(1)

	go handleRequest(func(c *gin.Context) {
		persons, err := models.GetPersons(20)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Hata": "Veritabanından kişiler alınamadı"})
			crudOperations.WithLabelValues("GET", "error").Inc()
			return
		}

		if persons == nil {
			c.JSON(http.StatusBadRequest, gin.H{"Hata": "Kayıt bulunamadı"})
			crudOperations.WithLabelValues("GET", "not_found").Inc()
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": persons})
		crudOperations.WithLabelValues("GET", "success").Inc()
	}, c, &wg)

	wg.Wait()

	duration := time.Since(start).Seconds()
	requestDuration.WithLabelValues("/api/v1/person", "GET").Observe(duration)
}

func getPersonById(c *gin.Context) {

	var wg sync.WaitGroup
	wg.Add(1)

	go handleRequest(func(c *gin.Context) {

		start := time.Now()

		id := c.Param("id")
		person, err := models.GetPersonById(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"HATA": "Veritabanında kişi aranırken bir hata oluştu"})
			crudOperations.WithLabelValues("getPersonById", "error").Inc()
			return
		}

		if person.FirstName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"Hata": "Kayıt bulunamadı"})
			crudOperations.WithLabelValues("getPersonById", "not_found").Inc()
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": person})
		crudOperations.WithLabelValues("getPersonById", "success").Inc()

		duration := time.Since(start).Seconds()
		requestDuration.WithLabelValues("/api/v1/person/:id", "GET").Observe(duration)
	}, c, &wg)

	wg.Wait()
}

func addPerson(c *gin.Context) {

	var wg sync.WaitGroup
	wg.Add(1)

	go handleRequest(func(c *gin.Context) {

		start := time.Now()

		var json models.Person

		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Hata": err.Error()})
			crudOperations.WithLabelValues("addPerson", "bad_request").Inc()
			return
		}

		if json.FirstName == "" || json.LastName == "" || json.Email == "" || json.IpAddress == "" {
			c.JSON(http.StatusBadRequest, gin.H{"Hata": "Geçersiz giriş verisi"})
			crudOperations.WithLabelValues("addPerson", "invalid_data").Inc()
			return
		}

		success, err := models.AddPerson(json)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Hata": "Kişi eklenirken bir hata oluştu"})
			crudOperations.WithLabelValues("addPerson", "error").Inc()
			return
		}

		if success {
			c.JSON(http.StatusOK, gin.H{"MSG": "BAŞARILI !!! PERSON EKLENDİ"})
			crudOperations.WithLabelValues("addPerson", "success").Inc()
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"HATA": "Kişi eklenemedi"})
			crudOperations.WithLabelValues("addPerson", "failed").Inc()
		}

		duration := time.Since(start).Seconds()
		requestDuration.WithLabelValues("/api/v1/person", "POST").Observe(duration)
	}, c, &wg)

	wg.Wait()
}

func updatePerson(c *gin.Context) {

	var wg sync.WaitGroup
	wg.Add(1)

	go handleRequest(func(c *gin.Context) {

		start := time.Now()

		var json models.Person

		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"HATA": err.Error()})
			crudOperations.WithLabelValues("updatePerson", "bad_request").Inc()
			return
		}

		personId, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"HATA": "GEÇERSİZ ID !"})
			crudOperations.WithLabelValues("updatePerson", "invalid_id").Inc()
		}

		success, err := models.UpdatePerson(json, personId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Hata": "Kişi güncellenirken bir hata oluştu"})
			crudOperations.WithLabelValues("updatePerson", "error").Inc()
			return
		}

		if success {
			c.JSON(http.StatusOK, gin.H{"MSG": "BAŞARILI !!! BİLGİLER DEĞİŞTİRİLDİ"})
			crudOperations.WithLabelValues("updatePerson", "success").Inc()
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"HATA": "Kişi bilgileri güncellenemedi"})
			crudOperations.WithLabelValues("updatePerson", "failed").Inc()
		}

		duration := time.Since(start).Seconds()
		requestDuration.WithLabelValues("/api/v1/person/:id", "PUT").Observe(duration)
	}, c, &wg)

	wg.Wait()
}

func deletePerson(c *gin.Context) {

	var wg sync.WaitGroup
	wg.Add(1)

	go handleRequest(func(c *gin.Context) {

		start := time.Now()

		personId, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"HATA": "GEÇERSİZ ID !"})
			crudOperations.WithLabelValues("deletePerson", "invalid_id").Inc()
		}

		success, err := models.DeletePerson(personId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Hata": "Kişi silinirken bir hata oluştu"})
			crudOperations.WithLabelValues("deletePerson", "error").Inc()
			return
		}

		if success {
			c.JSON(http.StatusOK, gin.H{"MSG": "BAŞARILI !!! BİLGİLER SİLİNDİ"})
			crudOperations.WithLabelValues("deletePerson", "success").Inc()
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"HATA": "Bilgiler silinemedi"})
			crudOperations.WithLabelValues("deletePerson", "failed").Inc()
		}

		duration := time.Since(start).Seconds()
		requestDuration.WithLabelValues("/api/v1/person/:id", "DELETE").Observe(duration)
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

		start := time.Now()

		secenekler := "200 OK\n" +
			"METOTLAR: GET,POST,PUT,DELETE,OPTIONS\n" +
			"HOST: http://localhost:8080\n"

		c.String(200, secenekler)
		crudOperations.WithLabelValues("options", "success").Inc()

		duration := time.Since(start).Seconds()
		requestDuration.WithLabelValues("/api/v1/person", "OPTIONS").Observe(duration)
	}, c, &wg)

	wg.Wait()
}
