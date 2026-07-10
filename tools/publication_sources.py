"""Find unsafe remote images in Markdown files published by SUMMARY.md."""

from __future__ import annotations

import html
import re
from pathlib import Path
from urllib.parse import urlparse


ALLOWED_REMOTE_IMAGE_HOSTS = frozenset({"img.shields.io"})
SUMMARY_ENTRY_RE = re.compile(r"(?m)^\s*[-*]\s+\[[^]]*\]\(([^)]+)\)")
MARKDOWN_IMAGE_RE = re.compile(r"!\[[^]]*\]\(([^)\s]+)")
HTML_IMAGE_RE = re.compile(
    r"<img\b[^>]*\bsrc\s*=\s*(?:\"([^\"]+)\"|'([^']+)'|([^\s>]+))",
    re.IGNORECASE,
)
FENCE_RE = re.compile(r"^\s{0,3}(`{3,}|~{3,})")


def _strip_fenced_blocks(text: str) -> str:
    output: list[str] = []
    open_fence: tuple[str, int] | None = None
    for line in text.splitlines():
        match = FENCE_RE.match(line)
        if match:
            marker = match.group(1)
            if open_fence is None:
                open_fence = (marker[0], len(marker))
            elif marker[0] == open_fence[0] and len(marker) >= open_fence[1]:
                open_fence = None
            output.append("")
            continue
        output.append("" if open_fence else line)
    return "\n".join(output)


def published_markdown_paths(source_root: Path) -> list[Path]:
    source_root = source_root.resolve()
    summary = source_root / "SUMMARY.md"
    if not summary.is_file():
        return []

    paths: list[Path] = []
    seen: set[Path] = set()
    for match in SUMMARY_ENTRY_RE.finditer(summary.read_text(encoding="utf-8")):
        relative = match.group(1).split("#", 1)[0].strip()
        if not relative.endswith(".md"):
            continue
        path = (source_root / relative).resolve()
        if source_root not in path.parents and path != source_root:
            raise ValueError(f"SUMMARY entry escapes source root: {relative}")
        if path.is_file() and path not in seen:
            seen.add(path)
            paths.append(path)
    return paths


def _is_unapproved_remote_image(raw_target: str) -> bool:
    target = html.unescape(raw_target.strip().strip("<>"))
    parsed = urlparse(target)
    if parsed.scheme.lower() not in {"http", "https"} and not parsed.netloc:
        return False
    return not (
        parsed.scheme.lower() == "https"
        and (parsed.hostname or "").lower() in ALLOWED_REMOTE_IMAGE_HOSTS
    )


def find_unapproved_remote_images(source_root: Path) -> list[tuple[Path, int, str]]:
    issues: list[tuple[Path, int, str]] = []
    for path in published_markdown_paths(source_root):
        body = _strip_fenced_blocks(path.read_text(encoding="utf-8", errors="ignore"))
        matches: list[tuple[int, str]] = []
        matches.extend(
            (match.start(), match.group(1))
            for match in MARKDOWN_IMAGE_RE.finditer(body)
        )
        for match in HTML_IMAGE_RE.finditer(body):
            target = next(group for group in match.groups() if group is not None)
            matches.append((match.start(), target))
        for start, target in matches:
            if _is_unapproved_remote_image(target):
                issues.append((path, body[:start].count("\n") + 1, target))
    return sorted(issues, key=lambda issue: (str(issue[0]), issue[1], issue[2]))
