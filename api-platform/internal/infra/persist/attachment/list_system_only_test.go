package attachmentpersist

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/attachment"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestListSystemOnlyScopedByEnnameAndType(t *testing.T) {
	dsn := os.Getenv("TEST_DSN")
	if dsn == "" {
		t.Skip("TEST_DSN not set")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	repo := NewRepo(db)
	ctx := context.Background()
	img := int8(0)
	vid := int8(1)

	_, totalAll, err := repo.List(ctx, 0, 0, false, &img, 1, 50)
	if err != nil {
		t.Fatal(err)
	}
	_, totalSys, err := repo.List(ctx, 0, 0, true, &img, 1, 50)
	if err != nil {
		t.Fatal(err)
	}
	if totalSys > totalAll {
		t.Fatalf("system image total %d > all %d", totalSys, totalAll)
	}

	// Insert a custom image then ensure systemOnly excludes it
	now := time.Now()
	custom := &attachment.Category{
		AttachmentCategoryName:   "临时自定义",
		AttachmentCategoryEnname: "tmp_custom_sys_test",
		Sort:                     1,
		MerID:                    0,
		IsSystem:                 0,
		CreateTime:               now,
	}
	if err := db.Create(custom).Error; err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = db.Where("attachment_category_id = ?", custom.AttachmentCategoryID).Delete(&attachment.Category{})
		_ = db.Where("attachment_name = ?", "tmp-custom-sys.png").Delete(&attachment.Attachment{})
	})
	asset := &attachment.Attachment{
		AttachmentCategoryID: custom.AttachmentCategoryID,
		AttachmentName:       "tmp-custom-sys.png",
		AttachmentSrc:        "/tmp/tmp-custom-sys.png",
		UploadType:           1,
		UserType:             0,
		UserID:               1,
		AttachmentType:       0,
		IsSystem:             0,
		CreateTime:           now,
	}
	if err := db.Create(asset).Error; err != nil {
		t.Fatal(err)
	}

	// 挂在系统分类但未标 is_system=1 的运营图，不得出现在「系统素材」
	opsInSystemCat := &attachment.Attachment{
		AttachmentCategoryID: 5106, // background_image
		AttachmentName:       "tmp-ops-bg.png",
		AttachmentSrc:        "/tmp/tmp-ops-bg.png",
		UploadType:           1,
		UserType:             0,
		UserID:               1,
		AttachmentType:       0,
		IsSystem:             0,
		CreateTime:           now,
	}
	if err := db.Create(opsInSystemCat).Error; err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = db.Where("attachment_name = ?", "tmp-ops-bg.png").Delete(&attachment.Attachment{})
	})

	_, totalAll2, err := repo.List(ctx, 0, 0, false, &img, 1, 100)
	if err != nil {
		t.Fatal(err)
	}
	sysRows, totalSys2, err := repo.List(ctx, 0, 0, true, &img, 1, 100)
	if err != nil {
		t.Fatal(err)
	}
	if totalSys2 >= totalAll2 {
		t.Fatalf("expected systemOnly < all after custom insert: sys=%d all=%d", totalSys2, totalAll2)
	}
	for _, row := range sysRows {
		if row.AttachmentCategoryID == custom.AttachmentCategoryID {
			t.Fatalf("custom category asset leaked into 系统素材")
		}
		if row.IsSystem != 1 {
			t.Fatalf("系统素材返回了 is_system!=1 的素材 id=%d", row.AttachmentID)
		}
		if row.AttachmentID == opsInSystemCat.AttachmentID {
			t.Fatalf("运营图 leaked into 系统素材")
		}
	}

	_, totalVidSys, err := repo.List(ctx, 0, 0, true, &vid, 1, 50)
	if err != nil {
		t.Fatal(err)
	}
	_ = totalVidSys
}
