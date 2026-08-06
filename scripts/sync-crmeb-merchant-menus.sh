#!/usr/bin/env bash
# 从 CRMEB 演示站拉「店铺菜单」全量并生成 sql/merchant/init_menu_crmeb_full.sql
# 用法：./scripts/sync-crmeb-merchant-menus.sh
# 依赖：curl、python3；演示账号见脚本内（平台后台，非商户）。
set -euo pipefail
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUT_JSON="${ROOT_DIR}/sql/merchant/crmeb_live_merchant_menu_lst.json"
OUT_SQL="${ROOT_DIR}/sql/merchant/init_menu_crmeb_full.sql"
REG_TS="${ROOT_DIR}/admin-merchant/src/views/ecrm/registry.ts"

ACCOUNT="${CRMEB_ADMIN_ACCOUNT:-18812341234}"
PASSWORD="${CRMEB_ADMIN_PASSWORD:-000000}"

TOKEN="$(curl -sS -m 20 -X POST 'https://mer.crmeb.net/sys/login' \
	-H 'Content-Type: application/json' \
	-d "{\"account\":\"${ACCOUNT}\",\"password\":\"${PASSWORD}\"}" \
	| python3 -c 'import sys,json; d=json.load(sys.stdin); assert d.get("status")==200, d; print(d["data"]["token"])')"

curl -sS -m 60 -H "X-Token: ${TOKEN}" -H 'Accept: application/json' \
	'https://mer.crmeb.net/sys/merchant/menu/lst' -o "${OUT_JSON}"

python3 - "${OUT_JSON}" "${OUT_SQL}" "${REG_TS}" <<'PY'
import json, re, sys
from pathlib import Path

src, out_sql, reg_ts = map(Path, sys.argv[1:4])
raw = json.loads(src.read_text(encoding="utf-8"))
assert raw.get("status") == 200, raw.get("message")

def flatten(nodes):
    out = []
    for n in nodes:
        kids = n.get("children") or []
        out.append({k: v for k, v in n.items() if k != "children"})
        out.extend(flatten(kids))
    return out

rows = flatten(raw["data"]["list"])
path_comp = dict(re.findall(r"'(/[^']*)'\s*:\s*'([^']+)'", reg_ts.read_text(encoding="utf-8")))
ICON_MAP = {
    "house": "lucide:home",
    "s-goods": "lucide:package",
    "goods": "lucide:package",
    "pie-chart": "lucide:pie-chart",
    "setting": "lucide:settings",
    "user": "lucide:users",
    "notebook-2": "lucide:notebook",
    "headset": "lucide:headset",
    "brush": "lucide:paintbrush",
    "shopping-bag-1": "lucide:shopping-bag",
    "s-order": "lucide:receipt",
    "s-marketing": "lucide:megaphone",
    "s-finance": "lucide:wallet",
}
ids = {int(r["menu_id"]) for r in rows}

def esc(s: str) -> str:
    return (s or "").replace("\\", "\\\\").replace("'", "\\'")

def route_to_code(route: str, menu_id: int, is_menu_crmeb: int) -> str:
    r = (route or "").strip()
    if not r or r == "self":
        return f"meru.self.{menu_id}" if r == "self" else f"meru.id.{menu_id}"
    if int(is_menu_crmeb) == 0 or (r and not r.startswith("/") and "/" not in r):
        code = re.sub(r"[^a-zA-Z0-9_.-]+", "_", r)
        return (code or f"meru.perm.{menu_id}")[:120]
    code = re.sub(r"[^a-zA-Z0-9_.-]+", "_", r.strip("/").replace("/", "."))
    return (code or f"meru.path.{menu_id}")[:120]

used = set()
out = []
for r in sorted(rows, key=lambda x: int(x["menu_id"])):
    mid = int(r["menu_id"])
    pid = int(r.get("pid") or 0)
    if pid and pid not in ids:
        pid = 0
    crmeb_is_menu = int(r.get("is_menu") or 0)
    is_menu = 1 if crmeb_is_menu == 1 else 2
    route = (r.get("route") or "").strip()
    path = (route if route and route != "self" else "") if is_menu == 1 else ""
    code = route_to_code(route, mid, crmeb_is_menu)
    base, n = code, 2
    while code.lower() in used:
        code = f"{base}.{mid}" if n == 2 else f"{base}.{mid}.{n}"
        n += 1
    used.add(code.lower())
    icon = r.get("icon") or ""
    if icon and ":" not in icon:
        icon = ICON_MAP.get(icon, f"lucide:{icon}")
    component = ""
    if is_menu == 1 and path.startswith("/"):
        p = path[len("/merchant") :] if path.startswith("/merchant/") else path
        if path == "/merchant":
            p = "/dashboard"
        elif not p.startswith("/"):
            p = "/" + p
        comp = path_comp.get(p) or path_comp.get(p.rstrip("/"))
        if comp:
            component = comp if comp.endswith(".vue") else f"views/{comp}.vue"
            if not component.startswith("views/"):
                component = f"views/{comp}.vue"
    out.append(
        {
            "id": mid,
            "parent_id": pid,
            "code": code,
            "name": r.get("menu_name") or "",
            "path": path[:255],
            "component": component[:255],
            "icon": icon[:128],
            "is_menu": is_menu,
            "is_route": 1 if is_menu == 1 and path.startswith("/") else 0,
            "sort": int(r.get("sort") or 0),
            "status": 1 if int(r.get("is_show") or 0) == 1 else 0,
            "created_at": r.get("create_time") or "2020-01-01 00:00:00",
        }
    )

lines = [
    "SET NAMES utf8mb4;",
    "USE `qixi_crm_merchant`;",
    "-- CRMEB live GET /sys/merchant/menu/lst → qixi_crm_m_menu",
    f"-- total {len(out)} (pages {sum(1 for x in out if x['is_menu']==1)} + buttons {sum(1 for x in out if x['is_menu']==2)})",
    "SET FOREIGN_KEY_CHECKS=0;",
    "DELETE FROM `qixi_crm_m_role_menu`;",
    "DELETE FROM `qixi_crm_m_menu`;",
    "SET FOREIGN_KEY_CHECKS=1;",
    "",
]
BATCH = 40
for i in range(0, len(out), BATCH):
    chunk = out[i : i + BATCH]
    lines.append(
        "INSERT INTO `qixi_crm_m_menu` (`id`,`parent_id`,`code`,`name`,`path`,`component`,`icon`,`is_menu`,`is_route`,`sort`,`status`,`created_at`) VALUES"
    )
    vals = [
        f"({x['id']},{x['parent_id']},'{esc(x['code'])}','{esc(x['name'])}','{esc(x['path'])}','{esc(x['component'])}','{esc(x['icon'])}',{x['is_menu']},{x['is_route']},{x['sort']},{x['status']},'{esc(x['created_at'])}')"
        for x in chunk
    ]
    lines.append(",\n".join(vals) + ";")
    lines.append("")
lines += [
    "INSERT IGNORE INTO `qixi_crm_m_role_menu` (`role_code`,`menu_id`)",
    "SELECT roles.role_code, menus.id FROM (SELECT 'owner' AS role_code UNION ALL SELECT 'manager') AS roles CROSS JOIN `qixi_crm_m_menu` AS menus;",
    "",
]
out_sql.write_text("\n".join(lines), encoding="utf-8")
print(f"wrote {out_sql} rows={len(out)}")
PY

echo "完成：${OUT_JSON} + ${OUT_SQL}"
