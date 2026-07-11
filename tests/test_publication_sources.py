from __future__ import annotations

import json
import shutil
import subprocess
import tempfile
import unittest
from pathlib import Path
from types import SimpleNamespace
from unittest import mock

from tools import publication_sources


PANDOC_READER = "markdown-simple_tables-multiline_tables-grid_tables-yaml_metadata_block"


def pandoc_image_targets(source: str) -> list[str]:
    result = subprocess.run(
        ["pandoc", "--from", PANDOC_READER, "--to", "json"],
        input=source,
        capture_output=True,
        text=True,
        check=False,
    )
    if result.returncode != 0:
        raise AssertionError(result.stderr)
    document = json.loads(result.stdout)
    targets: list[str] = []

    def visit(value):
        if isinstance(value, dict):
            if value.get("t") == "Image":
                targets.append(value["c"][-1][0])
            for child in value.values():
                visit(child)
        elif isinstance(value, list):
            for child in value:
                visit(child)

    visit(document)
    return targets


class PublicationSourceTests(unittest.TestCase):
    def scan(self, source: str) -> list[str]:
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            (root / "SUMMARY.md").write_text("* [Book](book.md)\n", encoding="utf-8")
            (root / "book.md").write_text(source, encoding="utf-8")
            return [target for _, _, target in publication_sources.find_unapproved_remote_images(root)]

    def test_publication_scan_matches_real_pandoc_image_nodes(self):
        self.assertIsNotNone(shutil.which("pandoc"), "real Pandoc is required for this contract")
        cases = {
            "escaped inline scheme": "![x](https\\://example.com/inline.svg)\n",
            "escaped reference scheme": (
                "![x][id]\n\n[id]: https\\://example.com/reference.svg\n"
            ),
            "blockquote definition": (
                "![x][id]\n\n> [id]: https://example.com/blockquote.svg\n"
            ),
            "closing bracket in alt code span": (
                "![a `]`][id]\n\n[id]: https://example.com/code-close.svg\n"
            ),
            "opening bracket in alt code span": (
                "![a `[`][id]\n\n[id]: https://example.com/code-open.svg\n"
            ),
            "full reference": "![x][id]\n\n[id]: https://example.com/full.svg\n",
            "collapsed reference": "![x][]\n\n[x]: https://example.com/collapsed.svg\n",
            "shortcut reference": "![x]\n\n[x]: https://example.com/shortcut.svg\n",
            "odd escaped marker": "\\![x][id]\n\n[id]: https://example.com/odd.svg\n",
            "even escaped marker": "\\\\![x][id]\n\n[id]: https://example.com/even.svg\n",
            "inline code": "`![x](https://example.com/inline-code.svg)`\n",
            "indented code": "    ![x](https://example.com/indented-code.svg)\n",
            "HTML comment": "<!-- ![x](https://example.com/comment.svg) -->\n",
        }
        for name, source in cases.items():
            with self.subTest(name=name):
                self.assertEqual(self.scan(source), pandoc_image_targets(source))

    def test_missing_pandoc_fails_closed_with_source_context(self):
        with mock.patch("shutil.which", return_value=None):
            with self.assertRaisesRegex(ValueError, "Pandoc.*publication"):
                self.scan("plain text\n")

    def test_pandoc_parse_failure_fails_closed_with_command_context(self):
        failure = SimpleNamespace(returncode=2, stdout="", stderr="reader failed")
        with mock.patch("subprocess.run", return_value=failure):
            with self.assertRaisesRegex(ValueError, "reader failed"):
                self.scan("plain text\n")

    def test_multiple_publication_sources_use_one_pandoc_batch(self):
        self.assertIsNotNone(shutil.which("pandoc"), "real Pandoc is required for this contract")
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            (root / "SUMMARY.md").write_text(
                "* [First](first.md)\n* [Second](second.md)\n", encoding="utf-8"
            )
            (root / "first.md").write_text("![x][shared]\n", encoding="utf-8")
            (root / "second.md").write_text(
                "[shared]: https://example.com/cross-file.svg\n", encoding="utf-8"
            )
            real_run = subprocess.run
            with mock.patch.object(
                publication_sources.subprocess, "run", wraps=real_run
            ) as pandoc_run:
                issues = publication_sources.find_unapproved_remote_images(root)

            self.assertEqual(pandoc_run.call_count, 1)
            self.assertEqual(
                [(path.name, target) for path, _, target in issues],
                [("first.md", "https://example.com/cross-file.svg")],
            )

    def test_markdown_image_srcset_checks_every_candidate(self):
        source = (
            '![local](local.png){srcset="local-small.png 1x, '
            'https://example.com/remote-one.png 2x, '
            'https://img.shields.io/badge/build-passing-green.svg 3x"}\n'
        )
        self.assertEqual(
            self.scan(source),
            ["https://example.com/remote-one.png"],
        )

    def test_raw_html_img_srcset_handles_quoted_and_escaped_urls(self):
        source = (
            "<img src='local.png' "
            "srcset='local-small.png 1x, "
            "https&#x3A;//example.com/remote-two.png 2x'>\n"
        )
        self.assertEqual(
            self.scan(source),
            ["https://example.com/remote-two.png"],
        )

    def test_picture_source_checks_src_and_srcset_without_splitting_data_commas(self):
        source = (
            '<picture><source src="https://example.com/source-direct.webp" '
            'srcset="data:text/plain,https://example.com/not-a-candidate 1x, '
            'local.webp 2x, https://example.com/remote-three.webp 3x">'
            '<img src="local.png"></picture>\n'
        )
        self.assertEqual(
            self.scan(source),
            [
                "https://example.com/remote-three.webp",
                "https://example.com/source-direct.webp",
            ],
        )


if __name__ == "__main__":
    unittest.main()
