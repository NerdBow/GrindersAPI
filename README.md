# ⏱️ GrindersAPI
---
## 🛑 **CAUTION** 🛑
Grinders is still in development and there is no frontend for it right now. 

It is basically not useable right now, unless you want to manually send in the HTTP request.
****

🚨 **(It is NOT the dating app! I repeat, it is NOT the dating app)** 🚨

This is the backend server for my time tracking application Grinders which is suppose to work in conjuction to the GrindersTUI.
The needs of Grinders will be ultimately tailored to what I want to put in as it is an application for myself.
I hope some people will also find the application useful as well.

## ⚙️ Features
- HTTP API
- Account system with JWT tokens
- Log creation, querying, updating, deleting

## ➡️ Plans
- Friend groups to share you logs with
- Add refresh tokens to auth

## 🧻 Documentation
### 🚧 In progress 🚧

## 🔧 Installation & Setup
```bash
# Clone the repo
git clone https://github.com/NerdBow/GrindersAPI

# Go int othe repo
cd GrindersAPI

# Create the database file
touch data/logs.db

# Create a .env for your settings
touch .env
```
### .env file (make sure to include everything)
```.env
# argon2id settings
HASHMEMORY=64
HASHLENGHT=50
HASHTIME=1
HASHTHREADS=4
SALTLENGTH=32

# JWT Settings
JWTSECRET=PASSWORDPASSWORDPASSWORDPASSWORD # Make sure to set this to something actually random if you plan to deploy this
JWTEXP=30 # This is in minutes

# The API only supports sqlite3 as of now so just keep it as this.
DATABASE=sqlite3
DATABASEFILE=./data/logs.db

# The port you want to the run the API on
PORT=8080
```
### Running with Go

```bash
go mod download
go run cmd/api/main.go
```

### Docker (I have never used docker before this. so if this sucks I am sorry 😭)

```bash
# Make the image
docker build -t grindersapi:latest ./

# Create a container from the image
# I would suggest to mount a volume for where the db is located
docker run -d -v ./data:/app/data --env-file .env -p 8080:8080 grindersapi:latest
```
