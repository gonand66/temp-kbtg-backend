package handlers

import (
	"database/sql"
	"fmt"
	"strconv"
	"temp-kbtg-backend/database"
	"temp-kbtg-backend/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CreateTransfer godoc
// @Summary Create points transfer
// @Description สร้างคำสั่งโอนแต้ม (ระบบจะสร้าง Idempotency-Key ให้อัตโนมัติ)
// @Tags Transfers
// @Accept json
// @Produce json
// @Param transfer body models.TransferCreateRequest true "Transfer data"
// @Success 201 {object} models.TransferCreateResponse
// @Failure 400 {object} map[string]interface{} "Validation error"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 409 {object} map[string]interface{} "Insufficient points"
// @Failure 422 {object} map[string]interface{} "Cannot transfer to yourself"
// @Router /transfers [post]
func CreateTransfer(c *fiber.Ctx) error {
	var req models.TransferCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "VALIDATION_ERROR",
			"message": "Invalid request body",
		})
	}

	// Validate required fields
	if req.FromUserID < 1 || req.ToUserID < 1 || req.Amount < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "VALIDATION_ERROR",
			"message": "fromUserId, toUserId, and amount must be greater than 0",
		})
	}

	// Check if trying to transfer to self
	if req.FromUserID == req.ToUserID {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error":   "BUSINESS_RULE_VIOLATION",
			"message": "Cannot transfer to yourself",
		})
	}

	// Generate idempotency key
	idemKey := uuid.New().String()

	// Start transaction
	tx, err := database.DB.Begin()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "INTERNAL_ERROR",
			"message": "Failed to start transaction",
		})
	}
	defer tx.Rollback()

	// Check if sender exists and has enough points
	var senderPoints int
	err = tx.QueryRow("SELECT points FROM users WHERE id = ?", req.FromUserID).Scan(&senderPoints)
	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "NOT_FOUND",
			"message": "Sender user not found",
		})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "INTERNAL_ERROR",
			"message": "Failed to check sender balance",
		})
	}

	if senderPoints < req.Amount {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":   "INSUFFICIENT_POINTS",
			"message": fmt.Sprintf("Insufficient points. Available: %d, Required: %d", senderPoints, req.Amount),
		})
	}

	// Check if receiver exists
	var receiverExists int
	err = tx.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", req.ToUserID).Scan(&receiverExists)
	if err != nil || receiverExists == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "NOT_FOUND",
			"message": "Receiver user not found",
		})
	}

	// Create transfer record
	now := time.Now().UTC().Format(time.RFC3339)
	result, err := tx.Exec(`
		INSERT INTO transfers (from_user_id, to_user_id, amount, status, note, idempotency_key, created_at, updated_at, completed_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, req.FromUserID, req.ToUserID, req.Amount, models.StatusCompleted, req.Note, idemKey, now, now, now)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "INTERNAL_ERROR",
			"message": "Failed to create transfer",
		})
	}

	transferID, _ := result.LastInsertId()

	// Update sender points
	newSenderBalance := senderPoints - req.Amount
	_, err = tx.Exec("UPDATE users SET points = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", newSenderBalance, req.FromUserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "INTERNAL_ERROR",
			"message": "Failed to deduct points from sender",
		})
	}

	// Get receiver current points
	var receiverPoints int
	tx.QueryRow("SELECT points FROM users WHERE id = ?", req.ToUserID).Scan(&receiverPoints)

	// Update receiver points
	newReceiverBalance := receiverPoints + req.Amount
	_, err = tx.Exec("UPDATE users SET points = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", newReceiverBalance, req.ToUserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "INTERNAL_ERROR",
			"message": "Failed to add points to receiver",
		})
	}

	// Record in ledger - sender (debit)
	_, err = tx.Exec(`
		INSERT INTO point_ledger (user_id, change, balance_after, event_type, transfer_id, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, req.FromUserID, -req.Amount, newSenderBalance, models.EventTransferOut, transferID, now)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "INTERNAL_ERROR",
			"message": "Failed to record sender ledger",
		})
	}

	// Record in ledger - receiver (credit)
	_, err = tx.Exec(`
		INSERT INTO point_ledger (user_id, change, balance_after, event_type, transfer_id, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, req.ToUserID, req.Amount, newReceiverBalance, models.EventTransferIn, transferID, now)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "INTERNAL_ERROR",
			"message": "Failed to record receiver ledger",
		})
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "INTERNAL_ERROR",
			"message": "Failed to commit transaction",
		})
	}

	// Fetch the created transfer
	transfer := fetchTransferByIdemKey(idemKey)

	// Set Idempotency-Key header
	c.Set("Idempotency-Key", idemKey)

	return c.Status(fiber.StatusCreated).JSON(models.TransferCreateResponse{
		Transfer: transfer,
	})
}

// GetTransferByID godoc
// @Summary Get transfer by ID
// @Description ดูสถานะคำสั่งโอน (ใช้ idemKey เป็น id)
// @Tags Transfers
// @Accept json
// @Produce json
// @Param id path string true "Idempotency Key (idemKey)"
// @Success 200 {object} models.TransferGetResponse
// @Failure 404 {object} map[string]interface{} "Transfer not found"
// @Router /transfers/{id} [get]
func GetTransferByID(c *fiber.Ctx) error {
	idemKey := c.Params("id")

	if idemKey == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "VALIDATION_ERROR",
			"message": "Transfer ID is required",
		})
	}

	transfer := fetchTransferByIdemKey(idemKey)
	if transfer.IdemKey == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "NOT_FOUND",
			"message": "Transfer not found",
		})
	}

	return c.JSON(models.TransferGetResponse{
		Transfer: transfer,
	})
}

// GetTransfers godoc
// @Summary Get transfer history
// @Description ค้น/ดูประวัติการโอน (กรองด้วย userId เท่านั้น)
// @Tags Transfers
// @Accept json
// @Produce json
// @Param userId query int true "User ID"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(20)
// @Success 200 {object} models.TransferListResponse
// @Failure 400 {object} map[string]interface{} "Validation error"
// @Router /transfers [get]
func GetTransfers(c *fiber.Ctx) error {
	userIDStr := c.Query("userId")
	if userIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "VALIDATION_ERROR",
			"message": "userId query parameter is required",
		})
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "VALIDATION_ERROR",
			"message": "userId must be a positive integer",
		})
	}

	// Parse pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.Query("pageSize", "20"))
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 200 {
		pageSize = 200
	}

	offset := (page - 1) * pageSize

	// Count total records
	var total int
	err = database.DB.QueryRow(`
		SELECT COUNT(*) FROM transfers 
		WHERE from_user_id = ? OR to_user_id = ?
	`, userID, userID).Scan(&total)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "INTERNAL_ERROR",
			"message": "Failed to count transfers",
		})
	}

	// Fetch transfers
	rows, err := database.DB.Query(`
		SELECT id, from_user_id, to_user_id, amount, status, note, idempotency_key, 
		       created_at, updated_at, completed_at, fail_reason
		FROM transfers 
		WHERE from_user_id = ? OR to_user_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, userID, userID, pageSize, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "INTERNAL_ERROR",
			"message": "Failed to fetch transfers",
		})
	}
	defer rows.Close()

	transfers := []models.Transfer{}
	for rows.Next() {
		var t models.Transfer
		var id int
		var note, completedAt, failReason sql.NullString
		var createdAt, updatedAt string

		err := rows.Scan(&id, &t.FromUserID, &t.ToUserID, &t.Amount, &t.Status, &note, &t.IdemKey,
			&createdAt, &updatedAt, &completedAt, &failReason)
		if err != nil {
			continue
		}

		t.TransferID = &id
		if note.Valid {
			t.Note = &note.String
		}
		if failReason.Valid {
			t.FailReason = &failReason.String
		}

		// Parse timestamps
		if parsedCreatedAt, err := time.Parse(time.RFC3339, createdAt); err == nil {
			t.CreatedAt = parsedCreatedAt
		}
		if parsedUpdatedAt, err := time.Parse(time.RFC3339, updatedAt); err == nil {
			t.UpdatedAt = parsedUpdatedAt
		}
		if completedAt.Valid {
			if parsedCompletedAt, err := time.Parse(time.RFC3339, completedAt.String); err == nil {
				t.CompletedAt = &parsedCompletedAt
			}
		}

		transfers = append(transfers, t)
	}

	return c.JSON(models.TransferListResponse{
		Data:     transfers,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	})
}

