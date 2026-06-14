#!/usr/bin/env python3
"""
小说下载器 - 将 novel543.com 上的小说下载并转为简体中文 txt
用法: python novel_downloader.py <目录页URL> [输出文件名]
示例: python novel_downloader.py https://www.novel543.com/0802619327/dir 冒姓琅琊.txt
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
    _converter = opencc.OpenCC('t2s')
    def to_simplified(text):
        return _converter.convert(text)
except ImportError:
    print("警告: 未安装 opencc，将跳过繁简转换。安装: pip install opencc-python-reimplemented")
    def to_simplified(text):
        return text


AD_PATTERNS = [
    r'本书由.{0,20}首发',
    r'最新章节.{0,30}小说网',
    r'请记住.{0,40}地址',
    r'手机版.{0,20}阅读',
    r'www\.[a-zA-Z0-9\-]+\.[a-zA-Z]{2,}',
    r'http[s]?://[^\s　-鿿]+',
    r'稷下[書书]院',
    r'novel543',
    r'如果您喜欢.{0,30}请.{0,10}收藏',
    r'求收藏.{0,20}求推荐',
    r'求月票.{0,20}打赏',
    r'(最新|最快)更新.{0,20}章节',
    r'温馨提示.{0,30}阅读',
    r'阅读.{0,20}记得收藏',
    r'小说.{0,10}免费.{0,10}阅读',
    r'章节错误.{0,30}举报',
    r'举报本章',
    r'上一[章页].{0,5}下一[章页]',
]

AD_RE = re.compile('|'.join(AD_PATTERNS), re.IGNORECASE)

HEADERS = {
    'User-Agent': (
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) '
        'AppleWebKit/537.36 (KHTML, like Gecko) '
        'Chrome/120.0.0.0 Safari/537.36'
    ),
    'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8',
    'Accept-Language': 'zh-CN,zh-TW;q=0.9,zh;q=0.8',
}


def fetch(url, retries=4):
    delays = [2, 4, 8, 16]
    for attempt in range(retries):
        try:
            resp = requests.get(url, headers=HEADERS, timeout=15)
            resp.raise_for_status()
            resp.encoding = resp.apparent_encoding or 'utf-8'
            return resp.text
        except Exception as e:
            if attempt < retries - 1:
                wait = delays[attempt] + random.uniform(0, 1)
                print(f"    重试 {attempt+2}/{retries}（{e}），等待 {wait:.1f}s...")
                time.sleep(wait)
            else:
                print(f"    跳过 {url}：{e}")
                return None


def clean_text(raw):
    """清理广告和无关内容，返回干净的正文字符串。"""
    lines = raw.splitlines()
    result = []
    blank_count = 0
    for line in lines:
        line = AD_RE.sub('', line).strip()
        # 过滤掉只有标点或极短的行（广告残留）
        core = re.sub(r'[\s　\W]+', '', line)
        if len(core) < 2:
            blank_count += 1
            if blank_count == 1:
                result.append('')
        else:
            blank_count = 0
            result.append(line)
    return '\n'.join(result).strip()


def extract_body(soup):
    """从 BeautifulSoup 对象中提取正文。"""
    # 先移除导航、广告等干扰标签
    for tag in soup.find_all(['nav', 'header', 'footer', 'script', 'style', 'aside']):
        tag.decompose()
    for tag in soup.find_all(class_=re.compile(r'(nav|ad|banner|toolbar|sidebar|foot|header)', re.I)):
        tag.decompose()

    # 按优先级尝试正文容器
    for sel in ['#content', '#chapter-content', '.chapter-content',
                '#chaptercontent', '.read-content', '#nr1', '.nr-book-text',
                'article', '.content', '#text']:
        node = soup.select_one(sel)
        if node and len(node.get_text(strip=True)) > 100:
            return node.get_text('\n')

    # 兜底：取最长 <div>
    divs = soup.find_all('div')
    if divs:
        return max((d.get_text('\n') for d in divs), key=len)
    return ''


def get_page_info(soup):
    """返回 (current_page, total_pages, next_page_href)。"""
    h1 = soup.find('h1')
    current_page = total_pages = 1
    if h1:
        m = re.search(r'\((\d+)/(\d+)\)', h1.get_text())
        if m:
            current_page, total_pages = int(m.group(1)), int(m.group(2))

    next_href = None
    # 查找"下一页"链接（区别于"下一章"）
    for a in soup.find_all('a', href=True):
        text = a.get_text(strip=True)
        if re.search(r'下一[页頁]', text):
            next_href = a['href']
            break

    # 如果没有明确的"下一页"按钮，但还有更多页，尝试从URL推断
    return current_page, total_pages, next_href


def chapter_url_base(url):
    """从 URL 提取章节基础ID，如 '8096_3' -> '8096_3'。"""
    name = url.split('/')[-1].replace('.html', '')
    # 去掉末尾的 _数字（页码后缀），保留章节ID
    return re.sub(r'(_\d+)+$', '', name)


def get_chapter_content(first_url, base):
    """抓取一个章节的所有分页，合并为完整正文。"""
    parts = []
    current_url = first_url
    ch_base = chapter_url_base(first_url)

    while current_url:
        html = fetch(current_url)
        if not html:
            break

        soup = BeautifulSoup(html, 'html.parser')
        raw = extract_body(soup)
        cleaned = clean_text(raw)
        if cleaned:
            parts.append(cleaned)

        cur_page, total_pages, next_href = get_page_info(soup)

        if cur_page >= total_pages:
            break  # 已是最后一页

        if next_href:
            next_url = urljoin(base, next_href)
            # 确保 next 还属于同一章节（URL 含相同章节基础ID）
            if ch_base in next_url:
                current_url = next_url
                time.sleep(0.3 + random.uniform(0, 0.3))
                continue

        # 若找不到下一页链接但页数未到头，尝试按规律拼接
        next_page_num = cur_page + 1
        # URL 规律: 8096_3.html -> 8096_3_2.html -> 8096_3_3.html
        guessed = urljoin(base, f"/{urlparse(first_url).path.rsplit('/', 1)[0].lstrip('/')}/{ch_base}_{next_page_num}.html")
        probe = requests.head(guessed, headers=HEADERS, timeout=8)
        if probe.status_code == 200:
            current_url = guessed
            time.sleep(0.3)
        else:
            break

    return '\n\n'.join(parts)


def get_chapters(dir_url):
    """获取目录页中所有章节的 (章节名, URL) 列表，按正序排列。"""
    parsed = urlparse(dir_url)
    base = f"{parsed.scheme}://{parsed.netloc}"

    html = fetch(dir_url)
    if not html:
        return [], ""

    soup = BeautifulSoup(html, 'html.parser')

    # 书名：去掉"章節列表"等后缀
    title = ""
    for sel in ['h1', '.book-name', '.title']:
        tag = soup.select_one(sel)
        if tag:
            title = re.sub(r'\s*(章節列表|章节列表|目录|目錄)\s*$', '', tag.get_text(strip=True))
            break
    if not title:
        title_tag = soup.find('title')
        if title_tag:
            title = re.sub(r'\s*[-_|–—].*$', '', title_tag.get_text(strip=True))

    # 章节链接（novel543 格式: /BOOKID/NUM_NUM.html）
    chapters = []
    seen = set()
    book_id = parsed.path.strip('/').split('/')[0]
    ch_pattern = re.compile(rf'/{re.escape(book_id)}/\d+_\d+\.html$')

    for a in soup.find_all('a', href=True):
        href = a['href']
        if ch_pattern.search(href):
            full_url = urljoin(base, href)
            if full_url not in seen:
                seen.add(full_url)
                chapters.append((a.get_text(strip=True), full_url))

    # 目录页通常倒序（最新在前），按 URL 中的数字正序排列
    def sort_key(item):
        nums = re.findall(r'\d+', item[1].split('/')[-1])
        return [int(n) for n in nums]

    chapters.sort(key=sort_key)
    return chapters, title


def main():
    parser = argparse.ArgumentParser(description='小说下载器（支持分页、繁→简转换、去广告）')
    parser.add_argument('url', help='目录页 URL')
    parser.add_argument('output', nargs='?', default='', help='输出文件名（默认用书名）')
    parser.add_argument('--delay', type=float, default=1.0, help='章节间延迟(秒)，默认 1.0')
    parser.add_argument('--start', type=int, default=1, help='从第几章开始，默认 1')
    parser.add_argument('--end', type=int, default=0, help='到第几章结束，默认 0（全部）')
    args = parser.parse_args()

    parsed = urlparse(args.url)
    base = f"{parsed.scheme}://{parsed.netloc}"

    print(f"正在获取目录: {args.url}")
    chapters, title = get_chapters(args.url)

    if not chapters:
        print("未找到章节链接，请检查 URL 或页面结构。")
        sys.exit(1)

    title_simple = to_simplified(title) if title else "novel"
    print(f"书名: {title_simple}")
    print(f"共找到 {len(chapters)} 章")

    output_file = args.output or f"{title_simple}.txt"
    end = args.end if args.end > 0 else len(chapters)
    to_fetch = chapters[args.start - 1:end]

    with open(output_file, 'w', encoding='utf-8') as f:
        if title_simple:
            f.write(title_simple + '\n')
            f.write('=' * 40 + '\n\n')

        for i, (ch_title, ch_url) in enumerate(to_fetch, start=args.start):
            print(f"[{i}/{len(chapters)}] {ch_title}  {ch_url}")
            content = get_chapter_content(ch_url, base)
            ch_title_s = to_simplified(ch_title)
            content_s = to_simplified(content)

            # 去掉正文中重复的章节标题行
            first_line = content_s.splitlines()[0] if content_s else ''
            if ch_title_s and first_line and ch_title_s.replace(' ', '') in first_line.replace(' ', ''):
                content_s = '\n'.join(content_s.splitlines()[1:]).lstrip()

            f.write(f"\n\n{'—' * 6} {ch_title_s} {'—' * 6}\n\n")
            f.write(content_s)
            f.write('\n')
            f.flush()

            if i < end:
                time.sleep(args.delay + random.uniform(0, 0.5))

    print(f"\n完成！已保存至: {output_file}")


if __name__ == '__main__':
    main()
