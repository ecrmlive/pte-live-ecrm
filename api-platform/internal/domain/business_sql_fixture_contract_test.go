package domain_test

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func TestAdminMenuSeedHasUniquePrimaryKeys(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot resolve current test file")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "../../.."))
	seed, err := os.ReadFile(filepath.Join(root, "sql/admin/init_data.sql"))
	if err != nil {
		t.Fatalf("read admin menu seed: %v", err)
	}
	const prefix = "INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`) VALUES"
	content := string(seed)
	start := strings.Index(content, prefix)
	if start < 0 {
		t.Fatal("admin menu seed insert not found")
	}
	section := content[start:]
	end := strings.Index(section, "ON DUPLICATE KEY UPDATE")
	if end < 0 {
		t.Fatal("admin menu seed upsert boundary not found")
	}
	ids := map[int]struct{}{}
	for _, match := range regexp.MustCompile(`(?m)^\s*\((\d+),`).FindAllStringSubmatch(section[:end], -1) {
		id, convErr := strconv.Atoi(match[1])
		if convErr != nil {
			t.Fatalf("parse menu id %q: %v", match[1], convErr)
		}
		if _, exists := ids[id]; exists {
			t.Fatalf("admin menu seed contains duplicate primary key %d", id)
		}
		ids[id] = struct{}{}
	}
	if len(ids) == 0 {
		t.Fatal("admin menu seed contains no menu rows")
	}
}

func TestBaseSeedDataDeclaresUTF8MB4(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot resolve current test file")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "../../.."))
	for _, relative := range []string{
		"sql/admin/init_data.sql",
		"sql/business/init_data.sql",
		"sql/merchant/init_data.sql",
	} {
		seed, err := os.ReadFile(filepath.Join(root, relative))
		if err != nil {
			t.Fatalf("read %s: %v", relative, err)
		}
		if !strings.Contains(string(seed), "SET NAMES utf8mb4;") {
			t.Fatalf("%s must declare SET NAMES utf8mb4 before inserting Chinese seed data", relative)
		}
	}
}

// A visible server menu must never fall back to the Vben placeholder.  This
// keeps the SQL menu seed and the concrete admin-platform component registry
// in lockstep without requiring a browser or any database credentials.
func TestAdminPageMenuHasConcreteVbenComponent(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot resolve current test file")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "../../.."))
	seed, err := os.ReadFile(filepath.Join(root, "sql/admin/init_data.sql"))
	if err != nil {
		t.Fatalf("read admin menu seed: %v", err)
	}
	registry, err := os.ReadFile(filepath.Join(root, "admin-platform/src/views/ecrm/registry.ts"))
	if err != nil {
		t.Fatalf("read Vben registry: %v", err)
	}
	paths := regexp.MustCompile(`\(\d+,\d+,'[^']+','[^']*','[^']*','([^']+)','page',\d+\)`).FindAllStringSubmatch(string(seed), -1)
	if len(paths) == 0 {
		t.Fatal("admin menu seed contains no page routes")
	}
	for _, match := range paths {
		path := match[1]
		if !strings.Contains(string(registry), "'"+path+"': 'ecrm/") {
			t.Fatalf("admin menu page %q has no concrete Vben component", path)
		}
	}
}