// Helper function to fetch transfer by idempotency key
func fetchTransferByIdemKey(idemKey string) models.Transfer {
	var t models.Transfer
	var id int
	var note, completedAt, failReason sql.NullString
	var createdAt, updatedAt string

	err := database.DB.QueryRow(`
		SELECT id, from_user_id, to_user_id, amount, status, note, idempotency_key,
		       created_at, updated_at, completed_at, fail_reason
		FROM transfers WHERE idempotency_key = ?
	`, idemKey).Scan(&id, &t.FromUserID, &t.ToUserID, &t.Amount, &t.Status, &note, &t.IdemKey,
		&createdAt, &updatedAt, &completedAt, &failReason)

	if err != nil {
		return models.Transfer{}
	}

	t.TransferID = &id
	if note.Valid {
		t.Note = &note.String
	}
	if failReason.Valid {
		t.FailReason = &failReason.String
	}

	// Parse timestamps
	if parsedCreatedAt, err := time.Parse(time.RFC3339, createdAt); err == nil {
		t.CreatedAt = parsedCreatedAt
	}
	if parsedUpdatedAt, err := time.Parse(time.RFC3339, updatedAt); err == nil {
		t.UpdatedAt = parsedUpdatedAt
	}
	if completedAt.Valid {
		if parsedCompletedAt, err := time.Parse(time.RFC3339, completedAt.String); err == nil {
			t.CompletedAt = &parsedCompletedAt
		}
	}

	return t
}
