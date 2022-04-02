package users

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/asdine/storm"
)

func Login(w http.ResponseWriter, r *http.Request) {
	email, password, ok := r.BasicAuth()
	if !ok || !CheckUsersPassword(email, password) {
		err := errors.New("неверный логин или пароль")
		log.Println(err)
		log.Println(email, password)
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}

	u, _ := store.getUserByEmail(email)
	res := struct {
		Token    string `json:"token"`
		AclGroup string `json:"aclGroup"`
	}{u.PasswordHash, u.AclGroup}

	token, _ := json.Marshal(&res)

	w.WriteHeader(http.StatusOK)
	w.Write(token)
}

func Register(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	req := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	if err := json.Unmarshal(body, &req); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Email == "" || req.Password == "" {
		err := errors.New("пустая почта или пароль")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := store.getUserByEmail(req.Email); err != storm.ErrNotFound {
		err = errors.New("уже существует")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// create new user
	u := User{}
	u.Email = req.Email
	u.AclGroup = "client"
	passHash, _ := HashPassword(req.Password)
	u.PasswordHash = passHash
	if err := u.save(); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 400)
		return
	}

	// send token and acl group to client
	res := struct {
		Token    string `json:"token"`
		AclGroup string `json:"aclGroup"`
	}{
		u.PasswordHash,
		u.AclGroup,
	}

	jsonToken, _ := json.Marshal(&res)

	w.WriteHeader(http.StatusOK)
	w.Write(jsonToken)
}
