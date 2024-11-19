package auth

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/storage/redis"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository - auth store APIs.
type Repository interface {
	createUser(user *User) (*User, error)
	createApiUser(user *User) (*User, error)
	createAuth(auth *Auth) (*Auth, error)
	findUserByEmail(email string) (*User, error)
	apifindUserByEmail(email string) (*User, error)
	storeSessionUser(token string, user *User) error
	getSessionUser(token string) (User, error)
}

// type repo struct {
// 	db *sqlx.DB
// 	rs *redis.Storage
// }

type repo struct {
	db *gorm.DB
	rs *redis.Storage
}

// NewRepository is a factory function of auth store.
// func NewRepository(db *sqlx.DB, rs *redis.Storage) Repository {
// 	return &repo{
// 		db: db,
// 		rs: rs,
// 	}
// }

func NewRepository(db *gorm.DB, rs *redis.Storage) Repository {
	return &repo{
		db: db,
		rs: rs,
	}
}

func (r repo) createUser(user *User) (*User, error) {
	query := `
		INSERT INTO users (email) VALUES ($1) RETURNING id
	`

	// Gunakan QueryRow untuk mengambil ID terakhir yang dimasukkan
	var lastInsertedID int64
	err := r.db.Raw(query, user.Email).Scan(&lastInsertedID).Error
	if err != nil {
		return nil, err
	}

	// Set ID pada user
	user.ID = lastInsertedID
	return user, nil
}

func (r repo) createApiUser(user *User) (*User, error) {
	query := `
		INSERT INTO users (name, email, username, password, role, nim, program_study, faculty, semester, social_media, emergency_contact)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
	`

	// Gunakan QueryRow untuk mengambil ID terakhir yang dimasukkan
	var lastInsertedID int64
	err := r.db.Raw(query, user.Name, user.Email, user.Username, user.Password, user.Role, user.NIM, user.ProgramStudy, user.Faculty, user.Semester, user.SocialMedia, user.EmergencyContact).Scan(&lastInsertedID).Error
	if err != nil {
		return nil, err
	}

	// Set ID pada user
	user.ID = lastInsertedID
	return user, nil
}

func (r repo) createAuth(auth *Auth) (*Auth, error) {
	// Validasi input auth
	if auth == nil || auth.UserID == 0 || auth.Type == "" || auth.Secret == "" {
		return nil, fmt.Errorf("invalid auth data")
	}

	query := `
		INSERT INTO auths (user_id, type, secret) VALUES ($1, $2, $3)
	`

	// Gunakan Exec untuk menyimpan data autentikasi
	err := r.db.Exec(query, auth.UserID, auth.Type, auth.Secret).Error
	if err != nil {
		return nil, err
	}

	return auth, nil
}

func (r repo) findUserByEmail(email string) (*User, error) {
	var queryUser struct {
		ID     int64  `db:"id"`
		UUID   string `db:"uuid"`
		Email  string `db:"email"`
		Secret string `db:"secret"`
	}
	query := `
		SELECT u.id, u.uuid, u.email, a.secret FROM users u JOIN auths a on u.id = a.user_id WHERE email = $1
	`
	// err := r.db.QueryRowx(query, email).StructScan(&queryUser)
	err := r.db.Raw(query, email).Scan(&queryUser).Error
	uid, _ := uuid.Parse(queryUser.UUID)
	u := &User{
		model: model{
			ID:   queryUser.ID,
			UUID: uid,
		},
		Email: queryUser.Email,
		Auth: &Auth{
			UserID: queryUser.ID,
			Type:   "email",
			Secret: queryUser.Secret,
		},
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r repo) apifindUserByEmail(email string) (*User, error) {
	// Struktur untuk memetakan hasil query
	var queryUser struct {
		ID       int64  `db:"id"`
		UUID     string `db:"uuid"`
		Email    string `db:"email"`
		Password string `db:"password"`
	}

	// Perbaiki query SQL
	query := `
		SELECT id, uuid, email, password 
		FROM users 
		WHERE email = $1
	`

	// Eksekusi query dan scan hasilnya
	err := r.db.Raw(query, email).Scan(&queryUser).Error
	if err != nil {
		return nil, err // Kembalikan error jika query gagal
	}

	// Konversi UUID dari string ke tipe UUID
	uid, err := uuid.Parse(queryUser.UUID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse UUID: %w", err)
	}

	// Buat objek User
	u := &User{
		model: model{
			ID:   queryUser.ID,
			UUID: uid,
		},
		Email: queryUser.Email,
		Auth: &Auth{
			UserID: queryUser.ID,
			Type:   "email",
			Secret: queryUser.Password, // Password diambil dari query
		},
	}

	return u, nil
}

func (r repo) storeSessionUser(token string, user *User) error {
	b, err := json.Marshal(user)
	if err != nil {
		return err
	}
	err = r.rs.Set("user--"+token, b, 0)
	if err != nil {
		return err
	}
	return nil
}

func (r repo) getSessionUser(token string) (User, error) {
	b, err := r.rs.Get("user--" + token)
	if err != nil {
		return User{}, err
	}
	var user User
	if err := json.Unmarshal(b, &user); err != nil {
		return User{}, err
	}
	return user, nil
}
