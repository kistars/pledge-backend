package models

import (
	"encoding/json"
	"errors"

	"github.com/kistars/pledge-backend/db"
	"gorm.io/gorm"
)

type TokenInfo struct {
	Id           int    `gorm:"column:id;primaryKey"`
	Logo         string `json:"logo" gorm:"column:logo"`
	Token        string `json:"token" gorm:"column:token"` // token合约地址
	Symbol       string `json:"symbol" gorm:"column:symbol"`
	ChainId      string `json:"chain_id" gorm:"column:chain_id"`
	Price        string `json:"price" gorm:"column:price"`
	Decimals     int    `json:"decimals" gorm:"column:decimals"`
	AbiFileExist int    `json:"abi_file_exist" gorm:"column:abi_file_exist"`
	CreatedAt    string `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    string `json:"updated_at" gorm:"column:updated_at"`
}

func NewTokenInfo() *TokenInfo {
	return &TokenInfo{}
}

// GetTokenInfo Get token information by token name
func (t *TokenInfo) GetTokenInfo(token, chainId string) (TokenInfo, error) {

	tokenInfo := TokenInfo{}
	redisKey := "token_info:" + chainId + ":" + token
	redisTokenInfoBytes, _ := db.RedisGet(redisKey)
	if len(redisTokenInfoBytes) <= 0 {
		err := db.Mysql.Table("token_info").Where("token=? and chain_id=?", token, chainId).First(&tokenInfo).Debug().Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return tokenInfo, nil
			} else {
				return tokenInfo, errors.New("record select err " + err.Error())
			}
		}
		_ = db.RedisSet(redisKey, RedisTokenInfo{
			Token:   token,
			ChainId: chainId,
			Price:   tokenInfo.Price,
			Logo:    tokenInfo.Logo,
			Symbol:  tokenInfo.Symbol,
		}, 0)
		return tokenInfo, nil
	} else {
		redisTokenInfo := RedisTokenInfo{}
		err := json.Unmarshal(redisTokenInfoBytes, &redisTokenInfo)
		if err != nil {
			return tokenInfo, errors.New("record Unmarshal err " + err.Error())
		}
		return TokenInfo{
			Logo:    redisTokenInfo.Logo,
			Token:   token,
			Symbol:  redisTokenInfo.Symbol,
			ChainId: chainId,
			Price:   redisTokenInfo.Price,
		}, nil
	}

}
