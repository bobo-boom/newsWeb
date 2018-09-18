package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"be_bobo1/models"
	"time"
)

type RegController struct {
	beego.Controller
}

func (this*RegController)ShowReg(){
	this.TplName = "register.html"
}

/*
1.拿到浏览器传递的数据



2.数据处理

3.插入数据库（数据库表User）

4.返回视图

*/

//注册业务
func(this*RegController)HandleReg(){
	//1.拿到浏览器传递的数据
name := this.GetString("userName")
passwd := this.GetString("password")
	//2.数据处理
	if name == "" || passwd == ""{
		beego.Info("用户名或者密码不能为空")
		this.TplName = "register.html"
		return
	}
	//3.插入数据库
		//1.获取ORM对象
		o := orm.NewOrm()

		//2.获取插入对象
		user := models.User{}
		//3.插入操作
		user.UserName = name
		user.Passwd = passwd

		_,err := o.Insert(&user)
		if err != nil{
			beego.Info("插入数据失败")
		}
		//4.返回登陆

		//this.TplName = "login.html"
		//this.Ctx.WriteString("注册成功")
		this.Redirect("/",302)
		//

}




type LoginController struct {
	beego.Controller
}

func (this*LoginController)ShowLogin(){


	userName:=this.Ctx.GetCookie("userName")
	if userName!=""{
		o:=orm.NewOrm()

		//获取查询对象
		user:=models.User{}
		user.UserName=userName
		//查询
		o.Read(&user,"userName")

		passwd:=user.Passwd

		this.Data["userName"]=userName
		this.Data["passwd"]=passwd

	} else{
		this.Data["userName"]=""
		this.Data["passswd"]=""
	}


	this.TplName = "login.html"
}
/*
1.拿到浏览器数据

2.数据处理

3.查找数据库

4.返回视图
 */

func (this*LoginController)HandleLogin(){
	//1.那数据
	name := this.GetString("userName")
	passwd := this.GetString("password")
	//beego.Info(name,passwd)
	//2.数据处理
	if name =="" || passwd ==""{
		beego.Info("用户名和密码不能为空")
		this.TplName = "login.html"
		return
	}

	//3.查找数据
		//1.获取orm对象
		o := orm.NewOrm()

		//2.获取查询对象
		user := models.User{}
		//3.查询
		user.UserName = name
		err:=o.Read(&user,"UserName")
		if err != nil{
			beego.Info("用户名失败")
			this.TplName = "login.html"
			return
		}

		//4.判断密码是否一直

		if user.Passwd != passwd{
			beego.Info("密码失败")
			this.TplName = "login.html"
			return
		}


		//5. 登陆成功记住用户名

	check:=this.GetString("remember")
	beego.Info(check)
	if check == "on"{
		this.Ctx.SetCookie("userName",name,time.Second*3600)
	}else {
		this.Ctx.SetCookie("userName","",-1)
	}

	this.SetSession("userName",name)

	//4.返回试图
		//this.Ctx.WriteString("登陆成功")
		this.Redirect("/Article/ShowArticle",302)


}