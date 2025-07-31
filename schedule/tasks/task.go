package tasks

import (
	"time"

	"github.com/jasonlvhit/gocron"
	"github.com/kistars/pledge-backend/db"
	"github.com/kistars/pledge-backend/log"
	"github.com/kistars/pledge-backend/schedule/common"
	"github.com/kistars/pledge-backend/schedule/services"
)

func Task() {
	// get environment variables
	common.GetEnv()

	// flush redis db
	err := db.RedisFlushDB()
	if err != nil {
		log.Logger.Error("clear redis err =" + err.Error())
		panic("clear redis error " + err.Error())
	}

	// run pool tasks
	s := gocron.NewScheduler()
	s.ChangeLoc(time.UTC)
	_ = s.Every(2).Minutes().From(gocron.NextTick()).Do(services.NewPoolService().UpdateAllPoolInfo)
	_ = s.Every(1).Minute().From(gocron.NextTick()).Do(services.NewTokenPrice().UpdateContractPrice)
	_ = s.Every(2).Hours().From(gocron.NextTick()).Do(services.NewTokenSymbol().UpdateContractSymbol)
	_ = s.Every(2).Hours().From(gocron.NextTick()).Do(services.NewTokenLogo().UpdateTokenLogo)
	_ = s.Every(30).Minutes().From(gocron.NextTick()).Do(services.NewBalanceMonitor().Monitor)
	//_ = s.Every(30).Minutes().From(gocron.NextTick()).Do(services.NewTokenPrice().SavePlgrPrice)
	_ = s.Every(30).Minutes().From(gocron.NextTick()).Do(services.NewTokenPrice().SavePlgrPriceTestNet)
	<-s.Start() // Start all the pending jobs
}
