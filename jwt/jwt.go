package jwt

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo"
	"github.com/twinj/uuid"
	"os"
	"strings"
	"time"
)

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type AccessDetails struct {
	AccessUuid string
	UserId     string
	UserRole   string
}

var client *redis.Client

func init() {

	os.Setenv("ACCESS_SECRET", "accesssecret")
	os.Setenv("REFRESH_SECRET", "refreshsecret")

	//Initializing redis
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err := client.Ping(client.Context()).Result()
	if err != nil {
		fmt.Println(err)
	}
}

// CreateToken create a access token and a refresh token with their uuids,expire time
func CreateToken(id, role string) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(150 * time.Minute).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(24 * 7 * time.Hour).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	var err error

	// Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = id
	atClaims["user_role"] = role
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	atClaims["user_id"] = id
	atClaims["user_role"] = role
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil

}

// CreateAuth insert the uuid and id with their expire time to redis db
func CreateAuth(id, role string, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := client.Set(client.Context(), td.AccessUuid, id+" "+role, at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := client.Set(client.Context(), td.RefreshUuid, id+" "+role, rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

// ExtractToken read the token from the request header of context
func ExtractToken(r echo.Context) string {
	bearToken := r.Request().Header.Get("Authorization")

	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

// VerifyToken get the token and check it then if it's valid it parses it to a jwt.token
func VerifyToken(r echo.Context) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SingingMethodHMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signig method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

// TokenValid validates the jwt.token
func TokenValid(r echo.Context) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

// ExtractTokenMetadata gives the uuid,id from the token
func ExtractTokenMetadata(r echo.Context) (*AccessDetails, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userid := claims["user_id"].(string)
		userrole := claims["user_role"].(string)
		return &AccessDetails{
			AccessUuid: accessUuid,
			UserId:     userid,
			UserRole:   userrole,
		}, nil
	}
	return nil, err
}

// FetchAuth give the id,role by existing uuid from redis
func FetchAuth(authD *AccessDetails) (id, role string, err error) {
	var ud string
	err = client.Get(client.Context(), authD.AccessUuid).Scan(&ud)
	if err != nil {
		return "", "", err
	}

	strArr := strings.Split(ud, " ")
	if len(strArr) == 2 {
		return strArr[0], strArr[1], nil
	}
	return "", "", errors.New("sth unusual happened")
}
