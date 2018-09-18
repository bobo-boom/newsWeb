package models

import ("github.com/astaxie/beego/orm"
_ "github.com/go-sql-driver/mysql"
	"time"
)

//表的设计


//和article多对多
type User struct {
	Id int
	UserName string
	Passwd string
	Articles  []*Article  `orm:"rel(m2m)"`
}

type ArticleType struct {
	Id int
	TypeName string `orm:"size(20)"`
	Articles []*Article `orm:"reverse(many)"`
}


type Article struct {
	Id2 int `orm:"pk;auto"`
	Title string`orm:"size(20)"`       //文章标题
	Content string`orm:"size(500)"`		//内容
	Img string	`orm:"size(50);null"`		//图片（路径）
	Time time.Time`orm:"type(datetime);auto_now_add"`		//发布时间
	Count int`orm:"default(0)"`//阅读量
	ArticleType  *ArticleType  `orm:"rel(fk)"`
	Users  []*User  `orm:"reverse(many)"`


	//num, err := dORM.QueryTable("post").Filter("Tags__Tag__Name", "golang").All(&posts)
	// select * from post where Tag.Name = golang
}

func init(){
	orm.RegisterDataBase("default","mysql","root:123456@tcp(127.0.0.1:3306)/newsTest?charset=utf8")
	orm.RegisterModel(new(User),new(Article),new(ArticleType))
	orm.RunSyncdb("default",false,true)
}
