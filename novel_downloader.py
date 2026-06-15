#!/usr/bin/env python3
"""
小说下载器 - novel543.com -> 简体 txt
用法: python novel_downloader.py <目录页URL> [输出文件名]
      python novel_downloader.py https://www.novel543.com/0802619327/dir
"""

import sys
import re
import time
import random
import argparse
from urllib.parse import urljoin, urlparse

import requests
from bs4 import BeautifulSoup

try:
    import opencc
    _cc = opencc.OpenCC('t2s')
    def to_s(text): return _cc.convert(text)
except ImportError:
    print("警告: 未安装 opencc，跳过繁简转换。pip install opencc-python-reimplemented")
    def to_s(text): return text

AD_RE = re.compile(
    r'本书由.{0,20}首发'
    r'|最新章节.{0,30}小说网'
    r'|请记住.{0,40}地址'
    r'|稷下[書书]院'
    r'|novel543'
    r'|如果您喜欢.{0,30}请.{0,10}收藏'
    r'|求收藏.{0,20}求推荐'
    r'|求月票.{0,20}打赏'
    r'|(最新|最快)更新.{0,20}章节'
    r'|章节错误.{0,30}举报'
    r'|举报本章'
    r'|温馨提示[：:].{0,80}'
    r'|登录用户.{0,80}'
    r'|建议大家登录使用'
    r'|跨设备永久保存'
    r'|站内信.{0,60}'
    r'|用户中心.{0,60}'
    r'|本站新增.{0,60}'
    r'|点击.{0,10}设置.{0,30}切换'
    r'|搜书名找不到.{0,60}'
    r'|www\.[a-zA-Z0-9\-]+\.[a-zA-Z]{2,}'
    r'|http[s]?://\S+',
    re.IGNORECASE,
)

HEADERS = {
    'User-Agent': (
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) '
        'AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36'
    ),
    'Accept-Language': 'zh-CN,zh-TW;q=0.9',
}


def fetch(url, retries=4):
    delays = [2, 4, 8, 16]
    for i in range(retries):
        try:
            r = requests.get(url, headers=HEADERS, timeout=15)
            r.raise_for_status()
            r.encoding = r.apparent_encoding or 'utf-8'
            return r.text
        except Exception as e:
            if i < retries - 1:
                w = delays[i] + random.uniform(0, 1)
                print(f"    重试({e}) 等待{w:.1f}s...")
                time.sleep(w)
            else:
                print(f"    跳过 {url}: {e}")
    return None


CHAPTER_TITLE_RE = re.compile(r'^第\d+章\s*.{0,20}\(\d+/\d+\)\s*$')

STOP_MARKERS = re.compile(r'温馨提示|溫馨提示|登录用户|登錄用戶|VIP会员|VIP會員|点击查看|點擊查看')

def clean_text(raw):
    lines = raw.splitlines()
    out, blanks = [], 0
    for line in lines:
        # 遇到站点提示/广告标记行，后续内容全部丢弃
        if STOP_MARKERS.search(line):
            break
        line = AD_RE.sub('', line).strip()
        # 去掉分页标题行，如"第1章 魂瓶 (2/2)"
        if CHAPTER_TITLE_RE.match(line):
            continue
        if len(re.sub(r'[\s　\W]', '', line)) < 2:
            blanks += 1
            if blanks == 1:
                out.append('')
        else:
            blanks = 0
            out.append(line)
    return '\n'.join(out).strip()


def extract_body(soup):
    for tag in soup.find_all(['nav', 'header', 'footer', 'script', 'style', 'aside']):
        tag.decompose()
    for tag in soup.find_all(class_=re.compile(r'nav|ad|banner|toolbar|sidebar|foot|header', re.I)):
        tag.decompose()
    for sel in ['#content', '#chapter-content', '.chapter-content',
                '#chaptercontent', '.read-content', '#nr1', 'article', '.content']:
        node = soup.select_one(sel)
        if node and len(node.get_text(strip=True)) > 100:
            return node.get_text('\n')
    divs = soup.find_all('div')
    return max((d.get_text('\n') for d in divs), key=len) if divs else ''


