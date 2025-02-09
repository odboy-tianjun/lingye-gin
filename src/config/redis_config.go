package config

import (
	"fmt"
	"github.com/go-redis/redis"
	"math/rand"
	"time"
)

var redisdb *redis.Client
var RedisClient *redis.Client

type RedisPool struct{}

func (p RedisPool) Init() {
	// 连接服务器
	redisdb = redis.NewClient(&redis.Options{
		Addr:     AppProps.Redis.Addr,
		Password: AppProps.Redis.Passwd,
		DB:       AppProps.Redis.Database,
	})

	// 开放
	RedisClient = redisdb

	// 心跳
	pong, err := redisdb.Ping().Result()

	// Output: PONG <nil>
	Logger.Info(pong, err)
}

func (p RedisPool) Test() {
	ExampleClient_String()
	ExampleClient_List()
	ExampleClient_Hash()
	ExampleClient_Set()
	ExampleClient_SortSet()
	ExampleClient_HyperLogLog()
	ExampleClient_CMD()
	ExampleClient_Scan()
	ExampleClient_Tx()
	ExampleClient_Script()
	ExampleClient_PubSub()
}

func ExampleClient_String() {
	Logger.Println("ExampleClient_String")
	defer Logger.Println("ExampleClient_String")

	//kv读写
	err := redisdb.Set("key", "value", 1*time.Second).Err()
	Logger.Println(err)

	//获取过期时间
	tm, err := redisdb.TTL("key").Result()
	Logger.Println(tm)

	val, err := redisdb.Get("key").Result()
	Logger.Println(val, err)

	val2, err := redisdb.Get("missing_key").Result()
	if err == redis.Nil {
		Logger.Println("missing_key does not exist")
	} else if err != nil {
		Logger.Println("missing_key", val2, err)
	}

	//不存在才设置 过期时间 nx ex
	value, err := redisdb.SetNX("counter", 0, 1*time.Second).Result()
	Logger.Println("setnx", value, err)

	//Incr
	result, err := redisdb.Incr("counter").Result()
	Logger.Println("Incr", result, err)
}

func ExampleClient_List() {
	Logger.Println("ExampleClient_List")
	defer Logger.Println("ExampleClient_List")

	//添加
	Logger.Println(redisdb.RPush("list_test", "message1").Err())
	Logger.Println(redisdb.RPush("list_test", "message2").Err())

	//设置
	Logger.Println(redisdb.LSet("list_test", 2, "message set").Err())

	//remove
	ret, err := redisdb.LRem("list_test", 3, "message1").Result()
	Logger.Println(ret, err)

	rLen, err := redisdb.LLen("list_test").Result()
	Logger.Println(rLen, err)

	//遍历
	lists, err := redisdb.LRange("list_test", 0, rLen-1).Result()
	Logger.Println("LRange", lists, err)

	//pop没有时阻塞
	result, err := redisdb.BLPop(1*time.Second, "list_test").Result()
	Logger.Println("result:", result, err, len(result))
}

func ExampleClient_Hash() {
	Logger.Println("ExampleClient_Hash")
	defer Logger.Println("ExampleClient_Hash")

	datas := map[string]interface{}{
		"name": "LI LEI",
		"sex":  1,
		"age":  28,
		"tel":  123445578,
	}

	//添加
	if err := redisdb.HMSet("hash_test", datas).Err(); err != nil {
		Logger.Fatal(err)
	}

	//获取
	rets, err := redisdb.HMGet("hash_test", "name", "sex").Result()
	Logger.Println("rets:", rets, err)

	//成员
	retAll, err := redisdb.HGetAll("hash_test").Result()
	Logger.Println("retAll", retAll, err)

	//存在
	bExist, err := redisdb.HExists("hash_test", "tel").Result()
	Logger.Println(bExist, err)

	bRet, err := redisdb.HSetNX("hash_test", "id", 100).Result()
	Logger.Println(bRet, err)

	//删除
	Logger.Println(redisdb.HDel("hash_test", "age").Result())
}

func ExampleClient_Set() {
	Logger.Println("ExampleClient_Set")
	defer Logger.Println("ExampleClient_Set")

	//添加
	ret, err := redisdb.SAdd("set_test", "11", "22", "33", "44").Result()
	Logger.Println(ret, err)

	//数量
	count, err := redisdb.SCard("set_test").Result()
	Logger.Println(count, err)

	//删除
	ret, err = redisdb.SRem("set_test", "11", "22").Result()
	Logger.Println(ret, err)

	//成员
	members, err := redisdb.SMembers("set_test").Result()
	Logger.Println(members, err)

	bret, err := redisdb.SIsMember("set_test", "33").Result()
	Logger.Println(bret, err)

	redisdb.SAdd("set_a", "11", "22", "33", "44")
	redisdb.SAdd("set_b", "11", "22", "33", "55", "66", "77")
	//差集
	diff, err := redisdb.SDiff("set_a", "set_b").Result()
	Logger.Println(diff, err)

	//交集
	inter, err := redisdb.SInter("set_a", "set_b").Result()
	Logger.Println(inter, err)

	//并集
	union, err := redisdb.SUnion("set_a", "set_b").Result()
	Logger.Println(union, err)

	ret, err = redisdb.SDiffStore("set_diff", "set_a", "set_b").Result()
	Logger.Println(ret, err)

	rets, err := redisdb.SMembers("set_diff").Result()
	Logger.Println(rets, err)
}

