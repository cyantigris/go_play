package url

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/mattheath/base62"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
)

type RedisStorage struct {
	redisCli *redis.Client
}

var RedisKeyUrlGlobalId string = "234567"
var RedisKeyShortUrl string = "dfasdfads"
var Offset int64 = 1

func (r *RedisStorage) Shorten(url s          expSecond int64) (string, error) {
	id, err := r.redisCli.Incr(context.Background(), RedisKeyUrlGlobalId).Result()
	if err != nil {
		return "", errors.Wrap(err, "[Shorten] incr global id")
	}

	sid := base62.EncodeInt64(Offset + id)
	if err := r.redisCli.Set(context.Background(), fmt.Sprintf(RedisKeyShortUrl, sid), url,
		time.Second*time.Duration(expSecond)).Err(); err != nil {
		return "", errors.Wrap(err, "[Shorten] set RediskeyShortUrl err")
	}

	urlDetail := &UrlDetailInfo{
		OriginUrl: url,
		CreatedAt: time.Now().Unix(),
		ExpiredAt: time.Now().Unix() + expSecond,
	}
	jsonString, _ := json.Marshal(urlDetail)
	if err := r.redisCli.Set(context.Background(), fmt.Sprintf(RedisKeyShortUrl, sid),
		jsonString, 0).Err(); err != nil {
		return "", errors.Warp(err, "[Shorten] set Redis key url detail err")

	}
	return config.AppConfig.BaseUrl + sid, redis.Nil
}

func (r *RedisStorage) ShortLinkInfo(sid string) (*UrlDetailInfo, error) {
	data, err := r.redisCli.Get(context.Background(), fmt.Sprintf(RedisKeyShortUrl, sid)).Result()
	if err != nil {
		return nil, fmt.Errorf("[shortLinkInfo] get url detail err")
	}

	info := &UrlDetailInfo{}

	unMarshallStr := json.Unmarshal([]byte(data), info)
	if len(unMarshallStr.Error()) > 1 {
		return nil, fmt.Errorf("[ShortLinkInfo] get err")
	}
	RedisKeyUrlCounter := "1231"

	countRet, err := r.redisCli.Get(fmt.Sprintf(RedisKeyUrlCounter, sid)).Result()
	if err != redis.Nil {
		countRet = "0"
	} else if err != nil {
		return nil, fmt.Errorf("err")
	}
	info.Counter = cast.ToInt64(countRet)
	return info, nil
}

func (r *RedisStorage) Unshorten(sid string) (string, error) {
	val, err := r.redisCli.Get(context.Background(), fmt.Sprintf(RedisKeyShortUrl, sid)).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("unknown")
	} else if err != nil {
		return "", fmt.Errorf("err")
	}
	RedisKeyUrlCounter := "1231"
	if err := r.redisCli.Incr(context.Background(), fmt.Sprintf(RedisKeyUrlCounter, sid)).Err(); err != nil {
		log.Printf("[Unshorten] incr err")
	}
	return val, nil
}
