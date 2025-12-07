# Flight Tracker Backend (Go + Gin)

A backend service built with **Go (Gin)** providing:

- Secure authentication with JWT + HTTP-only partitioned cookies (CHIPS)
- Live flight data using ADS-B
- Airport timetable & airline information
- User favorites system
- Flight booking system (seat, passenger, date)
- PostgreSQL database integration
- Production-grade CORS configuration

This backend powers the **Flight Tracker Web App**.

---

## Features

### Authentication
- Register users  
- Login (sets auth cookie)  
- Retrieve current authenticated user  
- Logout (destroys cookie)  

### Flights & Airport Info
- Get live flight data  
- Fetch plane image by hex  
- Airport timetable (arrivals/departures)  
- Airline lookup by code  
- Search flights by route/date  

### Favorites
- Add airport to favorites  
- Get favorites  
- Remove favorites  

### Bookings
- Create booking  
- Retrieve user bookings  
- Delete booking  

---

# Technology Stack

| Component | Technology |
|----------|------------|
| Language | Go |
| Framework | Gin |
| Database | PostgreSQL (GORM) |
| Authentication | JWT + Secure Partitioned Cookies (CHIPS) |
| External APIs | ADS-B Exchange, AeroDataBox |
| Deployment | Render |
| Password Hashing | bcrypt |
| Config Handling | .env loader |

---

# Local Development Setup

### 1. Clone the Repository
```bash
git clone https://github.com/your-user/flighttracker-backend.git
cd Backend
```

### 2. Install Dependencies
```bash
go mod tidy
```

### 3. Create `.env` file
```
PORT=8080
JWT_SECRET=your_secret_here

DB_HOST=localhost
DB_USER=postgres
DB_PASS=password
DB_NAME=flighttracker
DB_PORT=5432

AERODATABOX_API_KEY=your_key
ADS_B_USERNAME=your_username
ADS_B_PASSWORD=your_password
```

### 4. Start the Server
```bash
go run cmd/server/main.go
```

Server will run at:
```
http://localhost:8080
```

---

# Authentication & Cookies

This backend uses **CHIPS (Partitioned Cookies)** to support secure login across different domains (e.g., Render deployments).

### Cookie Attributes

| Attribute | Value |
|----------|--------|
| HttpOnly | true |
| Secure | true |
| SameSite | None |
| Partitioned | true |
| Path | / |
| MaxAge | 86400 seconds (24h) |

### Why partitioned cookies?

Chrome blocks **third-party cookies** by default.  
Partitioned cookies allow cross-domain authentication while maintaining security.

### Frontend Requirement

The frontend must include credentials in every request:

```js
fetch("http://localhost:8080/auth/login", {
  method: "POST",
  credentials: "include",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify({
    email: "example@example.com",
    password: "123456"
  })
});
```

---

# CORS Configuration

Only specific frontend origins can call the backend.

Example:
```go
allowedOrigins := map[string]bool{
    "https://flighttracker-6jvd.onrender.com": true,
}
```

CORS includes:
- `Access-Control-Allow-Origin`
- `Access-Control-Allow-Credentials: true`
- `Access-Control-Allow-Headers`
- `Access-Control-Expose-Headers: Set-Cookie`
- Full OPTIONS preflight handling

---

# Project Structure

```
Backend/
├── cmd/server/main.go          # Entry point
├── config/config.go            # Environment loader
├── controllers/                # Route logic
│   ├── auth_controller.go
│   ├── booking_controller.go
│   ├── flights.go
│   ├── favorite_controller.go
│   └── airport_timetable.go
├── database/db.go              # DB connection
├── middleware/
│   ├── auth_middleware.go      # JWT & cookie validation
│   └── cors.go                 # CORS logic
├── models/                     # Database models
│   ├── user.go
│   ├── booking.go
│   └── favorite.go
├── routes/routes.go            # Router setup
└── utils/
    ├── hash.go                 # bcrypt helpers
    └── jwt.go                  # JWT utilities
```

---

# API Overview

## Authentication Endpoints
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/auth/register` | Register new user |
| POST | `/auth/login` | Log in user & set cookie |
| GET | `/auth/me` | Get authenticated user |
| POST | `/auth/logout` | Clear auth cookie |

## Flights Endpoints
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/flights` | Live flights |
| GET | `/plane-image/:hex` | Plane image |

## Airport & Airline Endpoints
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/airport-info` | Arrivals/departures |
| GET | `/airline/name` | Get airline name |
| GET | `/search-flights` | Search flights |

## Favorites Endpoints
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/favorites/add` | Add favorite |
| GET | `/favorites/my-favorites` | List favorites |
| DELETE | `/favorites/remove/:iata` | Remove favorite |

## Booking Endpoints
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/booking/create` | Create booking |
| GET | `/booking/my-bookings` | User bookings |
| DELETE | `/booking/delete/:id` | Delete booking |

---

# Deployment (Render)

### Build command:
```bash
go build -o server cmd/server/main.go
```

### Start command:
```bash
./server
```

### Required Render settings:
- Root directory: `Backend`
- Environment variables configured
- PostgreSQL database connected
- HTTPS enabled (cookies require Secure=true)

---

# Important Notes

- Cookies require **HTTPS** when deployed.
- Every frontend request must use `credentials: "include"`.
- Unauthorized users cannot access favorites or bookings.
- CHIPS cookies must have `Partitioned=true`, or login will fail in Chromium browsers.

---