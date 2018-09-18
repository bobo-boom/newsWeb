package controllers

import (
	"github.com/astaxie/beego"
	"path"
	"time"
	"github.com/astaxie/beego/orm"
	"be_bobo1/models"
	"strconv"
	"math"
	"github.com/gomodule/redigo/redis"
	"encoding/gob"
	"bytes"
)

type ArticleController struct {
	beego.Controller
}

//文章列表页
func (this *ArticleController) ShowArticleList() {
	//1.查询
	//1.有一个orm对象
	o := orm.NewOrm()
	qs := o.QueryTable("Article")
	var articles []models.Article
	//qs.All(&articles)  //select * from Article

	//2.把数据传递给视图显示

	//3.分页功能

	//每一页数
	pageIndex := 1
	pageSize := 2

	//实现首页最后一页跳转
	pageIndex1 := this.GetString("pageIndex")
	if pageIndex1 != "" {
		pageIndex, _ = strconv.Atoi(pageIndex1)
	}
	//数据总数
	count, _ := qs.Count()
	//总页数
	pageCount := math.Ceil(float64(count) / float64(pageSize))
	//一页数据传递显示
	//第一个参数pageSize,设置获取多大数据量
	//第二个参数start,设置起始位置
	start:=pageSize*(pageIndex-1)
	qs.Limit(pageSize,start ).All(&articles)

	//上一页下一页实现
	//函数绑定

	//上一页下一页非使能
	FirstPage:= false
	if  pageIndex==1{
		FirstPage=true
	}

	//4.分类显示

	//获取类型数据
	var types []models.ArticleType




	//将类型与redis数据库进行联系

	conn,err:=redis.Dial("tcp",":6379")
	if  err!=nil{
		beego.Info("redis连接失败")
		return
	}
	defer conn.Close()
	typeBytes,err:=redis.Bytes(conn.Do("get","types" ))
	if  err!=nil{
		beego.Info("获取redis数据失败")

	}else {
		beego.Info("获取redis数据成功")
	}
	//解码
	deco:=gob.NewDecoder(bytes.NewReader(typeBytes))
	deco.Decode(&types)

	if  len(types)==0{
		o.QueryTable("ArticleType").All(&types)
		beego.Info("从mysql获取数据")
		var  buffer  bytes.Buffer
		eco:=gob.NewEncoder(&buffer)
		eco.Encode(types)
		conn.Do("set","types",buffer.Bytes())
	}
	this.Data["types"] = types
	//根据类型获取数据
	//1.接受数据
	typeName:=this.GetString("select")
	//beego.Info(typeName)
	//2.处理数据
	var articleswithtype []models.Article
	if typeName == ""{
		beego.Info("下拉框传递数据失败")
		//qs.Limit(pageSize,start).RelatedSel("ArticleType").All(&articleswithtype)
		qs.Limit(pageSize,start).RelatedSel("ArticleType").All(&articleswithtype)
	}else {
		//qs.Limit(pageSize,start).RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).All(&articleswithtype)

		qs.Limit(pageSize,start).RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).All(&articleswithtype)
	}

	this.Data["TypeName"]=typeName
	this.Data["FirstPage"]=FirstPage
	this.Data["pageIndex"] = pageIndex
	this.Data["pageCount"] = pageCount
	this.Data["count"] = count
	//this.Data["articles"] = articles
	this.Data["articles"] = articleswithtype
	this.TplName = "index.html"
}

func (this *ArticleController) ShowAddArticle() {
	//查询类型数据，传递到视图中
	o:=orm.NewOrm()
	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)
	this.Data["types"] = types
	this.Layout = "layout.html"
	this.TplName = "add.html"
}

/*
1.那数据
2.判断数据
3.插入数据
4.返回试图
 */

func (this *ArticleController) HandleAddArtcile() {
	//1.那数据
	//那标题
	artiName := this.GetString("articleName")
	artiContent := this.GetString("content")
	f, h, err := this.GetFile("uploadname")
	defer f.Close()
	//上传文件处理
	//1.判断文件格式
	ext := path.Ext(h.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		beego.Info("上传文件格式不正确")
		return
	}

	//2.文件大小
	if h.Size > 5000000 {
		beego.Info("文件太大，不允许上传")
		return
	}

	//3.不能重名
	fileName := time.Now().Format("2006-01-02 15:04:05")

	err2 := this.SaveToFile("uploadname", "./static/img/"+fileName+ext)
	if err != nil {
		beego.Info("上传文件失败")
		return
	}

	if err != nil {
		beego.Info("上传文件失败", err2)
		return
	}

	//3.插入数据
	//1.获取orm对象
	o := orm.NewOrm()
	//2.创建一个插入对象
	article := models.Article{}
	//3.赋值
	article.Title = artiName
	article.Content = artiContent
	article.Img = "/static/img/" + fileName + ext


	//4.返回试图


	//给article对象复制
	//获取到下拉框传递过来的类型数据
	typeName:=this.GetString("select")
	//类型判断
	if typeName == ""{
		beego.Info("下拉匡数据错误")
		return
	}
	//获取type对象
	var artiType models.ArticleType
	artiType.TypeName = typeName
	err=o.Read(&artiType,"TypeName")
	if err != nil{
		beego.Info("获取类型错误")
		return
	}
	article.ArticleType = &artiType


	//4.插入
	_,err = o.Insert(&article)
	if err != nil{
		beego.Info("插入数据失败")
		return
	}

	this.Redirect("/Article/ShowArticle", 302)

}

//显示文章详情
func (this *ArticleController) ShowContent() {
	//1.获取Id
	id := this.GetString("id")
	beego.Info(id)
	//2.查询数据
	//1.获取orm对象
	o := orm.NewOrm()
	//2.获取查询对象
	id2, _ := strconv.Atoi(id)
	article := models.Article{Id2: id2}
	//3.查询
	err := o.Read(&article)
	if err != nil {
		beego.Info("查询数据为空")
		return
	}
	article.Count += 1
	o.Update(&article) //没有指定更新哪一列，他会自己查

	//3.传递数据给视图
	this.Data["article"] = article
	this.Layout="layout.html"
	this.TplName = "content.html"
}

//1.URLchuanzhi
//2.执行delete操作

//删除文章
func (this *ArticleController) HandleDelete() {
	id, _ := this.GetInt("id")
	//1.orm对象
	o := orm.NewOrm()

	//要有删除对象
	article := models.Article{Id2: id}

	//3.删除
	o.Delete(&article)

	this.Redirect("/Article/ShowArticle", 302)
}


//下拉框
func (this*ArticleController)ShowAddType(){
	//1.读取类型表，显示数据
	o := orm.NewOrm()
	var artiTypes[]models.ArticleType
	//查询
	_,err:=o.QueryTable("ArticleType").All(&artiTypes)
	if err != nil{
		beego.Info("查询类型错误")
	}
	this.Data["types"] = artiTypes
	this.Layout="layout.html"
	this.TplName = "addType.html"
}

//处理添加类型业务
func (this*ArticleController)HandleAddType(){
	//1.获取数据
	typename:=this.GetString("typeName")
	//2.判断数据
	if typename == ""{
		beego.Info("添加类型数据为空")
		return
	}
	//3.执行插入操作
	o := orm.NewOrm()
	var artiType models.ArticleType
	artiType.TypeName = typename
	beego.Info(artiType.TypeName)
	_,err:=o.Insert(&artiType)
	if err != nil{
		beego.Info("插入失败",err)
		return
	}
	//4.展示视图？
	this.Redirect("/Article/AddArticleType",302)
}
