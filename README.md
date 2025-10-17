# KBTG Backend API

Backend API สำหรับจัดการข้อมูลผู้ใช้ (User Management) และระบบโอนแต้ม (Points Transfer) ที่พัฒนาด้วย Go + Fiber Framework และใช้ SQLite เป็นฐานข้อมูล

## 🚀 Features

- ✅ RESTful API สำหรับจัดการข้อมูลผู้ใช้ (CRUD)
- ✅ **Points Transfer System** - โอนแต้มระหว่างผู้ใช้แบบอะตอมมิก
- ✅ **Idempotency Support** - ป้องกันการโอนซ้ำด้วย Idempotency Key
- ✅ **Point Ledger** - บันทึกประวัติการเปลี่ยนแปลงแต้มทุกครั้ง (Audit Trail)
- ✅ **Transaction Safety** - ใช้ Database Transaction รับประกันความสอดคล้องของข้อมูล
- ✅ SQLite Database (ไม่ต้องติดตั้ง Database Server)
- ✅ Auto-generate Membership ID (LBK######)
- ✅ Middleware: CORS, Logger
- ✅ Sample Data พร้อมใช้งาน
- ✅ OpenAPI 3.1 Compliant (ตาม transfer.yml spec)

## 📋 Prerequisites

- Go 1.17 หรือสูงกว่า
- GCC compiler (สำหรับ SQLite driver)

## 🛠️ Installation

1. Clone หรือ download project นี้

2. ติดตั้ง dependencies:

```bash
go mod download
```

3. รัน application:

```bash
go run main.go
```

Server จะรันที่ `http://localhost:3000`

## 📁 Project Structure

```
temp-kbtg-backend/
├── main.go                    # Entry point & Routes
├── models/
│   ├── user.go               # User model & request structs
│   └── transfer.go           # Transfer & PointLedger models
├── database/
│   └── db.go                 # SQLite connection & initialization
├── handlers/
│   ├── user_handler.go       # User CRUD handlers
│   └── transfer_handler.go   # Transfer handlers
├── users.db                  # SQLite database (auto-created)
├── go.mod                    # Go module dependencies
└── README.md                 # คุณกำลังอ่านอยู่ตรงนี้
```

## 📊 User Model

ข้อมูลผู้ใช้ที่เก็บในระบบ:

| Field              | Type     | Description                                  |
| ------------------ | -------- | -------------------------------------------- |
| `id`               | Integer  | ID อัตโนมัติ (Primary Key)                   |
| `membership_id`    | String   | รหัสสมาชิก (เช่น LBK001234) - Auto-generated |
| `first_name`       | String   | ชื่อ                                         |
| `last_name`        | String   | นามสกุล                                      |
| `phone_number`     | String   | เบอร์โทรศัพท์                                |
| `email`            | String   | อีเมล (Unique)                               |
| `membership_level` | String   | ระดับสมาชิก (Gold/Silver/Bronze)             |
| `points`           | Integer  | แต้มคงเหลือ                                  |
| `joined_date`      | DateTime | วันที่สมัครสมาชิก                            |
| `created_at`       | DateTime | วันที่สร้างข้อมูล                            |
| `updated_at`       | DateTime | วันที่แก้ไขข้อมูลล่าสุด                      |

## 🔌 API Endpoints

### 1. Get All Users

ดึงรายการผู้ใช้ทั้งหมด

```http
GET /users
```

**Response Example:**

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "membership_id": "LBK001234",
      "first_name": "สมชาย",
      "last_name": "ใจดี",
      "phone_number": "081-234-5678",
      "email": "somchai@example.com",
      "membership_level": "Gold",
      "points": 15420,
      "joined_date": "2023-06-15T00:00:00Z",
      "created_at": "2025-10-17T13:46:28Z",
      "updated_at": "2025-10-17T13:46:28Z"
    }
  ],
  "total": 1
}
```

### 2. Get User by ID

ดึงข้อมูลผู้ใช้รายบุคคล

```http
GET /users/{id}
```

**Example:**

```bash
curl http://localhost:3000/users/1
```

**Response Example:**

```json
{
  "success": true,
  "data": {
    "id": 1,
    "membership_id": "LBK001234",
    "first_name": "สมชาย",
    "last_name": "ใจดี",
    "phone_number": "081-234-5678",
    "email": "somchai@example.com",
    "membership_level": "Gold",
    "points": 15420,
    "joined_date": "2023-06-15T00:00:00Z",
    "created_at": "2025-10-17T13:46:28Z",
    "updated_at": "2025-10-17T13:46:28Z"
  }
}
```

### 3. Create User

สร้างผู้ใช้ใหม่

```http
POST /users
Content-Type: application/json
```

**Request Body:**

```json
{
  "first_name": "ทดสอบ",
  "last_name": "ระบบ",
  "phone_number": "099-999-9999",
  "email": "test@example.com",
  "membership_level": "Gold"
}
```

**Required Fields:**

- `first_name` (required)
- `last_name` (required)
- `phone_number` (required)
- `email` (required)
- `membership_level` (optional, default: "Bronze")

**Example:**

```bash
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "ทดสอบ",
    "last_name": "ระบบ",
    "phone_number": "099-999-9999",
    "email": "test@example.com",
    "membership_level": "Gold"
  }'
