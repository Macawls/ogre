import satori from "satori";
import { Renderer } from "@takumi-rs/core";
import { container, text } from "@takumi-rs/helpers";
import { readdir, readFile } from "fs/promises";
import { join, basename } from "path";

const FIXTURES_DIR = join(import.meta.dir, "..", "fixtures");
const WIDTH = 1200;
const HEIGHT = 630;
const ITERATIONS = 20;

const fontRegular = Buffer.from(
  await fetch("https://cdn.jsdelivr.net/fontsource/fonts/inter@latest/latin-400-normal.ttf").then(r => r.arrayBuffer())
);
const fontBold = Buffer.from(
  await fetch("https://cdn.jsdelivr.net/fontsource/fonts/inter@latest/latin-700-normal.ttf").then(r => r.arrayBuffer())
);

const satoriFonts = [
  { name: "sans-serif", data: fontRegular, weight: 400 as const, style: "normal" as const },
  { name: "sans-serif", data: fontBold, weight: 700 as const, style: "normal" as const },
];

const takumiRenderer = new Renderer({
  fonts: [
    { name: "sans-serif", data: fontRegular, weight: 400, style: "normal" },
    { name: "sans-serif", data: fontBold, weight: 700, style: "normal" },
  ],
  loadDefaultFonts: false,
});

function parseStyle(s: string): Record<string, string | number> {
  const style: Record<string, string | number> = {};
  for (const decl of s.split(";")) {
    const idx = decl.indexOf(":");
    if (idx === -1) continue;
    const prop = decl.slice(0, idx).trim();
    const val = decl.slice(idx + 1).trim();
    const camel = prop.replace(/-([a-z])/g, (_, c) => c.toUpperCase());
    const num = Number(val);
    if (!isNaN(num) && val !== "") style[camel] = num;
    else if (val.endsWith("px")) { const n = Number(val.slice(0, -2)); style[camel] = isNaN(n) ? val : n; }
    else style[camel] = val;
  }
  return style;
}

type VNode = { type: string; props: { style?: any; children?: (VNode | string)[]; [k: string]: any } };

function htmlToVNode(html: string): VNode {
  const nodes = parseNodes(html.trim());
  if (nodes.length === 1 && typeof nodes[0] !== "string") return nodes[0];
  return { type: "div", props: { style: { display: "flex" }, children: nodes } };
}

function parseNodes(html: string): (VNode | string)[] {
  const nodes: (VNode | string)[] = [];
  let i = 0;
  while (i < html.length) {
    if (html[i] === "<") {
      if (html[i + 1] === "/") break;
      const tagEnd = html.indexOf(">", i);
      if (tagEnd === -1) break;
      const sc = html[tagEnd - 1] === "/";
      const tc = html.slice(i + 1, sc ? tagEnd - 1 : tagEnd).trim();
      const si = tc.indexOf(" ");
      const tag = si === -1 ? tc : tc.slice(0, si);
      const attrs = si === -1 ? "" : tc.slice(si + 1);
      const props: any = {};
      const sm = attrs.match(/style="([^"]*)"/);
      if (sm) props.style = parseStyle(sm[1]);
      if (sc || ["br", "img", "hr"].includes(tag)) { nodes.push({ type: tag, props }); i = tagEnd + 1; continue; }
      const is = tagEnd + 1;
      let d = 1, j = is;
      while (j < html.length && d > 0) {
        if (html[j] === "<") {
          if (html[j + 1] === "/") {
            const ce = html.indexOf(">", j);
            if (html.slice(j + 2, ce).trim() === tag) d--;
            if (d === 0) { const ch = parseNodes(html.slice(is, j)); if (ch.length) props.children = ch; nodes.push({ type: tag, props }); i = ce + 1; break; }
            j = ce + 1;
          } else { const ne = html.indexOf(">", j); const isc = html[ne - 1] === "/"; if (!isc) { const t2 = html.slice(j + 1, ne).trim(); const s2 = t2.indexOf(" "); const n2 = s2 === -1 ? t2 : t2.slice(0, s2); if (!["br", "img", "hr"].includes(n2)) d++; } j = ne + 1; }
        } else j++;
      }
      if (d > 0) i = html.length;
    } else {
      let end = html.indexOf("<", i);
      if (end === -1) end = html.length;
      const t = html.slice(i, end).trim();
      if (t) nodes.push(t);
      i = end;
    }
  }
  return nodes;
}

function vnodeToSatori(n: VNode | string): any {
  if (typeof n === "string") return n;
  const { children, ...rest } = n.props;
  return { type: n.type, props: { ...rest, children: children?.map(vnodeToSatori) } };
}

function vnodeToTakumi(n: VNode | string): any {
  if (typeof n === "string") return text({ style: {} }, n);
  const s = n.props.style || {};
  const ch = (n.props.children || []).map(vnodeToTakumi);
  return container({ style: s }, ch);
}

const files = (await readdir(FIXTURES_DIR)).filter(f => f.endsWith(".html")).sort();

console.log(`Benchmarking ${files.length} fixtures × ${ITERATIONS} iterations\n`);
console.log("| Fixture | Satori SVG (ms) | Takumi PNG (ms) | Takumi vs Satori |");
console.log("|---|---|---|---|");

for (const file of files) {
  const name = basename(file, ".html");
  const html = await readFile(join(FIXTURES_DIR, file), "utf-8");

  let satoriMs = -1;
  try {
    const el = vnodeToSatori(htmlToVNode(html));
    await satori(el, { width: WIDTH, height: HEIGHT, fonts: satoriFonts });
    const t0 = performance.now();
    for (let i = 0; i < ITERATIONS; i++) await satori(el, { width: WIDTH, height: HEIGHT, fonts: satoriFonts });
    satoriMs = (performance.now() - t0) / ITERATIONS;
  } catch { satoriMs = -1; }

  let takumiMs = -1;
  try {
    const tn = vnodeToTakumi(htmlToVNode(html));
    await takumiRenderer.render(tn, { width: WIDTH, height: HEIGHT, format: "png" });
    const t0 = performance.now();
    for (let i = 0; i < ITERATIONS; i++) await takumiRenderer.render(tn, { width: WIDTH, height: HEIGHT, format: "png" });
    takumiMs = (performance.now() - t0) / ITERATIONS;
  } catch { takumiMs = -1; }

  const ss = satoriMs > 0 ? satoriMs.toFixed(1) : "ERR";
  const ts = takumiMs > 0 ? takumiMs.toFixed(1) : "ERR";
  const cmp = satoriMs > 0 && takumiMs > 0 ? (takumiMs < satoriMs ? `${(satoriMs / takumiMs).toFixed(1)}x faster` : `${(takumiMs / satoriMs).toFixed(1)}x slower`) : "-";
  console.log(`| ${name} | ${ss} | ${ts} | ${cmp} |`);
}

console.log("\nDone.");