func TestBusinessSchemaAndChineseFixtureContract(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot resolve current test file")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "../../.."))
	schema, err := os.ReadFile(filepath.Join(root, "sql/business/init_table.sql"))
	if err != nil {
		t.Fatalf("read business schema: %v", err)
	}
	fixture, err := os.ReadFile(filepath.Join(root, "sql/business/init_test_data.sql"))
	if err != nil {
		t.Fatalf("read business fixture: %v", err)
	}

	for _, required := range []string{
		"CREATE TABLE IF NOT EXISTS `qixi_crm_b_seckill_time`",
		"CREATE TABLE IF NOT EXISTS `qixi_crm_b_seckill_active`",
		"CREATE TABLE IF NOT EXISTS `qixi_crm_b_combination_group`",
		"CREATE TABLE IF NOT EXISTS `qixi_crm_b_presell`",
	} {
		if !strings.Contains(string(schema), required) {
			t.Fatalf("business schema missing %q", required)
		}
	}
	for _, unsupported := range []string{"ADD COLUMN IF NOT EXISTS", "ADD UNIQUE INDEX IF NOT EXISTS"} {
		if strings.Contains(string(schema), unsupported) {
			t.Fatalf("business schema contains MySQL-incompatible %q", unsupported)
		}
	}
	for _, required := range []string{
		"SET NAMES utf8mb4;",
		"INSERT INTO `qixi_crm_b_seckill_active`",
		"轻奢羊绒针织衫限时秒杀",
		"INSERT INTO `qixi_crm_b_presell`",
		"晨间居家香氛套装",
	} {
		if !strings.Contains(string(fixture), required) {
			t.Fatalf("business fixture missing %q", required)
		}
	}
	for _, malformed := range []string{"((2102", "NOW()))))"} {
		if strings.Contains(string(fixture), malformed) {
			t.Fatalf("business fixture contains malformed SQL token %q", malformed)
		}
	}
	for _, required := range []string{
		"CREATE TABLE IF NOT EXISTS `qixi_crm_b_distribution_level`",
		"CREATE TABLE IF NOT EXISTS `qixi_crm_b_distribution_promoter`",
		"CREATE TABLE IF NOT EXISTS `qixi_crm_b_distribution_relation`",
		"CREATE TABLE IF NOT EXISTS `qixi_crm_b_commission_ledger`",
		"CREATE TABLE IF NOT EXISTS `qixi_crm_b_withdraw_bank`",
		"CREATE TABLE IF NOT EXISTS `qixi_crm_b_distribution_privilege`",
		"CREATE TABLE IF NOT EXISTS `qixi_crm_b_distribution_poster`",
		"INSERT INTO `qixi_crm_b_distribution_promoter`",
		"INSERT INTO `qixi_crm_b_withdraw_bank`",
		"分销夹具：推广资格由业务后台授权",
		"中国银行",
		"fixture-commission-9101-01",
	} {
		if !strings.Contains(string(schema), required) && !strings.Contains(string(fixture), required) {
			t.Fatalf("distribution schema or Chinese fixture missing %q", required)
		}
	}
	for _, required := range []string{
		"`refund_type` enum('money_only','return_and_refund')",
		"'awaiting_return','awaiting_receipt'",
		"CREATE TABLE IF NOT EXISTS `qixi_crm_b_refund_return_shipment`",
		"INSERT INTO qixi_crm_b_refund_return_shipment",
		"虚构退货物流夹具，仅用于本地闭环验收。",
		"CREATE TABLE IF NOT EXISTS `qixi_crm_b_refund_callback`",
		"INSERT INTO qixi_crm_b_payment_transaction",
		"mock-payment-CS-20260803-001",
	} {
		if !strings.Contains(string(schema), required) && !strings.Contains(string(fixture), required) {
			t.Fatalf("return-and-refund schema or Chinese fixture missing %q", required)
		}
	}
	for _, required := range []string{
		"`payout_idempotency_key` varchar(128) DEFAULT NULL",
		"`payout_reference` varchar(128) DEFAULT NULL",
		"UNIQUE KEY `uk_user_payout_key` (`user_id`,`payout_idempotency_key`)",
		"INSERT INTO `qixi_crm_b_withdrawal_application`",
		"中文演示审批通过，等待登记内部打款凭证。",
	} {
		if !strings.Contains(string(schema), required) && !strings.Contains(string(fixture), required) {
			t.Fatalf("withdrawal payout schema or Chinese fixture missing %q", required)
		}
	}
	for _, required := range []string{
		"CREATE TABLE IF NOT EXISTS `qixi_crm_b_customer_service_assignment_log`",
		"UNIQUE KEY `uk_binding_transfer_key` (`binding_id`,`idempotency_key`)",
		"客服转接审计夹具",
		"虚构演示：用户咨询居家商品，转交对应队列客服。",
	} {
		if !strings.Contains(string(schema), required) && !strings.Contains(string(fixture), required) {
			t.Fatalf("customer service transfer schema or Chinese fixture missing %q", required)
		}
	}
	for _, required := range []string{
		"CREATE TABLE IF NOT EXISTS `qixi_crm_b_product_sku_view`",
		"`merchant_sku_id` bigint unsigned NOT NULL DEFAULT 0",
		"CREATE TABLE IF NOT EXISTS `qixi_crm_b_stock_command_outbox`",
		"enum('reserve','confirm','release','restock')",
		"CREATE TABLE IF NOT EXISTS `qixi_crm_b_aftersale_item`",
		"INSERT INTO `qixi_crm_b_product_sku_view`",
		"(61001,1001,'61001'",
		"订单/库存闭环的业务侧 SKU 消费投影",
	} {
		if !strings.Contains(string(schema), required) && !strings.Contains(string(fixture), required) {
			t.Fatalf("sku inventory contract missing %q", required)
		}
	}
	for _, required := range []string{
		"CREATE TABLE IF NOT EXISTS `qixi_crm_b_settlement_command_outbox`",
		"`action` enum('accrue','reverse')",
	} {
		if !strings.Contains(string(schema), required) {
			t.Fatalf("settlement command schema missing %q", required)
		}
	}
	for _, relative := range []string{
		"api-business/internal/business/cart/handler.go",
		"api-business/internal/business/order/repository.go",
		"api-business/internal/business/order/create.go",
		"api-business/internal/business/order/stock.go",
		"api-business/internal/event/merchantstock/outbox.go",
		"api-merchant/internal/event/merchantstock/command.go",
		"api-platform/internal/platform/nativecatalog/product.go",
	} {
		content, err := os.ReadFile(filepath.Join(root, relative))
		if err != nil {
			t.Fatalf("read sku inventory source %s: %v", relative, err)
		}
		if !strings.Contains(string(content), "qixi_crm_b_product_sku_view") && !strings.Contains(string(content), "MerchantSKUID") {
			t.Fatalf("sku inventory source %s must preserve merchant SKU identity", relative)
		}
	}
	callback, err := os.ReadFile(filepath.Join(root, "api-business/internal/business/order/callback.go"))
	if err != nil {
		t.Fatalf("read refund stock callback: %v", err)
	}
	if strings.Count(string(callback), "enqueueStockRestockForRefund(tx, refund.ID)") != 2 {
		t.Fatal("both verified and sandbox refund callbacks must enqueue database-derived stock restock")
	}
	if strings.Count(string(callback), "enqueueSettlementReversalForRefund(tx, refund.ID)") != 2 {
		t.Fatal("both verified and sandbox refund callbacks must enqueue database-derived settlement reversal")
	}
	confirmReceipt, err := os.ReadFile(filepath.Join(root, "api-business/internal/business/order/handler.go"))
	if err != nil {
		t.Fatalf("read order receipt handler: %v", err)
	}
	if !strings.Contains(string(confirmReceipt), "enqueueSettlementAccrual(tx, order)") {
		t.Fatal("confirm receipt must write its settlement accrual in the order transaction")
	}
}

