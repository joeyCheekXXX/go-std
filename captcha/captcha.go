package captcha

import (
	"github.com/mojocn/base64Captcha"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"time"
)

// 验证码 支持单点及分布式

func init() {
	// init rand seed
	rand.NewSource(time.Now().UnixNano())
}

// Captcha basic information.
type Captcha struct {
	Driver base64Captcha.Driver
	Store  base64Captcha.Store
}

// NewCaptcha creates a captcha instance from driver and store
func NewCaptcha() *Captcha {
	driver := base64Captcha.NewDriverDigit(120, 220, 5, 0.7, 100)
	return &Captcha{Driver: driver, Store: base64Captcha.DefaultMemStore}
}

func (c *Captcha) SetRedisStore(redisStore redis.UniversalClient) {
	c.Store = NewDefaultRedisStore(redisStore)
}

// Generate generates a random id, base64 image string or an error if any
func (c *Captcha) Generate() (id, b64s string, err error) {
	id, content, answer := c.Driver.GenerateIdQuestionAnswer()
	item, err := c.Driver.DrawCaptcha(content)
	if err != nil {
		return "", "", err
	}
	c.Store.Set(id, answer)
	b64s = item.EncodeB64string()
	return
}

// Verify by a given id key and remove the captcha value in store,
// return boolean value.
// if you has multiple captcha instances which share a same store.
// You may want to call `store.Verify` method instead.
func (c *Captcha) Verify(id, answer string, clear bool) (match bool) {
	match = c.Store.Get(id, clear) == answer
	return
}
