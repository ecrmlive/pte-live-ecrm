#!/usr/bin/env python3
"""Patch native DIY model components: replace $parent.$parent with inject bridge."""
from __future__ import annotations

import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1] / "src/views/native/page/diy/model"

INJECT_BLOCK = """
inject: ['diyModel'],
methods: {
  diyEditer(index) {
    this.diyModel?.onEditer(index);
  },
  diyDeleteItem(index) {
    this.diyModel?.onDeleleItem(index);
  },
""".strip()

METHODS_PATCH = """
  diyEditer(index) {
    this.diyModel?.onEditer(index);
  },
  diyDeleteItem(index) {
    this.diyModel?.onDeleleItem(index);
  },
"""


def patch_file(path: Path) -> bool:
    text = path.read_text(encoding="utf-8")
    if "$parent.$parent.onEditer" not in text and "$parent.onEditer" not in text:
        return False

    text = text.replace("$parent.$parent.onEditer", "diyEditer")
    text = text.replace("$parent.$parent.onDeleleItem", "diyDeleteItem")
    text = text.replace("$parent.onEditer", "diyEditer")

    if "inject: ['diyModel']" in text:
        path.write_text(text, encoding="utf-8")
        return True

    if "export default {" in text:
        if re.search(r"\bmethods:\s*\{", text):
            text = re.sub(r"(\bmethods:\s*\{)", rf"\1{METHODS_PATCH}", text, count=1)
            text = text.replace(
                "export default {",
                "export default {\n  inject: ['diyModel'],",
                1,
            )
        else:
            text = text.replace(
                "export default {",
                f"export default {{\n  {INJECT_BLOCK},\n",
                1,
            )
    path.write_text(text, encoding="utf-8")
    return True


def main() -> None:
    count = 0
    for path in sorted(ROOT.glob("*.vue")):
        if patch_file(path):
            count += 1
            print(f"patched {path.name}")
    print(f"done: {count} files")


if __name__ == "__main__":
    main()
