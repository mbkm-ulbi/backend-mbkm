package auth

import (
	"fmt"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/crypto/bcrypt"
)

// Handler - HTTP auth handler.
type Handler struct {
	useCase UseCase
}

// NewHandler - a factory function of auth handler.
func NewHandler(useCase UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h Handler) apiLogin(c *fiber.Ctx) error {
	var cred Credentials

	// Parse input data
	err := c.BodyParser(&cred)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Validate user credentials
	sess := c.Locals("session").(*session.Session)
	sessionToken, ok := sess.Get("session_token").(string)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Forbidden"})
	}

	err = h.useCase.apilogin(cred, sessionToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	// Generate JWT token
	token, err := h.jwtGenerate(&cred)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"token":   token,
	})
}

func (h Handler) apiRegister(c *fiber.Ctx) error {

	var input struct {
		Name             string `json:"name"`
		Email            string `json:"email"`
		Username         string `json:"username"`
		Password         string `json:"password"`
		NIM              string `json:"nim"`
		ProgramStudy     string `json:"program_study"`
		Faculty          string `json:"faculty"`
		Semester         string `json:"semester"`
		SocialMedia      string `json:"social_media"`
		EmergencyContact string `json:"emergency_contact"`
	}

	// Parse input JSON
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Validasi input
	if input.Email == "" || input.Username == "" || input.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing required fields"})
	}

	// Parsing input form-data
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid form-data"})
	}

	// Data registrasi
	name := form.Value["name"][0]
	email := form.Value["email"][0]
	username := form.Value["username"][0]
	password := form.Value["password"][0]
	nim := form.Value["nim"][0]
	programStudy := form.Value["program_study"][0]
	faculty := form.Value["faculty"][0]
	semester := form.Value["semester"][0]
	socialMedia := form.Value["social_media"][0]
	emergencyContact := form.Value["emergency_contact"][0]

	// Handle file upload
	profilePicture := ""
	if files := form.File["profile_picture"]; len(files) > 0 {
		file := files[0]
		fileName := fmt.Sprintf("%d-%s", time.Now().Unix(), file.Filename)
		filePath := fmt.Sprintf("./uploads/profile/%s", fileName)
		if err := c.SaveFile(file, filePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save profile picture"})
		}
		profilePicture = fmt.Sprintf("/uploads/profile/%s", fileName)
	}

	// Hash password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// Simpan data user
	user := &User{
		Name:             name,
		Email:            email,
		Username:         username,
		Password:         string(hashedPassword),
		NIM:              nim,
		ProgramStudy:     programStudy,
		Faculty:          faculty,
		Semester:         semester,
		SocialMedia:      socialMedia,
		EmergencyContact: emergencyContact,
		ProfilePicture:   profilePicture,
	}

	// Simpan user ke database

	// err = h.useCase.register(user)
	// Panggil useCase.register
	// err = h.useCase.register(user)
	err = h.useCase.apiRegister(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// err = h.useCase.apiRegister(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to register user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"user": fiber.Map{
			"email":           user.Email,
			"username":        user.Username,
			"nim":             user.NIM,
			"profile_picture": user.ProfilePicture,
		},
	})
}

func (h Handler) apiProfile(c *fiber.Ctx) error {
	// Ambil klaim user dari context
	user := c.Locals("user").(jwt.MapClaims)

	// Kembalikan informasi user
	return c.JSON(fiber.Map{
		"message": "User profile",
		"user":    user,
	})
}

func (h Handler) apiLogout(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Logout successful. Please clear your token on the client side.",
	})
}

func (h Handler) getRegister(c *fiber.Ctx) error {
	return c.Render("pages/auth/register", fiber.Map{})
}

func (h Handler) postRegister(c *fiber.Ctx) error {
	var cred Credentials
	if err := c.BodyParser(&cred); err != nil {
		return c.Render("pages/auth/register", fiber.Map{"Error": err.Error()})
	}
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(cred.Password), 14)
	user := &User{
		Email: cred.Email,
		Auth: &Auth{
			Type:   "email",
			Secret: string(hashPassword),
		},
	}
	err := h.useCase.register(user)
	if err != nil {
		return c.Render("pages/auth/register", fiber.Map{"Error": err.Error()})
	}
	return c.Redirect("/auth/login")
}

func (h Handler) getLogin(c *fiber.Ctx) error {
	return c.Render("pages/auth/login", fiber.Map{})
}

func (h Handler) postLogin(c *fiber.Ctx) error {
	var cred Credentials
	err := c.BodyParser(&cred)
	if err != nil {
		return c.Render("pages/auth/login", fiber.Map{"Error": err.Error()})
	}
	sess := c.Locals("session").(*session.Session)
	sessionToken, ok := sess.Get("session_token").(string)
	if !ok {
		return c.Render("pages/auth/login", fiber.Map{"Error": "Forbidden"})
	}
	err = h.useCase.login(cred, sessionToken)
	if err != nil {
		return c.Render("pages/auth/login", fiber.Map{"Error": err.Error()})
	}
	token, err := h.jwtGenerate(&cred)
	if err != nil {
		return c.Render("pages/auth/login", fiber.Map{"Error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Success login", "data": token})
}

func (h Handler) jwtGenerate(cred *Credentials) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["identity"] = cred.Email
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", fiber.ErrInternalServerError
	}
	return t, nil
}

func (h Handler) getLogout(c *fiber.Ctx) error {
	sess := c.Locals("session").(*session.Session)
	err := sess.Destroy()
	if err != nil {
		return c.Render("pages/auth/login", fiber.Map{"Error": err.Error()})
	}
	return c.Redirect("/auth/login")
}

func (h Handler) getProfile(c *fiber.Ctx) error {
	//sess, err := session.Get("session", c)
	//if err != nil {
	//	return c.SendStatus(http.StatusForbidden)
	//}
	//log.Printf("%v", sess.Values)
	//token, ok := sess.Values["session_token"].(string)
	//if !ok {
	//	b := new(bytes.Buffer)
	//	t.ViewErrorForbiddenPage(b)
	//	return c.Stream(http.StatusForbidden, echo.MIMETextHTMLCharsetUTF8, b)
	//}
	//_, err = h.useCase.Profile(token)
	//if err != nil {
	//	return c.JSON(http.StatusForbidden, err)
	//}
	//return c.SendStatus(http.StatusForbidden)

	return c.Render("pages/post/index", fiber.Map{
		"Title": "Show post list page",
	}, "layouts/main")
}
