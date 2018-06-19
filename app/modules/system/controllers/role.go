package controllers

import (
	"strings"
	"mm-wiki/app/models"
	"mm-wiki/app/utils"
)

type RoleController struct {
	BaseController
}

func (this *RoleController) Add() {
	this.viewLayout("role/form", "default")
}

func (this *RoleController) Save() {

	name := strings.TrimSpace(this.GetString("name", ""))

	if name == "" {
		this.jsonError("角色名称不能为空！")
	}

	ok, err := models.RoleModel.HasRoleName(name)
	if err != nil {
		this.ErrorLog("添加角色失败："+err.Error())
		this.jsonError("添加角色失败！")
	}
	if ok {
		this.jsonError("角色名已经存在！")
	}

	roleId, err := models.RoleModel.Insert(map[string]interface{}{
		"name": name,
	})

	if err != nil {
		this.ErrorLog("添加角色失败：" + err.Error())
		this.jsonError("添加角色失败")
	}
	this.InfoLog("添加角色 "+utils.Convert.IntToString(roleId, 10)+" 成功")
	this.jsonSuccess("添加角色 "+utils.Convert.IntToString(roleId, 10)+" 成功", nil, "/system/role/list")
}

func (this *RoleController) List() {

	page, _ := this.GetInt("page", 1)
	keyword := strings.TrimSpace(this.GetString("keyword", ""))

	number := 20
	limit := (page - 1) * number
	var err error
	var count int64
	var roles []map[string]string
	if keyword != "" {
		count, err = models.RoleModel.CountRolesByKeyword(keyword)
		roles, err = models.RoleModel.GetRolesByKeywordAndLimit(keyword, limit, number)
	} else {
		count, err = models.RoleModel.CountRoles()
		roles, err = models.RoleModel.GetRolesByLimit(limit, number)
	}
	if err != nil {
		this.ErrorLog("获取角色列表失败: "+err.Error())
		this.ViewError("获取角色列表失败", "/system/main/index")
	}

	this.Data["roles"] = roles
	this.Data["keyword"] = keyword
	this.SetPaginator(number, count)
	this.viewLayout("role/list", "default")
}

func (this *RoleController) Edit() {

	roleId := this.GetString("role_id", "")
	if roleId == "" {
		this.ViewError("角色不存在", "/system/role/list")
	}

	role, err := models.RoleModel.GetRoleByRoleId(roleId)
	if err != nil {
		this.ViewError("角色不存在", "/system/role/list")
	}

	this.Data["role"] = role
	this.viewLayout("role/form", "default")
}

func (this *RoleController) Modify() {

	roleId := this.GetString("role_id", "")
	name := strings.TrimSpace(this.GetString("name", ""))

	if roleId == "" {
		this.jsonError("角色不存在！")
	}
	if name == "" {
		this.jsonError("角色名称不能为空！")
	}

	role, err := models.RoleModel.GetRoleByRoleId(roleId)
	if err != nil {
		this.ErrorLog("修改角色 "+roleId+" 失败: "+err.Error())
		this.jsonError("修改角色失败！")
	}
	if len(role) == 0 {
		this.jsonError("角色不存在！")
	}

	ok , _ := models.RoleModel.HasSameName(roleId, name)
	if ok {
		this.jsonError("角色名已经存在！")
	}
	_, err = models.RoleModel.Update(roleId, map[string]interface{}{
		"name": name,
	})

	if err != nil {
		this.ErrorLog("修改角色 "+roleId+" 失败：" + err.Error())
		this.jsonError("修改角色"+roleId+"失败")
	}
	this.InfoLog("修改角色 "+roleId+" 成功")
	this.jsonSuccess("修改角色成功", nil, "/system/role/list")
}