func ExampleClient_SortSet() {
	Logger.Println("ExampleClient_SortSet")
	defer Logger.Println("ExampleClient_SortSet")

	addArgs := make([]redis.Z, 100)
	for i := 1; i < 100; i++ {
		addArgs = append(addArgs, redis.Z{Score: float64(i), Member: fmt.Sprintf("a_%d", i)})
	}
	//Logger.Println(addArgs)

	Shuffle := func(slice []redis.Z) {
		r := rand.New(rand.NewSource(time.Now().Unix()))
		for len(slice) > 0 {
			n := len(slice)
			randIndex := r.Intn(n)
			slice[n-1], slice[randIndex] = slice[randIndex], slice[n-1]
			slice = slice[:n-1]
		}
	}

	//随机打乱
	Shuffle(addArgs)

	//添加
	ret, err := redisdb.ZAddNX("sortset_test", addArgs...).Result()
	Logger.Println(ret, err)

	//获取指定成员score
	score, err := redisdb.ZScore("sortset_test", "a_10").Result()
	Logger.Println(score, err)

	//获取制定成员的索引
	index, err := redisdb.ZRank("sortset_test", "a_50").Result()
	Logger.Println(index, err)

	count, err := redisdb.SCard("sortset_test").Result()
	Logger.Println(count, err)

	//返回有序集合指定区间内的成员
	rets, err := redisdb.ZRange("sortset_test", 10, 20).Result()
	Logger.Println(rets, err)

	//返回有序集合指定区间内的成员分数从高到低
	rets, err = redisdb.ZRevRange("sortset_test", 10, 20).Result()
	Logger.Println(rets, err)

	//指定分数区间的成员列表
	rets, err = redisdb.ZRangeByScore("sortset_test", redis.ZRangeBy{Min: "(30", Max: "(50", Offset: 1, Count: 10}).Result()
	Logger.Println(rets, err)
}

//用来做基数统计的算法，HyperLogLog 的优点是，在输入元素的数量或者体积非常非常大时，计算基数所需的空间总是固定 的、并且是很小的。
//每个 HyperLogLog 键只需要花费 12 KB 内存，就可以计算接近 2^64 个不同元素的基 数
func ExampleClient_HyperLogLog() {
	Logger.Println("ExampleClient_HyperLogLog")
	defer Logger.Println("ExampleClient_HyperLogLog")

	for i := 0; i < 10000; i++ {
		redisdb.PFAdd("pf_test_1", fmt.Sprintf("pfkey%d", i))
	}
	ret, err := redisdb.PFCount("pf_test_1").Result()
	Logger.Println(ret, err)

	for i := 0; i < 10000; i++ {
		redisdb.PFAdd("pf_test_2", fmt.Sprintf("pfkey%d", i))
	}
	ret, err = redisdb.PFCount("pf_test_2").Result()
	Logger.Println(ret, err)

	redisdb.PFMerge("pf_test", "pf_test_2", "pf_test_1")
	ret, err = redisdb.PFCount("pf_test").Result()
	Logger.Println(ret, err)
}

func ExampleClient_PubSub() {
	Logger.Println("ExampleClient_PubSub")
	defer Logger.Println("ExampleClient_PubSub")
	//发布订阅
	pubsub := redisdb.Subscribe("subkey")
	_, err := pubsub.Receive()
	if err != nil {
		Logger.Fatal("pubsub.Receive")
	}
	ch := pubsub.Channel()
	time.AfterFunc(1*time.Second, func() {
		Logger.Println("Publish")

		err = redisdb.Publish("subkey", "test publish 1").Err()
		if err != nil {
			Logger.Fatal("redisdb.Publish", err)
		}

		redisdb.Publish("subkey", "test publish 2")
	})
	for msg := range ch {
		Logger.Println("recv channel:", msg.Channel, msg.Pattern, msg.Payload)
	}
}

func ExampleClient_CMD() {
	Logger.Println("ExampleClient_CMD")
	defer Logger.Println("ExampleClient_CMD")

	//执行自定义redis命令
	Get := func(rdb *redis.Client, key string) *redis.StringCmd {
		cmd := redis.NewStringCmd("get", key)
		redisdb.Process(cmd)
		return cmd
	}

	v, err := Get(redisdb, "NewStringCmd").Result()
	Logger.Println("NewStringCmd", v, err)

	v, err = redisdb.Do("get", "redisdb.do").String()
	Logger.Println("redisdb.Do", v, err)
}

func ExampleClient_Scan() {
	Logger.Println("ExampleClient_Scan")
	defer Logger.Println("ExampleClient_Scan")

	//scan
	for i := 1; i < 1000; i++ {
		redisdb.Set(fmt.Sprintf("skey_%d", i), i, 0)
	}

	cusor := uint64(0)
	for {
		keys, retCusor, err := redisdb.Scan(cusor, "skey_*", int64(100)).Result()
		Logger.Println(keys, cusor, err)
		cusor = retCusor
		if cusor == 0 {
			break
		}
	}
}

func ExampleClient_Tx() {
	pipe := redisdb.TxPipeline()
	incr := pipe.Incr("tx_pipeline_counter")
	pipe.Expire("tx_pipeline_counter", time.Hour)

	// Execute
	//
	//     MULTI
	//     INCR pipeline_counter
	//     EXPIRE pipeline_counts 3600
	//     EXEC
	//
	// using one rdb-server roundtrip.
	_, err := pipe.Exec()
	fmt.Println(incr.Val(), err)
}

func ExampleClient_Script() {
	IncrByXX := redis.NewScript(`
        if redis.call("GET", KEYS[1]) ~= false then
            return redis.call("INCRBY", KEYS[1], ARGV[1])
        end
        return false
    `)

	n, err := IncrByXX.Run(redisdb, []string{"xx_counter"}, 2).Result()
	fmt.Println(n, err)

	err = redisdb.Set("xx_counter", "40", 0).Err()
	if err != nil {
		panic(err)
	}

	n, err = IncrByXX.Run(redisdb, []string{"xx_counter"}, 2).Result()
	fmt.Println(n, err)
}
