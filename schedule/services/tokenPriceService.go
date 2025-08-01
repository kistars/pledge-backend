package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kistars/pledge-backend/config"
	"github.com/kistars/pledge-backend/contract/bindings"
	"github.com/kistars/pledge-backend/db"
	"github.com/kistars/pledge-backend/log"
	serviceCommon "github.com/kistars/pledge-backend/schedule/common"
	"github.com/kistars/pledge-backend/schedule/models"
	"github.com/kistars/pledge-backend/utils"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type TokenPrice struct{}

func NewTokenPrice() *TokenPrice {
	return &TokenPrice{}
}

// UpdateContractPrice update contract price
func (s *TokenPrice) UpdateContractPrice() {
	var tokens []models.TokenInfo
	db.Mysql.Table("token_info").Find(&tokens) // read from db
	for _, t := range tokens {

		var err error
		var price int64 = 0

		if t.Token == "" {
			log.Logger.Sugar().Error("UpdateContractPrice token empty ", t.Symbol, t.ChainId)
			continue
		} else {
			switch t.ChainId {
			case config.Config.TestNet.ChainId:
				price, err = s.GetTestNetTokenPrice(t.Token) // 获取token价格
			case "56":
				// if strings.ToUpper(t.Token) == config.Config.MainNet.PlgrAddress { // get PLGR price from ku-coin(Only main network price)
				// 	priceStr, _ := db.RedisGetString("plgr_price")
				// 	priceF, _ := decimal.NewFromString(priceStr)
				// 	e8 := decimal.NewFromInt(100000000)
				// 	priceF = priceF.Mul(e8)
				// 	price = priceF.IntPart()
				// } else {
				// 	err, price = s.GetMainNetTokenPrice(t.Token)
				// }

				//err, price = s.GetMainNetTokenPrice(t.Token)

			}

			if err != nil {
				log.Logger.Sugar().Error("UpdateContractPrice err ", t.Symbol, t.ChainId, err)
				continue
			}
		}

		hasNewData, err := s.CheckPriceData(t.Token, t.ChainId, utils.Int64ToString(price))
		if err != nil {
			log.Logger.Sugar().Error("UpdateContractPrice CheckPriceData err ", err)
			continue
		}

		if hasNewData {
			err = s.SavePriceData(t.Token, t.ChainId, utils.Int64ToString(price)) // 写入mysql
			if err != nil {
				log.Logger.Sugar().Error("UpdateContractPrice SavePriceData err ", err)
				continue
			}
		}
	}
}

// GetMainNetTokenPrice get contract price on main net
func (s *TokenPrice) GetMainNetTokenPrice(token string) (int64, error) {
	ethereumConn, err := ethclient.Dial(config.Config.MainNet.NetUrl)
	if nil != err {
		log.Logger.Error(err.Error())
		return 0, err
	}

	bscPledgeOracleMainNetToken, err := bindings.NewBscPledgeOracleMainnetToken(common.HexToAddress(config.Config.MainNet.BscPledgeOracleToken), ethereumConn)
	if nil != err {
		log.Logger.Error(err.Error())
		return 0, err
	}

	price, err := bscPledgeOracleMainNetToken.GetPrice(nil, common.HexToAddress(token))
	if err != nil {
		log.Logger.Error(err.Error())
		return 0, err
	}

	return price.Int64(), nil
}

// GetTestNetTokenPrice get contract price on test net
func (s *TokenPrice) GetTestNetTokenPrice(token string) (int64, error) {
	ethereumConn, err := ethclient.Dial(config.Config.TestNet.NetUrl)
	if nil != err {
		log.Logger.Error(err.Error())
		return 0, err
	}

	// 获取合约实例
	bscPledgeOracleTestnetToken, err := bindings.NewBscPledgeOracleTestnetToken(common.HexToAddress(config.Config.TestNet.BscPledgeOracleToken), ethereumConn)
	if nil != err {
		log.Logger.Error(err.Error())
		return 0, err
	}

	price, err := bscPledgeOracleTestnetToken.GetPrice(nil, common.HexToAddress(token))
	if nil != err {
		log.Logger.Error(err.Error())
		return 0, err
	}

	return price.Int64(), nil
}

// CheckPriceData Saving price data to redis if it has new price
func (s *TokenPrice) CheckPriceData(token, chainId, price string) (bool, error) {
	redisKey := "token_info:" + chainId + ":" + token
	redisTokenInfoBytes, err := db.RedisGet(redisKey)
	if err != nil {
		log.Logger.Error(err.Error())
		return false, err
	}
	if len(redisTokenInfoBytes) <= 0 { // 缓存中没有token信息
		err = s.CheckTokenInfo(token, chainId) // 写入mysql
		if err != nil {
			log.Logger.Error(err.Error())
		}
		err = db.RedisSet(redisKey, models.RedisTokenInfo{
			Token:   token,
			ChainId: chainId,
			Price:   price,
		}, 0) // 写入redis
		if err != nil {
			log.Logger.Error(err.Error())
			return false, err
		}
	} else {
		redisTokenInfo := models.RedisTokenInfo{}
		err = json.Unmarshal(redisTokenInfoBytes, &redisTokenInfo)
		if err != nil {
			log.Logger.Error(err.Error())
			return false, err
		}

		if redisTokenInfo.Price == price {
			return false, nil
		}

		redisTokenInfo.Price = price
		err = db.RedisSet(redisKey, redisTokenInfo, 0) // 更新redis
		if err != nil {
			log.Logger.Error(err.Error())
			return true, err
		}
	}
	return true, nil
}

