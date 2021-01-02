package redis

import (
	"context"
	"fmt"
	"testing"
)

func TestGetAndSet(t *testing.T) {
	redis := New()
	err := redis.MSet(context.Background(), map[string]interface{}{"haha": 1111, "hehe": 2222})
	fmt.Println(err)
}
