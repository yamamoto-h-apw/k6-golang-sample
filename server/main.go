package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strings"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("secret") // セキュリティ上は環境変数等にすべき

type Credentials struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type Claims struct {
    Username string `json:"username"`
    jwt.RegisteredClaims
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
    var creds Credentials
    json.NewDecoder(r.Body).Decode(&creds)

    // パスワード認証（例：固定ユーザー）
    if creds.Username != "testuser" || creds.Password != "testpass" {
        http.Error(w, "unauthorized", http.StatusUnauthorized)
        return
    }

    // JWT 作成
    expirationTime := time.Now().Add(5 * time.Minute)
    claims := &Claims{
        Username: creds.Username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        http.Error(w, "internal error", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func privateHandler(w http.ResponseWriter, r *http.Request) {
    authHeader := r.Header.Get("Authorization")
    if !strings.HasPrefix(authHeader, "Bearer ") {
        http.Error(w, "missing token", http.StatusUnauthorized)
        return
    }

    tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
    claims := &Claims{}

    token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })

    if err != nil || !token.Valid {
        http.Error(w, "invalid token", http.StatusUnauthorized)
        return
    }

    fmt.Fprintf(w, "Hello, %s!", claims.Username)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "OK")
}

func main() {
    http.HandleFunc("/health", healthHandler)
    http.HandleFunc("/login", loginHandler)
    http.HandleFunc("/private", privateHandler)

    fmt.Println("Server running at http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}