# Hospital Backend API — Development Planning Documentation

---

## a. Project Structure

```
hospital-backend-api/
├── main.go                          # Entry point — โหลด .env, เชื่อม DB, เปิด server port 8080
├── go.mod                           # Go module dependencies
├── go.sum                           # Dependency checksums
├── .env                             # Environment variables (DB connection, JWT_SECRET)
├── .gitignore
├── .dockerignore
│
├── config/
│   └── database.go                  # เชื่อมต่อ PostgreSQL ผ่าน GORM, connection pool, auto-migrate
│
├── models/
│   ├── hospital.go                  # Hospital model — has many Staff, Patient
│   ├── staff.go                     # Staff model — belongs to Hospital
│   └── patient.go                   # Patient model — belongs to Hospital
│
├── handler/
│   ├── hospital.service.go          # Handler: CreateHospital, GetHospitals
│   ├── staff.service.go             # Handler: CreateStaff, LoginStaff
│   └── patient.service.go           # Handler: CreatePatient, SearchPatientById, SearchPatient
│
├── middleware/
│   └── auth.go                      # JWT Authentication middleware (AuthRequired)
│
├── routes/
│   ├── routes.go                    # SetupRoutes — รวม route ทั้งหมด
│   ├── hospital_routes.go           # /hospital/* routes
│   ├── staff_routes.go              # /staff/* routes
│   └── patient_routes.go            # /patient/* routes
│
├── tests/
│   ├── test_helper_test.go          # Helper functions สำหรับ test (setupDB, performRequest, etc.)
│   ├── hospital_test.go             # Unit tests สำหรับ Hospital endpoints
│   ├── staff_test.go                # Unit tests สำหรับ Staff endpoints
│   └── patient_test.go              # Unit tests สำหรับ Patient endpoints
│
├── Dockerfile                       # Docker image สำหรับ Go API
├── docker-compose.yml               # PostgreSQL + Go API + Nginx
└── nginx/
    └── nginx.conf                   # Reverse proxy → Go API port 8080
```

### Tech Stack

| Component        | Technology                   |
| ---------------- | ---------------------------- |
| Language         | Go 1.26                      |
| Web Framework    | Gin v1.12                    |
| ORM              | GORM v1.31                   |
| Database         | PostgreSQL 16 (Alpine)       |
| Authentication   | JWT (golang-jwt/v5) + bcrypt |
| Reverse Proxy    | Nginx (Alpine)               |
| Containerization | Docker + Docker Compose      |
| Testing          | Go testing + testify         |

### Infrastructure (Docker Compose)

| Service    | Image                   | Container Name    | Port            |
| ---------- | ----------------------- | ----------------- | --------------- |
| PostgreSQL | postgres:16-alpine      | hospital_postgres | 5432            |
| Go API     | (build from Dockerfile) | hospital_api      | 8080 (internal) |
| Nginx      | nginx:alpine            | hospital_nginx    | 80              |

---

## b. API Specification

### Base URL

- `http://localhost:80` (ผ่าน Nginx reverse proxy)

> หมายเหตุ: Go API (port 8080) ไม่ได้เปิด port ออกภายนอก — เข้าถึงได้เฉพาะผ่าน Nginx เท่านั้น

---

### Health Check

