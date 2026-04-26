# Hospital Backend API

ระบบ Backend API สำหรับจัดการข้อมูลโรงพยาบาล เจ้าหน้าที่ และผู้ป่วย

## Tech Stack

- **Go 1.26** + **Gin** (Web Framework)
- **GORM** (ORM) + **PostgreSQL** (Database)
- **JWT** + **bcrypt** (Authentication)
- **Nginx** (Reverse Proxy)
- **Docker Compose** (Containerization)

## Project Structure

```
├── main.go                  # Entry point
├── config/
│   └── database.go          # Database connection + auto-migrate
├── models/
│   ├── hospital.go          # Hospital model
│   ├── staff.go             # Staff model
│   └── patient.go           # Patient model
├── handler/
│   ├── hospital.service.go  # Hospital handlers
│   ├── staff.service.go     # Staff handlers
│   └── patient.service.go   # Patient handlers
├── middleware/
│   └── auth.go              # JWT auth middleware
├── routes/
│   ├── routes.go            # Route setup
│   ├── hospital_routes.go
│   ├── staff_routes.go
│   └── patient_routes.go
├── tests/                   # Unit tests
├── Dockerfile
├── docker-compose.yml
└── nginx/
    └── nginx.conf
```

## Getting Started

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) & Docker Compose

### Run with Docker Compose

```bash
# Clone the repo
git clone https://github.com/<your-username>/hospital-backend-api.git
cd hospital-backend-api

# Start all services (PostgreSQL + Go API + Nginx)
docker compose up --build -d

# Check status
docker compose ps

# View logs
docker compose logs -f api
```

API จะพร้อมใช้งานที่ `http://localhost`

### Stop

```bash
docker compose down
```

## API Endpoints

### Health Check

| Method | Path | Description     |
| ------ | ---- | --------------- |
| GET    | `/`  | Welcome message |

### Hospital

| Method | Path         | Description                 |
| ------ | ------------ | --------------------------- |
| POST   | `/hospital/` | สร้าง Hospital ใหม่         |
| GET    | `/hospital/` | ดึงรายชื่อ Hospital ทั้งหมด |

### Staff

| Method | Path            | Description                    |
| ------ | --------------- | ------------------------------ |
| POST   | `/staff/create` | สร้าง Staff ใหม่               |
| POST   | `/staff/login`  | เข้าสู่ระบบ (ได้รับ JWT token) |

### Patient

| Method | Path                  | Auth | Description                            |
| ------ | --------------------- | ---- | -------------------------------------- |
| POST   | `/patient/create`     | NO   | สร้างผู้ป่วยใหม่                       |
| GET    | `/patient/search/:id` | NO   | ค้นหาด้วย NationalID / PassportID      |
| GET    | `/patient/search`     | YES  | ค้นหาผู้ป่วย (เฉพาะ hospital เดียวกัน) |

> YES = ต้องส่ง `Authorization: Bearer <token>` ใน header

## Usage Examples

### สร้าง Hospital

```bash
curl -X POST http://localhost/hospital/ \
  -H "Content-Type: application/json" \
  -d '{"Name": "โรงพยาบาลศิริราช"}'
```

### สร้าง Staff

```bash
curl -X POST http://localhost/staff/create \
  -H "Content-Type: application/json" \
  -d '{"username": "nurse1", "password": "pass123", "hospital_id": 1}'
```

### Login

```bash
curl -X POST http://localhost/staff/login \
  -H "Content-Type: application/json" \
  -d '{"username": "nurse1", "password": "pass123", "hospital_id": 1}'
```

### ค้นหาผู้ป่วย (ต้อง login)

```bash
curl http://localhost/patient/search?first_name=สมชาย \
  -H "Authorization: Bearer <token>"
```

## Running Tests

```bash
cd tests
go test -v ./...
```

## Environment Variables

| Variable      | Description        | Default       |
| ------------- | ------------------ | ------------- |
| `DB_HOST`     | PostgreSQL host    | `postgres`    |
| `DB_PORT`     | PostgreSQL port    | `5432`        |
| `DB_USER`     | Database user      | `postgres`    |
| `DB_PASSWORD` | Database password  | -             |
| `DB_NAME`     | Database name      | `hospital_db` |
| `JWT_SECRET`  | Secret key for JWT | -             |

## Docker Services

| Service    | Image               | Port            |
| ---------- | ------------------- | --------------- |
| PostgreSQL | postgres:16-alpine  | 5432 (internal) |
| Go API     | Custom (Dockerfile) | 8080 (internal) |
| Nginx      | nginx:alpine        | 80 (public)     |

## ER Diagram

```
Hospital (1) ──── (N) Staff
Hospital (1) ──── (N) Patient
```

- **Hospital**: ID, Name (unique)
- **Staff**: ID, Username (unique), Password (hashed), HospitalID
- **Patient**: ID, Name (TH/EN), PatientHN (unique), NationalID (unique), PassportID (unique), HospitalID
