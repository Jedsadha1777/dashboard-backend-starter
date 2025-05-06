# Go Dashboard Backend Starter

โปรเจกต์เริ่มต้นสำหรับการพัฒนา Backend สำหรับ Admin Dashboard ที่มีความปลอดภัย, ขยายขนาดได้ และมีโครงสร้างเป็นโมดูล โดยใช้ Go

## คุณสมบัติ
- ระบบ Authentication แบบ JWT (github.com/golang-jwt/jwt/v5)
- โครงสร้างแบบ Clean Architecture (Controller/Service/Model)
- PostgreSQL + GORM (gorm.io/gorm)
- Middleware สำหรับการรับรองตัวตนและจัดการ Context
- การตั้งค่าผ่านไฟล์ `.env` ด้วย github.com/joho/godotenv
- Seeder สำหรับสร้าง Admin เริ่มต้น
- การจัดการด้วย Go Modules
- RESTful API ที่ได้มาตรฐาน
- การตรวจสอบและทำความสะอาดข้อมูลนำเข้า (Input Validation)
- Rate Limiting เพื่อป้องกันการโจมตี
- การจัดการ Transaction ในฐานข้อมูล
- Structured Logging
- ระบบ Graceful Shutdown

## โครงสร้างโฟลเดอร์
```
.
├── config/           # โหลด env และการตั้งค่าต่างๆ
├── controllers/      # จัดการตรรกะการรับคำขอ HTTP
├── db/               # การเชื่อมต่อฐานข้อมูล, Seeder, Transaction
├── middleware/       # Auth, CORS, Rate Limiting middlewares
├── models/           # GORM models และ Input structs
├── routes/           # การตั้งค่า Router
├── services/         # ตรรกะทางธุรกิจ
├── utils/            # JWT, Validation, Logging
├── main.go           # จุดเริ่มต้นโปรแกรม
├── .env              # (ไม่ควร commit) การตั้งค่าที่เป็นความลับ
├── go.mod / sum      # Go Modules
```

## ความต้องการของระบบ
- Go 1.19+
- PostgreSQL
- Git

## เริ่มต้นอย่างรวดเร็ว

### 1. Clone โปรเจกต์นี้
```bash
git clone https://github.com/jed1777/dashboard-starter.git
cd dashboard-starter
```

### 2. สร้างไฟล์ `.env`
และแก้ไขค่าต่างๆ ตามความเหมาะสม:
```
# Database Configuration
DB_USER=postgres
DB_PASSWORD=secure_password_here
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
```

### 3. ติดตั้ง Dependencies
```bash
go mod tidy
```

### 4. รันเซิร์ฟเวอร์
ทางเลือกที่ 1: รันแบบปกติ
```bash
go run main.go
```
เซิร์ฟเวอร์จะเริ่มทำงานที่ `http://localhost:8080` (หรือพอร์ตที่กำหนดใน .env)

ทางเลือกที่ 2: รันพร้อม Hot Reload ด้วย air

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
| POST   | /api/v1/auth/login | เข้าสู่ระบบ (รับ token) |
| POST   | /api/v1/auth/logout | ออกจากระบบ (invalidate token) |
| GET    | /api/v1/auth/profile | ดึงข้อมูลโปรไฟล์ admin |

### การจัดการผู้ใช้งาน
| Method | Endpoint | คำอธิบาย |
|--------|----------|---------|
| GET    | /api/v1/users | ดึงรายการผู้ใช้ (พร้อม pagination) |
| GET    | /api/v1/users/:id | ดึงข้อมูลผู้ใช้รายบุคคล |
| POST   | /api/v1/users | สร้างผู้ใช้ใหม่ |
| PUT    | /api/v1/users/:id | อัปเดตข้อมูลผู้ใช้ |
| DELETE | /api/v1/users/:id | ลบผู้ใช้ |

### Admin Dashboard
| Method | Endpoint | คำอธิบาย |
|--------|----------|---------|
| GET    | /api/v1/admin/dashboard | หน้า Dashboard สำหรับ admin |

## ตัวอย่างการใช้งาน API

### การเข้าสู่ระบบ
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@example.com", "password": "Admin@123!"}'
```

### การเข้าถึง Protected Route
```bash
curl http://localhost:8080/api/v1/admin/dashboard \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Seeder
Seeder จะสร้าง Admin เริ่มต้นโดยอัตโนมัติเมื่อรัน `main.go` ถ้ายังไม่มี Admin ในฐานข้อมูล

บัญชี Admin เริ่มต้น:
- Email: admin@example.com
- Password: Admin@123!
- **สำคัญ**: ควรเปลี่ยนรหัสผ่านทันทีหลังจาก login ครั้งแรก

## การพัฒนาต่อ
- เพิ่ม Unit Test และ Integration Test
- เพิ่ม API Documentation (เช่น Swagger)
- เพิ่ม Monitoring และ Metrics
- สร้าง Dockerfile สำหรับการ Deploy
- ตั้งค่า CI/CD Pipeline

## หมายเหตุ
- ไม่ควร Commit ไฟล์ `.env`
- สำหรับการใช้งานจริง: เปลี่ยน `JWT_SECRET`, ใช้ HTTPS, ใช้ Connection Pooling และ Enable TLS สำหรับ PostgreSQL
- ควรติดตั้ง Rate Limiting และ WAF เพิ่มเติมสำหรับการใช้งานในสภาพแวดล้อมจริง