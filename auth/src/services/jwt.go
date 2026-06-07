package services

import (
	"auth/logger"
	"auth/models"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// 秘密鍵
	JwtPrivateKey ed25519.PrivateKey

	// 公開鍵
	JwtPublicKey ed25519.PublicKey
)

func initJwt(certString string) {
	// PEMブロックの解析
	// PEM形式の秘密鍵を解析
	block, _ := pem.Decode([]byte(certString))
	if block == nil {
		logger.PrintErr("PEMデータの解析に失敗しました")
		return
	}

	// 秘密鍵のパース
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		logger.PrintErr("秘密鍵のパースに失敗しました: %v", err)
		return
	}

	// Ed25519秘密鍵の型アサーション
	edPrivateKey, ok := privateKey.(ed25519.PrivateKey)
	if !ok {
		logger.PrintErr("キーはEd25519秘密鍵ではありません")
		return
	}

	// 秘密鍵を保存
	JwtPrivateKey = edPrivateKey

	// 秘密鍵の情報を表示
	logger.Println("Ed25519秘密鍵の読み込みに成功しました\n")
	logger.Println("秘密鍵の長さ: %d バイト\n", len(edPrivateKey))

	// 対応する公開鍵を取得 (秘密鍵の後半32バイトが公開鍵)
	publicKey := edPrivateKey.Public().(ed25519.PublicKey)
	logger.Println("対応する公開鍵の長さ: %d バイト\n", len(publicKey))
	logger.Println("対応する公開鍵の値: %v\n", publicKey)

	// 公開鍵を保存
	JwtPublicKey = publicKey
}

const (
	tokenExpiry = time.Minute * 10
)

type AccessTokenClaim struct {
	UserID   string   // ユーザーID
	Name     string   // ユーザー名
	Email    string   // メールアドレス
	Labels   []string // ラベル
	ProvCode models.ProviderCode
	ProvUid  string
}

func AccessTokenJwt(args AccessTokenClaim) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims{
		"exp":      time.Now().Add(tokenExpiry).Unix(),
		"userID":   args.UserID,
		"name":     args.Name,
		"email":    args.Email,
		"labels":   args.Labels,
		"provCode": args.ProvCode,
		"provUid":  args.ProvUid,
	})

	tokenString, err := token.SignedString(JwtPrivateKey)
	return tokenString, err
}

type AccessTokenInfo struct {
	UserID   string
	Name     string
	Email    string
	Labels   []string
	ProvCode models.ProviderCode
	ProvUid  string
	Exp      int64
}

func ParseAccessToken(tokenString string) (*AccessTokenInfo, error) {
	if JwtPublicKey == nil {
		return nil, errors.New("JWT公開鍵が初期化されていません")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return JwtPublicKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodEdDSA.Alg()}))

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	userID, ok := claims["userID"].(string)
	if !ok || userID == "" {
		return nil, errors.New("トークンにuserIDクレームがありません")
	}

	info := &AccessTokenInfo{UserID: userID}

	if v, ok := claims["name"].(string); ok {
		info.Name = v
	}
	if v, ok := claims["email"].(string); ok {
		info.Email = v
	}
	if v, ok := claims["provCode"].(string); ok {
		info.ProvCode = models.ProviderCode(v)
	}
	if v, ok := claims["provUid"].(string); ok {
		info.ProvUid = v
	}
	if v, ok := claims["exp"].(float64); ok {
		info.Exp = int64(v)
	}

	if rawLabels, ok := claims["labels"].([]interface{}); ok {
		for _, l := range rawLabels {
			if s, ok := l.(string); ok {
				info.Labels = append(info.Labels, s)
			}
		}
	}

	return info, nil
}
