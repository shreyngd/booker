# Interview Slot Booking

A REST API written in Golang for booking interview slots for employees.


# Tools and libs
[Gin](https://github.com/gin-gonic/gin) 
[Mongo](https://github.com/mongodb/mongo-go-driver)
[JWT Auth](https://github.com/dgrijalva/jwt-go)
[OAuth2](https://github.com/golang/oauth2)
[WebSockets](https://github.com/gorilla/websocket)

# Getting Started

1) **Install Dependencies**
```bash
git clone https://github.com/shreyngd/booker.git
cd booker
make deps
cp .env.example .env
```
2) **Enter environment variables**
```
MONGO_URI="..."
SECRET_KEY="..."
GOOGLE_CLIENT_ID="..."
GOOGLE_CLIENT_SECRET="..."
GOOGLE_REDIRECT_URI="..."
....
```

3) **Run**
```
make run
```

