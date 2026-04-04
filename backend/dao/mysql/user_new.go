package mysql

import (
	"bluebell/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"

	"go.uber.org/zap"
)

// 把每一步数据库操作封装成函数
// 待logic层根据业务需求调用

const secret = "liwenzhou.com"

// CheckUserExist 检查指定用户名的用户是否存在
func CheckUserExist(username string) (err error) {
	sqlStr := `select count(id) from users where username = ?`
	var count int64
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return
}

// CheckEmailExist 检查邮箱是否存在
func CheckEmailExist(email string) (err error) {
	sqlStr := `select count(id) from users where email = ?`
	var count int64
	if err := db.Get(&count, sqlStr, email); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return
}

// CreateUser 创建用户
func CreateUser(user *models.User) error {
	// 对密码进行加密
	if user.Password != "" {
		user.Password = encryptPassword(user.Password)
	}

	sqlStr := `insert into users(
		username, email, password, role, nickname, bio,
		wechat_openid, github_id, ip_address, user_agent, extra
	) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// 处理可能为空的字段
	var email, nickname, bio, wechatOpenid, githubID interface{}
	email = nil
	if user.Email.Valid && user.Email.String != "" {
		email = user.Email.String
	}
	nickname = nil
	if user.Nickname.Valid && user.Nickname.String != "" {
		nickname = user.Nickname.String
	}
	bio = nil
	if user.Bio.Valid && user.Bio.String != "" {
		bio = user.Bio.String
	}
	wechatOpenid = nil
	if user.WechatOpenid != "" {
		wechatOpenid = user.WechatOpenid
	}
	githubID = nil
	if user.GithubID != "" {
		githubID = user.GithubID
	}

	result, err := db.Exec(sqlStr,
		user.Username, email, user.Password, user.Role,
		nickname, bio, wechatOpenid, githubID,
		user.IPAddress, user.UserAgent, user.Extra,
	)
	if err != nil {
		return err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = userID

	return nil
}

// encryptPassword 密码加密
func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	//hex.EncodeToString() 返回的是将[]byte类型的结果切片转换成16进制的数据
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

// Login 用户登录
func Login(user *models.User) (err error) {
	oPassword := user.Password // 因为后面会覆盖用户登录的密码 所以这里先保存
	sqlStr := `select id, username, email, password, role, status, nickname, avatar, bio, total_words, total_likes from users where username=?`
	err = db.Get(user, sqlStr, user.Username)
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}
	if err != nil {
		// 查询数据库失败
		return err
	}
	// 判断密码是否正确
	password := encryptPassword(oPassword)
	if password != user.Password {
		return ErrorInvalidPassword
	}
	return
}

// GetUserById 根据id获取用户信息
func GetUserById(uid int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select id, username, email, role, status, nickname, avatar, bio, total_words, total_likes, created_at from users where id = ?`
	err = db.Get(user, sqlStr, uid)
	return
}

// GetUserByUsername 根据用户名获取用户信息
func GetUserByUsername(username string) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select id, username, email, role, status, nickname, avatar, bio, total_words, total_likes, created_at from users where username = ?`
	err = db.Get(user, sqlStr, username)
	return
}

// GetUserByEmail 根据邮箱获取用户信息
func GetUserByEmail(email string) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select id, username, email, role, status, nickname, avatar, bio, total_words, total_likes, created_at from users where email = ?`
	err = db.Get(user, sqlStr, email)
	return
}

// GetUserByWechatOpenid 根据微信openid获取用户信息
func GetUserByWechatOpenid(openid string) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select id, username, email, role, status, nickname, avatar, bio, total_words, total_likes, created_at from users where wechat_openid = ?`
	err = db.Get(user, sqlStr, openid)
	return
}

// GetUserByGithubID 根据GitHub ID获取用户信息
func GetUserByGithubID(githubID string) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select id, username, email, role, status, nickname, avatar, bio, total_words, total_likes, created_at from users where github_id = ?`
	err = db.Get(user, sqlStr, githubID)
	return
}

