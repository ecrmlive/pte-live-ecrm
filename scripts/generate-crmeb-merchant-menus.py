#!/usr/bin/env python3
"""Generate sql/merchant/init_menu_crmeb_full.sql from CRMEB MER v4.0 install dump.

Source (read-only, outside repo):
  ~/Downloads/CRMEB多商户系统/CRMEB_MER_v4.0/install/crmeb_merchant.sql
  table eb_system_menu where is_mer=1
"""

from __future__ import annotations

import json
import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
DEFAULT_CRMEB = Path.home() / "Downloads/CRMEB多商户系统/CRMEB_MER_v4.0/install/crmeb_merchant.sql"
OUT = ROOT / "sql/merchant/init_menu_crmeb_full.sql"
REGISTRY = ROOT / "admin-merchant/src/views/ecrm/registry.ts"
SUMMARY = Path("/tmp/crmeb_mer_menu_summary.json")

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


def unq(s: str) -> str:
    s = s.strip()
    if s == "NULL":
        return ""
    if s.startswith("'") and s.endswith("'"):
        s = s[1:-1]
    return s.replace("\\'", "'").replace('\\"', '"')


def route_to_code(route: str, menu_id: int, is_menu: int) -> str:
    r = (route or "").strip()
    if not r or r == "self":
        return f"meru.self.{menu_id}" if r == "self" else f"meru.id.{menu_id}"
    if is_menu == 0 or (r and not r.startswith("/") and "/" not in r):
        code = re.sub(r"[^a-zA-Z0-9_.-]+", "_", r)
        return code[:120] or f"meru.perm.{menu_id}"
    code = r.strip("/").replace("/", ".")
    code = re.sub(r"[^a-zA-Z0-9_.-]+", "_", code)
    return (code or f"meru.path.{menu_id}")[:120]


def esc(s: str) -> str:
    return s.replace("\\", "\\\\").replace("'", "\\'")


def parse_rows(sql: str) -> list[dict]:
    blocks = re.findall(r"INSERT INTO `eb_system_menu` VALUES (.*?);", sql, re.S)
    text = ",".join(blocks)
    rows: list[dict] = []
    for m in re.finditer(
        r"\((\d+),(\d+),('(?:\\'|[^'])*'|NULL),('(?:\\'|[^'])*'|NULL),('(?:\\'|[^'])*'|NULL),"
        r"('(?:\\'|[^'])*'|NULL),('(?:\\'|[^'])*'|NULL),(-?\d+),(\d+),(\d+),(\d+),"
        r"('(?:\\'|[^'])*'|[^,]+),('(?:\\'|[^'])*'|[^,]+),(\d+)\)",
        text,
    ):
        rows.append(
            {
                "id": int(m.group(1)),
                "pid": int(m.group(2)),
                "icon": unq(m.group(4)),
                "name": unq(m.group(5)),
                "route": unq(m.group(6)),
                "sort": int(m.group(8)),
                "is_show": int(m.group(9)),
                "is_mer": int(m.group(10)),
                "is_menu": int(m.group(11)),
                "create_time": unq(m.group(12)),
            }
        )
    return rows


