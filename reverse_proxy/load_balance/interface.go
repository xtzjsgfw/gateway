package load_balance

type LoadBalance interface {
	Add(params ...string) error
	Get(string) (string, error)

	// 服务发现补充
	Update()
}
