#!/usr/bin/env python3
"""Generate docs/acceptance/PARITY-GAP-MATRIX.md from feature TSVs + registries + SQL seeds."""

from __future__ import annotations

import re
from collections import defaultdict
from datetime import date
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
OUT = ROOT / "docs/acceptance/PARITY-GAP-MATRIX.md"

PORTALS = {
    "平台后台": {
        "slug": "platform",
        "registry": ROOT / "admin-platform/src/views/ecrm/registry.ts",
        "sql": ROOT / "sql/admin/init_data.sql",
        "sql_path_re": re.compile(
            r"\(\d+,\d+,'[^']+','[^']*','[^']*','([^']+)','page',\d+\)"
        ),
        "feature_tsv": ROOT / "docs/generated/features-01-platform-admin.tsv",
        "baseline_ops": 1333,
    },
    "商户后台": {
        "slug": "merchant",
        "registry": ROOT / "admin-merchant/src/views/ecrm/registry.ts",
        "sql": ROOT / "sql/merchant/init_data.sql",
        "sql_path_re": re.compile(
            r"\(\d+,\d+,'[^']+','[^']+','(/[^']+)','views/ecrm/[^']*',"
        ),
        "feature_tsv": ROOT / "docs/generated/features-02-merchant-admin.tsv",
        "baseline_ops": 615,
    },
}


def load_registry(path: Path) -> set[str]:
    text = path.read_text(encoding="utf-8")
    paths = set(re.findall(r"'(/[^']*)'\s*:\s*'ecrm/", text))
    # also accept unquoted bare keys like brokerage
    for m in re.finditer(r"(?m)^\s*([A-Za-z_][\w]*)\s*:\s*'ecrm/", text):
        paths.add(m.group(1))
    return paths


def load_sql_pages(path: Path, pattern: re.Pattern[str]) -> set[str]:
    text = path.read_text(encoding="utf-8")
    return {m.group(1) for m in pattern.finditer(text)}


def load_feature_pages(tsv: Path) -> dict[str, dict]:
    """page_route -> {module, ops, title}"""
    pages: dict[str, dict] = {}
    if not tsv.exists():
        return pages
    with tsv.open(encoding="utf-8") as fh:
        header = fh.readline().rstrip("\n").split("\t")
        idx = {name: i for i, name in enumerate(header)}
        route_key = "page_route" if "page_route" in idx else "page"
        for line in fh:
            cols = line.rstrip("\n").split("\t")
            if len(cols) <= max(idx.values()):
                continue
            path = cols[idx[route_key]].strip().strip("`")
            if not path.startswith("/") and path != "brokerage":
                continue
            module = cols[idx["module_path"]] if "module_path" in idx else cols[idx.get("module", 1)]
            title = cols[idx["page"]] if "page" in idx else path
            entry = pages.setdefault(path, {"module": module, "ops": 0, "title": title})
            entry["ops"] += 1
    return pages


def is_registered(path: str, registry: set[str]) -> bool:
    if path in registry:
        return True
    normalized = path.rstrip("/")
    return any(r == normalized or r.rstrip("/") == normalized for r in registry)


def classify(path: str, registry: set[str], sql_pages: set[str]) -> str:
    if is_registered(path, registry):
        return "partial"
    if path in sql_pages:
        return "shell"
    return "missing"


