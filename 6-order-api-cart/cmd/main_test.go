package main

import (
	"bytes"
	"demo/order/6-order-api-cart/internal/auth"
	"demo/order/6-order-api-cart/internal/orders"
	"demo/order/6-order-api-cart/pkg/db"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func migrate() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	database, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	database.AutoMigrate(&db.Product{}, &db.User{}, &db.Order{}, &db.ProductOrder{})
}

func TestMigrationDb(t *testing.T) {
	migrate()
}

func initDb() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	database, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return database
}

func initData(database *gorm.DB) {
	database.Create(&db.Product{
		Name:        "Яблоки",
		Description: "Ранний урожай",
	})
}

func removeData(database *gorm.DB) {
	// database.Unscoped().Where("name = ?", "Яблоки").Delete(&db.Product{})
	database.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&db.Product{})
	database.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&db.ProductOrder{})
	database.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&db.Order{})
}

func TestLoginSuccess(t *testing.T) {
	db := initDb()
	initData(db)
	ts := httptest.NewServer(App())
	defer ts.Close()

	data, _ := json.Marshal(&auth.AuthRequest{
		Phone: "79991234567",
	})
	res, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected %d got %d", http.StatusOK, res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	var resp auth.AuthResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		t.Fatal(err)
	}
	if resp.SessionId == "" {
		t.Fatalf("Expected sessionId, got %s", resp.SessionId)
	}

	//test verify
	data, _ = json.Marshal(&auth.VerifyRequest{
		SessionId: resp.SessionId,
		Code:      3245,
	})
	vres, err := http.Post(ts.URL+"/auth/verify", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if vres.StatusCode != http.StatusOK {
		t.Fatalf("Expected %d got %d", http.StatusOK, res.StatusCode)
	}
	body, err = io.ReadAll(vres.Body)
	if err != nil {
		t.Fatal(err)
	}
	var vResp auth.VerifyResponse
	err = json.Unmarshal(body, &vResp)
	if err != nil {
		t.Fatal(err)
	}
	if vResp.Token == "" {
		t.Fatalf("Expected token, got %s", vResp.Token)
	}
	// t.Error(vResp.Token)
	// fmt.Println("token", vResp.Token)

}

func TestOrderSuccess(t *testing.T) {
	db := initDb()
	initData(db)
	ts := httptest.NewServer(App())
	defer ts.Close()

	ordData, _ := json.Marshal(&orders.OrderCreateRequest{
		Products: []string{"Яблоки"},
	})

	req, err := http.NewRequest(http.MethodPost, ts.URL+"/order", bytes.NewReader(ordData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwaG9uZSI6Ijc5OTkxMjM0NTY3In0.sMTFf4tgJZFWEOJluYSLFQocnAOVHxEWo1sdJSPEwvM")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected %d got %d", http.StatusCreated, resp.StatusCode)
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)

	var createResp orders.OrderCreateResponse
	err = json.Unmarshal(respBody, &createResp)
	if err != nil {
		t.Fatal(err)
	}
	if createResp.OrderId == 0 {
		t.Fatalf("Expected orderId, got %d", createResp.OrderId)
	}

	removeData(db)
}
