package jwt

import (
	"sso/internal/domain/models"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestNewToken(t *testing.T) {
	user := models.User{
		ID:    42,
		Email: "user@example.com",
	}
	app := models.App{
		ID:     1,
		Name:   "test-app",
		Secret: "test-secret-key",
	}
	duration := 24 * time.Hour

	tokenString, err := NewToken(user, app, duration)
	if err != nil {
		t.Fatalf("NewToken() error = %v", err)
	}
	if tokenString == "" {
		t.Fatal("NewToken() returned empty string")
	}

	// Парсим токен и проверяем claims
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(app.Secret), nil
	})
	if err != nil {
		t.Fatalf("Parse token: %v", err)
	}
	if !token.Valid {
		t.Fatal("token is not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("claims are not MapClaims")
	}

	if uid, ok := claims["uid"].(float64); !ok || int64(uid) != user.ID {
		t.Errorf("claim uid = %v, want %d", claims["uid"], user.ID)
	}
	if email, ok := claims["email"].(string); !ok || email != user.Email {
		t.Errorf("claim email = %v, want %s", claims["email"], user.Email)
	}
	if appID, ok := claims["app_id"].(float64); !ok || int(appID) != app.ID {
		t.Errorf("claim app_id = %v, want %d", claims["app_id"], app.ID)
	}
	if appName, ok := claims["app_name"].(string); !ok || appName != app.Name {
		t.Errorf("claim app_name = %v, want %s", claims["app_name"], app.Name)
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		t.Fatal("claim exp missing or wrong type")
	}
	expTime := time.Unix(int64(exp), 0)
	if expTime.Before(time.Now()) {
		t.Errorf("exp %v is in the past", expTime)
	}
	// exp должен быть примерно через duration (с допуском 2 секунды)
	expectedExp := time.Now().Add(duration)
	if expTime.Sub(expectedExp).Abs() > 2*time.Second {
		t.Errorf("exp = %v, want ~%v", expTime, expectedExp)
	}
}

func TestNewToken_InvalidSecret(t *testing.T) {
	user := models.User{ID: 1, Email: "a@b.c"}
	app := models.App{ID: 1, Name: "app", Secret: "secret"}
	tokenString, err := NewToken(user, app, time.Hour)
	if err != nil {
		t.Fatalf("NewToken() error = %v", err)
	}

	// Токен, подписанный другим секретом, не должен валидироваться
	_, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("wrong-secret"), nil
	})
	if err == nil {
		t.Error("expected error when parsing with wrong secret")
	}
}

func TestNewToken_EmptySecret(t *testing.T) {
	user := models.User{ID: 1, Email: "a@b.c"}
	app := models.App{ID: 1, Name: "app", Secret: ""}
	_, err := NewToken(user, app, time.Hour)
	if err != nil {
		t.Errorf("NewToken() with empty secret: unexpected error %v", err)
	}
}

func TestNewToken_ZeroDuration(t *testing.T) {
	user := models.User{ID: 1, Email: "a@b.c"}
	app := models.App{ID: 1, Name: "app", Secret: "x"}
	tokenString, err := NewToken(user, app, 0)
	if err != nil {
		t.Fatalf("NewToken() error = %v", err)
	}
	if tokenString == "" {
		t.Fatal("NewToken() returned empty string")
	}
	// При duration=0 токен сразу истекает — проверяем только что exp установлен в ~now
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		t.Fatalf("ParseUnverified: %v", err)
	}
	claims := token.Claims.(jwt.MapClaims)
	exp, ok := claims["exp"].(float64)
	if !ok {
		t.Fatal("claim exp missing")
	}
	expTime := time.Unix(int64(exp), 0)
	if expTime.After(time.Now().Add(time.Second)) {
		t.Errorf("exp with zero duration should be ~now, got %v", expTime)
	}
}
