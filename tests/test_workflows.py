from __future__ import annotations

import json
import os
import re
import subprocess
import tempfile
import textwrap
import unittest
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
WORKFLOW_DIR = ROOT / ".github" / "workflows"
WORKFLOWS = tuple(sorted(WORKFLOW_DIR.glob("*.y*ml")))
FULL_ACTION = re.compile(r"^[^@\s]+@[0-9a-f]{40}$")
REPO = "owner/repo"
SHA = "a" * 40
GET_REF = ["api", "--include", "--method", "GET", f"repos/{REPO}/git/ref/tags/preview-pdf"]
PATCH_REF = [
    "api", "--silent", "--method", "PATCH",
    f"repos/{REPO}/git/refs/tags/preview-pdf",
    "--raw-field", f"sha={SHA}", "--field", "force=true",
]
POST_REF = [
    "api", "--silent", "--method", "POST", f"repos/{REPO}/git/refs",
    "--raw-field", "ref=refs/tags/preview-pdf", "--raw-field", f"sha={SHA}",
]
VIEW_RELEASE = ["release", "view", "preview-pdf"]
EDIT_RELEASE = [
    "release", "edit", "preview-pdf", "--title", "Latest Preview Publications",
    "--notes-file", "dist/release-notes.md", "--prerelease",
]
CREATE_RELEASE = [
    "release", "create", "preview-pdf", "--title", "Latest Preview Publications",
    "--notes-file", "dist/release-notes.md", "--prerelease", "--latest=false",
    "--verify-tag",
]
UPLOAD_RELEASE = [
    "release", "upload", "preview-pdf", "dist/blockchain_guide.pdf",
    "dist/blockchain_guide.html", "dist/SHA256SUMS", "--clobber",
]


FAKE_GH = r'''#!/usr/bin/env python3
import json, os, sys
args = sys.argv[1:]
with open(os.environ["GH_LOG"], "a", encoding="utf-8") as stream:
    stream.write(json.dumps(args) + "\n")
scenario = os.environ["GH_SCENARIO"]
repo = "owner/repo"
sha = "a" * 40
get_ref = ["api", "--include", "--method", "GET", f"repos/{repo}/git/ref/tags/preview-pdf"]
patch_ref = ["api", "--silent", "--method", "PATCH", f"repos/{repo}/git/refs/tags/preview-pdf", "--raw-field", f"sha={sha}", "--field", "force=true"]
post_ref = ["api", "--silent", "--method", "POST", f"repos/{repo}/git/refs", "--raw-field", "ref=refs/tags/preview-pdf", "--raw-field", f"sha={sha}"]
view_release = ["release", "view", "preview-pdf"]
edit_release = ["release", "edit", "preview-pdf", "--title", "Latest Preview Publications", "--notes-file", "dist/release-notes.md", "--prerelease"]
create_release = ["release", "create", "preview-pdf", "--title", "Latest Preview Publications", "--notes-file", "dist/release-notes.md", "--prerelease", "--latest=false", "--verify-tag"]
upload_release = ["release", "upload", "preview-pdf", "dist/blockchain_guide.pdf", "dist/blockchain_guide.html", "dist/SHA256SUMS", "--clobber"]
if os.environ.get("GH_REPO") != repo:
    print("explicit GH_REPO required", file=sys.stderr); raise SystemExit(2)
if args == get_ref:
    if scenario == "network":
        print("network failure", file=sys.stderr); raise SystemExit(1)
    code = scenario.split("_", 1)[0]
    if code in {"401", "403", "404", "429", "503"}:
        print(f"HTTP/2.0 {code} status"); print(f"HTTP {code}", file=sys.stderr); raise SystemExit(1)
    print("HTTP/2.0 200 OK\n\n{}"); raise SystemExit(0)
if args in (patch_ref, post_ref, edit_release, create_release, upload_release):
    raise SystemExit(0)
if args == view_release:
    if scenario.endswith("release_missing"):
        print("release not found", file=sys.stderr); raise SystemExit(1)
    if scenario.endswith("release_error"):
        print("release network error", file=sys.stderr); raise SystemExit(1)
    raise SystemExit(0)
print(f"unexpected argv: {args!r}", file=sys.stderr); raise SystemExit(2)
'''


def step_script(text: str, name: str) -> str:
    marker = f"      - name: {name}\n"
    start = text.index(marker) + len(marker)
    run = text.index("        run: |\n", start) + len("        run: |\n")
    end = text.find("\n      - name:", run)
    return textwrap.dedent(text[run : len(text) if end < 0 else end])


