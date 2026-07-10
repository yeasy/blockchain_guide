"""Find unsafe remote images in published Markdown via Pandoc's production AST."""

from __future__ import annotations

import html
import json
import re
import secrets
import shutil
import string
import subprocess
from html.parser import HTMLParser
from pathlib import Path
from typing import Iterator
from urllib.parse import urlparse


ALLOWED_REMOTE_IMAGE_HOSTS = frozenset({"img.shields.io"})
SUMMARY_ENTRY_RE = re.compile(r"(?m)^\s*[-*]\s+\[[^]]*\]\(([^)]+)\)")
PANDOC_MARKDOWN_READER = (
    "markdown-simple_tables-multiline_tables-grid_tables-yaml_metadata_block"
)
PANDOC_TIMEOUT_SECONDS = 60
ESCAPABLE_PUNCTUATION = frozenset(string.punctuation)


class _HTMLImageParser(HTMLParser):
    """Extract img src values from HTML fragments already identified by Pandoc."""

    def __init__(self) -> None:
        super().__init__(convert_charrefs=True)
        self.targets: list[str] = []

    def _record_image(
        self, tag: str, attrs: list[tuple[str, str | None]]
    ) -> None:
        if tag.casefold() != "img":
            return
        for name, value in attrs:
            if name.casefold() == "src" and value is not None:
                self.targets.append(value)
                return

    def handle_starttag(
        self, tag: str, attrs: list[tuple[str, str | None]]
    ) -> None:
        self._record_image(tag, attrs)

    def handle_startendtag(
        self, tag: str, attrs: list[tuple[str, str | None]]
    ) -> None:
        self._record_image(tag, attrs)


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


def _pandoc_document(source: str, source_root: Path, source_count: int) -> dict:
    pandoc = shutil.which("pandoc")
    if pandoc is None:
        raise ValueError(
            f"Pandoc is required to inspect publication images under {source_root}"
        )

    command = [
        pandoc,
        "--from",
        PANDOC_MARKDOWN_READER,
        "--to",
        "json",
    ]
    try:
        result = subprocess.run(
            command,
            input=source,
            capture_output=True,
            text=True,
            check=False,
            timeout=PANDOC_TIMEOUT_SECONDS,
        )
    except (OSError, subprocess.TimeoutExpired) as error:
        raise ValueError(
            f"Pandoc publication image parse could not run for {source_count} "
            f"source(s) under {source_root} ({' '.join(command)}): {error}"
        ) from error

    if result.returncode != 0:
        detail = (result.stderr or result.stdout).strip() or "no diagnostic output"
        raise ValueError(
            f"Pandoc publication image parse failed for {source_count} source(s) "
            f"under {source_root} ({' '.join(command)}): {detail}"
        )

    try:
        document = json.loads(result.stdout)
    except (TypeError, json.JSONDecodeError) as error:
        raise ValueError(
            f"Pandoc returned invalid JSON while inspecting {source_count} source(s) "
            f"under {source_root}: {error}"
        ) from error
    if not isinstance(document, dict) or not isinstance(document.get("blocks"), list):
        raise ValueError(
            f"Pandoc returned an invalid document AST while inspecting "
            f"{source_count} source(s) under {source_root}"
        )
    return document


def _source_marker_index(block: object, marker_prefix: str) -> int | None:
    if not isinstance(block, dict) or block.get("t") not in {"Para", "Plain"}:
        return None
    content = block.get("c")
    if not isinstance(content, list) or len(content) != 1:
        return None
    token = content[0]
    if not isinstance(token, dict) or token.get("t") != "Str":
        return None
    value = token.get("c")
    if not isinstance(value, str):
        return None
    match = re.fullmatch(re.escape(marker_prefix) + r"(\d+)BOUNDARY", value)
    return int(match.group(1)) if match else None


def _walk_ast(value: object) -> Iterator[dict]:
    if isinstance(value, dict):
        yield value
        for child in value.values():
            yield from _walk_ast(child)
    elif isinstance(value, list):
        for child in value:
            yield from _walk_ast(child)


def _image_targets(node: dict, source_root: Path) -> list[str]:
    node_type = node.get("t")
    content = node.get("c")
    if node_type == "Image":
        if (
            not isinstance(content, list)
            or not content
            or not isinstance(content[-1], list)
            or not content[-1]
            or not isinstance(content[-1][0], str)
        ):
            raise ValueError(
                f"Pandoc returned a malformed Image node while inspecting {source_root}"
            )
        return [content[-1][0]]

    if node_type not in {"RawInline", "RawBlock"}:
        return []
    if (
        not isinstance(content, list)
        or len(content) != 2
        or not isinstance(content[0], str)
        or not isinstance(content[1], str)
        or not content[0].casefold().startswith("html")
    ):
        return []
    parser = _HTMLImageParser()
    parser.feed(content[1])
    parser.close()
    return parser.targets


def _markdown_unescape(text: str) -> str:
    """Decode target spelling only for best-effort diagnostic line lookup."""

    output: list[str] = []
    cursor = 0
    while cursor < len(text):
        if (
            text[cursor] == "\\"
            and cursor + 1 < len(text)
            and text[cursor + 1] in ESCAPABLE_PUNCTUATION
        ):
            cursor += 1
        output.append(text[cursor])
        cursor += 1
    return html.unescape("".join(output))


def _target_line(source: str, target: str) -> int:
    """Locate a resolved target for diagnostics; AST detection remains authoritative."""

    decoded_source = _markdown_unescape(source)
    offset = decoded_source.find(target)
    if offset < 0:
        return 1
    return decoded_source[:offset].count("\n") + 1


def find_unapproved_remote_images(source_root: Path) -> list[tuple[Path, int, str]]:
    source_root = source_root.resolve()
    paths = published_markdown_paths(source_root)
    if not paths:
        return []

    bodies: list[str] = []
    for path in paths:
        try:
            bodies.append(path.read_text(encoding="utf-8"))
        except (OSError, UnicodeError) as error:
            raise ValueError(
                f"Could not read publication source {path} while inspecting images: {error}"
            ) from error

    marker_prefix = f"PANDOCSOURCE{secrets.token_hex(16).upper()}ZZ"
    chunks = [
        f"\n\n{marker_prefix}{index}BOUNDARY\n\n{body}\n"
        for index, body in enumerate(bodies)
    ]
    document = _pandoc_document("".join(chunks), source_root, len(paths))

    issues: list[tuple[Path, int, str]] = []
    current_index: int | None = None
    next_marker = 0
    for block in document["blocks"]:
        marker_index = _source_marker_index(block, marker_prefix)
        if marker_index is not None:
            if marker_index != next_marker or marker_index >= len(paths):
                raise ValueError(
                    f"Pandoc source boundaries were invalid while inspecting {source_root}"
                )
            current_index = marker_index
            next_marker += 1
            continue

        for node in _walk_ast(block):
            targets = _image_targets(node, source_root)
            if targets and current_index is None:
                raise ValueError(
                    f"Pandoc found an image outside a publication source boundary under "
                    f"{source_root}"
                )
            for target in targets:
                if _is_unapproved_remote_image(target):
                    assert current_index is not None
                    issues.append(
                        (
                            paths[current_index],
                            _target_line(bodies[current_index], target),
                            target,
                        )
                    )

    if next_marker != len(paths):
        raise ValueError(
            f"Pandoc did not preserve all {len(paths)} publication source boundaries "
            f"under {source_root}; refusing an incomplete image check"
        )
    return sorted(issues, key=lambda issue: (str(issue[0]), issue[1], issue[2]))