def get_chapter_content(first_url, base):
    """抓取章节所有分页，合并正文。利用 URL 前缀判断是否仍在本章。"""
    # 章节基础名：如 8096_3（目录给的 URL 永远是第一页，直接用文件名去掉 .html）
    name0 = re.sub(r'\.html$', '', first_url.split('/')[-1])
    ch_base = name0  # e.g. "8096_3"；页码 URL 形如 "8096_3_2"，以 ch_base+"_" 开头

    parts = []
    current_url = first_url

    while current_url:
        html = fetch(current_url)
        if not html:
            break

        soup = BeautifulSoup(html, 'html.parser')

        # 先找导航链接（extract_body 会 decompose 掉 nav 标签，必须在此之前处理）
        next_url = None
        for a in soup.find_all('a', href=True):
            if re.search(r'下一[章页頁]', a.get_text(strip=True)):
                candidate = urljoin(base, a['href'])
                cname = re.sub(r'\.html$', '', candidate.split('/')[-1])
                # e.g. ch_base="8096_3", cname="8096_3_2" → same chapter
                if cname == ch_base or cname.startswith(ch_base + '_'):
                    next_url = candidate
                break

        parts.append(clean_text(to_s(extract_body(soup))))

        if next_url:
            current_url = next_url
            time.sleep(0.4 + random.uniform(0, 0.2))
        else:
            break

    return '\n\n'.join(p for p in parts if p)


def get_chapters(dir_url):
    parsed = urlparse(dir_url)
    base = f"{parsed.scheme}://{parsed.netloc}"
    book_id = parsed.path.strip('/').split('/')[0]

    html = fetch(dir_url)
    if not html:
        return [], ""

    soup = BeautifulSoup(html, 'html.parser')

    title = ""
    for sel in ['h1', '.book-name', '.title']:
        tag = soup.select_one(sel)
        if tag:
            title = re.sub(r'\s*(章節列表|章节列表|目录|目錄)\s*$', '', tag.get_text(strip=True))
            break
    if not title:
        t = soup.find('title')
        if t:
            title = re.sub(r'\s*[-|–—].*$', '', t.get_text(strip=True))

    ch_pat = re.compile(rf'/{re.escape(book_id)}/\d+_\d+\.html$')
    seen, chapters = set(), []
    for a in soup.find_all('a', href=True):
        if ch_pat.search(a['href']):
            url = urljoin(base, a['href'])
            if url not in seen:
                seen.add(url)
                chapters.append((a.get_text(strip=True), url))

    # 目录页通常倒序，按 URL 中数字正序排列
    chapters.sort(key=lambda x: [int(n) for n in re.findall(r'\d+', x[1].split('/')[-1])])
    return chapters, title


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('url')
    parser.add_argument('output', nargs='?', default='')
    parser.add_argument('--delay', type=float, default=0.8, help='章节间延迟(秒)')
    parser.add_argument('--start', type=int, default=1)
    parser.add_argument('--end', type=int, default=0)
    args = parser.parse_args()

    parsed = urlparse(args.url)
    base = f"{parsed.scheme}://{parsed.netloc}"

    print(f"获取目录: {args.url}")
    chapters, title = get_chapters(args.url)
    if not chapters:
        print("未找到章节，请检查 URL。")
        sys.exit(1)

    title_s = to_s(title) or "novel"
    print(f"书名: {title_s}  共 {len(chapters)} 章")

    out_file = args.output or f"{title_s}.txt"
    end = args.end if args.end > 0 else len(chapters)
    to_fetch = chapters[args.start - 1:end]

    with open(out_file, 'w', encoding='utf-8') as f:
        f.write(title_s + '\n' + '=' * 40 + '\n\n')

        for i, (ch_title, ch_url) in enumerate(to_fetch, start=args.start):
            print(f"[{i}/{len(chapters)}] {ch_title}")
            content = get_chapter_content(ch_url, base)
            ch_s = to_s(ch_title)
            body = to_s(content)

            # 去掉正文第一行重复的章节标题
            lines = body.splitlines()
            if lines and ch_s.replace(' ', '') in lines[0].replace(' ', ''):
                body = '\n'.join(lines[1:]).lstrip()

            f.write(f"\n\n{ch_s}\n\n{body}\n")
            f.flush()

            if i < end:
                time.sleep(args.delay + random.uniform(0, 0.4))

    print(f"\n完成！→ {out_file}")


if __name__ == '__main__':
    main()
