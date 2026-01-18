package handlers

import (
	"encoding/csv"
	"fmt"
	"io"
	"mbkm-go/internal/database"
	"mbkm-go/internal/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// --- Settings (Bobot Nilai) ---

func GetBobotNilai(c *fiber.Ctx) error {
	// Logic from Laravel: get_id_program_studi() usually fetches from auth user's prodi if applicable, or general.
	// We'll simplisticly list all or filter by query `prodi_id`.
	prodiID, _ := strconv.Atoi(c.Query("prodi_id"))

	var bobot models.BobotNilai
	query := database.DB.Model(&models.BobotNilai{})

	if prodiID != 0 {
		query = query.Where("id_program_studi = ?", prodiID)
	} else {
		// Try fetch where id_program_studi is NOT NULL if no specific filter (mimicking Laravel's behavior often picking first available)
		query = query.Where("id_program_studi IS NOT NULL")
	}

	if err := query.First(&bobot).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"data": nil}) // Return null data if not set
	}

	return c.JSON(fiber.Map{"data": bobot})
}

func UpdateBobotNilai(c *fiber.Ctx) error {
	type BobotInput struct {
		IDProgramStudi       uint    `json:"id_program_studi"`
		BobotNilaiPerusahaan float64 `json:"bobot_nilai_perusahaan"`
		BobotNilaiPembimbing float64 `json:"bobot_nilai_pembimbing"`
		BobotNilaiPenguji    float64 `json:"bobot_nilai_penguji"`
	}
	input := new(BobotInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if input.IDProgramStudi == 0 {
		// Attempt to get from Auth user if missing? Or require it.
		// For now require it or default to 1 for testing.
		return c.Status(400).JSON(fiber.Map{"error": "id_program_studi is required"})
	}

	var bobot models.BobotNilai
	if err := database.DB.Where("id_program_studi = ?", input.IDProgramStudi).First(&bobot).Error; err != nil {
		// Create
		bobot = models.BobotNilai{
			IDProgramStudi:       input.IDProgramStudi,
			BobotNilaiPerusahaan: input.BobotNilaiPerusahaan,
			BobotNilaiPembimbing: input.BobotNilaiPembimbing,
			BobotNilaiPenguji:    input.BobotNilaiPenguji,
		}
		database.DB.Create(&bobot)
		return c.Status(201).JSON(fiber.Map{"status": true, "data": bobot})
	} else {
		// Update
		database.DB.Model(&bobot).Updates(map[string]interface{}{
			"bobot_nilai_perusahaan": input.BobotNilaiPerusahaan,
			"bobot_nilai_pembimbing": input.BobotNilaiPembimbing,
			"bobot_nilai_penguji":    input.BobotNilaiPenguji,
		})
		return c.Status(200).JSON(fiber.Map{"status": true, "data": bobot})
	}
}

// --- Import Handlers ---

func ImportStudents(c *fiber.Ctx) error {
	// Parse file
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "File required"})
	}

	f, err := file.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to open file"})
	}
	defer f.Close()

	reader := csv.NewReader(f)
	var users []models.User
	// var roles []models.RoleUserStub

	// Fetch 'student' role ID
	var studentRole models.Role
	database.DB.Where("title = ?", "student").First(&studentRole)
	if studentRole.ID == 0 {
		// Fallback id 2 based on seeder assumption
		studentRole.ID = 2
	}

	rowNum := 0
	inserted := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
		rowNum++
		if rowNum == 1 { // Skip header
			continue
		}

		if len(record) < 5 {
			continue
		}

		// CSV Format: Name, NIM, Birthdate (m/d/Y), ProgramStudy, Status
		name := record[0]
		nim := record[1]
		birthdateStr := record[2]
		prodi := record[3]
		status := record[4]

		birthdate, _ := time.Parse("01/02/2006", birthdateStr)
		if birthdate.IsZero() {
			// Try other format?
			birthdate, _ = time.Parse("2006-01-02", birthdateStr)
		}

		// Password: YYYYMMDD from birthdate
		passwordRaw := birthdate.Format("20060102")
		hash, _ := bcrypt.GenerateFromPassword([]byte(passwordRaw), bcrypt.DefaultCost)

		email := fmt.Sprintf("%s@mbkm.ulbi.ac.id", nim)

		user := models.User{
			Name:         name,
			Email:        email,
			Username:     nim,
			Password:     string(hash),
			NIM:          &nim,
			ProgramStudy: &prodi,
			Status:       &status,
			Role:         "student",
			Birthdate:    &birthdate,
			Verified:     false,
			Approved:     false,
		}

		if err := database.DB.Create(&user).Error; err == nil {
			users = append(users, user)
			inserted++

			// Attach Role
			// Direct DB insert for Pivot might be faster or use Association
			// Using Association:
			database.DB.Model(&user).Association("Roles").Append(&models.Role{ID: studentRole.ID})
		}
	}

	return c.JSON(fiber.Map{
		"status":       inserted > 0,
		"total_import": inserted,
	})
}

// Temporary Type for Pivot if needed, though GORM handles it via Association
type RoleUserStub struct {
	UserID uint
	RoleID uint
}