func TestProductCommentVirtualFixtureContract(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot resolve current test file")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "../../.."))
	schema, err := os.ReadFile(filepath.Join(root, "sql/business/init_table.sql"))
	if err != nil {
		t.Fatalf("read business schema: %v", err)
	}
	fixture, err := os.ReadFile(filepath.Join(root, "sql/business/init_test_data.sql"))
	if err != nil {
		t.Fatalf("read business fixture: %v", err)
	}
	for _, required := range []string{"qixi_crm_b_product_comment_moderation_audit", "virtual_author_name", "idx_product_visible_sort", "虚构中文虚拟评论", "演示用户小满"} {
		if !strings.Contains(string(schema), required) && !strings.Contains(string(fixture), required) {
			t.Fatalf("product comment schema or fixture missing %q", required)
		}
	}
}

func TestMerchantSettlementLedgerFixtureContract(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot resolve current test file")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "../../.."))
	schema, err := os.ReadFile(filepath.Join(root, "sql/merchant/init_table.sql"))
	if err != nil {
		t.Fatalf("read merchant schema: %v", err)
	}
	fixture, err := os.ReadFile(filepath.Join(root, "sql/merchant/init_test_data.sql"))
	if err != nil {
		t.Fatalf("read merchant fixture: %v", err)
	}
	for _, required := range []string{
		"CREATE TABLE IF NOT EXISTS `qixi_crm_m_settlement_entry`",
		"enum('order_accrual','refund_reversal')",
		"INSERT INTO `qixi_crm_m_settlement_entry`",
		"settlement:accrue:99001",
		"settlement:reverse:88001",
	} {
		if !strings.Contains(string(schema), required) && !strings.Contains(string(fixture), required) {
			t.Fatalf("merchant settlement ledger contract missing %q", required)
		}
	}
	for _, relative := range []string{
		"api-business/internal/event/merchantledger/outbox.go",
		"api-merchant/internal/event/merchantledger/command.go",
	} {
		content, err := os.ReadFile(filepath.Join(root, relative))
		if err != nil {
			t.Fatalf("read settlement ledger source %s: %v", relative, err)
		}
		if !strings.Contains(string(content), "qixi_crm_m_settlement_entry") && !strings.Contains(string(content), "qixi_crm_b_settlement_command_outbox") {
			t.Fatalf("settlement ledger source %s must use an explicit ledger or command outbox", relative)
		}
	}
	freezer, err := os.ReadFile(filepath.Join(root, "api-merchant/internal/event/merchantledger/freezer.go"))
	if err != nil {
		t.Fatalf("read settlement freezer: %v", err)
	}
	if !strings.Contains(string(freezer), "qixi_crm_m_settlement_bill") || !strings.Contains(string(freezer), "bill_frozen") {
		t.Fatal("settlement freezer must conditionally freeze merchant pending bills")
	}
}

func TestAdminMerchantSettlementProjectionFixtureContract(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot resolve current test file")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "../../.."))
	schema, err := os.ReadFile(filepath.Join(root, "sql/admin/init_table.sql"))
	if err != nil {
		t.Fatalf("read admin schema: %v", err)
	}
	fixture, err := os.ReadFile(filepath.Join(root, "sql/admin/init_test_data.sql"))
	if err != nil {
		t.Fatalf("read admin fixture: %v", err)
	}
	for _, required := range []string{
		"CREATE TABLE IF NOT EXISTS `qixi_crm_a_merchant_settlement_view`",
		"`status` enum('bill_pending','bill_frozen','withdraw_applied','approved','paid','rejected','cancelled')",
		"INSERT INTO `qixi_crm_a_merchant_settlement_view`",
		"七禧演示店铺",
		"1280.50,'withdraw_applied'",
	} {
		if !strings.Contains(string(schema), required) && !strings.Contains(string(fixture), required) {
			t.Fatalf("admin settlement projection contract missing %q", required)
		}
	}
}

