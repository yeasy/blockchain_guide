from __future__ import annotations

import contextlib
import io
import tempfile
import unittest
from pathlib import Path
from unittest import mock

import check_project_rules as checker


class ProjectRuleTests(unittest.TestCase):
    def run_checker(self, source: str) -> tuple[int, str]:
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            (root / "SUMMARY.md").write_text("* [Book](README.md)\n", encoding="utf-8")
            (root / "README.md").write_text(source, encoding="utf-8")
            with (
                mock.patch.object(checker, "ROOT", root),
                contextlib.redirect_stdout(io.StringIO()) as stdout,
            ):
                result = checker.main()
            return result, stdout.getvalue()

    def test_checker_rejects_unapproved_remote_images_in_publication_sources(self):
        result, output = self.run_checker(
            "![dynamic chart](https://example.com/dynamic.svg)\n"
        )
        self.assertEqual(result, 1)
        self.assertIn("remote image", output)

    def test_checker_allows_explicitly_whitelisted_badges(self):
        result, output = self.run_checker(
            "[![build](https://img.shields.io/badge/build-passing-green.svg)](https://example.com)\n"
        )
        self.assertEqual(result, 0, output)

    def test_checker_rejects_full_reference_image_with_normalized_multiline_label(self):
        result, output = self.run_checker(
            "![dynamic chart][Chart\n   ID]\n\n"
            "[  chart\n"
            "   id ]:\n"
            "  <https://example.com/dynamic.svg>\n"
            '  "Dynamic chart"\n'
        )
        self.assertEqual(result, 1)
        self.assertIn("README.md:6: remote image", output)
        self.assertIn("https://example.com/dynamic.svg", output)

    def test_checker_rejects_collapsed_and_shortcut_reference_images(self):
        cases = (
            "![dynamic chart][]\n\n[dynamic chart]: https://example.com/collapsed.svg\n",
            "![dynamic chart]\n\n[dynamic chart]: https://example.com/shortcut.svg\n",
            "![   ]\n\n[   ]: https://example.com/empty.svg\n",
        )
        for source in cases:
            with self.subTest(source=source.splitlines()[0]):
                result, output = self.run_checker(source)
                self.assertEqual(result, 1)
                self.assertIn("remote image", output)

    def test_checker_rejects_pandoc_reference_labels_with_escaped_punctuation(self):
        cases = (
            (
                "![dynamic\\] chart](https://example.com/inline-escaped.svg)\n",
                "inline-escaped.svg",
            ),
            (
                "![dynamic chart][Chart\\]   ID]\n\n"
                "[chart\\] id]: https://example.com/full-escaped.svg\n",
                "full-escaped.svg",
            ),
            (
                "![chart\\]][]\n\n"
                "[chart\\]]: https://example.com/collapsed-escaped.svg\n",
                "collapsed-escaped.svg",
            ),
            (
                "![chart\\]]\n\n"
                "[chart\\]]: https://example.com/shortcut-escaped.svg\n",
                "shortcut-escaped.svg",
            ),
            (
                "\\\\![chart\\]]\n\n"
                "[chart\\]]: https://example.com/even-backslashes.svg\n",
                "even-backslashes.svg",
            ),
            (
                "![chart\\\\][]\n\n"
                "[chart\\\\]: https://example.com/escaped-backslash.svg\n",
                "escaped-backslash.svg",
            ),
            (
                "![chart](https\\://example.com/escaped-scheme.svg)\n",
                "escaped-scheme.svg",
            ),
            (
                "![chart][id]\n\n> [id]: https://example.com/blockquote.svg\n",
                "blockquote.svg",
            ),
            (
                "![a `]`][id]\n\n[id]: https://example.com/code-span.svg\n",
                "code-span.svg",
            ),
        )
        for source, target in cases:
            with self.subTest(source=source.splitlines()[0]):
                result, output = self.run_checker(source)
                self.assertEqual(result, 1)
                self.assertIn(target, output)

    def test_checker_applies_badge_allowlist_to_reference_images(self):
        result, output = self.run_checker(
            "[![build][badge]][project]\n\n"
            "[badge]: https://img.shields.io/badge/build-passing-green.svg\n"
            "[project]: https://example.com\n"
        )
        self.assertEqual(result, 0, output)

    def test_checker_ignores_non_image_markdown_syntax(self):
        cases = (
            "\\![dynamic chart]\n\n[dynamic chart]: https://example.com/escaped.svg\n",
            "\\![chart\\]]\n\n[chart\\]]: https://example.com/escaped-bracket.svg\n",
            "`![code](https://example.com/inline-code.svg)`\n",
            "    ![code](https://example.com/indented-code.svg)\n",
            "<!-- ![comment](https://example.com/comment.svg) -->\n",
        )
        for source in cases:
            with self.subTest(source=source.splitlines()[0]):
                result, output = self.run_checker(source)
                self.assertEqual(result, 0, output)


if __name__ == "__main__":
    unittest.main()
