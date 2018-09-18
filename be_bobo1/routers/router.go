package routers

import (
	"be_bobo1/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
    //beego.Router("/", &controllers.MainController{})
	//beego.InsertFilter("/Article/*",beego.BeforeRouter,FilterFunc)
	beego.InsertFilter("/Article/*",beego.BeforeRouter,FilterFunc)
    beego.Router("/register",&controllers.RegController{},"get:ShowReg;post:HandleReg")
	beego.Router("/",&controllers.LoginController{},"get:ShowLogin;post:HandleLogin")
	beego.Router("/Article/ShowArticle",&controllers.ArticleController{},"get:ShowArticleList")
	beego.Router("/Article/AddArticle",&controllers.ArticleController{},"get:ShowAddArticle;post:HandleAddArtcile")
	beego.Router("/Article/ArticleContent",&controllers.ArticleController{},"get:ShowContent")
	beego.Router("/Article/DeleteArticle",&controllers.ArticleController{},"get:HandleDelete")
	beego.Router("/Article/AddArticleType",&controllers.ArticleController{},"get:ShowAddType;post:HandleAddType")


	}

//
//var FilterFunc = func(ctx *context.Context) {
//	userName := ctx.Input.Session("userName")
//	if userName == nil{
//		ctx.Redirect(302,"/")//如果有输出不再往下执行
//	}
//}

var  FilterFunc= func(ctx * context.Context) {
	userName:=ctx.Input.Session("userName")
	if  userName==nil{
		ctx.Redirect(302,"/")
	}
}