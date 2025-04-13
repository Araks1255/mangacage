package handlers

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/handlers/chapters"
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
)

func TestCreateChapter(t *testing.T) {
	key := viper.Get("SECRET_KEY").(string)

	chaptersOnModerationCollection := mongoDB.Collection("chapters_on_moderation_pages")
	chaptersPagesCollection := mongoDB.Collection("chapters_pages")

	defer func() {
		chaptersOnModerationCollection.DeleteMany(nil, bson.M{})
		chaptersPagesCollection.DeleteMany(nil, bson.M{})
	}()

	h := chapters.NewHandler(db, chaptersOnModerationCollection, chaptersPagesCollection)

	r := gin.New()
	r.Use(middlewares.AuthMiddleware(key))

	r.POST("/volume/:id/chapters", h.CreateChapter)

	var userID uint
	db.Raw("SELECT id FROM users WHERE user_name = 'user_test'").Scan(&userID)
	if userID == 0 {
		t.Fatal("Юзер не найден")
	}

	claims := models.Claims{
		ID: userID,
		StandardClaims: jwt.StandardClaims{
			Subject:   "test",
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		t.Fatal(err)
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	if err = writer.WriteField("name", "testChapter"); err != nil {
		t.Fatal(err)
	}
	if err = writer.WriteField("description", "someDescription"); err != nil {
		t.Fatal(err)
	}

	part, err := writer.CreateFormFile("pages", "file")
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile("./test_data/test_chapter_page.png")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = part.Write(data); err != nil {
		t.Fatal(err)
	}

	writer.Close()

	var volumeID uint
	db.Raw("SELECT id FROM volumes WHERE name = 'volume_test'").Scan(&volumeID)
	if volumeID == 0 {
		t.Fatal("Том не найден")
	}

	url := fmt.Sprintf("/volume/%d/chapters", volumeID)

	req := httptest.NewRequest("POST", url, &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	req.AddCookie(&http.Cookie{
		Name:  "mangacage_token",
		Value: tokenString,
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 201 {
		t.Fatal(w.Body.String())
	}
}