func TestMerchantStatusAndIntentionMutationsHaveRBACContract(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot resolve current test file")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "../../.."))
	menuSeed, err := os.ReadFile(filepath.Join(root, "sql/admin/init_data.sql"))
	if err != nil {
		t.Fatalf("read admin menu seed: %v", err)
	}
	for _, required := range []string{
		"merchant.status.manage",
		"merchant.intention.audit",
		"merchant.intention.assign_region",
		"merchant.intention.delete",
	} {
		if !strings.Contains(string(menuSeed), required) {
			t.Fatalf("merchant mutation menu contract missing %q", required)
		}
	}
	handler, err := os.ReadFile(filepath.Join(root, "api-platform/internal/platform/merchant/handler.go"))
	if err != nil {
		t.Fatalf("read merchant handler: %v", err)
	}
	for _, required := range []string{
		"RequireAdminMenu(h.adminDB, \"merchant.status.manage\")",
		"RequireAdminMenu(h.adminDB, \"merchant.intention.audit\")",
		"RequireAdminMenu(h.adminDB, \"merchant.intention.assign_region\")",
		"RequireAdminMenu(h.adminDB, \"merchant.intention.delete\")",
	} {
		if !strings.Contains(string(handler), required) {
			t.Fatalf("merchant mutation handler missing RBAC guard %q", required)
		}
	}
}

func TestMerchantCategoryHasIndependentTablePageAndRBAC(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot resolve current test file")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "../../.."))
	adminSchema, err := os.ReadFile(filepath.Join(root, "sql/admin/init_table.sql"))
	if err != nil {
		t.Fatalf("read admin schema: %v", err)
	}
	merchantRepo, err := os.ReadFile(filepath.Join(root, "api-platform/internal/infra/persist/merchant/repo.go"))
	if err != nil {
		t.Fatalf("read merchant repo: %v", err)
	}
	menuSeed, err := os.ReadFile(filepath.Join(root, "sql/admin/init_data.sql"))
	if err != nil {
		t.Fatalf("read admin menu seed: %v", err)
	}
	registry, err := os.ReadFile(filepath.Join(root, "admin-platform/src/views/ecrm/registry.ts"))
	if err != nil {
		t.Fatalf("read ecrm registry: %v", err)
	}
	for _, required := range []string{
		"CREATE TABLE IF NOT EXISTS `qixi_crm_a_merchant_category`",
		"`commission_rate` decimal(5,2)",
	} {
		if !strings.Contains(string(adminSchema), required) {
			t.Fatalf("merchant category schema missing %q", required)
		}
	}
	for _, required := range []string{
		"qixi_crm_a_merchant_category",
		"merchant.category.manage",
		"'/merchant/categories': 'ecrm/merchant/categories'",
	} {
		if !strings.Contains(string(merchantRepo), required) && !strings.Contains(string(menuSeed), required) && !strings.Contains(string(registry), required) {
			t.Fatalf("merchant category production contract missing %q", required)
		}
	}
}

func TestStoreGroupHasTreeMembershipTemplateAndRBACContract(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot resolve current test file")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "../../.."))
	adminSchema, err := os.ReadFile(filepath.Join(root, "sql/admin/init_table.sql"))
	if err != nil {
		t.Fatalf("read admin schema: %v", err)
	}
	fixture, err := os.ReadFile(filepath.Join(root, "sql/admin/init_test_data.sql"))
	if err != nil {
		t.Fatalf("read admin fixture: %v", err)
	}
	menuSeed, err := os.ReadFile(filepath.Join(root, "sql/admin/init_data.sql"))
	if err != nil {
		t.Fatalf("read admin menu seed: %v", err)
	}
	handler, err := os.ReadFile(filepath.Join(root, "api-platform/internal/platform/storegroup/handler.go"))
	if err != nil {
		t.Fatalf("read store group handler: %v", err)
	}
	registry, err := os.ReadFile(filepath.Join(root, "admin-platform/src/views/ecrm/registry.ts"))
	if err != nil {
		t.Fatalf("read ecrm registry: %v", err)
	}
	for _, required := range []string{
		"CREATE TABLE IF NOT EXISTS `qixi_crm_a_store_group`",
		"CREATE TABLE IF NOT EXISTS `qixi_crm_a_store_group_merchant`",
		"`path` varchar(500)",
		"`diy_page_id` bigint unsigned",
		"七禧中文演示商圈",
		"merchant.group.manage",
		"/store-groups/:id/template",
		"errGroupCycle",
		"errGroupTooDeep",
		"errMerchantAbsent",
		"'/merchant/grouping': 'ecrm/merchant/grouping'",
	} {
		joined := string(adminSchema) + string(fixture) + string(menuSeed) + string(handler) + string(registry)
		if !strings.Contains(joined, required) {
			t.Fatalf("store group production contract missing %q", required)
		}
	}
}