def main() -> None:
    lines: list[str] = [
        "# CRMEB 双后台对照缺口矩阵",
        "",
        f"> 自动生成：`scripts/generate-parity-gap-matrix.py` · {date.today().isoformat()}",
        ">",
        "> 状态：`missing` 无页面 · `shell` SQL 有菜单但 registry 无组件 · `partial` 已注册待按钮级闭环 · `done` 需人工在验收台账标关闭后回写。",
        ">",
        "> 基线路径来自 `docs/generated/features-0{1,2}-*.tsv` 的 `page_route`；完成以 TSV 操作行 + `crmeb-vben-parity.md` 六条为准。",
        "",
    ]

    for portal, cfg in PORTALS.items():
        registry = load_registry(cfg["registry"])
        sql_pages = load_sql_pages(cfg["sql"], cfg["sql_path_re"])
        feature_pages = load_feature_pages(cfg["feature_tsv"])

        # Baseline = CRMEB feature routes; also surface SQL/registry-only extras
        paths = sorted(set(feature_pages) | sql_pages | registry)

        by_status: dict[str, list[str]] = defaultdict(list)
        by_module: dict[str, dict[str, int]] = defaultdict(lambda: defaultdict(int))
        covered_ops = 0
        total_ops = sum(v["ops"] for v in feature_pages.values())

        for path in paths:
            status = classify(path, registry, sql_pages)
            by_status[status].append(path)
            mod = (feature_pages.get(path) or {}).get("module") or (
                path.split("/")[1] if path.startswith("/") and "/" in path[1:] else "其他"
            )
            if isinstance(mod, str) and " / " in mod:
                mod = mod.split(" / ")[0]
            by_module[str(mod) or "其他"][status] += 1
            if status == "partial" and path in feature_pages:
                covered_ops += feature_pages[path]["ops"]

        lines += [
            f"## {portal}（{cfg['slug']}）",
            "",
            f"| 指标 | 数量 |",
            f"| --- | ---: |",
            f"| 基线操作数 | {cfg['baseline_ops']} |",
            f"| 特征表页面路由 | {len(feature_pages)} |",
            f"| SQL 叶子 page | {len(sql_pages)} |",
            f"| registry 路径 | {len(registry)} |",
            f"| 特征表 ops（已挂 registry） | {covered_ops} / {total_ops} |",
            f"| missing | {len(by_status['missing'])} |",
            f"| shell | {len(by_status['shell'])} |",
            f"| partial | {len(by_status['partial'])} |",
            "",
            "### 按模块状态计数",
            "",
            "| 模块 | missing | shell | partial |",
            "| --- | ---: | ---: | ---: |",
        ]
        for mod in sorted(by_module):
            c = by_module[mod]
            lines.append(
                f"| {mod} | {c['missing']} | {c['shell']} | {c['partial']} |"
            )

        lines += ["", "### missing（优先补齐）", ""]
        missing = [p for p in by_status["missing"] if p in feature_pages]
        extras = [p for p in by_status["missing"] if p not in feature_pages]
        show = missing or by_status["missing"]
        if not show:
            lines.append("_无_")
        else:
            for path in show[:150]:
                ops = (feature_pages.get(path) or {}).get("ops", 0)
                title = (feature_pages.get(path) or {}).get("title", "")
                suffix = f" · {title}" if title else ""
                suffix += f" · {ops} ops" if ops else ""
                lines.append(f"- `{path}`{suffix}")
            if len(show) > 150:
                lines.append(f"- … 另有 {len(show) - 150} 条")
            if extras and missing:
                lines.append(f"- （另有 {len(extras)} 条非特征表路径）")

        lines += ["", "### shell（有菜单无组件）", ""]
        shell = by_status["shell"]
        if not shell:
            lines.append("_无_")
        else:
            for path in shell:
                lines.append(f"- `{path}`")

        lines += ["", "### partial（已注册，待按钮级 / 布局闭环）", ""]
        for path in sorted(registry)[:100]:
            ops = (feature_pages.get(path) or {}).get("ops", 0)
            lines.append(f"- `{path}`" + (f" · {ops} ops" if ops else ""))
        if len(registry) > 100:
            lines.append(f"- … 另有 {len(registry) - 100} 条 registry 路径")
        lines.append("")

    OUT.parent.mkdir(parents=True, exist_ok=True)
    OUT.write_text("\n".join(lines) + "\n", encoding="utf-8")
    print(f"wrote {OUT.relative_to(ROOT)}")


if __name__ == "__main__":
    main()
