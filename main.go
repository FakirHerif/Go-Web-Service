package main

import (
	"database/sql"
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

	"example.com/webservice/auth"
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

	r.POST("/login", auth.Login)
	r.GET("/secured", auth.TokenAuthMiddleware(), auth.SecuredEndpoint) // TOKEN ÖRNEĞİ: İSTENİLEN ENDPOINT İÇİN auth.TokenAuthMiddleware() KULLANILIR ÖRNEK: v1.GET("person", auth.TokenAuthMiddleware(), getPersons)

	v1 := r.Group("/api/v1")

	{
		v1.GET("person", auth.TokenAuthMiddleware(), getPersons)
		v1.GET("person/:id", auth.TokenAuthMiddleware(), getPersonById)
		v1.POST("person", auth.TokenAuthMiddleware(), addPerson)
		v1.PUT("person/:id", auth.TokenAuthMiddleware(), updatePerson)
		v1.DELETE("person/:id", auth.TokenAuthMiddleware(), deletePerson)
		v1.OPTIONS("person", auth.TokenAuthMiddleware(), options)
		v1.GET("/user", auth.TokenAuthMiddleware(), getUsers)
		v1.GET("/user/:id", auth.TokenAuthMiddleware(), getUserByID)
		v1.POST("/user", auth.TokenAuthMiddleware(), addUser)
		v1.PUT("/user/:id", auth.TokenAuthMiddleware(), updateUser)
		v1.DELETE("/user/:id", auth.TokenAuthMiddleware(), deleteUser)
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
	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(1)

	go handleRequest(func(c *gin.Context) {
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
	}, c, &wg)

	wg.Wait()

	duration := time.Since(start).Seconds()
	requestDuration.WithLabelValues("/api/v1/person/:id", "GET").Observe(duration)
}

func addPerson(c *gin.Context) {
	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(1)

	go handleRequest(func(c *gin.Context) {
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
	}, c, &wg)

	wg.Wait()

	duration := time.Since(start).Seconds()
	requestDuration.WithLabelValues("/api/v1/person", "POST").Observe(duration)
}

func updatePerson(c *gin.Context) {
	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(1)

	go handleRequest(func(c *gin.Context) {
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
	}, c, &wg)

	wg.Wait()

	duration := time.Since(start).Seconds()
	requestDuration.WithLabelValues("/api/v1/person/:id", "PUT").Observe(duration)
}

func deletePerson(c *gin.Context) {
	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(1)

	go handleRequest(func(c *gin.Context) {
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
	}, c, &wg)

	wg.Wait()

	duration := time.Since(start).Seconds()
	requestDuration.WithLabelValues("/api/v1/person/:id", "DELETE").Observe(duration)
}

// @Summary Get available options
// @Description Get available options for the API
// @Tags persons
// @Produce plain
// @Success 200 {string} string "Available options for the API"
// @Router /api/v1/person [options]
func options(c *gin.Context) {
	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(1)

	go handleRequest(func(c *gin.Context) {

		secenekler := "200 OK\n" +
			"METOTLAR: GET,POST,PUT,DELETE,OPTIONS\n" +
			"HOST: http://localhost:8080\n"

		c.String(200, secenekler)
		crudOperations.WithLabelValues("options", "success").Inc()
	}, c, &wg)

	wg.Wait()

	duration := time.Since(start).Seconds()
	requestDuration.WithLabelValues("/api/v1/person", "OPTIONS").Observe(duration)
}

func getUsers(c *gin.Context) {
	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(1)

	go handleRequest(func(c *gin.Context) {
		users, err := models.GetUsers()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Hata": "Kullanıcılar alınamadı"})
			crudOperations.WithLabelValues("GET", "error").Inc()
			return
		}

		if len(users) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"Hata": "Kayıt bulunamadı"})
			crudOperations.WithLabelValues("GET", "not_found").Inc()
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": users})
		crudOperations.WithLabelValues("GET", "success").Inc()
	}, c, &wg)

	wg.Wait()

	duration := time.Since(start).Seconds()
	requestDuration.WithLabelValues("/api/v1/user", "GET").Observe(duration)
}

func getUserByID(c *gin.Context) {
	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(1)

	go handleRequest(func(c *gin.Context) {
		userIDString := c.Param("id")
		userID, err := strconv.Atoi(userIDString)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz Kullanıcı ID'si"})
			return
		}

		user, err := models.GetUserByID(userID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Kullanıcı Bulunamadı"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Kullanıcı çağırılırken hata oluştu"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": user})
	}, c, &wg)

	wg.Wait()

	duration := time.Since(start).Seconds()
	requestDuration.WithLabelValues("/api/v1/user/:id", "GET").Observe(duration)
}

func addUser(c *gin.Context) {
	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(1)

	go handleRequest(func(c *gin.Context) {
		var user models.User

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Hata": err.Error()})
			crudOperations.WithLabelValues("addUser", "bad_request").Inc()
			return
		}

		id, err := models.CreateUser(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Hata": "Kullanıcı eklenemedi"})
			crudOperations.WithLabelValues("addUser", "error").Inc()
			return
		}

		if id != 0 {
			c.JSON(http.StatusOK, gin.H{"message": "Kullanıcı başarıyla eklendi", "id": id})
			crudOperations.WithLabelValues("addUser", "success").Inc()
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"Hata": "Kullanıcı eklenemedi"})
			crudOperations.WithLabelValues("addUser", "failed").Inc()
		}
	}, c, &wg)

	wg.Wait()

	duration := time.Since(start).Seconds()
	requestDuration.WithLabelValues("/api/v1/user", "POST").Observe(duration)
}

func updateUser(c *gin.Context) {
	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(1)

	go handleRequest(func(c *gin.Context) {
		userIDStr := c.Param("id")

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz Kullanıcı ID'si"})
			crudOperations.WithLabelValues("updateUser", "bad_request").Inc()
			return
		}

		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			crudOperations.WithLabelValues("updateUser", "bad_request").Inc()
			return
		}

		user.ID = userID

		err = models.UpdateUser(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Kullanıcı güncellenemedi"})
			crudOperations.WithLabelValues("updateUser", "error").Inc()
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Kullanıcı başarıyla güncellendi"})
		crudOperations.WithLabelValues("updateUser", "success").Inc()
	}, c, &wg)

	wg.Wait()

	duration := time.Since(start).Seconds()
	requestDuration.WithLabelValues("/api/v1/user/:id", "PUT").Observe(duration)
}

func deleteUser(c *gin.Context) {
	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(1)

	go handleRequest(func(c *gin.Context) {
		userID := c.Param("id")
		id, _ := strconv.Atoi(userID)

		err := models.DeleteUser(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Hata": "Kullanıcı silinemedi"})
			crudOperations.WithLabelValues("deleteUser", "error").Inc()
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Kullanıcı başarıyla silindi"})
		crudOperations.WithLabelValues("deleteUser", "success").Inc()
	}, c, &wg)

	wg.Wait()

	duration := time.Since(start).Seconds()
	requestDuration.WithLabelValues("/api/v1/user/:id", "DELETE").Observe(duration)
}
