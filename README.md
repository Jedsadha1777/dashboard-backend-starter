# Go Dashboard Backend Starter

โปรเจกต์เริ่มต้นสำหรับการพัฒนา Backend สำหรับ Admin Dashboard โดยใช้ Go และ PostgreSQL

## คุณสมบัติหลัก

- **Authentication & Authorization**
  - ระบบ Authentication แบบ JWT พร้อม Token Versioning (github.com/golang-jwt/jwt/v5)
  - ระบบ Refresh Token แยกต่างหาก (อายุ 1 ปี ตรงนี้เป็นค่าตั้งต้น ถ้าเอาไปใช้บน prod ควรแก้ไข)
  - การจัดการเข้าสู่ระบบของ Admin, User และ IoT Device
  - การรีเซ็ต API key และ token invalidation

- **โครงสร้างแบบ Clean Architecture**
  - แยกส่วน Controller/Service/Repository/Model
  - การจัดโครงสร้างไฟล์ที่เป็นระเบียบและบำรุงรักษาง่าย
  - Generic Repository Pattern เพื่อลดโค้ดที่ซ้ำซ้อน

- **ฐานข้อมูลและการจัดการ**
  - PostgreSQL + GORM (gorm.io/gorm)
  - การจัดการ Transaction อย่างปลอดภัย
  - Migrations และ Seeder อัตโนมัติ
  - Soft Delete สำหรับการกู้คืนข้อมูล

- **ระบบรักษาความปลอดภัย**
  - IP-based Rate Limiting เพื่อป้องกันการโจมตี
  - API Key Generation สำหรับ IoT Device
  - การกรองข้อมูลนำเข้าและการตรวจสอบความถูกต้อง
  - การเข้ารหัสรหัสผ่านด้วย bcrypt

- **ความสะดวกในการพัฒนา**
  - การตั้งค่าผ่านไฟล์ `.env` ด้วย github.com/joho/godotenv
  - Structured Logging แบบปรับแต่งได้
  - Graceful Shutdown เพื่อป้องกันการสูญเสียข้อมูล
  - RESTful API ที่สอดคล้องกับมาตรฐาน

- **ความพร้อมสำหรับการขยาย**
  - การจัดการ Pagination สำหรับ endpoints ที่เกี่ยวข้องกับข้อมูลจำนวนมาก
  - การจัดการ Caching และ Database Connection Pool
  - โครงสร้างแบบโมดูลาร์เพื่อรองรับการขยาย
  - การเพิ่มฟีเจอร์ใหม่ได้ง่าย

## โครงสร้างโฟลเดอร์

```
.
├── cmd/
│   ├── migrate/       # เครื่องมือสำหรับการ migration
│   └── seed/          # เครื่องมือสำหรับการเพิ่มข้อมูลตั้งต้น
├── config/            # การตั้งค่าแอปพลิเคชันและการโหลด env
├── controllers/       # จัดการการรับ request และส่ง response
├── db/                # การเชื่อมต่อฐานข้อมูล, Repository, Seeder, Transaction
├── middleware/        # Auth, CORS, Rate Limiting middlewares
├── models/            # GORM models และ Input validation structs
├── routes/            # การกำหนด Router และการจัดกลุ่ม endpoints
├── services/          # ตรรกะทางธุรกิจและการดำเนินการข้อมูล
├── utils/             # JWT, Validation, Logging, Pagination
├── main.go            # จุดเริ่มต้นแอปพลิเคชัน
├── .env.example       # ตัวอย่างไฟล์การตั้งค่า environment
├── go.mod / go.sum    # Go Modules
└── README.md          # เอกสารโปรเจกต์
```

## ความต้องการของระบบ

- Go 1.20+
- PostgreSQL 13+
- Git

## เริ่มต้นอย่างรวดเร็ว

### 1. Clone โปรเจกต์

```bash
git clone https://github.com/yourusername/dashboard-starter.git
cd dashboard-starter
```

### 2. สร้างไฟล์ `.env`

สร้างไฟล์ `.env` จาก `.env.example` และแก้ไขค่าต่างๆ ตามความเหมาะสม:

```bash
cp .env.example .env
# แก้ไขไฟล์ .env ตามความเหมาะสม
```

รายละเอียดตัวแปรสภาพแวดล้อมที่สำคัญ:

```
# Database Configuration
DB_USER=postgres
DB_PASSWORD=your_secure_password_here
DB_NAME=dashboard
DB_PORT=5432
DB_HOST=localhost
DB_TIMEZONE=UTC
DB_SSLMODE=disable

# Server Configuration
SERVER_PORT=8080
SERVER_READ_TIMEOUT=10
SERVER_WRITE_TIMEOUT=10
# ตั้งค่า trusted proxies (ว่างเปล่า = ไม่เชื่อถือ proxy ใดๆ)
TRUSTED_PROXIES=

# JWT Configuration
JWT_SECRET=your_strong_random_jwt_secret_key_here
JWT_EXPIRY_MINUTES=1440

# Logging Configuration
LOG_LEVEL=info
LOG_TO_FILE=false

# Rate Limiting
RATE_LIMIT_REQUESTS_PER_MINUTE=60
RATE_LIMIT_PATHS=/api/v1/auth/login,/api/v1/auth/register,/api/v1/users
```