// CheckTokenInfo  Insert token information if it was not in mysql
func (s *TokenPrice) CheckTokenInfo(token, chainId string) error {
	tokenInfo := models.TokenInfo{}
	err := db.Mysql.Table("token_info").Where("token=? and chain_id=?", token, chainId).First(&tokenInfo).Debug().Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tokenInfo = models.TokenInfo{}
			nowDateTime := utils.GetCurDateTimeFormat()
			tokenInfo.Token = token
			tokenInfo.ChainId = chainId
			tokenInfo.UpdatedAt = nowDateTime
			tokenInfo.CreatedAt = nowDateTime
			err = db.Mysql.Table("token_info").Create(tokenInfo).Debug().Error
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

// SavePriceData Saving price data to mysql if it has new price
func (s *TokenPrice) SavePriceData(token, chainId, price string) error {

	nowDateTime := utils.GetCurDateTimeFormat()

	err := db.Mysql.Table("token_info").Where("token=? and chain_id=? ", token, chainId).Updates(map[string]interface{}{
		"price":      price,
		"updated_at": nowDateTime,
	}).Debug().Error
	if err != nil {
		log.Logger.Sugar().Error("UpdateContractPrice SavePriceData err ", err)
		return err
	}

	return nil
}

// SavePlgrPrice Saving price data to mysql if it has new price
func (s *TokenPrice) SavePlgrPrice() {
	priceStr, _ := db.RedisGetString("plgr_price")
	priceF, _ := decimal.NewFromString(priceStr)
	e8 := decimal.NewFromInt(100000000)
	priceF = priceF.Mul(e8)
	price := priceF.IntPart()

	ethereumConn, err := ethclient.Dial(config.Config.MainNet.NetUrl)
	if nil != err {
		log.Logger.Error(err.Error())
		return
	}
	bscPledgeOracleMainNetToken, err := bindings.NewBscPledgeOracleMainnetToken(common.HexToAddress(config.Config.MainNet.BscPledgeOracleToken), ethereumConn)
	if nil != err {
		log.Logger.Error(err.Error())
		return
	}

	privateKeyEcdsa, err := crypto.HexToECDSA(serviceCommon.PlgrAdminPrivateKey)
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKeyEcdsa, big.NewInt(utils.StringToInt64(config.Config.MainNet.ChainId)))
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	transactOpts := bind.TransactOpts{
		From:      auth.From,
		Nonce:     nil,
		Signer:    auth.Signer, // Method to use for signing the transaction (mandatory)
		Value:     big.NewInt(0),
		GasPrice:  nil,
		GasFeeCap: nil,
		GasTipCap: nil,
		GasLimit:  0,
		Context:   ctx,
		NoSend:    false, // Do all transact steps but do not send the transaction
	}

	_, err = bscPledgeOracleMainNetToken.SetPrice(&transactOpts, common.HexToAddress(config.Config.MainNet.PlgrAddress), big.NewInt(price))

	log.Logger.Sugar().Info("SavePlgrPrice ", err)

	a, d := s.GetMainNetTokenPrice(config.Config.MainNet.PlgrAddress)
	log.Logger.Sugar().Info("GetMainNetTokenPrice ", a, d)
}

// SavePlgrPriceTestNet  Saving price data to mysql if it has new price
func (s *TokenPrice) SavePlgrPriceTestNet() {

	price := 22222
	ethereumConn, err := ethclient.Dial(config.Config.TestNet.NetUrl)
	if nil != err {
		log.Logger.Error(err.Error())
		return
	}
	// 获取合约实例
	bscPledgeOracleTestNetToken, err := bindings.NewBscPledgeOracleMainnetToken(common.HexToAddress(config.Config.TestNet.BscPledgeOracleToken), ethereumConn)
	if nil != err {
		log.Logger.Error(err.Error())
		return
	}

	// 从环境变量中读取私钥
	privateKeyEcdsa, err := crypto.HexToECDSA(serviceCommon.PlgrAdminPrivateKey)
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKeyEcdsa, big.NewInt(utils.StringToInt64(config.Config.TestNet.ChainId)))
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	transactOpts := bind.TransactOpts{
		From:      auth.From,
		Nonce:     nil,
		Signer:    auth.Signer, // Method to use for signing the transaction (mandatory)
		Value:     big.NewInt(0),
		GasPrice:  nil,
		GasFeeCap: nil,
		GasTipCap: nil,
		GasLimit:  0,
		Context:   ctx,
		NoSend:    false, // Do all transact steps but do not send the transaction
	}

	// 更新链上价格数据
	_, err = bscPledgeOracleTestNetToken.SetPrice(&transactOpts, common.HexToAddress(config.Config.TestNet.PlgrAddress), big.NewInt(int64(price)))

	log.Logger.Sugar().Info("SavePlgrPrice ", err)

	a, d := s.GetTestNetTokenPrice(config.Config.TestNet.PlgrAddress)
	fmt.Println(a, d, 5555)
}