def main() -> int:
    crmeb = Path(sys.argv[1]) if len(sys.argv) > 1 else DEFAULT_CRMEB
    if not crmeb.is_file():
        print(f"missing CRMEB SQL: {crmeb}", file=sys.stderr)
        return 1

    mer = [r for r in parse_rows(crmeb.read_text(errors="ignore")) if r["is_mer"] == 1]
    ids = {r["id"] for r in mer}
    path_comp = dict(re.findall(r"'(/[^']*)'\s*:\s*'([^']+)'", REGISTRY.read_text(encoding="utf-8")))

    used_codes: set[str] = set()
    out_rows: list[dict] = []
    for r in sorted(mer, key=lambda x: x["id"]):
        is_menu = 1 if r["is_menu"] == 1 else 2
        route = r["route"]
        # CRMEB admin list hides path and shows route as 菜单地址 (incl. special "self").
        path = route or ""
        code = route_to_code(route, r["id"], r["is_menu"])
        base = code
        n = 2
        while code.lower() in used_codes:
            code = f"{base}.{r['id']}" if n == 2 else f"{base}.{r['id']}.{n}"
            n += 1
        used_codes.add(code.lower())

        icon = r["icon"]
        if icon and ":" not in icon:
            icon = ICON_MAP.get(icon, f"lucide:{icon}" if icon else "")

        component = ""
        if is_menu == 1 and path.startswith("/"):
            comp = path_comp.get(path) or path_comp.get(path.rstrip("/")) or path_comp.get(path + "/")
            if comp:
                component = comp if comp.startswith("views/") else f"views/{comp}.vue"
                if not component.endswith(".vue"):
                    component = f"views/{comp}.vue"

        create_time = r["create_time"] or "2020-01-01 00:00:00"
        if create_time.startswith("'"):
            create_time = create_time.strip("'")

        out_rows.append(
            {
                "id": r["id"],
                "parent_id": r["pid"] if r["pid"] in ids else 0,
                "code": code,
                "name": r["name"],
                "path": path[:255],
                "component": component[:255],
                "icon": icon[:128],
                "is_menu": is_menu,
                "is_route": 1 if is_menu == 1 and path.startswith("/") else 0,
                "sort": r["sort"],
                "status": 1 if r["is_show"] == 1 else 0,
                "created_at": create_time,
            }
        )

    lines = [
        "SET NAMES utf8mb4;",
        "USE `qixi_crm_merchant`;",
        "-- Full CRMEB merchant menus (eb_system_menu.is_mer=1) → qixi_crm_m_menu",
        "-- Regenerate: python3 scripts/generate-crmeb-merchant-menus.py",
        "SET FOREIGN_KEY_CHECKS=0;",
        "DELETE FROM `qixi_crm_m_role_menu`;",
        "DELETE FROM `qixi_crm_m_menu`;",
        "SET FOREIGN_KEY_CHECKS=1;",
        "",
    ]
    batch = 40
    for i in range(0, len(out_rows), batch):
        chunk = out_rows[i : i + batch]
        lines.append(
            "INSERT INTO `qixi_crm_m_menu` "
            "(`id`,`parent_id`,`code`,`name`,`path`,`component`,`icon`,`is_menu`,`is_route`,`sort`,`status`,`created_at`) VALUES"
        )
        vals = [
            f"({x['id']},{x['parent_id']},'{esc(x['code'])}','{esc(x['name'])}','{esc(x['path'])}',"
            f"'{esc(x['component'])}','{esc(x['icon'])}',{x['is_menu']},{x['is_route']},{x['sort']},"
            f"{x['status']},'{esc(x['created_at'])}')"
            for x in chunk
        ]
        lines.append(",\n".join(vals) + ";")
        lines.append("")

    lines.append("INSERT IGNORE INTO `qixi_crm_m_role_menu` (`role_code`,`menu_id`)")
    lines.append(
        "SELECT roles.role_code, menus.id FROM "
        "(SELECT 'owner' AS role_code UNION ALL SELECT 'manager') AS roles "
        "CROSS JOIN `qixi_crm_m_menu` AS menus;"
    )
    lines.append("")
    OUT.write_text("\n".join(lines), encoding="utf-8")

    summary = {
        "total": len(out_rows),
        "pages": sum(1 for x in out_rows if x["is_menu"] == 1),
        "buttons": sum(1 for x in out_rows if x["is_menu"] == 2),
        "tops": [(x["id"], x["name"], x["path"]) for x in out_rows if x["parent_id"] == 0],
        "out": str(OUT),
    }
    SUMMARY.write_text(json.dumps(summary, ensure_ascii=False, indent=2), encoding="utf-8")
    print(json.dumps(summary, ensure_ascii=False, indent=2))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