### 3. ติดตั้ง Dependencies

```bash
go mod tidy
```

### 4. สร้างและเริ่มต้น Database

```bash
# สร้างฐานข้อมูลใน PostgreSQL
createdb dashboard

# ทำ migration และ seed
go run cmd/seed/main.go
```

### 5. รันเซิร์ฟเวอร์

**ทางเลือกที่ 1: รันแบบปกติ**

```bash
go run main.go
```

เซิร์ฟเวอร์จะเริ่มทำงานที่ `http://localhost:8080` (หรือพอร์ตที่กำหนดใน .env)

**ทางเลือกที่ 2: รันพร้อม Hot Reload ด้วย air**

ติดตั้ง air (เพียงครั้งเดียว):

```bash
go install github.com/cosmtrek/air@latest
```

ตรวจสอบให้แน่ใจว่า $GOPATH/bin หรือ $HOME/go/bin อยู่ใน $PATH

รัน:

```bash
air
```

เซิร์ฟเวอร์จะรีโหลดโดยอัตโนมัติเมื่อไฟล์มีการเปลี่ยนแปลง

## การตั้งค่า Trusted Proxies

เพื่อแก้ไขคำเตือน "You trusted all proxies, this is NOT safe" ให้ตั้งค่า `TRUSTED_PROXIES` ในไฟล์ `.env`:

1. สำหรับ Development ในเครื่องท้องถิ่น (ไม่เชื่อถือ proxy ใดๆ):
```
TRUSTED_PROXIES=
```

2. สำหรับใช้กับ Reverse Proxy เช่น Nginx:
```
TRUSTED_PROXIES=127.0.0.1,10.0.0.1
```

3. สำหรับใช้ในเครือข่ายภายใน:
```
TRUSTED_PROXIES=192.168.0.0/16,10.0.0.0/8
```

## API Endpoints

### การจัดการ Authentication

| Method | Endpoint | คำอธิบาย |
|--------|----------|---------|
| POST   | /api/v1/auth/login | เข้าสู่ระบบ admin (รับ token) |
| POST   | /api/v1/auth/logout | ออกจากระบบ (invalidate token) |
| POST   | /api/v1/auth/refresh | รีเฟรช access token ด้วย refresh token |
| POST   | /api/v1/auth/device | ยืนยันตัวตนสำหรับอุปกรณ์ IoT |
| GET    | /api/v1/auth/profile | ดึงข้อมูลโปรไฟล์ผู้ใช้งาน |

### การจัดการผู้ใช้งาน

| Method | Endpoint | คำอธิบาย |
|--------|----------|---------|
| GET    | /api/v1/admin/users | ดึงรายการผู้ใช้ (พร้อม pagination) |
| GET    | /api/v1/admin/users/:id | ดึงข้อมูลผู้ใช้รายบุคคล |
| POST   | /api/v1/admin/users | สร้างผู้ใช้ใหม่ |
| PUT    | /api/v1/admin/users/:id | อัปเดตข้อมูลผู้ใช้ |
| DELETE | /api/v1/admin/users/:id | ลบผู้ใช้ |

### การจัดการอุปกรณ์ IoT

| Method | Endpoint | คำอธิบาย |
|--------|----------|---------|
| GET    | /api/v1/admin/devices | ดึงรายการอุปกรณ์ (พร้อม pagination) |
| GET    | /api/v1/admin/devices/:id | ดึงข้อมูลอุปกรณ์เฉพาะ |
| POST   | /api/v1/admin/devices | ลงทะเบียนอุปกรณ์ใหม่ |
| PUT    | /api/v1/admin/devices/:id | อัปเดตข้อมูลอุปกรณ์ |
| DELETE | /api/v1/admin/devices/:id | ลบอุปกรณ์ |
| POST   | /api/v1/admin/devices/:id/reset-key | รีเซ็ท API key ของอุปกรณ์ |

### การจัดการบทความ

| Method | Endpoint | คำอธิบาย |
|--------|----------|---------|
| GET    | /api/v1/admin/articles | ดึงรายการบทความ (พร้อม pagination และการค้นหา) |
| GET    | /api/v1/admin/articles/:id | ดึงข้อมูลบทความเฉพาะ |
| POST   | /api/v1/admin/articles | สร้างบทความใหม่ |
| PUT    | /api/v1/admin/articles/:id | อัปเดตบทความ |
| DELETE | /api/v1/admin/articles/:id | ลบบทความ |
| POST   | /api/v1/admin/articles/:id/publish | เผยแพร่บทความ |

