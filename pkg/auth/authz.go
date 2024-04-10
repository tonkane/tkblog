package auth

import (
	"time"

	casbin "github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	adapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

const (
	// casbin 访问控制模型 Access Control List
	// 允许管理员定义一系列的访问控制规则，这些规则指定了哪些用户或用户组（角色）可以对哪些资源执行哪些操作
	// 请求定义 （例如：用户对资源的读取）/ 策略定义 （定义和前者相同，规则部分）/  策略效果（p.eft 指的是策略效果）
	// 匹配器 用于判断请求是否匹配一个策略 keyMatch 可以用来处理通配符  regexMatch 允许正则
	aclmodel = `[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act)`
)

// 定义授权器 Enforcer-执行者
type Authz struct {
	*casbin.SyncedCachedEnforcer
}

// 初始化
func NewAuthz(db *gorm.DB) (*Authz, error) {
	adapter, err := adapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	m, _ := model.NewModelFromString(aclmodel)

	enforcer, err := casbin.NewSyncedCachedEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}

	if err := enforcer.LoadPolicy(); err != nil {
		return nil, err
	}
	// 每5秒同步一次策略
	enforcer.StartAutoLoadPolicy(5 * time.Second)

	a := &Authz{enforcer}

	return a, nil
}

// 授权操作
func (a *Authz) Authorize(sub, obj, act string) (bool, error) {
	return a.Enforce(sub, obj, act)
}