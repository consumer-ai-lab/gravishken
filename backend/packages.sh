#!/bin/bash

echo "Creating init file for go"
go mod init gravtest
echo "Created go.mod file"

echo "Installing all dependencies for go"

go get go.mongodb.org/mongo-driver/mongo
go get github.com/joho/godotenv
go get golang.org/x/crypto/bcrypt
go get github.com/golang-jwt/jwt/v5


echo "Installation finished"
