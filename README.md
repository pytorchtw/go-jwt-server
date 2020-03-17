# Golang JWT Login Server

This is a starter project that demonstrates the features of the following Golang related stack. Testing code included

## Backend
- golang-migrate for sql migrations - [github.com/golang-migrate/migrate](https://github.com/golang-migrate/migrate)
- sqlboiler for db orm code generation - [github.com/volatiletech/sqlboiler](https://github.com/volatiletech/sqlboiler)
- swagger and go-swagger for api code generation - [github.com/go-swagger/go-swagger](https://github.com/go-swagger/go-swagger)
- jwt-go - [github.com/dgrijalva/jwt-go](https://github.com/dgrijalva/jwt-go)
- postgresql in docker container
- docker, docker-compose

## Frontend
- react, react-hooks
- material-ui
- axios
- react-router-dom

# React Frontend Installation
```
cd react_frontend/react-material
npm install
```

# Usage

start postgresql db in container 
```
docker-compose up postgresql
```

start golang server via script
```
bash start_testserver.sh # ./start_testserver.sh
```

compile and start frontend server
```
cd react_frontend/react-materail
npm start
```

# Screenshots

![](/screenshots/screenshot_1.png?raw=true)
