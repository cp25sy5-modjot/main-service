package auth
import (
  "context"
  "encoding/json"
  "net/http"
  "os"
  "time"
  "google.golang.org/api/idtoken"
  "github.com/golang-jwt/jwt/v5"
)

type loginReq struct{ IDToken string `json:"id_token"` }

func GoogleLogin(w http.ResponseWriter, r *http.Request) {
  var req loginReq
  if err := json.NewDecoder(r.Body).Decode(&req); err != nil { http.Error(w,"bad request",400); return }

  v, _ := idtoken.NewValidator(context.Background())
  // AUDIENCE MUST MATCH the client ID that produced the token (iOS/Android/Web)
  payload, err := v.Validate(r.Context(), req.IDToken, "YOUR_CLIENT_ID.apps.googleusercontent.com")
  if err != nil { http.Error(w,"invalid token",401); return }

  sub := payload.Claims["sub"].(string)
  email := payload.Claims["email"].(string)

  // upsert user by `sub`...

  claims := jwt.MapClaims{
    "sub": sub, "email": email,
    "aud": "yourapp-api", "iss": "https://api.yourapp.local",
    "iat": time.Now().Unix(), "exp": time.Now().Add(15*time.Minute).Unix(),
  }
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  signed, _ := token.SignedString([]byte(os.Getenv("APP_JWT_SECRET")))
  json.NewEncoder(w).Encode(map[string]string{"access_token": signed})
}
