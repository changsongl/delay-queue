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

Outer:
	for i := 0; i < 10; i++ {
		fmt.Println(i)
		if i == 1 {
			fmt.Println("outer")
			continue Outer
		}
	}
	fmt.Println("end")
}