class WorkflowTests(unittest.TestCase):
    def text(self, name: str) -> str:
        path = WORKFLOW_DIR / name
        self.assertTrue(path.is_file(), path)
        return path.read_text(encoding="utf-8")

    def test_all_actions_are_full_sha_pinned_with_version_comments(self):
        failures = []
        for path in WORKFLOWS:
            for line_no, line in enumerate(path.read_text(encoding="utf-8").splitlines(), 1):
                match = re.search(r"\buses:\s*([^\s#]+)(?:\s+#\s*(v\S+))?", line)
                if match and (not FULL_ACTION.fullmatch(match.group(1)) or not match.group(2)):
                    failures.append(f"{path.name}:{line_no}: {line.strip()}")
        self.assertEqual(failures, [])

    def test_checkout_drops_credentials_and_permissions_are_minimal(self):
        for path in WORKFLOWS:
            text = path.read_text(encoding="utf-8")
            self.assertIn("\npermissions: {}\n", text, path.name)
            lines = text.splitlines()
            for index, line in enumerate(lines):
                if "uses: actions/checkout@" in line:
                    self.assertIn("persist-credentials: false", "\n".join(lines[index:index+10]), path.name)
        for name in ("ci.yaml", "auto-release.yml", "preview-pdf.yml"):
            self.assertRegex(self.text(name), r"(?ms)^  build:\n(?:.*?\n)*?    permissions:\n      contents: read\b")
        self.assertRegex(self.text("auto-release.yml"), r"(?ms)^  release:.*?permissions:\n      contents: write\n      id-token: write\n      attestations: write\b")
        self.assertRegex(self.text("preview-pdf.yml"), r"(?ms)^  publish:.*?permissions:\n      contents: write\b")

    def test_downloads_mermaid_and_artifacts_have_integrity_gates(self):
        for name in ("ci.yaml", "auto-release.yml", "preview-pdf.yml"):
            text = self.text(name)
            for marker in ("MDPRESS_SHA256", "PANDOC_SHA256", "sha256sum -c -", "tools/verify_release_artifacts.py", "SHA256SUMS"):
                self.assertIn(marker, text, f"{name}: {marker}")
            self.assertIn("npm ci --prefix tools/mermaid --ignore-scripts", text, name)
            self.assertNotIn("continue-on-error", text, name)
        package = json.loads((ROOT / "tools/mermaid/package.json").read_text(encoding="utf-8"))
        lock = json.loads((ROOT / "tools/mermaid/package-lock.json").read_text(encoding="utf-8"))
        self.assertEqual(package["dependencies"]["@mermaid-js/mermaid-cli"], "10.9.1")
        self.assertEqual(lock["packages"][""]["dependencies"]["@mermaid-js/mermaid-cli"], "10.9.1")

    def test_every_publication_workflow_builds_pdf_html_and_checksums(self):
        for name in ("ci.yaml", "auto-release.yml", "preview-pdf.yml"):
            text = self.text(name)
            for marker in ("mdpress build", "tools/render_mermaid.py", "tools/build_html_reader.py", "--pdf", "--html", "--checksums"):
                self.assertIn(marker, text, f"{name}: {marker}")
            self.assertIn("if-no-files-found: error", text, name)
        auto = self.text("auto-release.yml")
        self.assertIn("actions/attest-build-provenance@0f67c3f4856b2e3261c31976d6725780e5e4c373 # v4.1.1", auto)
        self.assertRegex(auto, r"(?s)subject-path:.*?\.pdf.*?\.html.*?SHA256SUMS")
        self.assertRegex(auto, r"(?s)files:.*?\.pdf.*?\.html.*?SHA256SUMS")

    def test_ci_runs_each_chaincode_package_as_a_matrix(self):
        ci = self.text("ci.yaml")
        for example in ("example01", "example02", "example03", "example04", "example05", "example06"):
            self.assertIn(example, ci)
        self.assertIn("matrix:", ci)
        self.assertIn('go test "./examples/${{ matrix.example }}" -count=1', ci)
        self.assertIn("working-directory: 11_app_dev", ci)

    def test_renderer_does_not_kill_unrelated_processes_and_can_fail_closed(self):
        source = (ROOT / "tools/render_mermaid.py").read_text(encoding="utf-8")
        self.assertNotIn("pkill", source)
        self.assertIn("--require-all", source)
        for name in ("ci.yaml", "auto-release.yml", "preview-pdf.yml"):
            self.assertIn("--require-all", self.text(name), name)

    def run_preview(self, scenario: str, *, repo: str = REPO, sha: str = SHA):
        text = self.text("preview-pdf.yml")
        scripts = [step_script(text, name) for name in (
            "Synchronize mutable preview tag", "Create or update preview release", "Replace preview assets"
        )]
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            gh = root / "gh"
            gh.write_text(FAKE_GH, encoding="utf-8")
            gh.chmod(0o755)
            log = root / "log.jsonl"
            env = os.environ.copy()
            env.update({"PATH": f"{root}{os.pathsep}{env.get('PATH','')}", "GH_LOG": str(log), "GH_SCENARIO": scenario, "GH_REPO": repo, "GH_TOKEN": "x", "GITHUB_SHA": sha})
            result = None
            for script in scripts:
                result = subprocess.run(["/bin/bash", "-c", script], cwd=ROOT, env=env, text=True, capture_output=True, check=False)
                if result.returncode:
                    break
            commands = [json.loads(line) for line in log.read_text(encoding="utf-8").splitlines()] if log.exists() else []
            return result, commands

    def test_preview_existing_and_missing_paths_use_exact_order(self):
        result, commands = self.run_preview("200_release_exists")
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual(commands, [GET_REF, PATCH_REF, VIEW_RELEASE, EDIT_RELEASE, UPLOAD_RELEASE])
        result, commands = self.run_preview("404_release_missing")
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual(commands, [GET_REF, POST_REF, VIEW_RELEASE, CREATE_RELEASE, UPLOAD_RELEASE])

    def test_preview_fails_closed_before_publish(self):
        for scenario in ("401", "403", "429", "503", "network"):
            with self.subTest(scenario=scenario):
                result, commands = self.run_preview(scenario)
                self.assertNotEqual(result.returncode, 0)
                self.assertEqual(commands, [GET_REF])
        result, commands = self.run_preview("200_release_error")
        self.assertNotEqual(result.returncode, 0)
        self.assertEqual(commands, [GET_REF, PATCH_REF, VIEW_RELEASE])
        self.assertNotIn(UPLOAD_RELEASE, commands)
        for repo, sha in (("bad/repo/extra", SHA), (REPO, "a" * 39)):
            result, commands = self.run_preview("200_release_exists", repo=repo, sha=sha)
            self.assertNotEqual(result.returncode, 0)
            self.assertEqual(commands, [])


if __name__ == "__main__":
    unittest.main()