func TestAdminPlatformDoesNotRegisterPlaceholderAsBusinessPage(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot resolve current test file")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "../../.."))
	menuAdapter, err := os.ReadFile(filepath.Join(root, "admin-platform/src/utils/ecrm-menu.ts"))
	if err != nil {
		t.Fatalf("read admin menu adapter: %v", err)
	}
	if strings.Contains(string(menuAdapter), "ecrm/placeholder/index") {
		t.Fatal("business menu adapter must not map unknown server pages to placeholder")
	}
	for _, required := range []string{
		".map(mapNode)",
		"未注册真实组件的叶子不能降级到",
	} {
		if !strings.Contains(string(menuAdapter), required) {
			t.Fatalf("business menu adapter missing production route guard %q", required)
		}
	}
	registry, err := os.ReadFile(filepath.Join(root, "admin-platform/src/views/ecrm/registry.ts"))
	if err != nil {
		t.Fatalf("read admin page registry: %v", err)
	}
	seed, err := os.ReadFile(filepath.Join(root, "sql/admin/init_data.sql"))
	if err != nil {
		t.Fatalf("read admin menu seed: %v", err)
	}
	pageRoutes := regexp.MustCompile(`'(/[^']*)','page'`).FindAllStringSubmatch(string(seed), -1)
	if len(pageRoutes) == 0 {
		t.Fatal("admin menu seed has no page routes")
	}
	for _, match := range pageRoutes {
		if !strings.Contains(string(registry), "'"+match[1]+"':") {
			t.Fatalf("seeded page %q has no real Vben component registration", match[1])
		}
	}
}

func TestAssistDomainUsesBusinessPrefixAndFixture(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot resolve current test file")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "../../.."))
	for _, relative := range []string{
		"api-platform/internal/domain/assist/model.go",
		"api-platform/internal/infra/persist/assist/repo.go",
	} {
		content, err := os.ReadFile(filepath.Join(root, relative))
		if err != nil {
			t.Fatalf("read assist source %s: %v", relative, err)
		}
		if strings.Contains(string(content), "qixi_m_") {
			t.Fatalf("assist source %s must not access legacy qixi_m_ tables", relative)
		}
	}
	mainSource, err := os.ReadFile(filepath.Join(root, "api-platform/cmd/main.go"))
	if err != nil {
		t.Fatalf("read platform main: %v", err)
	}
	if !strings.Contains(string(mainSource), "assistpersist.NewRepo(businessDB)") {
		t.Fatal("assist repository must be connected to the business database")
	}
	fixture, err := os.ReadFile(filepath.Join(root, "sql/business/init_test_data.sql"))
	if err != nil {
		t.Fatalf("read business fixture: %v", err)
	}
	if !strings.Contains(string(fixture), "INSERT INTO `qixi_crm_b_assist`") || !strings.Contains(string(fixture), "好友助力满员后可按助力价下单。") {
		t.Fatal("business fixture must contain the Chinese assist acceptance data")
	}
}

func TestBroadcastDomainUsesBusinessPrefixAndDoesNotSerializeSensitiveFields(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot resolve current test file")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "../../.."))
	for _, relative := range []string{
		"api-platform/internal/domain/broadcast/model.go",
		"api-platform/internal/infra/persist/broadcast/repo.go",
	} {
		content, err := os.ReadFile(filepath.Join(root, relative))
		if err != nil {
			t.Fatalf("read broadcast source %s: %v", relative, err)
		}
		if strings.Contains(string(content), "qixi_m_") {
			t.Fatalf("broadcast source %s must not access legacy qixi_m_ tables", relative)
		}
	}
	model, err := os.ReadFile(filepath.Join(root, "api-platform/internal/domain/broadcast/model.go"))
	if err != nil {
		t.Fatalf("read broadcast model: %v", err)
	}
	for _, required := range []string{
		"gorm:\"column:push_url\" json:\"-\"",
		"gorm:\"column:phone\" json:\"-\"",
		"qixi_crm_b_broadcast_room",
	} {
		if !strings.Contains(string(model), required) {
			t.Fatalf("broadcast model missing sensitive-field or business-table guard %q", required)
		}
	}
	mainSource, err := os.ReadFile(filepath.Join(root, "api-platform/cmd/main.go"))
	if err != nil {
		t.Fatalf("read platform main: %v", err)
	}
	if !strings.Contains(string(mainSource), "broadcastpersist.NewRepo(businessDB)") {
		t.Fatal("broadcast repository must be connected to the business database")
	}
	schema, err := os.ReadFile(filepath.Join(root, "sql/business/init_table.sql"))
	if err != nil {
		t.Fatalf("read business schema: %v", err)
	}
	fixture, err := os.ReadFile(filepath.Join(root, "sql/business/init_test_data.sql"))
	if err != nil {
		t.Fatalf("read business fixture: %v", err)
	}
	for _, required := range []string{
		"CREATE TABLE IF NOT EXISTS `qixi_crm_b_broadcast_room`",
		"CREATE TABLE IF NOT EXISTS `qixi_crm_b_broadcast_room_goods`",
		"INSERT INTO `qixi_crm_b_broadcast_room`",
		"CRM Live服饰秋日穿搭直播间",
		"推流地址、主播手机号始终为空",
	} {
		if !strings.Contains(string(schema), required) && !strings.Contains(string(fixture), required) {
			t.Fatalf("broadcast business schema or Chinese fixture missing %q", required)
		}
	}
}

