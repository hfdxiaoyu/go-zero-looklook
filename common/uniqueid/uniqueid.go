package uniqueid

import (
	"github.com/sony/sonyflake"
	"github.com/zeromicro/go-zero/core/logx"
)

// 生成分布式的唯一id
var flake *sonyflake.Sonyflake

func init() {
	// 实例化生成分布式id的方法
	flake = sonyflake.NewSonyflake(sonyflake.Settings{})
}

func GenId() int64 {
	// NextID生成下一个唯一ID。在Sonyflake时间溢出后，NextID返回一个错误。
	id, err := flake.NextID()
	if err != nil {
		logx.Severef("flake NextID failed with %s \n", err)
		panic(err)
	}

	return int64(id)
}
