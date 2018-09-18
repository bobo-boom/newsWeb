package main

import (
	_ "be_bobo1/routers"
	"github.com/astaxie/beego"
	_"be_bobo1/models"

	"strconv"
)

func main() {
	beego.AddFuncMap("ShowPrePage",HandlePrePage)
	beego.Run()

}

func  HandlePrePage(data  int )(int){
	pageIndex:=data-1
	return   pageIndex

}

//
//func  HandlePrePage(data  int  )(string  ){
//
//
//	pageIndex:=data-1
//	pageIndex1:=strconv.Itoa(pageIndex)
//	return  pageIndex1
//}

func  HandleAftPage(data int )(string){

	pageIndex:=data+1
	pageIndex1:=strconv.Itoa(pageIndex)
	return  pageIndex1
}