func TestCommunityDomainUsesBusinessTablesAndChineseFixture(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot resolve current test file")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "../../.."))
	for _, relative := range []string{
		"api-platform/internal/domain/community/model.go",
		"api-platform/internal/infra/persist/community/repo.go",
	} {
		content, err := os.ReadFile(filepath.Join(root, relative))
		if err != nil {
			t.Fatalf("read community source %s: %v", relative, err)
		}
		if strings.Contains(string(content), "qixi_m_") {
			t.Fatalf("community source %s must not access legacy qixi_m_ tables", relative)
		}
	}
	mainSource, err := os.ReadFile(filepath.Join(root, "api-platform/cmd/main.go"))
	if err != nil {
		t.Fatalf("read platform main: %v", err)
	}
	if !strings.Contains(string(mainSource), "communitypersist.NewRepo(businessDB)") {
		t.Fatal("community repository must be connected to the business database")
	}
	schema, err := os.ReadFile(filepath.Join(root, "sql/business/init_table.sql"))
	if err != nil {
		t.Fatalf("read business schema: %v", err)
	}
	fixture, err := os.ReadFile(filepath.Join(root, "sql/business/init_test_data.sql"))
	if err != nil {
		t.Fatalf("read business fixture: %v", err)
	}
	for _, required := range []string{
		"CREATE TABLE IF NOT EXISTS `qixi_crm_b_social_category`",
		"CREATE TABLE IF NOT EXISTS `qixi_crm_b_social_topic`",
		"CREATE TABLE IF NOT EXISTS `qixi_crm_b_social_post`",
		"CREATE TABLE IF NOT EXISTS `qixi_crm_b_social_reply`",
		"INSERT INTO `qixi_crm_b_social_post`",
		"通勤针织衫的三种叠穿思路",
	} {
		if !strings.Contains(string(schema), required) && !strings.Contains(string(fixture), required) {
			t.Fatalf("community business schema or Chinese fixture missing %q", required)
		}
	}
}

func TestUserTagDomainUsesBusinessTablesChineseFixtureAndRBAC(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot resolve current test file")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "../../.."))
	for _, relative := range []string{
		"api-platform/internal/domain/usertag/model.go",
		"api-platform/internal/infra/persist/usertag/repo.go",
	} {
		content, err := os.ReadFile(filepath.Join(root, relative))
		if err != nil {
			t.Fatalf("read usertag source %s: %v", relative, err)
		}
		if strings.Contains(string(content), "qixi_m_") {
			t.Fatalf("usertag source %s must not access legacy qixi_m_ tables", relative)
		}
	}
	mainSource, err := os.ReadFile(filepath.Join(root, "api-platform/cmd/main.go"))
	if err != nil {
		t.Fatalf("read platform main: %v", err)
	}
	if !strings.Contains(string(mainSource), "usertagpersist.NewRepo(businessDB)") {
		t.Fatal("user tag repository must be connected to the business database")
	}
	handler, err := os.ReadFile(filepath.Join(root, "api-platform/internal/platform/usertag/handler.go"))
	if err != nil {
		t.Fatalf("read usertag handler: %v", err)
	}
	for _, required := range []string{
		"labelManage := middleware.RequirePlatformMenu",
		"groupManage := middleware.RequirePlatformMenu",
		"r.GET(\"/user/labels\", labelManage, h.ListLabels)",
		"r.GET(\"/user/groups\", groupManage, h.ListGroups)",
	} {
		if !strings.Contains(string(handler), required) {
			t.Fatalf("usertag handler missing RBAC guard %q", required)
		}
	}
	schema, err := os.ReadFile(filepath.Join(root, "sql/business/init_table.sql"))
	if err != nil {
		t.Fatalf("read business schema: %v", err)
	}
	fixture, err := os.ReadFile(filepath.Join(root, "sql/business/init_test_data.sql"))
	if err != nil {
		t.Fatalf("read business fixture: %v", err)
	}
	for _, required := range []string{
		"CREATE TABLE IF NOT EXISTS `qixi_crm_b_user_label`",
		"CREATE TABLE IF NOT EXISTS `qixi_crm_b_user_group`",
		"CREATE TABLE IF NOT EXISTS `qixi_crm_b_user_label_relation`",
		"INSERT INTO `qixi_crm_b_user_label`",
		"高频复购用户",
		"CRM Live精选会员",
	} {
		if !strings.Contains(string(schema), required) && !strings.Contains(string(fixture), required) {
			t.Fatalf("usertag business schema or Chinese fixture missing %q", required)
		}
	}
}