// UpdateUser 更新用户信息
func UpdateUser(id int64, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	// 如果更新密码，需要加密
	if password, ok := updates["password"]; ok {
		updates["password"] = encryptPassword(password.(string))
	}

	sqlStr := "UPDATE users SET "
	args := []interface{}{}
	i := 0
	for field, value := range updates {
		if i > 0 {
			sqlStr += ", "
		}
		sqlStr += field + " = ?"
		args = append(args, value)
		i++
	}
	sqlStr += ", updated_at = NOW() WHERE id = ?"
	args = append(args, id)

	result, err := db.Exec(sqlStr, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrorUserNotExist
	}

	return nil
}

// UpdateUserLoginInfo 更新用户登录信息
func UpdateUserLoginInfo(id int64, ip string) error {
	sqlStr := "UPDATE users SET last_login_at = NOW(), ip_address = ? WHERE id = ?"
	_, err := db.Exec(sqlStr, ip, id)
	return err
}

// GetUserList 获取用户列表（管理员用）
func GetUserList(role string, status string, page, size int) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	conditions := []interface{}{}
	sqlWhere := ""
	if role != "" && role != "all" {
		sqlWhere += " WHERE role = ?"
		conditions = append(conditions, role)
	}
	if status != "" && status != "all" {
		if sqlWhere == "" {
			sqlWhere += " WHERE status = ?"
		} else {
			sqlWhere += " AND status = ?"
		}
		conditions = append(conditions, status)
	}

	// 计算总数
	countSql := "SELECT COUNT(*) FROM users" + sqlWhere
	err := db.Get(&total, countSql, conditions...)
	if err != nil {
		return nil, 0, err
	}

	// 查询数据
	querySql := "SELECT id, username, email, role, status, nickname, avatar, bio, total_words, total_likes, created_at, last_login_at FROM users" +
		sqlWhere + " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	conditions = append(conditions, size, (page-1)*size)

	err = db.Select(&users, querySql, conditions...)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// UpdateUserStatus 更新用户状态（管理员用）
func UpdateUserStatus(id int64, status string) error {
	sqlStr := "UPDATE users SET status = ?, updated_at = NOW() WHERE id = ?"
	result, err := db.Exec(sqlStr, status, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrorUserNotExist
	}

	return nil
}

// UpdateUserStats 更新用户统计数据（总字数、总点赞数）
func UpdateUserStats(id int64, totalWords, totalLikes int64) error {
	sqlStr := "UPDATE users SET total_words = ?, total_likes = ?, updated_at = NOW() WHERE id = ?"
	_, err := db.Exec(sqlStr, totalWords, totalLikes, id)
	return err
}

// GetUserByID 根据ID获取用户信息
func GetUserByID(id int64) (*models.User, error) {
	user := new(models.User)
	sqlStr := `SELECT id, username, email, password, role, nickname, bio, avatar,
			   total_words, total_likes, status, created_at, updated_at
			   FROM users WHERE id = ?`
	err := db.Get(user, sqlStr, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUserRole 更新用户角色
func UpdateUserRole(id int64, role string) error {
	sqlStr := "UPDATE users SET role = ?, updated_at = NOW() WHERE id = ?"
	result, err := db.Exec(sqlStr, role, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrorUserNotExist
	}

	return nil
}

// BatchUpdateUserStatus 批量更新用户状态
func BatchUpdateUserStatus(userIDs []int64, status string) error {
	// 使用简单的循环更新，避免复杂的SQL构建
	for _, userID := range userIDs {
		sqlStr := "UPDATE users SET status = ?, updated_at = NOW() WHERE id = ?"
		result, err := db.Exec(sqlStr, status, userID)
		if err != nil {
			return err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			// 记录日志但不中断操作
			zap.L().Warn("user not found when batch update status", zap.Int64("user_id", userID))
		}
	}

	return nil
}