from __future__ import annotations

import contextlib
import io
import tempfile
import unittest
from pathlib import Path
from types import SimpleNamespace
from unittest import mock


try:
    from tools import verify_release_artifacts as verifier
except ImportError:
    verifier = None


class ReleaseArtifactTests(unittest.TestCase):
    def module(self):
        self.assertIsNotNone(verifier, "tools/verify_release_artifacts.py must be importable")
        return verifier

    def test_html_title_and_every_published_mermaid_are_verified(self):
        module = self.module()
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / "book.html"
            valid = (
                "<!doctype html><html><head><title>区块链技术指南</title></head><body>"
                '<figure class="diagram"><svg></svg></figure></body></html>'
            )
            path.write_text(valid, encoding="utf-8")
            module.verify_html(path, "区块链技术指南", expected_mermaid_count=1)
            invalid = (
                valid.replace("区块链技术指南", "错误标题", 1),
                valid.replace("<svg></svg>", "fallback"),
                valid.replace("</body>", '<pre class="diagram-fallback">source</pre></body>'),
                valid.replace("</body>", "MERMAIDZZ0ZZ</body>"),
            )
            for text in invalid:
                path.write_text(text, encoding="utf-8")
                with self.assertRaises(module.ArtifactError):
                    module.verify_html(path, "区块链技术指南", expected_mermaid_count=1)

    def test_pdf_signature_and_title_are_mandatory(self):
        module = self.module()
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / "book.pdf"
            path.write_bytes(b"%PDF-1.7\nfixture")
            with (
                mock.patch.object(module.shutil, "which", return_value="/usr/bin/tool"),
                mock.patch.object(module, "command_output", return_value="Title: 区块链技术指南\n"),
            ):
                module.verify_pdf(path, "区块链技术指南")
            path.write_bytes(b"invalid")
            with self.assertRaises(module.ArtifactError):
                module.verify_pdf(path, "区块链技术指南")

    def test_manifest_covers_both_artifacts_and_detects_tampering(self):
        module = self.module()
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            pdf, html = root / "book.pdf", root / "book.html"
            manifest = root / "SHA256SUMS"
            pdf.write_bytes(b"pdf")
            html.write_bytes(b"html")
            module.write_checksums([pdf, html], manifest)
            module.verify_checksums(manifest)
            self.assertEqual(len(manifest.read_text(encoding="utf-8").splitlines()), 2)
            html.write_bytes(b"tampered")
            with self.assertRaises(module.ArtifactError):
                module.verify_checksums(manifest)

    def test_mermaid_count_only_uses_unique_summary_chapters(self):
        module = self.module()
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            (root / "SUMMARY.md").write_text("* [A](a.md)\n* [again](a.md)\n", encoding="utf-8")
            (root / "a.md").write_text("```mermaid\ngraph TD\n```\n", encoding="utf-8")
            (root / "unused.md").write_text("```mermaid\ngraph LR\n```\n", encoding="utf-8")
            self.assertEqual(module.count_summary_mermaid(root), 1)

    def test_artifact_main_rejects_unapproved_remote_publication_images(self):
        module = self.module()
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            pdf, html = root / "book.pdf", root / "book.html"
            manifest = root / "SHA256SUMS"
            (root / "SUMMARY.md").write_text("* [Book](book.md)\n", encoding="utf-8")
            source = root / "book.md"
            source.write_text(
                '<img src="https://example.com/dynamic.svg" alt="dynamic">\n',
                encoding="utf-8",
            )
            args = SimpleNamespace(
                title="Book",
                pdf=pdf,
                html=html,
                source_root=root,
                checksums=manifest,
            )
            with (
                mock.patch.object(module, "parse_args", return_value=args),
                mock.patch.object(module, "verify_pdf"),
                mock.patch.object(module, "verify_html"),
                mock.patch.object(module, "write_checksums"),
                mock.patch.object(module, "verify_checksums"),
                contextlib.redirect_stderr(io.StringIO()) as stderr,
            ):
                self.assertEqual(module.main(), 1)
            self.assertIn("remote image", stderr.getvalue())

            source.write_text(
                "![dynamic chart][chart]\n\n"
                "[chart]: https://example.com/reference.svg\n",
                encoding="utf-8",
            )
            with (
                mock.patch.object(module, "parse_args", return_value=args),
                mock.patch.object(module, "verify_pdf"),
                mock.patch.object(module, "verify_html"),
                mock.patch.object(module, "write_checksums"),
                mock.patch.object(module, "verify_checksums"),
                contextlib.redirect_stderr(io.StringIO()) as stderr,
            ):
                self.assertEqual(module.main(), 1)
            self.assertIn("reference.svg", stderr.getvalue())

            source.write_text(
                "![dynamic chart][chart\\]]\n\n"
                "[chart\\]]: https://example.com/escaped-reference.svg\n",
                encoding="utf-8",
            )
            with (
                mock.patch.object(module, "parse_args", return_value=args),
                mock.patch.object(module, "verify_pdf"),
                mock.patch.object(module, "verify_html"),
                mock.patch.object(module, "write_checksums"),
                mock.patch.object(module, "verify_checksums"),
                contextlib.redirect_stderr(io.StringIO()) as stderr,
            ):
                self.assertEqual(module.main(), 1)
            self.assertIn("escaped-reference.svg", stderr.getvalue())

            source.write_text(
                "[![build](https://img.shields.io/badge/build-passing-green.svg)](https://example.com)\n",
                encoding="utf-8",
            )
            with (
                mock.patch.object(module, "parse_args", return_value=args),
                mock.patch.object(module, "verify_pdf"),
                mock.patch.object(module, "verify_html"),
                mock.patch.object(module, "write_checksums"),
                mock.patch.object(module, "verify_checksums"),
                contextlib.redirect_stdout(io.StringIO()),
            ):
                self.assertEqual(module.main(), 0)

    def test_html_reader_disables_ambiguous_pandoc_block_syntax(self):
        source = (Path(__file__).resolve().parents[1] / "tools" / "build_html_reader.py").read_text(
            encoding="utf-8"
        )
        self.assertIn(
            'reader = "markdown-simple_tables-multiline_tables-grid_tables-yaml_metadata_block"',
            source,
        )
        self.assertIn('"-f", reader', source)


if __name__ == "__main__":
    unittest.main()
