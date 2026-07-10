"""Find unsafe remote images in Markdown files published by SUMMARY.md."""

from __future__ import annotations

import html
import re
from pathlib import Path
from urllib.parse import urlparse


ALLOWED_REMOTE_IMAGE_HOSTS = frozenset({"img.shields.io"})
SUMMARY_ENTRY_RE = re.compile(r"(?m)^\s*[-*]\s+\[[^]]*\]\(([^)]+)\)")
MARKDOWN_IMAGE_RE = re.compile(r"!\[[^]]*\]\(([^)\s]+)")
FULL_REFERENCE_IMAGE_RE = re.compile(r"!\[([^]]*)\]\[([^]]*)\]")
SHORTCUT_REFERENCE_IMAGE_RE = re.compile(r"!\[([^]]+)\](?![\[(])")
REFERENCE_DEFINITION_RE = re.compile(
    r"^[ \t]{0,3}\[((?:[^\[\]\n]|\n(?![ \t]*\n)){1,999})\]:[ \t]*(.*)$",
    re.MULTILINE,
)
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


def _normalize_reference_label(label: str) -> str:
    return re.sub(r"\s+", " ", label.strip()).casefold()


def _is_escaped(body: str, start: int) -> bool:
    backslashes = 0
    while start > 0 and body[start - 1] == "\\":
        backslashes += 1
        start -= 1
    return backslashes % 2 == 1


def _parse_reference_destination(value: str) -> str:
    value = value.lstrip()
    if not value:
        return ""
    if value.startswith("<"):
        end = value.find(">", 1)
        return value[1:end] if end >= 0 else ""

    destination: list[str] = []
    depth = 0
    escaped = False
    for char in value:
        if escaped:
            destination.append(char)
            escaped = False
            continue
        if char == "\\":
            destination.append(char)
            escaped = True
            continue
        if char.isspace() and depth == 0:
            break
        if char == "(":
            depth += 1
        elif char == ")":
            if depth == 0:
                break
            depth -= 1
        destination.append(char)
    return "".join(destination)


def _reference_definitions(body: str) -> dict[str, str]:
    definitions: dict[str, str] = {}
    for match in REFERENCE_DEFINITION_RE.finditer(body):
        label, remainder = match.groups()
        normalized_label = _normalize_reference_label(label)
        if not normalized_label:
            continue
        destination = _parse_reference_destination(remainder)
        if not destination and body.startswith("\n", match.end()):
            next_line = body[match.end() + 1 :].split("\n", 1)[0]
            destination = _parse_reference_destination(next_line)
        if destination:
            definitions.setdefault(normalized_label, destination)
    return definitions


def find_unapproved_remote_images(source_root: Path) -> list[tuple[Path, int, str]]:
    issues: list[tuple[Path, int, str]] = []
    for path in published_markdown_paths(source_root):
        body = _strip_fenced_blocks(path.read_text(encoding="utf-8", errors="ignore"))
        definitions = _reference_definitions(body)
        matches: list[tuple[int, str]] = []
        matches.extend(
            (match.start(), match.group(1))
            for match in MARKDOWN_IMAGE_RE.finditer(body)
            if not _is_escaped(body, match.start())
        )
        for match in FULL_REFERENCE_IMAGE_RE.finditer(body):
            if _is_escaped(body, match.start()):
                continue
            alt, explicit_label = match.groups()
            label = explicit_label or alt
            target = definitions.get(_normalize_reference_label(label))
            if target is not None:
                matches.append((match.start(), target))
        for match in SHORTCUT_REFERENCE_IMAGE_RE.finditer(body):
            if _is_escaped(body, match.start()):
                continue
            normalized_label = _normalize_reference_label(match.group(1))
            if not normalized_label:
                continue
            target = definitions.get(normalized_label)
            if target is not None:
                matches.append((match.start(), target))
        for match in HTML_IMAGE_RE.finditer(body):
            target = next(group for group in match.groups() if group is not None)
            matches.append((match.start(), target))
        for start, target in matches:
            if _is_unapproved_remote_image(target):
                issues.append((path, body[:start].count("\n") + 1, target))
    return sorted(issues, key=lambda issue: (str(issue[0]), issue[1], issue[2]))
