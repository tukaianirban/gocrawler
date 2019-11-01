package urlcache

import (
	"github.com/go-redis/redis"
	"log"
)

type RedisClient struct {
	redisclient			*redis.Client

	urlHitsHashName		string
}

func NewRedisCache(urlHitsHashName string) (*RedisClient, error) {

	r := RedisClient{
		redisclient:		redis.NewClient(&redis.Options{
								Addr:		"localhost:6379",
								Password:	"",
								DB:			0,
							}),
		urlHitsHashName:	urlHitsHashName,
	}

	_, err := r.redisclient.Ping().Result()

	if err!=nil {
		return nil, err
	}

	// remove any existing entries in the referred-to hashmap
	err = r.redisclient.Del(urlHitsHashName).Err()

	return &r, err
}

func (self *RedisClient)StoreURL (url string) error {

	return self.redisclient.HSet(self.urlHitsHashName, url, 1).Err()

}

func (self *RedisClient)CheckURLExists(url string) bool {

	resultcmd := self.redisclient.HGet(self.urlHitsHashName, url)

	return resultcmd.Err() == nil

}

func (self *RedisClient)SetURLHit(url string) (bool, error) {

	resultcmd := self.redisclient.HGet(self.urlHitsHashName, url)
	if resultcmd.Err() != nil {

		return false, self.StoreURL(url)
	}

	hitcount, err := resultcmd.Int()
	if err != nil {

		log.Printf("error getting hitcounts for url:%s", url)
		return true, self.StoreURL(url)
	}

	return true, self.redisclient.HSet(self.urlHitsHashName, url, hitcount+1).Err()
}

func (self *RedisClient)GetCacheSize() int64 {

	return self.redisclient.HLen(self.urlHitsHashName).Val()
}
