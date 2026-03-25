# Transaction-Processing-Engine# 💳 Transaction Processing Engine (Go)

This project is a **simple payment authorization engine** built using Go — kind of like a mini version of what happens behind the scenes when you swipe your card.

Instead of relying on frameworks or external libraries, everything here is built using Go’s **standard library only**. The goal? Keep it clean, fast, and easy to understand.

---

## 🚀 What This Project Does

* Validates card details (number + PIN)
* Handles transactions like:

  * Withdrawals
  * Top-ups
* Checks balance before approving transactions
* Stores transaction history
* Simulates real-world payment responses (success, failure, etc.)

Everything runs **in-memory**, so no database setup needed.

---

## 📁 Project Structure (Simple Breakdown)

```
transaction-engine/
├── main.go         → Starts the server & routes
├── models/         → Data structures (Card, Transaction)
├── store/          → In-memory storage + seed data
└── handlers/       → API logic (core functionality)
```

No unnecessary complexity. Each folder has a clear responsibility.

---

## ⚙️ How to Run

### Requirements

* Go 1.21+

### Steps

```bash
cd transaction-engine
go run main.go
```

Server will start at:

```
http://localhost:8080
```

---

## 🧪 Test Cards (Already Loaded)

You don’t need to create anything — just use these:

| Card Number      | PIN  | Balance | Status  |
| ---------------- | ---- | ------- | ------- |
| 4123456789012345 | 1234 | 1000    | ACTIVE  |
| 4123456789012346 | 5678 | 500     | ACTIVE  |
| 4123456789012347 | 9999 | 250     | BLOCKED |

---

## 🔌 API Endpoints

### 1. Make a Transaction

**POST** `/api/transaction`

Example request:

```json
{
  "cardNumber": "4123456789012345",
  "pin": "1234",
  "type": "withdraw",
  "amount": 200
}
```

### Possible Responses

| Case               | Code | Result  |
| ------------------ | ---- | ------- |
| Success            | 00   | SUCCESS |
| Invalid card       | 05   | FAILED  |
| Wrong PIN          | 06   | FAILED  |
| Not enough balance | 99   | FAILED  |

---

### 2. Check Balance

**GET**

```
/api/card/balance/{cardNumber}
```

---

### 3. Transaction History

**GET**

```
/api/card/transactions/{cardNumber}
```

---

## 🧪 Quick Testing (cURL)

### Withdraw

```bash
curl -X POST http://localhost:8080/api/transaction \
-H "Content-Type: application/json" \
-d '{"cardNumber":"4123456789012345","pin":"1234","type":"withdraw","amount":200}'
```

### Top-up

```bash
curl -X POST http://localhost:8080/api/transaction \
-H "Content-Type: application/json" \
-d '{"cardNumber":"4123456789012345","pin":"1234","type":"topup","amount":500}'
```

---

## 🔐 Security Notes

* PINs are **hashed using SHA-256**
* No plaintext storage (basic but important)
* Thread-safe using `sync.RWMutex`

---

## ⚠️ Limitations (Be Honest — This Matters)

* No database (everything resets on restart)
* No authentication beyond PIN
* No logging/monitoring
* Not production-ready — this is a learning project

---

## 🎯 Why This Project Matters

This isn’t just a CRUD app.

It shows:

* Understanding of **backend flow (request → validation → response)**
* Handling of **real-world edge cases**
* Clean separation of concerns
* Ability to build without frameworks

---

## 📌 Next Improvements (If You’re Serious)

If you want this to actually impress recruiters, do this:

* Add PostgreSQL / MySQL
* Add JWT authentication
* Add transaction logs
* Dockerize it
* Write unit tests

Right now, it's a **good base** — not a standout project.

---

## 📬 Postman

Import your collection and start testing instantly.

---

## Final Thought

This project proves you understand backend fundamentals.

But don’t stop here — this is **step 1**, not the finish line.
