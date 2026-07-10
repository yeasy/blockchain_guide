#!/usr/bin/env python3
"""Verify PDF/HTML publications and write a portable SHA-256 manifest."""

from __future__ import annotations

import argparse
import hashlib
import html
import re
import shutil
import subprocess
import sys
from pathlib import Path


class ArtifactError(ValueError):
    """A publication is missing, malformed, or inconsistent with its source."""


def normalized_title(value: str) -> str:
    return " ".join(html.unescape(value).split())


def require_file(path: Path) -> None:
    if not path.is_file():
        raise ArtifactError(f"{path} does not exist or is not a file")
    if path.stat().st_size == 0:
        raise ArtifactError(f"{path} is empty")


def command_output(command: list[str]) -> str:
    result = subprocess.run(command, capture_output=True, text=True, check=False)
    if result.returncode != 0:
        detail = (result.stderr or result.stdout).strip()
        raise ArtifactError(f"command failed ({' '.join(command)}): {detail}")
    return result.stdout


def verify_pdf(path: Path, expected_title: str) -> None:
    require_file(path)
    with path.open("rb") as source:
        if source.read(5) != b"%PDF-":
            raise ArtifactError(f"{path} does not have a PDF signature")

    if shutil.which("pdfinfo") is None or shutil.which("pdftotext") is None:
        raise ArtifactError("pdfinfo and pdftotext are required for PDF title verification")

    metadata = command_output(["pdfinfo", str(path)])
    title_match = re.search(r"(?m)^Title:\s*(.*)$", metadata)
    metadata_title = normalized_title(title_match.group(1)) if title_match else ""
    expected = normalized_title(expected_title)
    if metadata_title == expected:
        return

    cover_text = normalized_title(
        command_output(["pdftotext", "-f", "1", "-l", "2", str(path), "-"])
    )
    if expected not in cover_text:
        raise ArtifactError(
            f"{path} title mismatch: expected {expected_title!r}; "
            f"PDF metadata title was {metadata_title!r}"
        )


def count_summary_mermaid(source_root: Path) -> int:
    """Count Mermaid fences in unique Markdown files published by SUMMARY.md."""

    source_root = source_root.resolve()
    summary = source_root / "SUMMARY.md"
    require_file(summary)
    paths: list[Path] = []
    seen: set[Path] = set()
    for match in re.finditer(
        r"(?m)^\s*[-*]\s+\[[^]]*\]\(([^)]+)\)",
        summary.read_text(encoding="utf-8"),
    ):
        relative = match.group(1).split("#", 1)[0].strip()
        if not relative.endswith(".md"):
            continue
        path = (source_root / relative).resolve()
        if source_root not in path.parents and path != source_root:
            raise ArtifactError(f"SUMMARY entry escapes source root: {relative}")
        if path in seen:
            continue
        require_file(path)
        seen.add(path)
        paths.append(path)

    return sum(
        len(re.findall(r"(?m)^[ \t]*```mermaid[ \t]*$", path.read_text(encoding="utf-8")))
        for path in paths
    )


def verify_html(path: Path, expected_title: str, *, expected_mermaid_count: int) -> None:
    require_file(path)
    text = path.read_text(encoding="utf-8")
    title_match = re.search(
        r"<title(?:\s[^>]*)?>(.*?)</title>", text, re.IGNORECASE | re.DOTALL
    )
    actual_title = normalized_title(title_match.group(1)) if title_match else ""
    expected = normalized_title(expected_title)
    if actual_title != expected:
        raise ArtifactError(
            f"{path} title mismatch: expected {expected_title!r}, got {actual_title!r}"
        )

    placeholder = re.search(
        r"MERMAIDZZ\d+ZZ|PGBKZZ|"
        r"<pre\b[^>]*class=[\"'][^\"']*\bdiagram-fallback\b[^\"']*[\"'][^>]*>",
        text,
        re.IGNORECASE,
    )
    if placeholder:
        raise ArtifactError(
            f"{path} contains unresolved reader placeholder: {placeholder.group(0)}"
        )

    figures = re.findall(
        r"<figure\b[^>]*\bclass=[\"'][^\"']*\bdiagram\b[^\"']*[\"'][^>]*>"
        r".*?</figure>",
        text,
        re.IGNORECASE | re.DOTALL,
    )
    if any("<svg" not in figure.lower() for figure in figures):
        raise ArtifactError(f"{path} contains a Mermaid figure without an inline SVG")
    if len(figures) != expected_mermaid_count:
        raise ArtifactError(
            f"{path} Mermaid count mismatch: expected {expected_mermaid_count}, "
            f"got {len(figures)}"
        )


def sha256_file(path: Path) -> str:
    digest = hashlib.sha256()
    with path.open("rb") as source:
        for chunk in iter(lambda: source.read(1024 * 1024), b""):
            digest.update(chunk)
    return digest.hexdigest()


def write_checksums(paths: list[Path], destination: Path) -> None:
    destination.parent.mkdir(parents=True, exist_ok=True)
    lines: list[str] = []
    names: set[str] = set()
    for path in sorted(paths, key=lambda item: item.name):
        require_file(path)
        if path.parent.resolve() != destination.parent.resolve():
            raise ArtifactError(f"{path} must be beside checksum manifest {destination}")
        if path.name in names:
            raise ArtifactError(f"duplicate artifact name: {path.name}")
        names.add(path.name)
        lines.append(f"{sha256_file(path)}  {path.name}\n")
    if not lines:
        raise ArtifactError("cannot write an empty checksum manifest")
    destination.write_text("".join(lines), encoding="utf-8")


def verify_checksums(manifest: Path) -> None:
    require_file(manifest)
    entries = manifest.read_text(encoding="utf-8").splitlines()
    if not entries:
        raise ArtifactError(f"{manifest} is empty")
    names: set[str] = set()
    for line_number, line in enumerate(entries, 1):
        match = re.fullmatch(r"([0-9a-f]{64})  ([^/\\]+)", line)
        if match is None:
            raise ArtifactError(f"{manifest}:{line_number}: malformed SHA-256 entry")
        expected, name = match.groups()
        if name == manifest.name or name in names:
            raise ArtifactError(
                f"{manifest}:{line_number}: duplicate or recursive entry {name!r}"
            )
        names.add(name)
        path = manifest.parent / name
        require_file(path)
        actual = sha256_file(path)
        if actual != expected:
            raise ArtifactError(
                f"checksum mismatch for {path}: expected {expected}, got {actual}"
            )


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--title", required=True)
    parser.add_argument("--pdf", type=Path, required=True)
    parser.add_argument("--html", type=Path, required=True)
    parser.add_argument("--source-root", type=Path, required=True)
    parser.add_argument("--checksums", type=Path, required=True)
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    try:
        verify_pdf(args.pdf, args.title)
        expected_mermaid_count = count_summary_mermaid(args.source_root)
        verify_html(args.html, args.title, expected_mermaid_count=expected_mermaid_count)
        write_checksums([args.pdf, args.html], args.checksums)
        verify_checksums(args.checksums)
    except ArtifactError as error:
        print(f"artifact verification failed: {error}", file=sys.stderr)
        return 1

    print(f"verified artifacts: {args.pdf}, {args.html}")
    print(f"verified Mermaid diagrams: {expected_mermaid_count}")
    print(f"verified checksums: {args.checksums}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
