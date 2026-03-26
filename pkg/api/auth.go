package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// SigninRequest — структура входящего запроса на авторизацию
// содержит пароль, введённый пользователем
type SigninRequest struct {
	Password string `json:"password"`
}

// SigninResponse — структура ответа сервера
// либо возвращает токен, либо ошибку
type SigninResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

// Claims — структура данных внутри JWT токена
// хранит хэш пароля и стандартные JWT-поля (например срок жизни)
type Claims struct {
	Hash string `json:"hash"`
	jwt.RegisteredClaims
}

// passwordHash — вычисляет SHA256-хэш пароля
// используется для хранения в токене вместо открытого пароля
func passwordHash(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}

// makeJWT — создаёт JWT-токен
// внутрь кладёт хэш пароля и срок действия (8 часов)
// подписывает токен самим паролем
func makeJWT(password string) (string, error) {
	claims := Claims{
		Hash: passwordHash(password),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(password))
}

// validJWT — проверяет валидность JWT-токена
// проверяет подпись, срок действия и совпадение хэша пароля
func validJWT(tokenString string, password string) bool {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(password), nil
	})
	if err != nil {
		return false
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return false
	}

	return claims.Hash == passwordHash(password)
}

// signinHandler — обработчик POST /api/signin
// принимает пароль, сравнивает с TODO_PASSWORD
// если совпадает — возвращает JWT токен
func signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, SigninResponse{
			Error: "method not allowed",
		})
		return
	}

	pass := os.Getenv("TODO_PASSWORD")
	if pass == "" {
		writeJSON(w, http.StatusInternalServerError, SigninResponse{
			Error: "Пароль не задан",
		})
		return
	}

	var req SigninRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, SigninResponse{
			Error: err.Error(),
		})
		return
	}

	if req.Password != pass {
		writeJSON(w, http.StatusUnauthorized, SigninResponse{
			Error: "Неверный пароль",
		})
		return
	}

	token, err := makeJWT(pass)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, SigninResponse{
			Error: err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, SigninResponse{
		Token: token,
	})
}

// auth — middleware для проверки авторизации
// если TODO_PASSWORD задан:
//   - берёт JWT токен из cookie
//   - проверяет его валидность
//   - если всё ок → пускает дальше
//   - иначе → возвращает 401 Unauthorized
func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if pass == "" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil || cookie.Value == "" {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		if !validJWT(cookie.Value, pass) {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
