// admin-init 仅用于 local/test 的受控初始化。
// 管理员明文密码只读取被 Git 忽略的 YAML，数据库中仅保存 bcrypt hash。
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/crmlive/qixi-live-ecrm/api-platform/internal/pkg/config"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type seedConfig struct {
	Admins []seedAdmin `yaml:"admins"`
}

type seedAdmin struct {
	DisplayName string   `yaml:"display_name"`
	Password    string   `yaml:"password"`
	Roles       []string `yaml:"roles"`
	Username    string   `yaml:"username"`
}

type adminUser struct {
	ID           uint64 `gorm:"column:id;primaryKey"`
	Username     string `gorm:"column:username"`
	PasswordHash string `gorm:"column:password_hash"`
	DisplayName  string `gorm:"column:display_name"`
	Status       int8   `gorm:"column:status"`
}

func (adminUser) TableName() string { return "qixi_crm_a_admin_user" }

type adminUserRole struct {
	AdminUserID uint64 `gorm:"column:admin_user_id"`
	RoleID      uint64 `gorm:"column:role_id"`
}

func (adminUserRole) TableName() string { return "qixi_crm_a_admin_user_role" }

type role struct {
	ID   uint64 `gorm:"column:id"`
	Code string `gorm:"column:code"`
}

func (role) TableName() string { return "qixi_crm_a_role" }

func main() {
	appConfigPath := flag.String("app-config", "", "api-platform app.yaml")
	seedConfigPath := flag.String("seed-config", "", "管理员初始化 YAML")
	flag.Parse()
	if *appConfigPath == "" || *seedConfigPath == "" {
		fail("必须指定 -app-config 与 -seed-config")
	}

	appCfg, err := config.Load(*appConfigPath)
	if err != nil {
		fail("读取应用配置失败: %v", err)
	}
	dsn, err := appCfg.DSNFor(config.DatabaseAdmin)
	if err != nil {
		fail("读取 admin 数据库失败: %v", err)
	}
	// 本命令运行在宿主机，Docker 内服务名需映射到 compose 固定的本机端口。
	dsn = strings.Replace(dsn, "pte_live_mysql:3306", "127.0.0.1:13306", 1)

	raw, err := os.ReadFile(*seedConfigPath)
	if err != nil {
		fail("读取管理员 YAML 失败: %v", err)
	}
	var seed seedConfig
	if err := yaml.Unmarshal(raw, &seed); err != nil {
		fail("解析管理员 YAML 失败: %v", err)
	}
	if len(seed.Admins) == 0 {
		fail("管理员 YAML 至少需要一个 admins 条目")
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		fail("连接 admin 数据库失败: %v", err)
	}
	created := 0
	for _, item := range seed.Admins {
		if err := initializeAdmin(db, item); err != nil {
			fail("初始化管理员失败: %v", err)
		}
		created++
	}
	fmt.Printf("已初始化 %d 个后台管理员（密码未输出）\n", created)
}

func initializeAdmin(db *gorm.DB, input seedAdmin) error {
	input.Username = strings.TrimSpace(input.Username)
	input.DisplayName = strings.TrimSpace(input.DisplayName)
	// 本地初始化兼容项目约定的 demo 账号；线上后台通过改密接口仍强制至少 8 位。
	if input.Username == "" || len(input.Password) < 6 || len(input.Roles) == 0 {
		return fmt.Errorf("账号、至少 6 位密码和至少一个角色必填")
	}
	roles := uniqueRoles(input.Roles)
	var roleRows []role
	if err := db.Where("code IN ? AND status = 1", roles).Find(&roleRows).Error; err != nil {
		return err
	}
	if len(roleRows) != len(roles) {
		return fmt.Errorf("包含不存在或禁用的后台角色")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		var account adminUser
		err := tx.Where("username = ?", input.Username).First(&account).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if err == gorm.ErrRecordNotFound {
			account = adminUser{Username: input.Username, PasswordHash: string(hash), DisplayName: input.DisplayName, Status: 1}
			if err := tx.Create(&account).Error; err != nil {
				return err
			}
		} else if err := tx.Model(&account).Updates(map[string]any{"password_hash": string(hash), "display_name": input.DisplayName, "status": 1, "auth_version": gorm.Expr("auth_version + 1")}).Error; err != nil {
			return err
		}
		if err := tx.Where("admin_user_id = ?", account.ID).Delete(&adminUserRole{}).Error; err != nil {
			return err
		}
		bindings := make([]adminUserRole, 0, len(roleRows))
		for _, row := range roleRows {
			bindings = append(bindings, adminUserRole{AdminUserID: account.ID, RoleID: row.ID})
		}
		return tx.Create(&bindings).Error
	})
}

func uniqueRoles(input []string) []string {
	seen := make(map[string]struct{}, len(input))
	result := make([]string, 0, len(input))
	for _, value := range input {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; !ok {
			seen[value] = struct{}{}
			result = append(result, value)
		}
	}
	return result
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
