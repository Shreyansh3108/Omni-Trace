# OmniTrace API 

** Live Deployment:** [https://omnitrace-api.onrender.com]

Hey there! Welcome to the repository for OmniTrace. I built this API to showcase my approach to writing clean, production-ready backend code in Go. 

When planning this project, my main goal was to avoid heavy abstractions. I wanted to build something fast, safe, and easy to maintain, which drove my decisions for the tech stack below.

##  Why I Chose This Stack

* **Go & GoFiber:** I love Go for its concurrency and simplicity. I decided to use GoFiber as the web framework because it handles routing cleanly and is incredibly fast under the hood (thanks to `fasthttp`), without feeling overly bloated.
* **PostgreSQL & pgxpool:** I used Neon to host a serverless Postgres database. To make sure the server handles multiple requests smoothly, I wired it up using `pgxpool` for concurrent connection pooling.
* **SQLC (Instead of an ORM):** This is probably my favorite part of the architecture. Instead of letting a heavy ORM guess how to write my queries, I wrote raw SQL and used SQLC to compile it into type-safe Go code. It catches database errors at compile time instead of runtime, which gives me a lot of peace of mind.
* **Uber Zap:** Standard library logging is fine for small scripts, but for a real API, I wanted structured JSON logging. Zap makes it really easy to track down issues if things break.
* **Go-Playground Validator:** Added this at the router level. It validates the incoming JSON payloads so that bad data (like an empty email or a tiny password) gets rejected before it ever touches the database layer.

## 🛣 Endpoints

### 1. Health Check
Just a simple route to make sure the server is breathing.
* **GET** `/health`
* **Response:** `200 OK`

    {
      "message": "OmniTrace system is fully operational",
      "status": "success"
    }

### 2. Create User
Expects a JSON body. It will fail if the email isn't a real email format, or if the password is under 6 characters.
* **POST** `/users`
* **Payload:**

    {
      "email": "user@example.com",
      "password_hash": "secure123",
      "full_name": "Jane Doe"
    }

* **Response:** `201 Created` (Returns the new user record with generated UUID and timestamps).

### 3. Get User
* **GET** `/users/:email`
* **Response:** `200 OK` 

##  How to Run It Locally

If you want to spin this up on your own machine, it just takes a few steps:

1. **Clone the repo:**

    git clone https://github.com/Shreyansh3108/Omni-Trace.git
    cd Omni-Trace

2. **Grab the dependencies:**

    go mod tidy

3. **Set up the environment:**
    Create a `.env` file in the root folder and drop your Postgres connection string in there:

    PORT=3000
    DB_SOURCE=postgresql://username:password@your-db-url...

4. **Run the server:**

    go run main.go

##  What I'd Do Next
Since this was built within a tight timeframe, there are definitely things I would add before pushing this to an actual production environment:
1. Actually hash the incoming passwords using `bcrypt` (right now it just stores the hash string from the JSON payload).
2. Add JWT-based authentication to secure the user endpoints.
3. Write standard unit tests for the handler functions. 

Thanks for taking the time to review my code! Let me know if you have any questions about the architecture.
