"""Find unsafe remote images in Markdown files published by SUMMARY.md."""

from __future__ import annotations

import html
import re
import string
from pathlib import Path
from urllib.parse import urlparse


ALLOWED_REMOTE_IMAGE_HOSTS = frozenset({"img.shields.io"})
SUMMARY_ENTRY_RE = re.compile(r"(?m)^\s*[-*]\s+\[[^]]*\]\(([^)]+)\)")
HTML_IMAGE_RE = re.compile(
    r"<img\b[^>]*\bsrc\s*=\s*(?:\"([^\"]+)\"|'([^']+)'|([^\s>]+))",
    re.IGNORECASE,
)
FENCE_RE = re.compile(r"^\s{0,3}(`{3,}|~{3,})")
ESCAPABLE_PUNCTUATION = frozenset(string.punctuation)
MAX_REFERENCE_LABEL_LENGTH = 999


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


def _scan_bracket_content(
    body: str,
    start: int,
    *,
    allow_nested: bool,
    max_length: int | None = None,
    reject_blank_lines: bool = False,
) -> tuple[str, int, int] | None:
    """Return decoded content, the next offset, and raw content length."""
    if start >= len(body) or body[start] != "[":
        return None

    content: list[str] = []
    cursor = start + 1
    nesting = 0
    while cursor < len(body):
        char = body[cursor]
        raw_length = cursor - start
        if char == "]" and nesting == 0:
            return "".join(content), cursor + 1, raw_length - 1
        if max_length is not None and raw_length > max_length:
            return None
        if char == "\\" and cursor + 1 < len(body):
            escaped = body[cursor + 1]
            if escaped in ESCAPABLE_PUNCTUATION:
                if max_length is not None and raw_length + 1 > max_length:
                    return None
                content.append(escaped)
                cursor += 2
                continue
        if char == "[":
            if not allow_nested:
                return None
            nesting += 1
        elif char == "]":
            nesting -= 1
        elif char == "\n" and reject_blank_lines:
            lookahead = cursor + 1
            while lookahead < len(body) and body[lookahead] in " \t":
                lookahead += 1
            if lookahead < len(body) and body[lookahead] == "\n":
                return None
        content.append(char)
        cursor += 1
    return None


def _scan_reference_label(body: str, start: int) -> tuple[str, int, int] | None:
    return _scan_bracket_content(
        body,
        start,
        allow_nested=False,
        max_length=MAX_REFERENCE_LABEL_LENGTH,
        reject_blank_lines=True,
    )


def _closing_brackets(body: str) -> dict[int, int]:
    closing_brackets: dict[int, int] = {}
    stack: list[int] = []
    backslashes = 0
    for offset, char in enumerate(body):
        if char == "\\":
            backslashes += 1
            continue
        escaped = backslashes % 2 == 1
        backslashes = 0
        if escaped:
            continue
        if char == "[":
            stack.append(offset)
        elif char == "]" and stack:
            closing_brackets[stack.pop()] = offset
    return closing_brackets


def _decode_bracket_content(body: str, start: int, end: int) -> str:
    content: list[str] = []
    cursor = start + 1
    while cursor < end:
        char = body[cursor]
        if char == "\\" and cursor + 1 < end:
            escaped = body[cursor + 1]
            if escaped in ESCAPABLE_PUNCTUATION:
                content.append(escaped)
                cursor += 2
                continue
        content.append(char)
        cursor += 1
    return "".join(content)


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
    line_start = 0
    while line_start < len(body):
        line_end = body.find("\n", line_start)
        if line_end < 0:
            line_end = len(body)

        label_start = line_start
        while label_start < line_end and body[label_start] in " \t":
            label_start += 1
        if label_start - line_start > 3 or label_start >= line_end:
            line_start = line_end + 1
            continue

        scanned = _scan_reference_label(body, label_start)
        if scanned is None:
            line_start = line_end + 1
            continue
        label, label_end, _ = scanned
        if not label or label_end >= len(body) or body[label_end] != ":":
            line_start = line_end + 1
            continue

        destination_line_end = body.find("\n", label_end + 1)
        if destination_line_end < 0:
            destination_line_end = len(body)
        remainder = body[label_end + 1 : destination_line_end]
        normalized_label = _normalize_reference_label(label)
        if not normalized_label:
            line_start = line_end + 1
            continue
        destination = _parse_reference_destination(remainder)
        if not destination and destination_line_end < len(body):
            next_line = body[destination_line_end + 1 :].split("\n", 1)[0]
            destination = _parse_reference_destination(next_line)
        if destination:
            definitions.setdefault(normalized_label, destination)
        line_start = line_end + 1
    return definitions


def _markdown_image_matches(
    body: str, definitions: dict[str, str]
) -> list[tuple[int, str]]:
    matches: list[tuple[int, str]] = []
    closing_brackets = _closing_brackets(body)
    cursor = 0
    while cursor < len(body):
        marker = body.find("![", cursor)
        if marker < 0:
            break
        if _is_escaped(body, marker):
            cursor = marker + 2
            continue

        alt_start = marker + 1
        alt_end = closing_brackets.get(alt_start)
        if alt_end is None:
            cursor = marker + 2
            continue
        alt = _decode_bracket_content(body, alt_start, alt_end)
        image_end = alt_end + 1
        alt_length = alt_end - alt_start - 1
        cursor = image_end

        if image_end < len(body) and body[image_end] == "(":
            target = _parse_reference_destination(body[image_end + 1 :])
            if target:
                matches.append((marker, target))
            continue

        label = alt
        label_length = alt_length
        if image_end < len(body) and body[image_end] == "[":
            label_scanned = _scan_reference_label(body, image_end)
            if label_scanned is None:
                continue
            explicit_label, cursor, explicit_length = label_scanned
            if explicit_label:
                label = explicit_label
                label_length = explicit_length
        if not label or label_length > MAX_REFERENCE_LABEL_LENGTH:
            continue
        normalized_label = _normalize_reference_label(label)
        if not normalized_label:
            continue
        target = definitions.get(normalized_label)
        if target is not None:
            matches.append((marker, target))
    return matches


def find_unapproved_remote_images(source_root: Path) -> list[tuple[Path, int, str]]:
    issues: list[tuple[Path, int, str]] = []
    for path in published_markdown_paths(source_root):
        body = _strip_fenced_blocks(path.read_text(encoding="utf-8", errors="ignore"))
        definitions = _reference_definitions(body)
        matches = _markdown_image_matches(body, definitions)
        for match in HTML_IMAGE_RE.finditer(body):
            target = next(group for group in match.groups() if group is not None)
            matches.append((match.start(), target))
        for start, target in matches:
            if _is_unapproved_remote_image(target):
                issues.append((path, body[:start].count("\n") + 1, target))
    return sorted(issues, key=lambda issue: (str(issue[0]), issue[1], issue[2]))