**GET /**

| Item     | Detail                                           |
| -------- | ------------------------------------------------ |
| Response | 200 — `{ "message": "Welcome to Hospital API" }` |

---

### Hospital Endpoints

#### POST /hospital/

สร้าง Hospital ใหม่

**Request Body:**

```json
{
  "Name": "โรงพยาบาล ABC"
}
```

**Responses:**

| Status Code               | Response Body                                                                                                    |
| ------------------------- | ---------------------------------------------------------------------------------------------------------------- |
| 200 OK                    | `{ "ID": 1, "Name": "โรงพยาบาล ABC", "CreatedAt": "...", "UpdatedAt": "...", "Staffs": null, "Patients": null }` |
| 400 Bad Request           | `{ "error": "<validation message>" }`                                                                            |
| 500 Internal Server Error | `{ "error": "<database error>" }`                                                                                |

---

#### GET /hospital/

ดึงรายชื่อ Hospital ทั้งหมด

**Responses:**

| Status Code               | Response Body                                                                      |
| ------------------------- | ---------------------------------------------------------------------------------- |
| 200 OK                    | `[ { "ID": 1, "Name": "...", "CreatedAt": "...", "UpdatedAt": "...", ... }, ... ]` |
| 500 Internal Server Error | `{ "error": "<database error>" }`                                                  |

---

### Staff Endpoints

#### POST /staff/create

สร้าง Staff ใหม่ (password ถูก hash ด้วย bcrypt ก่อนเก็บลง DB)

**Request Body:**

```json
{
  "username": "nurse1",
  "password": "pass123",
  "hospital_id": 1
}
```

**Responses:**

| Status Code               | Response Body                                                                                    |
| ------------------------- | ------------------------------------------------------------------------------------------------ |
| 201 Created               | `{ "id": 1, "username": "nurse1", "hospital": { "ID": 1, "Name": "..." }, "created_at": "..." }` |
| 400 Bad Request           | `{ "error": "<validation message>" }`                                                            |
| 409 Conflict              | `{ "error": "username นี้มีอยู่ในระบบแล้ว" }`                                                    |
| 500 Internal Server Error | `{ "error": "<error>" }`                                                                         |

**หมายเหตุ:** password จะไม่ถูกส่งกลับมาใน response (ใช้ `json:"-"` ใน Staff model)

---

#### POST /staff/login

เข้าสู่ระบบ — ได้รับ JWT token กลับมา (หมดอายุ 24 ชั่วโมง)

**Request Body:**

```json
{
  "username": "nurse1",
  "password": "pass123",
  "hospital_id": 1
}
```

**Responses:**

| Status Code      | Response Body                                                                                                             |
| ---------------- | ------------------------------------------------------------------------------------------------------------------------- |
| 200 OK           | `{ "message": "เข้าสู่ระบบสำเร็จ", "token": "<JWT>", "data": { "staff_id": 1, "username": "nurse1", "hospital_id": 1 } }` |
| 400 Bad Request  | `{ "error": "<validation message>" }`                                                                                     |
| 401 Unauthorized | `{ "error": "username หรือ password ไม่ถูกต้อง" }`                                                                        |

**JWT Token Claims:**

| Claim       | Type   | Description        |
| ----------- | ------ | ------------------ |
| staff_id    | uint   | ID ของ Staff       |
| hospital_id | uint   | ID ของ Hospital    |
| username    | string | ชื่อผู้ใช้         |
| exp         | int64  | หมดอายุ (24 ชม.)   |
| iat         | int64  | เวลาที่สร้าง token |

---

### Patient Endpoints

#### POST /patient/create

สร้างผู้ป่วยใหม่

**Request Body:**

```json
{
  "FirstNameTH": "สมชาย",
  "MiddleNameTH": "",
  "LastNameTH": "ใจดี",
  "FirstNameEN": "Somchai",
  "MiddleNameEN": "",
  "LastNameEN": "Jaidee",
  "DateOfBirth": "1990-01-15",
  "PatientHN": "HN001",
  "NationalID": "1234567890123",
  "PassportID": "",
  "PhoneNumber": "0812345678",
  "Email": "somchai@email.com",
  "Gender": "M",
  "HospitalID": 1
}
```

**Responses:**

| Status Code               | Response Body                                                                    |
| ------------------------- | -------------------------------------------------------------------------------- |
| 200 OK                    | Patient object พร้อม Hospital data (ไม่มี nested Patients list)                  |
| 400 Bad Request           | `{ "error": "<validation message>" }`                                            |
| 409 Conflict              | `{ "error": "ข้อมูลซ้ำ: ไม่สามารถบันทึกได้เนื่องจากมีข้อมูลนี้อยู่ในระบบแล้ว" }` |
| 500 Internal Server Error | `{ "error": "<error>" }`                                                         |

---

#### GET /patient/search/:id

ค้นหาผู้ป่วยด้วย NationalID หรือ PassportID (ไม่ต้อง login)

**Path Parameter:** `:id` — เลขบัตรประชาชน หรือ เลขพาสปอร์ต

**Responses:**

| Status Code   | Response Body                                                             |
| ------------- | ------------------------------------------------------------------------- |
| 200 OK        | `{ "FirstNameTH": "...", "LastNameTH": "...", "NationalID": "...", ... }` |
| 404 Not Found | `{ "error": "ไม่พบข้อมูลผู้ป่วย" }`                                       |

---

#### GET /patient/search (ต้อง Login)

ค้นหาผู้ป่วยเฉพาะ hospital เดียวกับ staff ที่ login — ต้องมี JWT token

**Headers:**

```
Authorization: Bearer <JWT token>
```

**Query Parameters (ทุก field optional):**

| Parameter     | Description           | Match Type      |
| ------------- | --------------------- | --------------- |
| national_id   | รหัสบัตรประชาชน       | exact match     |
| passport_id   | รหัสพาสปอร์ต          | exact match     |
| first_name    | ชื่อ (TH หรือ EN)     | ILIKE (partial) |
| middle_name   | ชื่อกลาง (TH หรือ EN) | ILIKE (partial) |
| last_name     | นามสกุล (TH หรือ EN)  | ILIKE (partial) |
| date_of_birth | วันเกิด               | exact match     |
| phone_number  | เบอร์โทรศัพท์         | exact match     |
| email         | อีเมล                 | ILIKE (partial) |

**Responses:**

| Status Code               | Response Body                                                                       |
| ------------------------- | ----------------------------------------------------------------------------------- |
| 200 OK                    | `{ "count": 2, "data": [ { "national_id": "...", "first_name_th": "...", ... } ] }` |
| 401 Unauthorized          | `{ "error": "กรุณาเข้าสู่ระบบก่อน (ไม่พบ token)" }`                                 |
| 500 Internal Server Error | `{ "error": "<error>" }`                                                            |

---

### Authentication Middleware (AuthRequired)

ใช้กับ endpoint ที่ต้อง login — ตรวจ JWT token จาก `Authorization: Bearer <token>`

**ผลลัพธ์ที่เป็นไปได้:**

| กรณี                       | Status | Response                                                           |
| -------------------------- | ------ | ------------------------------------------------------------------ |
| ไม่มี Authorization header | 401    | `{ "error": "กรุณาเข้าสู่ระบบก่อน (ไม่พบ token)" }`                |
| format ไม่ถูกต้อง          | 401    | `{ "error": "รูปแบบ token ไม่ถูกต้อง (ต้องเป็น Bearer <token>)" }` |
| token หมดอายุ/ไม่ valid    | 401    | `{ "error": "token ไม่ถูกต้องหรือหมดอายุ" }`                       |
| ผ่าน                       | -      | แนบ staff_id, hospital_id, username ใน context                     |

---

## c. ER Diagram

### Entity Relationship Diagram

```
┌────────────────────────────────────────┐
│              HOSPITAL                  │
├────────────────────────────────────────┤
│  ID          uint        PK           │
│  Name        string      NOT NULL, UK │
│  CreatedAt   time.Time                │
│  UpdatedAt   time.Time                │
├────────────────────────────────────────┤
│  Staffs      []Staff     1:N          │
│  Patients    []Patient   1:N          │
└──────────────┬─────────────┬──────────┘
               │             │
          1:N  │             │  1:N
               │             │
┌──────────────▼──────┐ ┌────▼─────────────────────────┐
│       STAFF         │ │          PATIENT              │
├─────────────────────┤ ├──────────────────────────────┤
│ ID         uint  PK │ │ ID            uint       PK  │
│ HospitalID uint  FK │ │ FirstNameTH   string         │
│ Username   string UK │ │ MiddleNameTH  string         │
│ Password   string    │ │ LastNameTH    string         │
│ CreatedAt  time      │ │ FirstNameEN   string         │
│ UpdatedAt  time      │ │ MiddleNameEN  string         │
├─────────────────────┤ │ LastNameEN    string         │
│ Hospital   Hospital  │ │ DateOfBirth   string         │
└─────────────────────┘ │ PatientHN     string    UK   │
                        │ NationalID    string    UK   │
                        │ PassportID    string    UK   │
                        │ PhoneNumber   string         │
                        │ Email         string         │
                        │ Gender        string         │
                        │ HospitalID    uint       FK  │
                        │ CreatedAt     time            │
                        │ UpdatedAt     time            │
                        ├──────────────────────────────┤
                        │ Hospital      Hospital        │
                        └──────────────────────────────┘
```

### Relationships

| Relation           | Type        | Foreign Key                       | Description                       |
| ------------------ | ----------- | --------------------------------- | --------------------------------- |
| Hospital → Staff   | One-to-Many | staff.hospital_id → hospital.id   | โรงพยาบาล 1 แห่ง มีได้หลาย staff  |
| Hospital → Patient | One-to-Many | patient.hospital_id → hospital.id | โรงพยาบาล 1 แห่ง มีได้หลายผู้ป่วย |

### Unique Constraints

| Table    | Column     | Description              |
| -------- | ---------- | ------------------------ |
| Hospital | Name       | ชื่อโรงพยาบาลห้ามซ้ำ     |
| Staff    | Username   | username ห้ามซ้ำ         |
| Patient  | PatientHN  | รหัสผู้ป่วย (HN) ห้ามซ้ำ |
| Patient  | NationalID | เลขบัตรประชาชนห้ามซ้ำ    |
| Patient  | PassportID | เลขพาสปอร์ตห้ามซ้ำ       |

### Field Details

#### Hospital

| Field     | Type      | Size | Constraints      |
| --------- | --------- | ---- | ---------------- |
| ID        | uint      | -    | Primary Key      |
| Name      | string    | -    | NOT NULL, UNIQUE |
| CreatedAt | time.Time | -    | Auto             |
| UpdatedAt | time.Time | -    | Auto             |

#### Staff

| Field      | Type      | Size | Constraints               |
| ---------- | --------- | ---- | ------------------------- |
| ID         | uint      | -    | Primary Key               |
| HospitalID | uint      | -    | NOT NULL, INDEX, FK       |
| Username   | string    | 100  | NOT NULL, UNIQUE          |
| Password   | string    | 255  | NOT NULL, JSON hidden (-) |
| CreatedAt  | time.Time | -    | Auto                      |
| UpdatedAt  | time.Time | -    | Auto                      |

#### Patient

| Field        | Type      | Size | Constraints         |
| ------------ | --------- | ---- | ------------------- |
| ID           | uint      | -    | Primary Key         |
| FirstNameTH  | string    | 100  | -                   |
| MiddleNameTH | string    | 100  | -                   |
| LastNameTH   | string    | 100  | -                   |
| FirstNameEN  | string    | 100  | -                   |
| MiddleNameEN | string    | 100  | -                   |
| LastNameEN   | string    | 100  | -                   |
| DateOfBirth  | string    | 100  | -                   |
| PatientHN    | string    | 50   | UNIQUE, INDEX       |
| NationalID   | string    | 20   | UNIQUE, INDEX       |
| PassportID   | string    | 20   | UNIQUE, INDEX       |
| PhoneNumber  | string    | 20   | -                   |
| Email        | string    | 100  | -                   |
| Gender       | string    | 1    | -                   |
| HospitalID   | uint      | -    | NOT NULL, INDEX, FK |
| CreatedAt    | time.Time | -    | Auto                |
| UpdatedAt    | time.Time | -    | Auto                |
