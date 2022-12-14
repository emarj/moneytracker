package sqlite

import "golang.org/x/crypto/bcrypt"

func (s *SQLiteStore) Login(user string, passwordHash []byte) (bool, error) {
	//JUST A PLACEHOLDER
	err := bcrypt.CompareHashAndPassword(passwordHash, []byte("secret"))
	if err != nil {
		return false, err
	}
	if user == "joe" {
		return true, nil
	}

	return false, nil
}
