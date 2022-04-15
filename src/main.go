package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"
)

func GenerateOTP() int {
	return rand.Intn(1000000)
}

type UserData struct {
	UserName        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	FullName        string `json:"fullname"`
	OneTimeKey      int    `json:"otp"`
}

type UserLogin struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

var AllUsers = make(map[string]UserData)

func Home(w http.ResponseWriter, r *http.Request) {
	log.Print("Inside Home")
	fmt.Fprintln(w, "Welcome to homepage!")
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	log.Print("Inside SignUp")

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintln(w, "Error in SignUp: Reading request body of SignUp call failed ERROR: ", err)
		log.Print("Error in SignUp: Reading request body of SignUp call failed ERROR: ", err)
		return
	}

	recvdUserData := &UserData{}
	err = json.Unmarshal(reqBody, recvdUserData)
	if err != nil {
		fmt.Fprintln(w, "Error in SignUp: unmarshal on request body failed, ERROR: ", err)
		log.Print("Error in SignUp: unmarshal on request body failed, ERROR: ", err)
		return
	}

	if recvdUserData.UserName == "" || recvdUserData.Password == "" || recvdUserData.ConfirmPassword == "" || recvdUserData.FullName == "" {
		fmt.Fprintln(w, "Please provide non-empty values for userName, password, confirm_password, fullName")
		log.Print("Please provide non-empty values for userName, password, confirm_password, fullName")
		return
	}

	if recvdUserData.Password != recvdUserData.ConfirmPassword {
		fmt.Fprintln(w, "Password and confirm_password values not matching. Both values should be equal")
		log.Print("Password and confirm_password values not matching. Both values should be equal")
		return
	}

	if _, found := AllUsers[recvdUserData.UserName]; found {
		fmt.Fprintln(w, "User with userName: ", recvdUserData.UserName, " already exist. Please try with another username.")
		log.Print("User with userName: ", recvdUserData.UserName, " already exist. Please try with another username.")
		return
	}

	recvdUserData.OneTimeKey = -1
	AllUsers[recvdUserData.UserName] = *recvdUserData

	fmt.Fprintln(w, "User successfully created! Now try login before accessing Profile page")
	log.Print("User successfully created! Now try login before accessing Profile page")
}

func Login(w http.ResponseWriter, r *http.Request) {
	log.Print("Inside Login")
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintln(w, "Error in Login: Reading request body of Login call failed ERROR: ", err)
		log.Print("Error in Login: Reading request body of Login call failed ERROR: ", err)
		return
	}

	recvdUserLogin := &UserLogin{}
	err = json.Unmarshal(reqBody, recvdUserLogin)
	if err != nil {
		fmt.Fprintln(w, "Error in Login: unmarshal on request body failed, ERROR: ", err)
		log.Print("Error in Login: unmarshal on request body failed, ERROR: ", err)
		return
	}

	if user, found := AllUsers[recvdUserLogin.UserName]; !found {
		fmt.Fprintln(w, "User with userName: ", recvdUserLogin.UserName, " does not exist. Please try with another username.")
		log.Print("User with userName: ", recvdUserLogin.UserName, " does not exist. Please try with another username.")
		return
	} else if recvdUserLogin.Password == user.Password {
		user.OneTimeKey = GenerateOTP()
		AllUsers[recvdUserLogin.UserName] = user

		fmt.Fprintln(w, "Login successfully! Use this OTP for accessing Profile page: ", user.OneTimeKey)
		log.Print("Login successfully! Use this OTP for accessing Profile page: ", user.OneTimeKey)
		return
	} else {
		fmt.Fprintln(w, "Combination of this username and password does not exist. Please try again.")
		log.Print("Combination of this username and password does not exist. Please try again.")
		return
	}
}

func Profile(w http.ResponseWriter, r *http.Request) {
	log.Print("Inside Profile")

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintln(w, "Error in Profile: Reading request body of Profile call failed ERROR: ", err)
		log.Print("Error in Profile: Reading request body of Profile call failed ERROR: ", err)
		return
	}

	recvdUser := &UserData{}
	err = json.Unmarshal(reqBody, recvdUser)
	if err != nil {
		fmt.Fprintln(w, "Error in Profile: unmarshal on request body failed, ERROR: ", err)
		log.Print("Error in Profile: unmarshal on request body failed, ERROR: ", err)
		return
	}

	if user, found := AllUsers[recvdUser.UserName]; !found {
		fmt.Fprintln(w, "User with userName: ", recvdUser.UserName, " does not exist.")
		log.Print("User with userName: ", recvdUser.UserName, " does not exist.")
		return
	} else if recvdUser.OneTimeKey == user.OneTimeKey && recvdUser.OneTimeKey >= 0 {
		fmt.Fprintln(w, "OTP matched. You can see the Profile page of: ", recvdUser.UserName)
		log.Print("OTP matched. You can see the Profile page of: ", recvdUser.UserName)
		return
	} else {
		fmt.Fprintln(w, "OTP does not matched. Please try again with correct OTP.")
		log.Print("OTP does not matched. Please try again with correct OTP.")
		return
	}

}

func Logout(w http.ResponseWriter, r *http.Request) {
	log.Print("Inside Logout")

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintln(w, "Error in Logout: Reading request body of Logout call failed ERROR: ", err)
		log.Print("Error in Logout: Reading request body of Logout call failed ERROR: ", err)
		return
	}

	recvdUserLogin := &UserLogin{}
	err = json.Unmarshal(reqBody, recvdUserLogin)
	if err != nil {
		fmt.Fprintln(w, "Error in Logout: unmarshal on request body failed, ERROR: ", err)
		log.Print("Error in Logout: unmarshal on request body failed, ERROR: ", err)
		return
	}

	if user, found := AllUsers[recvdUserLogin.UserName]; !found {
		fmt.Fprintln(w, "User with userName: ", recvdUserLogin.UserName, " does not exist.")
		log.Print("User with userName: ", recvdUserLogin.UserName, " does not exist.")
		return
	} else {
		user.OneTimeKey = -1
		AllUsers[recvdUserLogin.UserName] = user

		fmt.Fprintln(w, "Logout successfully!")
		log.Print("Logout successfully!")
		return
	}
}

func HandleRequests(routerObj *mux.Router) {
	routerObj.HandleFunc("/", Home).Methods("GET")
	routerObj.HandleFunc("/signup", SignUp).Methods("POST")
	routerObj.HandleFunc("/login", Login).Methods("POST")
	routerObj.HandleFunc("/profile", Profile).Methods("POST")
	routerObj.HandleFunc("/logout", Logout).Methods("POST")
}

func main() {
	routerObj := mux.NewRouter()

	HandleRequests(routerObj)

	log.Print("Listening at https://localhost:3333")
	err := http.ListenAndServe(":3333", routerObj)
	if err != nil {
		log.Fatal("Failed to start server, ERROR: ", err)
	}
}
