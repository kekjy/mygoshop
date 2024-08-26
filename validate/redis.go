package main

import (
	"context"
	"fmt"
)

// lua 脚本
var script = `
	-- redis-cli --eval product.lua product-1 3
	-- 获取商品库存信息
	local counts = redis.call("HMGET", KEYS[1], "total", "ordered")
	-- 将总库存转换为数值
	local total = tonumber(counts[1])
	-- 将已被秒杀的库存转换为数值
	local ordered = tonumber(counts[2])
	-- 获得订单购买数量
	local delta = KEYS[2]
	-- 如果当前请求的库存量加上已被秒杀的库存量仍然小于总库存量，就可以更新库存
	if ordered + delta <= total then
		-- 更新已秒杀的库存量
		redis.call("HINCRBY", KEYS[1], "ordered", delta)
		return delta
	end
	return 0
`

// 获取秒杀商品
func GetProductFromRedis(productName string, productNum string) bool {
	ctx := context.Background()
	fmt.Println("productName: ", productName)
	fmt.Println("productNum: ", productNum)
	result, err := rdb.Eval(ctx, script, []string{productName, productNum}).Result()
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}

	if result == productNum {
		fmt.Println("success")
		return true
	} else {
		fmt.Println("failed")
		return false
	}
}
