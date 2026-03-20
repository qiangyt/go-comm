#!/usr/bin/env python3
"""
检查 Go 文件中使用的 i18n entry 是否在 locale 文件中定义。
"""

import os
import re
import sys
from pathlib import Path
from collections import defaultdict


def find_i18n_keys_in_go(go_files: list[str]) -> dict[str, set[str]]:
    """
    从 Go 文件中提取所有 i18n 调用的 key。
    返回 {file: {keys}}
    """
    # 匹配 localize("key") 或 T("key") 格式
    pattern = re.compile(r'(?:localize|T)\s*\(\s*["\']([^"\']+)["\']')
    file_keys = defaultdict(set)
    file_keys = defaultdict(set)

    for filepath in go_files:
        try:
            content = Path(filepath).read_text()
            keys = pattern.findall(content)
            file_keys[filepath].update(keys)
        except Exception as e:
            print(f"Warning: Failed to read {filepath}: {e}", file=sys.stderr)

    return file_keys


def find_i18n_keys_in_locales(locales_dir: str) -> set[str]:
    """
    从 locale 文件中提取所有定义的 i18n key。
    返回 {keys}
    """
    keys = set()
    pattern = re.compile(r"^\[([^\]]+)\]", re.MULTILINE)

    for filepath in Path(locales_dir).glob("*.toml"):
        try:
            content = filepath.read_text()
            found_keys = pattern.findall(content)
            keys.update(found_keys)
        except Exception as e:
            print(f"Warning: Failed to read {filepath}: {e}", file=sys.stderr)

    return keys


def main():
    # 找到项目根目录
    script_dir = Path(__file__).parent
    go_comm_dir = script_dir.parent

    # 路径 - 排除测试文件
    go_files = [f for f in go_comm_dir.glob("*.go") if not f.name.endswith("_test.go")]
    locales_dir = go_comm_dir / "q18n/locales"

    if not locales_dir.exists():
        print(f"Error: Locales directory not found: {locales_dir}", file=sys.stderr)
        sys.exit(1)

    # 提取 i18n keys
    file_keys = find_i18n_keys_in_go([str(f) for f in go_files])
    locale_keys = find_i18n_keys_in_locales(str(locales_dir))

    # 检查
    errors = []
    for filepath, keys in file_keys.items():
        for key in keys:
            if key not in locale_keys:
                errors.append(f"{filepath}: i18n key '{key}' not found in locales")

    # 输出结果
    if errors:
        print("❌ i18n validation FAILED")
        print(f"Found {len(errors)} missing i18n entries:\n")
        for error in sorted(errors):
            print(f"  - {error}")
        sys.exit(1)
    else:
        print("✅ i18n validation PASSED")
        total_keys = sum(len(keys) for keys in file_keys.values())
        print(f"  - Checked {len(go_files)} Go files")
        print(f"  - Found {total_keys} i18n entries")
        print(f"  - All {len(locale_keys)} locale keys are valid")
        sys.exit(0)


if __name__ == "__main__":
    main()