func TestBusinessZoneUsesUnifiedAdminTablesRBACAndDoesNotSerializePaymentData(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot resolve current test file")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "../../.."))
	for _, relative := range []string{
		"api-platform/internal/domain/circle/model.go",
		"api-platform/internal/infra/persist/circle/repo.go",
	} {
		content, err := os.ReadFile(filepath.Join(root, relative))
		if err != nil {
			t.Fatalf("read business-zone source %s: %v", relative, err)
		}
		if strings.Contains(string(content), "qixi_m_") {
			t.Fatalf("business-zone source %s must not access legacy qixi_m_ tables", relative)
		}
	}
	model, err := os.ReadFile(filepath.Join(root, "api-platform/internal/domain/circle/model.go"))
	if err != nil {
		t.Fatalf("read business-zone model: %v", err)
	}
	for _, required := range []string{
		"qixi_crm_a_business_zone",
		"qixi_crm_a_business_zone_agent",
		"gorm:\"column:payment_account\" json:\"-\"",
		"gorm:\"column:payment_bank\" json:\"-\"",
		"gorm:\"column:payment_qr_img\" json:\"-\"",
	} {
		if !strings.Contains(string(model), required) {
			t.Fatalf("business-zone model missing unified-table or payment guard %q", required)
		}
	}
	mainSource, err := os.ReadFile(filepath.Join(root, "api-platform/cmd/main.go"))
	if err != nil {
		t.Fatalf("read platform main: %v", err)
	}
	if !strings.Contains(string(mainSource), "platformcircle.NewHandler(circleSvc, gdb, businessDB)") {
		t.Fatal("business-zone handler must use unified admin database RBAC")
	}
	handler, err := os.ReadFile(filepath.Join(root, "api-platform/internal/platform/circle/handler.go"))
	if err != nil {
		t.Fatalf("read business-zone handler: %v", err)
	}
	for _, required := range []string{
		"middleware.RequireAdminRoles(\"platform\")",
		"region.zone.manage",
		"region.agent.manage",
		"region.agent.review",
	} {
		if !strings.Contains(string(handler), required) {
			t.Fatalf("business-zone handler missing RBAC guard %q", required)
		}
	}
	schema, err := os.ReadFile(filepath.Join(root, "sql/admin/init_table.sql"))
	if err != nil {
		t.Fatalf("read admin schema: %v", err)
	}
	seed, err := os.ReadFile(filepath.Join(root, "sql/admin/init_data.sql"))
	if err != nil {
		t.Fatalf("read admin menu seed: %v", err)
	}
	fixture, err := os.ReadFile(filepath.Join(root, "sql/admin/init_test_data.sql"))
	if err != nil {
		t.Fatalf("read admin fixture: %v", err)
	}
	for _, required := range []string{
		"CREATE TABLE IF NOT EXISTS `qixi_crm_a_business_zone`",
		"CREATE TABLE IF NOT EXISTS `qixi_crm_a_business_zone_agent`",
		"region.zone.manage",
		"region.agent.manage",
		"region.agent.review",
		"INSERT INTO `qixi_crm_a_business_zone`",
		"华东中文演示区域",
		"结算账号刻意为空",
	} {
		if !strings.Contains(string(schema), required) && !strings.Contains(string(seed), required) && !strings.Contains(string(fixture), required) {
			t.Fatalf("business-zone schema, RBAC seed, or Chinese fixture missing %q", required)
		}
	}
}

func TestLogisticsUsesUnifiedAdminTablesRBACAndChineseFixture(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot resolve current test file")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "../../.."))
	model, err := os.ReadFile(filepath.Join(root, "api-platform/internal/domain/logistics/model.go"))
	if err != nil {
		t.Fatalf("read logistics model: %v", err)
	}
	if strings.Contains(string(model), "qixi_m_") {
		t.Fatal("logistics model must not access legacy qixi_m_ tables")
	}
	for _, required := range []string{"qixi_crm_a_express", "qixi_crm_a_city"} {
		if !strings.Contains(string(model), required) {
			t.Fatalf("logistics model missing unified admin table %q", required)
		}
	}
	mainSource, err := os.ReadFile(filepath.Join(root, "api-platform/cmd/main.go"))
	if err != nil {
		t.Fatalf("read platform main: %v", err)
	}
	if !strings.Contains(string(mainSource), "platformlogistics.NewHandler(logisticsSvc, gdb)") {
		t.Fatal("logistics handler must use unified admin database RBAC")
	}
	handler, err := os.ReadFile(filepath.Join(root, "api-platform/internal/platform/logistics/handler.go"))
	if err != nil {
		t.Fatalf("read logistics handler: %v", err)
	}
	for _, required := range []string{"middleware.RequireAdminRoles(\"platform\")", "freight.express.manage"} {
		if !strings.Contains(string(handler), required) {
			t.Fatalf("logistics handler missing RBAC guard %q", required)
		}
	}
	schema, err := os.ReadFile(filepath.Join(root, "sql/admin/init_table.sql"))
	if err != nil {
		t.Fatalf("read admin schema: %v", err)
	}
	seed, err := os.ReadFile(filepath.Join(root, "sql/admin/init_data.sql"))
	if err != nil {
		t.Fatalf("read admin menu seed: %v", err)
	}
	fixture, err := os.ReadFile(filepath.Join(root, "sql/admin/init_test_data.sql"))
	if err != nil {
		t.Fatalf("read admin fixture: %v", err)
	}
	for _, required := range []string{
		"CREATE TABLE IF NOT EXISTS `qixi_crm_a_express`",
		"CREATE TABLE IF NOT EXISTS `qixi_crm_a_city`",
		"freight.express.manage",
		"INSERT INTO `qixi_crm_a_express`",
		"七禧演示快递",
		"上海市区",
		"黄浦区",
	} {
		if !strings.Contains(string(schema), required) && !strings.Contains(string(seed), required) && !strings.Contains(string(fixture), required) {
			t.Fatalf("logistics schema, RBAC seed, or Chinese fixture missing %q", required)
		}
	}
}