```

**Response Example:**

```json
{
  "success": true,
  "message": "User created successfully",
  "data": {
    "id": 4,
    "membership_id": "LBK000004",
    "first_name": "ทดสอบ",
    "last_name": "ระบบ",
    "phone_number": "099-999-9999",
    "email": "test@example.com",
    "membership_level": "Gold",
    "points": 0,
    "joined_date": "2025-10-17T13:50:00Z",
    "created_at": "2025-10-17T13:50:00Z",
    "updated_at": "2025-10-17T13:50:00Z"
  }
}
```

### 4. Update User

แก้ไขข้อมูลผู้ใช้

```http
PUT /users/{id}
Content-Type: application/json
```

**Request Body (ส่งเฉพาะ field ที่ต้องการแก้ไข):**

```json
{
  "first_name": "สมชาย",
  "last_name": "ใจดี",
  "phone_number": "081-234-5678",
  "email": "somchai@example.com",
  "membership_level": "Platinum",
  "points": 20000
}
```

**Example:**

```bash
curl -X PUT http://localhost:3000/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "points": 20000,
    "membership_level": "Platinum"
  }'
```

**Response Example:**

```json
{
  "success": true,
  "message": "User updated successfully",
  "data": {
    "id": 1,
    "membership_id": "LBK001234",
    "first_name": "สมชาย",
    "last_name": "ใจดี",
    "phone_number": "081-234-5678",
    "email": "somchai@example.com",
    "membership_level": "Platinum",
    "points": 20000,
    "joined_date": "2023-06-15T00:00:00Z",
    "created_at": "2025-10-17T13:46:28Z",
    "updated_at": "2025-10-17T13:52:00Z"
  }
}
```

### 5. Delete User

ลบผู้ใช้

```http
DELETE /users/{id}
```

**Example:**

```bash
curl -X DELETE http://localhost:3000/users/1
```

**Response Example:**

```json
{
  "success": true,
  "message": "User deleted successfully"
}
```

## 🎁 Sample Data

เมื่อรัน application ครั้งแรก ระบบจะสร้างข้อมูลตัวอย่าง 3 รายการให้อัตโนมัติ:

1. **สมชาย ใจดี** - Gold Member (15,420 แต้ม)
2. **สมหญิง รักดี** - Silver Member (8,500 แต้ม)
3. **สมศักดิ์ มีสุข** - Bronze Member (2,100 แต้ม)

## 🔧 Technologies

- **Go** - Programming Language
- **Fiber v2** - Web Framework (Express-like for Go)
- **SQLite** - Embedded Database
- **go-sqlite3** - SQLite Driver for Go

## 📝 Response Format

ทุก API response จะมีรูปแบบดังนี้:

**Success Response:**

```json
{
  "success": true,
  "message": "ข้อความ (ถ้ามี)",
  "data": { ... },
  "total": 0
}
```

**Error Response:**

```json
{
  "success": false,
  "message": "Error message",
  "error": "Detailed error (ถ้ามี)"
}
```

## 🚨 Error Codes

| Status Code | Description                         |
| ----------- | ----------------------------------- |
| 200         | Success                             |
| 201         | Created                             |
| 400         | Bad Request (ข้อมูลไม่ถูกต้อง)      |
| 404         | Not Found (ไม่พบข้อมูล)             |
| 409         | Conflict (ข้อมูลซ้ำ เช่น Email ซ้ำ) |
| 500         | Internal Server Error               |

## 🧪 Testing

### ทดสอบด้วย curl:

```bash
# ดูรายการผู้ใช้ทั้งหมด
curl http://localhost:3000/users

# ดูผู้ใช้รายบุคคล
curl http://localhost:3000/users/1

# สร้างผู้ใช้ใหม่
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{"first_name":"ทดสอบ","last_name":"ระบบ","phone_number":"099-999-9999","email":"test@example.com"}'

# แก้ไขข้อมูลผู้ใช้
curl -X PUT http://localhost:3000/users/1 \
  -H "Content-Type: application/json" \
  -d '{"points":25000}'

# ลบผู้ใช้
curl -X DELETE http://localhost:3000/users/1
```

### ทดสอบด้วย Postman:

1. Import collection จาก API endpoints ด้านบน
2. Set base URL: `http://localhost:3000`
3. ทดสอบแต่ละ endpoint ตามต้องการ

