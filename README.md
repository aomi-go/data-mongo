# Golang & MongoDB

基于MongoDB Drive 封装的数据访问层

### 快速开始

1. 安装

```
go get github.com/aomi-go/data-mongo
```

2. 初始化数据源环境变量

```
export MONGODB_URI=mongodb://username:pwd@127.0.0.1/db
```

3. 创建结构体和一个单例模式的数据仓库

```go
package main

import (
	"github.com/aomi-go/data-mongo/datasource"
	"github.com/aomi-go/data-mongo/repository"
	"reflect"
	"sync"
)

type User struct {
	Id   string `json:"id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name,omitempty" bson:"name,omitempty"`
}

type UserRepository struct {
	repository.MongoRepository // 继承MongoRepository, 内置了常用的方法
}

// 创建一个单例模式的UserRepository实例
var (
	instance *UserRepository
	once     sync.Once
)

// UserRepoInstance
// datasource.Database() 获取一个MongoDB 数据库实例
// reflect.TypeOf(User{}) 传入结构体的类型 
func UserRepoInstance() *UserRepository {
	once.Do(func() {
		instance = &UserRepository{
			repository.NewMongoRepo(datasource.Database(), reflect.TypeOf(User{}), ""),
		}
	})
	return instance
}

func main() {
	// 保存User
	UserRepoInstance().Save(User{Name: "Sean"})
	// 根据id查找
	UserRepoInstance().FindById("62414682643781d27a7022ce")
	// 判断id对应的数据是否存在
	UserRepoInstance().ExistsById("62414682643781d27a7022ce")
	// 根据id删除
	UserRepoInstance().DeleteById("62414682643781d27a7022ce")
}

```