func TestArticleDomainUsesUnifiedAdminTablesAndChineseFixture(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot resolve current test file")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "../../.."))
	model, err := os.ReadFile(filepath.Join(root, "api-platform/internal/domain/article/model.go"))
	if err != nil {
		t.Fatalf("read article model: %v", err)
	}
	if strings.Contains(string(model), "qixi_m_") {
		t.Fatal("article model must not use legacy qixi_m_ tables")
	}
	for _, required := range []string{"qixi_crm_a_article_category", "qixi_crm_a_article"} {
		if !strings.Contains(string(model), required) {
			t.Fatalf("article model missing unified admin table %q", required)
		}
	}
	schema, err := os.ReadFile(filepath.Join(root, "sql/admin/init_table.sql"))
	if err != nil {
		t.Fatalf("read admin schema: %v", err)
	}
	fixture, err := os.ReadFile(filepath.Join(root, "sql/admin/init_test_data.sql"))
	if err != nil {
		t.Fatalf("read admin fixture: %v", err)
	}
	for _, required := range []string{
		"CREATE TABLE IF NOT EXISTS `qixi_crm_a_article_category`",
		"CREATE TABLE IF NOT EXISTS `qixi_crm_a_article`",
		"INSERT INTO `qixi_crm_a_article_category`",
		"七禧商城秋季服务公告",
		"居家香氛选购小贴士",
	} {
		if !strings.Contains(string(schema), required) && !strings.Contains(string(fixture), required) {
			t.Fatalf("article admin schema or Chinese fixture missing %q", required)
		}
	}
}

func TestContentDomainUsesUnifiedAdminTablesAndChineseFixture(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot resolve current test file")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "../../.."))
	for _, relative := range []string{
		"api-platform/internal/domain/content/model.go",
		"api-platform/internal/domain/content/service.go",
	} {
		content, err := os.ReadFile(filepath.Join(root, relative))
		if err != nil {
			t.Fatalf("read content source %s: %v", relative, err)
		}
		if strings.Contains(string(content), "qixi_m_") {
			t.Fatalf("content source %s must not use legacy qixi_m_ tables", relative)
		}
	}
	schema, err := os.ReadFile(filepath.Join(root, "sql/admin/init_table.sql"))
	if err != nil {
		t.Fatalf("read admin schema: %v", err)
	}
	fixture, err := os.ReadFile(filepath.Join(root, "sql/admin/init_test_data.sql"))
	if err != nil {
		t.Fatalf("read admin fixture: %v", err)
	}
	for _, required := range []string{
		"CREATE TABLE IF NOT EXISTS `qixi_crm_a_notice`",
		"CREATE TABLE IF NOT EXISTS `qixi_crm_a_setting_cache`",
		"CREATE TABLE IF NOT EXISTS `qixi_crm_a_marketing_decor`",
		"CREATE TABLE IF NOT EXISTS `qixi_crm_a_config_item`",
		"INSERT INTO `qixi_crm_a_notice`",
		"INSERT INTO `qixi_crm_a_marketing_decor`",
		"INSERT INTO `qixi_crm_a_config_item`",
		"七禧商城本地验收公告",
		"夏日焕新氛围图",
		"夏日香氛",
		"本地验收未配置通道",
	} {
		if !strings.Contains(string(schema), required) && !strings.Contains(string(fixture), required) {
			t.Fatalf("content admin schema or Chinese fixture missing %q", required)
		}
	}
}

func TestAttachmentDomainUsesUnifiedAdminTablesAndChineseFixture(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot resolve current test file")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "../../.."))
	model, err := os.ReadFile(filepath.Join(root, "api-platform/internal/domain/attachment/model.go"))
	if err != nil {
		t.Fatalf("read attachment model: %v", err)
	}
	if strings.Contains(string(model), "qixi_m_") {
		t.Fatal("attachment model must not use legacy qixi_m_ tables")
	}
	for _, required := range []string{"qixi_crm_a_attachment_category", "qixi_crm_a_attachment_asset"} {
		if !strings.Contains(string(model), required) {
			t.Fatalf("attachment model missing unified admin table %q", required)
		}
	}
	schema, err := os.ReadFile(filepath.Join(root, "sql/admin/init_table.sql"))
	if err != nil {
		t.Fatalf("read admin schema: %v", err)
	}
	fixture, err := os.ReadFile(filepath.Join(root, "sql/admin/init_test_data.sql"))
	if err != nil {
		t.Fatalf("read admin fixture: %v", err)
	}
	for _, required := range []string{
		"CREATE TABLE IF NOT EXISTS `qixi_crm_a_attachment_category`",
		"CREATE TABLE IF NOT EXISTS `qixi_crm_a_attachment_asset`",
		"`is_system`",
		"store_cover",
		"other_image",
		"store_video",
		"product_video",
		"other_video",
	} {
		if !strings.Contains(string(schema), required) && !strings.Contains(string(fixture), required) && !strings.Contains(string(model), required) {
			dataSQL, _ := os.ReadFile(filepath.Join(root, "sql/admin/init_data.sql"))
			if !strings.Contains(string(dataSQL), required) {
				t.Fatalf("attachment admin schema/fixture/data missing %q", required)
			}
		}
	}
}