## 💡 Tips

- **Database Location**: ไฟล์ `users.db` จะถูกสร้างใน root directory ของโปรเจค
- **Membership ID Format**: LBK + 6 หลัก (เช่น LBK001234)
- **Default Membership Level**: Bronze (ถ้าไม่ระบุตอนสร้าง)
- **CORS**: Enable ทุก origins (`*`) สำหรับ development
- **Idempotency Key**: ระบบจะ auto-generate UUID สำหรับทุกการโอนแต้ม
- **Transaction Safety**: ทุกการโอนแต้มใช้ Database Transaction เพื่อความปลอดภัย

---

## 🔄 Transfer API (Points Transfer System)

### Database Schema

#### Transfers Table

เก็บคำสั่งโอนแต้ม พร้อม idempotency key สำหรับค้นหา

| Field             | Type     | Description                                                    |
| ----------------- | -------- | -------------------------------------------------------------- |
| `id`              | Integer  | ID ภายในระบบ (Auto-increment)                                  |
| `from_user_id`    | Integer  | ID ของผู้โอน                                                   |
| `to_user_id`      | Integer  | ID ของผู้รับ                                                   |
| `amount`          | Integer  | จำนวนแต้มที่โอน (> 0)                                          |
| `status`          | String   | สถานะ (pending/processing/completed/failed/cancelled/reversed) |
| `note`            | String   | หมายเหตุ (optional)                                            |
| `idempotency_key` | String   | Unique key สำหรับค้นหาและป้องกันการโอนซ้ำ                      |
| `created_at`      | DateTime | วันที่สร้าง                                                    |
| `updated_at`      | DateTime | วันที่อัปเดตล่าสุด                                             |
| `completed_at`    | DateTime | วันที่ทำสำเร็จ                                                 |
| `fail_reason`     | String   | เหตุผลที่ล้มเหลว (ถ้ามี)                                       |

#### Point Ledger Table

สมุดบัญชีแต้ม - บันทึกทุกการเปลี่ยนแปลงแต้ม (Append-only)

| Field           | Type     | Description                                          |
| --------------- | -------- | ---------------------------------------------------- |
| `id`            | Integer  | ID ภายในระบบ                                         |
| `user_id`       | Integer  | ID ของผู้ใช้                                         |
| `change`        | Integer  | จำนวนที่เปลี่ยนแปลง (+รับ / -โอนออก)                 |
| `balance_after` | Integer  | ยอดคงเหลือหลังทำรายการ                               |
| `event_type`    | String   | ประเภท (transfer_out/transfer_in/adjust/earn/redeem) |
| `transfer_id`   | Integer  | อ้างอิงถึง transfers.id                              |
| `reference`     | String   | ข้อมูลอ้างอิงเพิ่มเติม                               |
| `metadata`      | String   | JSON metadata                                        |
| `created_at`    | DateTime | วันที่สร้าง                                          |

### Transfer API Endpoints

#### 1. Create Transfer (POST /transfers)

สร้างคำสั่งโอนแต้ม - ระบบจะ generate Idempotency-Key ให้อัตโนมัติ

```http
POST /transfers
Content-Type: application/json
```

**Request Body:**

```json
{
  "fromUserId": 1,
  "toUserId": 2,
  "amount": 250,
  "note": "ขอบคุณสำหรับช่วยงาน"
}
```

**Response (201 Created):**

```json
{
  "transfer": {
    "idemKey": "5d1f8c7a-2b5b-4b1f-9f2a-8f50b0a8d9f3",
    "transferId": 1,
    "fromUserId": 1,
    "toUserId": 2,
    "amount": 250,
    "status": "completed",
    "note": "ขอบคุณสำหรับช่วยงาน",
    "createdAt": "2025-10-17T14:03:12Z",
    "updatedAt": "2025-10-17T14:03:12Z",
    "completedAt": "2025-10-17T14:03:12Z"
  }
}
```

**Response Headers:**

```
Idempotency-Key: 5d1f8c7a-2b5b-4b1f-9f2a-8f50b0a8d9f3
```

**Example:**

```bash
curl -X POST http://localhost:3000/transfers \
  -H "Content-Type: application/json" \
  -d '{
    "fromUserId": 1,
    "toUserId": 2,
    "amount": 250,
    "note": "ขอบคุณสำหรับช่วยงาน"
  }'
```

**Error Responses:**

- **400 Bad Request**: ข้อมูลไม่ถูกต้อง

```json
{
  "error": "VALIDATION_ERROR",
  "message": "amount must be > 0"
}
```

- **404 Not Found**: ไม่พบผู้ใช้

