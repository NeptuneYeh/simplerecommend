package controllers

import (
	"encoding/json"
	"github.com/NeptuneYeh/simplerecommend/init/auth"
	"github.com/NeptuneYeh/simplerecommend/init/config"
	RedisModule "github.com/NeptuneYeh/simplerecommend/init/redis"
	"github.com/NeptuneYeh/simplerecommend/init/store"
	"github.com/NeptuneYeh/simplerecommend/internal/app/requests/recommendProductRequests"
	db "github.com/NeptuneYeh/simplerecommend/internal/infra/database/mysql/sqlc"
	"github.com/NeptuneYeh/simplerecommend/pkg/token"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

type RecommendProductController struct {
	store       db.Store
	auth        token.Maker
	config      *config.Module
	redisClient *redis.Client
}

func NewRecommendProductController() *RecommendProductController {
	return &RecommendProductController{
		store:       *store.MyStore,
		auth:        *auth.MyAuth,
		config:      config.MyConfig,
		redisClient: RedisModule.MyRedisClient,
	}
}

func (c *RecommendProductController) ListProducts(ctx *gin.Context) {
	// check redis
	val, err := c.redisClient.Get(ctx, "recommend_products").Result()
	if err == redis.Nil {
		// "recommend_products" 不存在
	} else if err != nil {
		// TODO 寄一些緊急的 logger 例如 sentry, 失敗不應該影響流程, 但必須馬上知道
	} else {
		var respProducts []recommendProductRequests.RecommendProductResponse
		err = json.Unmarshal([]byte(val), &respProducts)
		if err != nil {
			// TODO 寄一些緊急的 logger 例如 sentry, 失敗不應該影響流程, 但必須馬上知道
		}
		ctx.JSON(http.StatusOK, respProducts)
		return
	}

	products, err := c.store.ListRecommendProducts(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}
	var respProducts []recommendProductRequests.RecommendProductResponse
	for _, v := range products {
		respProduct := recommendProductRequests.RecommendProductResponse{
			ID:    v.ID.Int64,
			Name:  v.Name.String,
			Price: v.Price.String,
		}
		respProducts = append(respProducts, respProduct)
	}
	// cache to redis
	jsonData, err := json.Marshal(respProducts)
	if err != nil {
		// TODO 紀錄 轉 json 失敗不應該影響流程
	}
	if jsonData != nil {
		err = c.redisClient.Set(ctx, "recommend_products", jsonData, 10*time.Minute).Err()
		if err != nil {
			// TODO 寄一些緊急的 logger 例如 sentry, 轉存 Redis 失敗不應該影響流程, 但必須馬上知道
		}
	}
	// 假裝要取一段時間
	time.Sleep(2 * time.Second)
	ctx.JSON(http.StatusOK, respProducts)
}
