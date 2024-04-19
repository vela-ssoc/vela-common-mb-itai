package accord

const (
	PathPrefix        = "/api/v1"
	PathTaskLoad      = "/task/load"  // 加载指定配置
	PathTaskSync      = "/task/sync"  // 同步配置
	PathTaskTable     = "/task/table" // 同步配置
	PathThirdDiff     = "/third/diff"
	PathElasticReset  = "/elastic/reset"
	PathEmcReset      = "/emc/reset"
	PathEmailReset    = "/email/reset"
	PathStoreReset    = "/store/reset"
	PathNotifierReset = "/notifier/reset"
	PathStartup       = "/startup"
	PathUpgrade       = "/upgrade"
	PathCommand       = "/command"
	FPTaskLoad        = PathPrefix + PathTaskLoad
	FPTaskSync        = PathPrefix + PathTaskSync
	FPTaskTable       = PathPrefix + PathTaskTable
	FPThirdDiff       = PathPrefix + PathThirdDiff
	FPElasticReset    = PathPrefix + PathElasticReset
	FPEmcReset        = PathPrefix + PathEmcReset
	FPEmailReset      = PathPrefix + PathEmailReset
	FPStoreReset      = PathPrefix + PathStoreReset
	FPNotifierReset   = PathPrefix + PathNotifierReset
	FPStartup         = PathPrefix + PathStartup
	FPUpgrade         = PathPrefix + PathUpgrade
	FPCommand         = PathPrefix + PathCommand
)

type TaskLoadRequest struct {
	MinionID    int64  `json:"minion_id"`
	SubstanceID int64  `json:"substance_id"`
	Inet        string `json:"inet"`
}

type StoreRestRequest struct {
	ID string `json:"id"`
}

type ThirdDiff struct {
	Name  string `json:"name"`  // 三方文件名
	Event string `json:"event"` // 事件类型：delete-文件删除 update-文件更新
}

const (
	ThirdUpdate = "update"
	ThirdDelete = "delete"
)

type TaskTable struct {
	TaskID int64 `json:"task_id"`
}

type Startup struct {
	ID int64 `json:"id"`
}

type Upgrade struct {
	ID         []int64 `json:"id"`
	Semver     string  `json:"semver"`     // 二进制版本号
	Customized string  `json:"customized"` // 定制版本类型
}

type Command struct {
	ID  []int64 `json:"id"`
	Cmd string  `json:"cmd"`
}

type IDs struct {
	ID []int64 `json:"id"`
}
