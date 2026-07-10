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


if __name__ == "__main__":
    unittest.main()