```json
{
  "error": "NOT_FOUND",
  "message": "Sender user not found"
}
```

- **409 Conflict**: แต้มไม่พอ

```json
{
  "error": "INSUFFICIENT_POINTS",
  "message": "Insufficient points. Available: 100, Required: 250"
}
```

- **422 Unprocessable Entity**: โอนให้ตัวเอง

```json
{
  "error": "BUSINESS_RULE_VIOLATION",
  "message": "Cannot transfer to yourself"
}
```

#### 2. Get Transfer by ID (GET /transfers/{id})

ดูสถานะคำสั่งโอน - ใช้ `idemKey` เป็น id

```http
GET /transfers/{idemKey}
```

**Example:**

```bash
curl http://localhost:3000/transfers/5d1f8c7a-2b5b-4b1f-9f2a-8f50b0a8d9f3
```

**Response (200 OK):**

```json
{
  "transfer": {
    "idemKey": "5d1f8c7a-2b5b-4b1f-9f2a-8f50b0a8d9f3",
    "transferId": 1,
    "fromUserId": 1,
    "toUserId": 2,
    "amount": 250,
    "status": "completed",
    "note": "ขอบคุณสำหรับช่วยงาน",
    "createdAt": "2025-10-17T14:03:12Z",
    "updatedAt": "2025-10-17T14:03:12Z",
    "completedAt": "2025-10-17T14:03:12Z"
  }
}
```

**Error Response (404 Not Found):**

```json
{
  "error": "NOT_FOUND",
  "message": "Transfer not found"
}
```

#### 3. Get Transfer History (GET /transfers)

ค้นหา/ดูประวัติการโอน - กรองด้วย userId (แสดงทั้งโอนออกและรับเข้า)

```http
GET /transfers?userId={userId}&page={page}&pageSize={pageSize}
```

**Query Parameters:**

- `userId` (required): ID ของผู้ใช้ที่ต้องการดูประวัติ
- `page` (optional, default=1): หน้าที่ต้องการ
- `pageSize` (optional, default=20, max=200): จำนวนรายการต่อหน้า

**Example:**

```bash
curl "http://localhost:3000/transfers?userId=1&page=1&pageSize=20"
```

**Response (200 OK):**

```json
{
  "data": [
    {
      "idemKey": "5d1f8c7a-2b5b-4b1f-9f2a-8f50b0a8d9f3",
      "transferId": 1,
      "fromUserId": 1,
      "toUserId": 2,
      "amount": 250,
      "status": "completed",
      "note": "ขอบคุณสำหรับช่วยงาน",
      "createdAt": "2025-10-17T14:03:12Z",
      "updatedAt": "2025-10-17T14:03:12Z",
      "completedAt": "2025-10-17T14:03:12Z"
    },
    {
      "idemKey": "a8b4f2e0-5562-4f1c-9b62-2a2f2f4c9b10",
      "transferId": 2,
      "fromUserId": 3,
      "toUserId": 1,
      "amount": 100,
      "status": "completed",
      "createdAt": "2025-10-17T10:00:00Z",
      "updatedAt": "2025-10-17T10:00:00Z",
      "completedAt": "2025-10-17T10:00:00Z"
    }
  ],
  "page": 1,
  "pageSize": 20,
  "total": 2
}
```

### Transfer Status Values

| Status       | Description    |
| ------------ | -------------- |
| `pending`    | รอดำเนินการ    |
| `processing` | กำลังดำเนินการ |
| `completed`  | สำเร็จ         |
| `failed`     | ล้มเหลว        |
| `cancelled`  | ยกเลิก         |
| `reversed`   | ย้อนกลับ       |

### Event Types (Point Ledger)

| Event Type     | Description            |
| -------------- | ---------------------- |
| `transfer_out` | โอนแต้มออก (ลบแต้ม)    |
| `transfer_in`  | รับโอนแต้ม (เพิ่มแต้ม) |
| `adjust`       | ปรับปรุงแต้ม           |
| `earn`         | ได้รับแต้ม             |
| `redeem`       | แลกแต้ม                |

### Business Rules

1. **ไม่สามารถโอนแต้มให้ตัวเองได้**
2. **ผู้โอนต้องมีแต้มเพียงพอ** (จำนวนแต้ม >= จำนวนที่ต้องการโอน)
3. **ทุกการโอนใช้ Database Transaction** เพื่อความปลอดภัย
4. **บันทึกทุกการเปลี่ยนแปลงใน Point Ledger** (Audit Trail)
5. **Idempotency Key เป็น UUID** ที่ unique สำหรับแต่ละรายการโอน

---

## 📄 License

MIT

## 👨‍💻 Development

สร้างโดย: KBTG Team  
วันที่: October 17, 2025

---

**Happy Coding! 🚀**
