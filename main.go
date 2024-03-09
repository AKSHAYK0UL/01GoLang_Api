package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Course struct {
	CourseId    string  `json:"courseid"`
	CourseName  string  `json:"coursename"`
	CoursePrice int     `json:"courseprice"`
	Author      *Author `json:"author"`
}
type Author struct {
	FullName string `json:"fullname"`
	Website  string `json:"website"`
}

// helper method
func (c *Course) IsEmpty() bool { //(c *course) means it is the part of "Course Struct"
	//return c.CourseId == "" && c.CourseName == ""
	return c.CourseName == ""

}

// Fake database
var courses []Course

func main() {
	//seeding
	courses = append(courses, Course{"101", "Flutter", 599, &Author{"Alex", "FreeCodeCamp"}}, Course{"102", "GoLang", 499, &Author{"Ali", "Udemy"}})
	fmt.Println("Building Api in GoLang")
	r := mux.NewRouter()
	r.HandleFunc("/", serveHome).Methods("GET")
	r.HandleFunc("/all", getAllCourses).Methods("GET")
	r.HandleFunc("/course/{id}", getOneCourse).Methods("GET")
	r.HandleFunc("/create", CreateOneCourse).Methods("POST")
	r.HandleFunc("/update/{id}", update).Methods("PUT")
	r.HandleFunc("/delall", deleteAll).Methods("DELETE")
	r.HandleFunc("/del/{id}", deleteOneCourse).Methods("DELETE")
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}
	log.Fatal(http.ListenAndServe(":0.0.0.0:"+port, r))
}

// On Home Page
func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Welcome To Goland backend</h1>"))
}

// Get all Courses
func getAllCourses(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get all Courses")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

// get one course by id
func getOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get One Course")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, cour := range courses {
		if cour.CourseId == params["id"] {
			json.NewEncoder(w).Encode(cour)
			return
		}
	}
	json.NewEncoder(w).Encode("No courese found with this Id ")

}

// Create one Course
func CreateOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create One Course")
	w.Header().Set("Content-Type", "application/json")
	//body is empty
	if r.Body == nil {
		json.NewEncoder(w).Encode("Please send some data")
	}
	//what if {}
	var course Course
	_ = json.NewDecoder(r.Body).Decode(&course) //get the request body and store in the course and check if its empty or not
	if course.IsEmpty() {
		json.NewEncoder(w).Encode("No data inside Json")

		return
	}
	//generate Id
	rand.Seed(time.Now().UnixNano())
	course.CourseId = strconv.Itoa(rand.Intn(100))
	courses = append(courses, course)
	json.NewEncoder(w).Encode(courses)

}

// update Course
func update(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update Course")
	w.Header().Set("Content-Type", "application/json")
	paramas := mux.Vars(r)
	for index, cour := range courses {
		if cour.CourseId == paramas["id"] {
			courses = append(courses[:index], courses[index+1:]...)

			var course Course
			_ = json.NewDecoder(r.Body).Decode(&course)
			course.CourseId = paramas["id"]
			courses = append(courses, course)
			json.NewEncoder(w).Encode(course)

		}
	}

}

// delete all courses
func deleteAll(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete all courses")
	w.Header().Set("Content-Type", "application/json")
	courses = courses[:0]
	json.NewEncoder(w).Encode("All courses deleted")
}

// delete one Course
func deleteOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete One Course")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, cour := range courses {
		if cour.CourseId == params["id"] {
			courses = append(courses[:index], courses[index+1:]...)
			json.NewEncoder(w).Encode(&courses)
			break

		}
	}
}
