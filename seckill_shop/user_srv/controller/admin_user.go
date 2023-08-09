package controller

import (
	"context"
	"zhiliao_user_srv/proto/admin_user"
	"fmt"
	"zhiliao_user_srv/models"
	"zhiliao_user_srv/data_source"
	"zhiliao_user_srv/utils"
	"github.com/pkg/errors"
	"strconv"
)
type AdminUser struct {}

func (a *AdminUser)AdminUserLogin(ctx context.Context, in *zhiliao_user_srv.AdminUserRequest, out *zhiliao_user_srv.AdminUserResponse) error {

	user_name := in.Username
	password := in.Password

	fmt.Println("++++++++++++++")
	fmt.Println(user_name)
	fmt.Println(password)

	md5_password := utils.Md5pwd(password)

	admin_user := models.AdminUser{}

	result := data_source.Db.Where("user_name = ?",user_name).Where("password =?",md5_password).Find(&admin_user)

	if result.Error != nil {
		out.Code = 500
		out.Msg = "用户名或密码错误"
		return errors.New("用户名或密码错误")
	}

	out.Code = 200
	out.Msg = "登录成功"
	out.UserName = admin_user.UserName
	return nil


}

func (a *AdminUser)FrontUserList(ctx context.Context, in *zhiliao_user_srv.FrontUsersRequest, out *zhiliao_user_srv.FrontUsersResponse) error  {

	currentPage := in.CurrentPage
	pageSize := in.Pagesize

	/*
		current offset limit
		1       0        2       2 * (1 - 1)
		2       2        2		 2 * (2 - 1)
		3       4         2		2 * (3 -1 )

		offset = limit * (current - 1)
	 */

	 offsetNum := pageSize * (currentPage - 1)

	users := []models.FrontUser{}
	result := data_source.Db.Limit(pageSize).Offset(offsetNum).Find(&users)
	fmt.Println(len(users))


	if result.Error != nil {
		out.Code = 500
		out.Msg = "没有查询到数据"
	}

	var count int32
	users_count := []models.FrontUser{}
	data_source.Db.Find(&users_count).Count(&count)

	front_users_rep := []*zhiliao_user_srv.FrontUser{}

	for _,front_user := range users{
		front_user_rep := zhiliao_user_srv.FrontUser{}
		front_user_rep.Email = front_user.Email
		front_user_rep.Desc = front_user.Desc
		front_user_rep.Status = strconv.FormatInt(int64(front_user.Status),10)
		front_user_rep.CreateTime = front_user.CreateTime.Format("2006-01-02 15:04:05")
		front_users_rep = append(front_users_rep,&front_user_rep)
	}
	out.Code = 200
	out.Msg = "成功"
	out.FrontUsers = front_users_rep
	out.Total = count
	out.Current = currentPage
	out.PageSize = pageSize
	return nil

}