package users

import (
	"log"

	"github.com/asdine/storm"
)

type Store struct {
	db *storm.DB
}

var store = Store{
	initUsersDB(),
}

func initUsersDB() *storm.DB {
	db, err := storm.Open("storage/users.db")
	if err != nil {
		log.Println(err)
		return nil
	}
	return db
}

func CloseDB() {
	store.db.Close()
}

func (u *User) save() error {
	if err := store.db.Save(u); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (u *User) update() error {
	du, err := store.getUserByEmail(u.Email)
	if err != nil {
		log.Println(err)
		return err
	}
	u.ID = du.ID

	if err := store.db.Update(u); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (u *User) updatePassword(newPassword string) (newHash string, err error) {

	hash, _ := HashPassword(newPassword)

	u.PasswordHash = hash
	if err := u.update(); err != nil {
		log.Println(err)
		return "", err
	}
	return u.PasswordHash, nil
}

func (s *Store) getUserByEmail(email string) (u User, err error) {
	if err := s.db.One("Email", email, &u); err != nil {
		log.Println(err)
		return u, err
	}
	log.Println(u)
	return u, nil
}

func (s *Store) getAllUsers() (users []User, err error) {

	if err := s.db.All(&users); err != nil {
		log.Println(err)
		return users, err
	}
	log.Println(users)
	return users, nil
}

func GetUserByEmail(email string) (u User, err error) {
	return store.getUserByEmail(email)
}
