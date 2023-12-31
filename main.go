package main

import (
	"errors"
	"fmt"
	_ "fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

// creating struct
type Student struct {
	FirstName string `json:"firstname,omitempty"`
	LastName  string `json:"LastName,omitempty"`
	Roll_no   int    `json:"roll_no,omitempty"`
	Class     int    `json:"class,omitempty"`
	Phone_No  int    `json:"Phone_No,omitempty"`
}

// declaring struct values
var Students = []Student{
	{FirstName: "Rahul", LastName: "Pandita", Roll_no: 17, Class: 12, Phone_No: 78897},
	{FirstName: "john ", LastName: "Doe", Roll_no: 18, Class: 12, Phone_No: 97887},
	{FirstName: "Db", LastName: "Cooper", Roll_no: 19, Class: 11, Phone_No: 9419},
	{FirstName: "Tony", LastName: "Montana", Roll_no: 21, Class: 11, Phone_No: 2250},
	{FirstName: "papa", LastName: "Franku", Roll_no: 23, Class: 12, Phone_No: 9682},
	{FirstName: "Mah", LastName: "Dank", Roll_no: 32, Class: 12, Phone_No: 1234},
}

// GET request for students
func get_student_data(c *gin.Context) {
	Request.Println("Get student data")
	c.IndentedJSON(http.StatusOK, Students)
	Output.Printf("%v \n \n ", Students)
}

// POST request (add new items)
func add_sudent_data(c *gin.Context) {
	Request.Println("Add student data")
	// create a new student
	var new_student Student

	//using Bindjson to bind the received json
	err := c.BindJSON(&new_student) // returns an error

	if err != nil {
		Errors.Printf("%v \n \n", err)
		return
	}

	// if roll number already exixts it gives an error
	for _, j := range Students {
		if new_student.Roll_no == j.Roll_no {
			c.IndentedJSON(http.StatusConflict, gin.H{"message ": "roll number already exists"})
			Errors.Println("roll number already exists \n ")
			return
		}
	}

	// add new student to the student slice

	Students = append(Students, new_student)
	c.IndentedJSON(http.StatusCreated, new_student)
	Output.Printf("%v \n \n ", new_student)
}

// GET request to receive specific item
func get_specific_student(c *gin.Context) {
	Request.Println("Get specific student")
	roll_no := c.Param("roll_no")

	roll_no_conv, _ := strconv.ParseInt(roll_no, 10, 32)
	roll_no_int := int(roll_no_conv)

	for _, s := range Students {

		if s.Roll_no == roll_no_int {
			c.IndentedJSON(http.StatusFound, s)
			Output.Printf("%v \n \n ", s)
			return
		}
	}
	Errors.Println(" student not found \n ")
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "student not found"})
}

// get student for delete student
func get_student_by_rollno(roll_no int) (*Student, error) {

	for i, s := range Students {
		if s.Roll_no == roll_no {
			return &Students[i], nil
		}
	}
	Errors.Println(" student not found \n ")
	return nil, errors.New(" student not found")
}

// Delete request for deleting specific student
func delete_specific_student(c *gin.Context) {
	// returns query and boolean value
	roll_no, ok := c.GetQuery("id")

	Request.Println("Delete student : ", roll_no)

	roll_no_conv, _ := strconv.ParseInt(roll_no, 10, 32)
	roll_no_int := int(roll_no_conv)

	// if getquery returns false we raise an error
	if !ok {
		fmt.Println(roll_no)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "roll number not found"})
		Errors.Println("roll number not found \n ")
		return
	}

	Std, err := get_student_by_rollno(roll_no_int)

	if err != nil {
		return
	}

	set_all_null(Std)
	c.IndentedJSON(http.StatusOK, gin.H{"message": "student deleted"})
	Output.Printf("student deleted : %v\n ", roll_no)
}

// sets all values to null (used by delete function)
func set_all_null(std *Student) {
	std.Roll_no = 0
	std.Class = 0
	std.FirstName = ""
	std.LastName = ""
	std.Phone_No = 0
}

var (
	SessionTracker *log.Logger
	Errors         *log.Logger
	Request        *log.Logger
	Output         *log.Logger
)

func init() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		panic(err)
	}

	SessionTracker = log.New(file, "Session : ", log.Lshortfile|log.Ltime)
	Errors = log.New(file, "Errors : ", log.Lshortfile|log.Ltime)
	Request = log.New(file, "Request : ", log.Lshortfile|log.Ltime)
	Output = log.New(file, "Output : ", log.Lshortfile|log.Ltime)

	SessionTracker.Println("######################### \n ")
}

func main() {
	println("ho \n")
	println("hi")
	SessionTracker.Println("New session started")
	router := gin.Default()
	router.GET("/students", get_student_data)
	router.GET("/students/:roll_no", get_specific_student)
	router.POST("/add", add_sudent_data)
	router.PATCH("/delete", delete_specific_student)
	router.Run("localhost:8080")

}
