package handlers

import (
	"fmt"
	"strings"
	"temp-kbtg-backend/database"
	"temp-kbtg-backend/models"

	"github.com/gofiber/fiber/v2"
)

// GetAllUsers godoc
// @Summary Get all users
// @Description ดึงรายการผู้ใช้ทั้งหมด
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "success, data, total"
// @Failure 500 {object} map[string]interface{} "error response"
// @Router /users [get]
func GetAllUsers(c *fiber.Ctx) error {
	rows, err := database.DB.Query(`
		SELECT id, membership_id, first_name, last_name, phone_number, email, 
		       membership_level, points, joined_date, created_at, updated_at 
		FROM users
		ORDER BY id DESC
	`)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to fetch users",
			"error":   err.Error(),
		})
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.MembershipID,
			&user.FirstName,
			&user.LastName,
			&user.PhoneNumber,
			&user.Email,
			&user.MembershipLevel,
			&user.Points,
			&user.JoinedDate,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Failed to scan user",
				"error":   err.Error(),
			})
		}
		users = append(users, user)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    users,
		"total":   len(users),
	})
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description ดึงข้อมูลผู้ใช้รายบุคคล
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{} "success, data"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /users/{id} [get]
func GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var user models.User
	err := database.DB.QueryRow(`
		SELECT id, membership_id, first_name, last_name, phone_number, email, 
		       membership_level, points, joined_date, created_at, updated_at 
		FROM users WHERE id = ?
	`, id).Scan(
		&user.ID,
		&user.MembershipID,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
		&user.Email,
		&user.MembershipLevel,
		&user.Points,
		&user.JoinedDate,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    user,
	})
}

// CreateUser godoc
// @Summary Create a new user
// @Description สร้างผู้ใช้ใหม่
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "User data"
// @Success 201 {object} map[string]interface{} "success, message, data"
// @Failure 400 {object} map[string]interface{} "Validation error"
// @Failure 409 {object} map[string]interface{} "Email already exists"
// @Router /users [post]
func CreateUser(c *fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Validate required fields
	if req.FirstName == "" || req.LastName == "" || req.PhoneNumber == "" || req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "First name, last name, phone number, and email are required",
		})
	}

	// Set default membership level if not provided
	if req.MembershipLevel == "" {
		req.MembershipLevel = "Bronze"
	}

	// Generate membership ID
	membershipID, err := generateMembershipID()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to generate membership ID",
			"error":   err.Error(),
		})
	}

	// Insert user
	result, err := database.DB.Exec(`
		INSERT INTO users (membership_id, first_name, last_name, phone_number, email, membership_level, points, joined_date)
		VALUES (?, ?, ?, ?, ?, ?, 0, CURRENT_TIMESTAMP)
	`, membershipID, req.FirstName, req.LastName, req.PhoneNumber, req.Email, req.MembershipLevel)

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"success": false,
				"message": "Email or membership ID already exists",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create user",
			"error":   err.Error(),
		})
	}

	userID, _ := result.LastInsertId()

	// Fetch the created user
	var user models.User
	database.DB.QueryRow(`
		SELECT id, membership_id, first_name, last_name, phone_number, email, 
		       membership_level, points, joined_date, created_at, updated_at 
		FROM users WHERE id = ?
	`, userID).Scan(
		&user.ID,
		&user.MembershipID,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
		&user.Email,
		&user.MembershipLevel,
		&user.Points,
		&user.JoinedDate,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "User created successfully",
		"data":    user,
	})
}

// UpdateUser godoc
// @Summary Update user
// @Description แก้ไขข้อมูลผู้ใช้
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body models.UpdateUserRequest true "User update data"
// @Success 200 {object} map[string]interface{} "success, message, data"
// @Failure 400 {object} map[string]interface{} "Validation error"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /users/{id} [put]
func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")

	var req models.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Check if user exists
	var exists int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", id).Scan(&exists)
	if err != nil || exists == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "User not found",
		})
	}

	// Build update query dynamically
	updates := []string{}
	args := []interface{}{}

	if req.FirstName != "" {
		updates = append(updates, "first_name = ?")
		args = append(args, req.FirstName)
	}
	if req.LastName != "" {
		updates = append(updates, "last_name = ?")
		args = append(args, req.LastName)
	}
	if req.PhoneNumber != "" {
		updates = append(updates, "phone_number = ?")
		args = append(args, req.PhoneNumber)
	}
	if req.Email != "" {
		updates = append(updates, "email = ?")
		args = append(args, req.Email)
	}
	if req.MembershipLevel != "" {
		updates = append(updates, "membership_level = ?")
		args = append(args, req.MembershipLevel)
	}
	if req.Points != nil {
		updates = append(updates, "points = ?")
		args = append(args, *req.Points)
	}

	if len(updates) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "No fields to update",
		})
	}

	updates = append(updates, "updated_at = CURRENT_TIMESTAMP")
	args = append(args, id)

	query := fmt.Sprintf("UPDATE users SET %s WHERE id = ?", strings.Join(updates, ", "))
	_, err = database.DB.Exec(query, args...)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to update user",
			"error":   err.Error(),
		})
	}

	// Fetch updated user
	var user models.User
	database.DB.QueryRow(`
		SELECT id, membership_id, first_name, last_name, phone_number, email, 
		       membership_level, points, joined_date, created_at, updated_at 
		FROM users WHERE id = ?
	`, id).Scan(
		&user.ID,
		&user.MembershipID,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
		&user.Email,
		&user.MembershipLevel,
		&user.Points,
		&user.JoinedDate,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "User updated successfully",
		"data":    user,
	})
}

// DeleteUser godoc
// @Summary Delete user
// @Description ลบผู้ใช้
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{} "success, message"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /users/{id} [delete]
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	// Check if user exists
	var exists int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", id).Scan(&exists)
	if err != nil || exists == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "User not found",
		})
	}

	// Delete user
	_, err = database.DB.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to delete user",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "User deleted successfully",
	})
}

// Helper function to generate membership ID
func generateMembershipID() (string, error) {
	var lastID int
	err := database.DB.QueryRow("SELECT COALESCE(MAX(id), 0) FROM users").Scan(&lastID)
	if err != nil {
		return "", err
	}

	nextID := lastID + 1
	membershipID := fmt.Sprintf("LBK%06d", nextID)
	return membershipID, nil
}