### Admin Dashboard

| Method | Endpoint | คำอธิบาย |
|--------|----------|---------|
| GET    | /api/v1/admin/dashboard | ข้อมูลสรุปสำหรับ admin dashboard |

## ตัวอย่างการใช้งาน API

### การเข้าสู่ระบบ Admin

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@example.com", "password": "Admin@123!"}'
```

ตัวอย่างการตอบกลับ:

```json
{
  "success": true,
  "data": {
    "token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9...",
    "refresh_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9...",
    "expires_at": "2025-05-09T10:30:00Z",
    "user_id": 1,
    "user_type": "admin"
  }
}
```

### การสร้างบทความใหม่

```bash
curl -X POST http://localhost:8080/api/v1/admin/articles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "บทความใหม่",
    "content": "เนื้อหาของบทความใหม่",
    "slug": "new-article",
    "summary": "สรุปย่อของบทความใหม่",
    "status": "draft"
  }'
```

ตัวอย่างการตอบกลับ:

```json
{
  "success": true,
  "data": {
    "id": 1,
    "title": "บทความใหม่",
    "content": "เนื้อหาของบทความใหม่",
    "slug": "new-article",
    "summary": "สรุปย่อของบทความใหม่",
    "status": "draft",
    "published_at": null,
    "admin_id": 1,
    "admin": {
      "id": 1,
      "email": "admin@example.com"
    },
    "created_at": "2025-05-08T15:30:00Z",
    "updated_at": "2025-05-08T15:30:00Z"
  }
}
```

## การตั้งค่า Rate Limiting

เพื่อป้องกัน API จากการใช้งานมากเกินไปหรือการโจมตี แอปพลิเคชันมีระบบ rate limiting ที่สามารถปรับแต่งได้ผ่านตัวแปรสภาพแวดล้อมต่อไปนี้:

| ตัวแปร | คำอธิบาย | ค่าเริ่มต้น |
|--------|----------|---------|
| `RATE_LIMIT_REQUESTS_PER_MINUTE` | จำนวนคำขอสูงสุดที่อนุญาตต่อนาทีสำหรับ path ที่มีการจำกัด | 60 |
| `RATE_LIMIT_PATHS` | รายการ API path ที่ควรมีการจำกัดอัตรา (คั่นด้วยเครื่องหมายจุลภาค) | `/api/v1/auth/login` |

### ตัวอย่างการตั้งค่า

```
# Rate Limiting
RATE_LIMIT_REQUESTS_PER_MINUTE=60
RATE_LIMIT_PATHS=/api/v1/auth/login,/api/v1/auth/register,/api/v1/admin/articles
```

เมื่อตั้งค่านี้:
- path `/api/v1/auth/login`, `/api/v1/auth/register`, และ `/api/v1/admin/articles` จะถูกจำกัดอัตรา
- แต่ละ path จะอนุญาตให้มีการขอข้อมูลสูงสุด 60 ครั้งต่อนาที
- เมื่อเกินขีดจำกัด API จะส่งกลับสถานะ 429 Too Many Requests

## Seeder และข้อมูลตั้งต้น

เมื่อทำการรัน `main.go` หรือ `go run cmd/seed/main.go` โปรแกรมจะสร้างข้อมูลตั้งต้นโดยอัตโนมัติหากยังไม่มีข้อมูลในฐานข้อมูล:

### บัญชี Admin เริ่มต้น

- **Email**: admin@example.com
- **Password**: Admin@123!
- **ข้อควรระวัง**: ควรเปลี่ยนรหัสผ่านทันทีหลังจาก login ครั้งแรกในสภาพแวดล้อมการผลิต

### ข้อมูลตัวอย่าง

ระบบจะสร้างข้อมูลผู้ใช้ตัวอย่าง 5 รายการสำหรับการทดสอบในสภาพแวดล้อมการพัฒนา

## การทดสอบ

สามารถรันการทดสอบได้ด้วยคำสั่ง:

```bash
go test ./...
```

## การดำเนินงานในสภาพแวดล้อมการผลิต

สำหรับการใช้งานในสภาพแวดล้อมการผลิต ควรพิจารณาขั้นตอนต่อไปนี้:

1. ตั้งค่า `JWT_SECRET` ที่ซับซ้อนและไม่คาดเดา
2. เปิดใช้งาน SSL/TLS
3. เปลี่ยนรหัสผ่าน admin เริ่มต้น
4. กำหนดค่า `TRUSTED_PROXIES` อย่างเหมาะสม
5. ตั้งค่า Rate Limiting ให้เหมาะสมกับการใช้งาน
6. พิจารณาใช้ Docker สำหรับการ deploy
