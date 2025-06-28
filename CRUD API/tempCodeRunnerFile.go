package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Patient model
type Patient struct {
	ID      uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Gender  string `json:"gender"`
	Disease string `json:"disease"`
}

var db *gorm.DB
var err error

func main() {
	// PostgreSQL connection string
	dsn := "host=localhost user=postgres password=itu@2003 dbname=hospitaldbb port=5432 sslmode=disable"

	// Connect to PostgreSQL
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("❌ Failed to connect to the database")
	}

	// Debug check - Confirm connected database
	var currentDB string
	db.Raw("SELECT current_database()").Scan(&currentDB)
	fmt.Println("✅ Connected to database:", currentDB)

	// Auto-migrate the Patient model (creates table if not exists)
	err = db.AutoMigrate(&Patient{})
	if err != nil {
		panic("❌ Failed to migrate database schema")
	}

	// Initialize Gin router
	router := gin.Default()

	// Routes
	router.GET("/patients", getAllPatients)
	router.GET("/patients/:id", getPatientByID)
	router.POST("/patients", addPatient)
	router.PUT("/patients/:id", updatePatient)
	router.DELETE("/patients/:id", deletePatient)

	// Start server
	fmt.Println("🚀 Server is running at http://localhost:8080")
	router.Run(":8080")
}

////////////////////////////
// Handler Functions
////////////////////////////

// Get all patients
func getAllPatients(c *gin.Context) {
	var patients []Patient
	result := db.Find(&patients)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, patients)
}

// Get patient by ID
func getPatientByID(c *gin.Context) {
	id := c.Param("id")
	var patient Patient
	result := db.First(&patient, id)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Patient not found"})
		return
	}

	c.JSON(http.StatusOK, patient)
}

// Add new patient
func addPatient(c *gin.Context) {
	var newPatient Patient
	if err := c.ShouldBindJSON(&newPatient); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := db.Create(&newPatient)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, newPatient)
}

// Update patient
func updatePatient(c *gin.Context) {
	id := c.Param("id")
	var patient Patient

	// Check if patient exists
	if err := db.First(&patient, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Patient not found"})
		return
	}

	// Create a separate variable to bind updated data (without overwriting ID)
	var updatedData Patient
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	patient.Name = updatedData.Name
	patient.Age = updatedData.Age
	patient.Gender = updatedData.Gender
	patient.Disease = updatedData.Disease

	db.Save(&patient)

	c.JSON(http.StatusOK, patient)
}

// Delete patient
func deletePatient(c *gin.Context) {
	id := c.Param("id")
	var patient Patient

	if err := db.First(&patient, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Patient not found"})
		return
	}

	db.Delete(&patient)

	c.JSON(http.StatusOK, gin.H{"message": "Patient deleted successfully"})
}
