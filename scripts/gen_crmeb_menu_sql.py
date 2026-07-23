#!/usr/bin/env python3
"""Generate sql/043_crmeb_system_menu_full.sql from CRMEB MER v4.0 dump."""
from __future__ import annotations

import os
import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
DEFAULT_SRC = Path(
    os.path.expanduser("~/Downloads/CRMEB多商户系统/CRMEB_MER_v4.0/install/crmeb_merchant.sql")
)
DEFAULT_BACKUP = Path(
    os.path.expanduser("~/Downloads/CRMEB多商户系统/CRMEB_MER_v4.0/backup/update_4_0.sql")
)
OUT = ROOT / "sql" / "043_crmeb_system_menu_full.sql"
RESERVED = set(range(1, 173))
ID_OFFSET = 20000


def sql_str(s: str) -> str:
    return "'" + s.replace("\\", "\\\\").replace("'", "\\'") + "'"


def unquote(f: str) -> str:
    f = f.strip()
    if len(f) >= 2 and f[0] == "'" and f[-1] == "'":
        return f[1:-1].replace("\\'", "'").replace("\\\\", "\\")
    return f


def parse_tuples(blob: str) -> list[list[str]]:
    idx = blob.find("VALUES")
    if idx < 0:
        return []
    data = blob[idx + 6 :]
    rows: list[list[str]] = []
    i, n = 0, len(data)
    while i < n:
        while i < n and data[i] in " \t\r\n,":
            i += 1
        if i >= n or data[i] != "(":
            break
        i += 1
        fields: list[str] = []
        cur: list[str] = []
        in_str = False
        esc = False
        while i < n:
            ch = data[i]
            if in_str:
                cur.append(ch)
                if esc:
                    esc = False
                elif ch == "\\":
                    esc = True
                elif ch == "'":
                    in_str = False
                i += 1
                continue
            if ch == "'":
                in_str = True
                cur.append(ch)
                i += 1
                continue
            if ch == ",":
                fields.append("".join(cur).strip())
                cur = []
                i += 1
                continue
            if ch == ")":
                fields.append("".join(cur).strip())
                i += 1
                break
            cur.append(ch)
            i += 1
        if len(fields) >= 14:
            rows.append(fields)
        while i < n and data[i] not in ",;":
            i += 1
        if i < n and data[i] == ";":
            break
        if i < n and data[i] == ",":
            i += 1
    return rows


def load_inserts(path: Path) -> list[list[str]]:
    if not path.exists():
        return []
    text = path.read_text(encoding="utf-8", errors="ignore")
    rows: list[list[str]] = []
    for m in re.finditer(r"INSERT INTO `eb_system_menu`[^;]+;", text, re.I | re.S):
        rows.extend(parse_tuples(m.group(0)))
    return rows


def main() -> None:
    src = Path(os.environ.get("CRMEB_MENU_SQL", DEFAULT_SRC))
    backup = Path(os.environ.get("CRMEB_MENU_SQL_BACKUP", DEFAULT_BACKUP))
    all_rows = load_inserts(src) + load_inserts(backup)
    by_id: dict[int, list[str]] = {}
    for fields in all_rows:
        try:
            mid = int(fields[0])
        except Exception:
            continue
        by_id[mid] = fields
    remap = {mid: (mid + ID_OFFSET if mid in RESERVED else mid) for mid in by_id}

    lines = [
        "-- 全量导入 CRMEB MER v4.0 eb_system_menu → qixi_system_menu",
        f"-- 来源: {src.name}" + (f" + {backup.name}" if backup.exists() else ""),
        "-- 变换: is_mer 0→1 / 1→2; is_menu 0→2 / 1→1; 与本仓自定义 menu_id 1–172 冲突的 CRMEB 行偏移 +20000",
        "-- 生成: python3 scripts/gen_crmeb_menu_sql.py",
        "USE `qixi_mergers`;",
        "",
    ]
    for mid, fields in sorted(by_id.items()):
        new_id = remap[mid]
        pid = int(fields[1] or 0)
        new_pid = remap.get(pid, pid) if pid else 0
        path = unquote(fields[2])
        icon = unquote(fields[3])
        name = unquote(fields[4])
        route = unquote(fields[5])
        params = unquote(fields[6])
        sort = int(fields[7] or 0)
        is_show = int(fields[8] or 0)
        crmeb_is_mer = int(fields[9] or 0)
        crmeb_is_menu = int(fields[10] or 0)
        is_agent = int(fields[13] or 0)
        qixi_is_mer = 1 if crmeb_is_mer == 0 else 2
        qixi_is_menu = 1 if crmeb_is_menu == 1 else 2
        lines.append(
            "INSERT INTO `qixi_system_menu` (`menu_id`,`pid`,`path`,`icon`,`menu_name`,`route`,`params`,`sort`,`is_show`,`is_mer`,`is_menu`,`is_agent`)\n"
            f"SELECT {new_id}, {new_pid}, {sql_str(path)}, {sql_str(icon)}, {sql_str(name)}, {sql_str(route)}, {sql_str(params)}, {sort}, {is_show}, {qixi_is_mer}, {qixi_is_menu}, {is_agent}\n"
            f"WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = {new_id});"
        )

    lines.extend(
        [
            "",
            "-- 角色 rules：按导入后的全表重建（平台超管 / 商户模板角色）",
            "SET SESSION group_concat_max_len = 1024*1024;",
            "",
            "UPDATE `qixi_system_role` r",
            "SET r.`rules` = (",
            "  SELECT GROUP_CONCAT(m.`menu_id` ORDER BY m.`menu_id` SEPARATOR ',')",
            "  FROM `qixi_system_menu` m WHERE m.`is_mer` = 1",
            ")",
            "WHERE r.`role_id` = 1 AND r.`mer_id` = 0;",
            "",
            "UPDATE `qixi_system_role` r",
            "SET r.`rules` = (",
            "  SELECT GROUP_CONCAT(m.`menu_id` ORDER BY m.`menu_id` SEPARATOR ',')",
            "  FROM `qixi_system_menu` m WHERE m.`is_mer` = 2",
            ")",
            "WHERE r.`role_id` = 2;",
            "",
            "INSERT INTO `qixi_schema_meta` (`version`, `note`)",
            "SELECT 'crmeb-menu-043', '全量 CRMEB 菜单/按钮导入（冲突 id +20000；rules 按 is_mer 重建）'",
            "WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'crmeb-menu-043');",
            "",
        ]
    )
    OUT.write_text("\n".join(lines), encoding="utf-8")
    print(f"wrote {OUT} rows={len(by_id)} remapped={sum(1 for m in by_id if m in RESERVED)}")


if __name__ == "__main__":
    main()
