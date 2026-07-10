#!/usr/bin/env python3
"""Render Mermaid diagrams from SUMMARY chapters to numbered SVG files."""

from __future__ import annotations

import argparse
import glob
import json
import os
import re
import shutil
import subprocess
import sys
from pathlib import Path


SAFETY_MESSAGE = "must be an independent directory outside protected source trees"
STANDARD_CHROME_PATHS = (
    "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
    "/Applications/Chromium.app/Contents/MacOS/Chromium",
)


def paths_overlap(left: Path, right: Path) -> bool:
    return left == right or left in right.parents or right in left.parents


def validate_output_directory(book_dir: str, svg_out: str) -> tuple[Path, Path]:
    book = Path(book_dir).expanduser().resolve()
    output = Path(svg_out).expanduser().resolve()
    repository = Path(__file__).resolve().parents[1]
    protected_exact = {Path(output.anchor).resolve(), Path.home().resolve(), Path.cwd().resolve()}
    if output in protected_exact or paths_overlap(output, book) or paths_overlap(output, repository):
        raise ValueError(f"--svg-out {output} {SAFETY_MESSAGE}")
    return book, output


def clean_generated_files(output: Path) -> None:
    generated_svg = re.compile(r"^(?:d-\d+|_c(?:-\d+)?)\.svg$")
    generated_names = {"_chunk.md", "_pptr.json", "_rc.json"}
    for entry in output.iterdir():
        if entry.name in generated_names or generated_svg.fullmatch(entry.name):
            if entry.is_file() or entry.is_symlink():
                entry.unlink()


def find_chrome() -> str | None:
    configured = os.environ.get("CHROME_BIN")
    candidates = (
        [configured]
        if configured
        else [
            *(shutil.which(name) for name in (
                "google-chrome-stable", "google-chrome", "chromium-browser", "chromium", "chrome"
            )),
            *STANDARD_CHROME_PATHS,
        ]
    )
    return next(
        (str(Path(path).resolve()) for path in candidates if path and Path(path).is_file()),
        None,
    )


def summary_sources(book: Path) -> list[str]:
    summary = book / "SUMMARY.md"
    if not summary.is_file():
        raise ValueError(f"SUMMARY.md not found under {book}")
    sources: list[str] = []
    seen: set[str] = set()
    for line in summary.read_text(encoding="utf-8").splitlines():
        match = re.match(r"^\s*[-*]\s+\[.*?\]\(([^)]+?)\)", line)
        if not match:
            continue
        relative = match.group(1).split("#", 1)[0].strip()
        chapter = (book / relative).resolve()
        if not relative.endswith(".md") or relative in seen or not chapter.is_file():
            continue
        if book not in chapter.parents:
            raise ValueError(f"SUMMARY entry escapes book directory: {relative}")
        seen.add(relative)
        text = chapter.read_text(encoding="utf-8")
        sources.extend(
            match.group(1)
            for match in re.finditer(r"```mermaid[ \t]*\n(.*?)\n[ \t]*```", text, re.DOTALL)
        )
    return sources


parser = argparse.ArgumentParser(description=__doc__)
parser.add_argument("--book-dir", default=".")
parser.add_argument("--svg-out", required=True)
parser.add_argument("--chunk", type=int, default=25)
mode = parser.add_mutually_exclusive_group()
mode.add_argument(
    "--require-all",
    action="store_true",
    help="fail if Chrome, mmdc, or any rendered SVG is missing",
)
mode.add_argument(
    "--allow-fallback",
    action="store_true",
    help="explicitly allow missing SVGs to fall back to Mermaid source",
)
args = parser.parse_args()

try:
    book_path, svg_path = validate_output_directory(args.book_dir, args.svg_out)
    sources = summary_sources(book_path)
except ValueError as error:
    print(f"Mermaid rendering failed: {error}", file=sys.stderr)
    raise SystemExit(2)

svg_path.mkdir(parents=True, exist_ok=True)
clean_generated_files(svg_path)
total = len(sources)
print(f"mermaid diagrams found: {total}")
if total == 0:
    raise SystemExit(0)

chrome = find_chrome()
if not chrome:
    message = "no Chrome executable found"
    if args.require_all:
        print(f"Mermaid rendering failed: {message}", file=sys.stderr)
        raise SystemExit(1)
    print(f"WARNING: {message} -> all diagrams will fall back to source")
    raise SystemExit(0)

mmdc = shutil.which("mmdc")
if not mmdc:
    message = "mmdc is not on PATH"
    if args.require_all:
        print(f"Mermaid rendering failed: {message}", file=sys.stderr)
        raise SystemExit(1)
    print(f"WARNING: {message} -> all diagrams will fall back to source")
    raise SystemExit(0)

print(f"using Chrome: {chrome}")
puppeteer_config = svg_path / "_pptr.json"
render_config = svg_path / "_rc.json"
puppeteer_config.write_text(
    json.dumps({
        "executablePath": chrome,
        "args": ["--no-sandbox", "--disable-gpu", "--disable-dev-shm-usage"],
    }),
    encoding="utf-8",
)
render_config.write_text(json.dumps({"theme": "default"}), encoding="utf-8")


def rendered_count() -> int:
    return sum((svg_path / f"d-{index + 1}.svg").is_file() for index in range(total))


def render(indices: list[int]) -> None:
    chunk = svg_path / "_chunk.md"
    chunk.write_text(
        "\n".join(f"```mermaid\n{sources[index]}\n```\n" for index in indices),
        encoding="utf-8",
    )
    for stale in glob.glob(str(svg_path / "_c*.svg")):
        Path(stale).unlink()
    result = subprocess.run(
        [
            mmdc, "-i", str(chunk), "-o", str(svg_path / "_c.svg"),
            "-p", str(puppeteer_config), "-c", str(render_config), "-b", "transparent",
        ],
        capture_output=True,
        text=True,
        check=False,
    )
    if result.returncode != 0:
        print(result.stderr.strip() or result.stdout.strip(), file=sys.stderr)
    for rendered_index, source_index in enumerate(indices, 1):
        candidate = svg_path / f"_c-{rendered_index}.svg"
        if len(indices) == 1 and not candidate.is_file():
            candidate = svg_path / "_c.svg"
        if candidate.is_file() and candidate.stat().st_size > 0:
            candidate.replace(svg_path / f"d-{source_index + 1}.svg")
    for stale in glob.glob(str(svg_path / "_c*.svg")):
        Path(stale).unlink()


for chunk_index, start in enumerate(range(0, total, args.chunk), 1):
    render(list(range(start, min(start + args.chunk, total))))
    print(f"  chunk {chunk_index}: {rendered_count()}/{total}", flush=True)

for attempt in range(4):
    missing = [index for index in range(total) if not (svg_path / f"d-{index + 1}.svg").is_file()]
    if not missing:
        break
    print(f"  retry {attempt + 1}: {len(missing)} missing", flush=True)
    for start in range(0, len(missing), 8):
        render(missing[start:start + 8])

for temporary in (puppeteer_config, render_config, svg_path / "_chunk.md"):
    temporary.unlink(missing_ok=True)
print(f"RENDERED {rendered_count()}/{total} diagrams")
missing = [index + 1 for index in range(total) if not (svg_path / f"d-{index + 1}.svg").is_file()]
if args.require_all and missing:
    print(f"Mermaid rendering failed for diagrams: {missing}", file=sys.stderr)
    raise SystemExit(1)
