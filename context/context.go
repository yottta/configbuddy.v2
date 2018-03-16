package context

var globalCtx Context

const (
	ParsedConfigDataKey = "PARSED_CONFIG"
)

type Context interface {
	init()
	ConfigsPaths() []string
	StoreData(string, interface{})
	GetData(string) interface{}
}

type ApplicationContext struct {
	Configs []string
	Data    map[string]interface{}
}

func GetContext() Context {
	return globalCtx
}

func InitContext(ctx Context) {
	if globalCtx != nil {
		panic("Context cannot be initialized twice")
	}
	globalCtx = ctx
}

func (a *ApplicationContext) ConfigsPaths() []string {
	return a.Configs
}

func (a *ApplicationContext) StoreData(key string, val interface{}) {
	a.Data[key] = val
}

func (a *ApplicationContext) GetData(key string) interface{} {
	return a.Data[key]
}

func (a *ApplicationContext) init() {

}
