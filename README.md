# KBTG Backend API

Backend API สำหรับจัดการข้อมูลผู้ใช้ (User Management) ที่พัฒนาด้วย Go + Fiber Framework และใช้ SQLite เป็นฐานข้อมูล

## 🚀 Features

- ✅ RESTful API สำหรับจัดการข้อมูลผู้ใช้ (CRUD)
- ✅ SQLite Database (ไม่ต้องติดตั้ง Database Server)
- ✅ Auto-generate Membership ID (LBK######)
- ✅ Middleware: CORS, Logger
- ✅ Sample Data พร้อมใช้งาน

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
│   └── user.go               # User model & request structs
├── database/
│   └── db.go                 # SQLite connection & initialization
├── handlers/
│   └── user_handler.go       # CRUD handlers
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

## 📄 License

MIT

## 👨‍💻 Development

สร้างโดย: KBTG Team  
วันที่: October 17, 2025

---

**Happy Coding! 🚀**
