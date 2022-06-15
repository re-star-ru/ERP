package users

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/asdine/storm"
)

type User struct {
	ID           int    `storm:"id,increment"`
	Email        string `json:"email" storm:"unique"`
	Name         string `json:"name"`
	PasswordHash string `json:"passwordHash,omitempty"`
	AclGroup     string `json:"aclGroup,omitempty"`
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value("email").(string)

	u, err := store.getUserByEmail(email)
	if err != nil && err == storm.ErrNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println(r.Context().Value("aclGroup"))

	res, _ := json.Marshal(&u)

	w.Write(res)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := store.getAllUsers()
	if err != nil && err == storm.ErrNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, _ := json.Marshal(&users)

	w.Write(res)
}

func Update(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)
	u := User{}
	if err := json.Unmarshal(body, &u); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	uu := User{}
	uu.Email = r.Context().Value("email").(string)
	uu.AclGroup = u.AclGroup
	uu.Name = u.Name
	if err := uu.update(); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)

	req := struct {
		NewPassword string `json:"newPassword"`
	}{}
	json.Unmarshal(body, &req)

	email := r.Context().Value("email").(string)
	u, err := store.getUserByEmail(email)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	newHash, err := u.updatePassword(req.NewPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(newHash))
}
