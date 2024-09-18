package snowflake

import (
	"github.com/yitter/idgenerator-go/idgen"
)

var (
	// 工作区ID 用于分布式区分服务，保证并发时生成的ID唯一性
	workerID uint16 = 1
	// 是否已初始化?
	hasInit bool
)

func SetWorkerId(workerId uint16) {
	workerID = workerId
}

// initWorker
//
//	@Description: 初始化工作区 需要在调用前设定workerID
func initWorker() {
	// 创建 IdGeneratorOptions 对象，可在构造函数中输入 WorkerId：
	var options = idgen.NewIdGeneratorOptions(workerID)
	// options.WorkerIdBitLength = 10  // 默认值6，限定 WorkerId 最大值为2^6-1，即默认最多支持64个节点。
	// options.SeqBitLength = 6; // 默认值6，限制每毫秒生成的ID个数。若生成速度超过5万个/秒，建议加大 SeqBitLength 到 10。
	// options.BaseTime = Your_Base_Time // 如果要兼容老系统的雪花算法，此处应设置为老系统的BaseTime。
	// ...... 其它参数参考 IdGeneratorOptions 定义。

	// 保存参数（务必调用，否则参数设置不生效）：
	idgen.SetIdGenerator(options)

	// 以上过程只需全局一次，且应在生成ID之前完成。

	hasInit = true

}

func GenerateID() int64 {

	if hasInit == false {
		initWorker()
	}

	return idgen.NextId()
}
