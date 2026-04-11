import satori from "satori";
import { Resvg } from "@resvg/resvg-js";
import { Renderer } from "@takumi-rs/core";
import { container, text } from "@takumi-rs/helpers";

const WIDTH = 1200;
const HEIGHT = 630;

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

let takumiRenderer: InstanceType<typeof Renderer> | null = null;
try {
  takumiRenderer = new Renderer({
    fonts: [
      { name: "sans-serif", data: fontRegular, weight: 400, style: "normal" },
      { name: "sans-serif", data: fontBold, weight: 700, style: "normal" },
    ],
    loadDefaultFonts: false,
  });
} catch {}

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
      const is2 = tagEnd + 1;
      let d = 1, j = is2;
      while (j < html.length && d > 0) {
        if (html[j] === "<") {
          if (html[j + 1] === "/") {
            const ce = html.indexOf(">", j);
            if (html.slice(j + 2, ce).trim() === tag) d--;
            if (d === 0) { const ch = parseNodes(html.slice(is2, j)); if (ch.length) props.children = ch; nodes.push({ type: tag, props }); i = ce + 1; break; }
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

function htmlToVNode(html: string): VNode {
  const nodes = parseNodes(html.trim());
  if (nodes.length === 1 && typeof nodes[0] !== "string") return nodes[0];
  return { type: "div", props: { style: { display: "flex" }, children: nodes } };
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

const html = await Bun.stdin.text();
const vnode = htmlToVNode(html);
const result: any = { satori: null, takumi: null };

try {
  const el = vnodeToSatori(vnode);
  const t0 = performance.now();
  const svg = await satori(el, { width: WIDTH, height: HEIGHT, fonts: satoriFonts });
  const satoriMs = performance.now() - t0;

  const resvg = new Resvg(svg, { fitTo: { mode: "width" as const, value: WIDTH } });
  const png = resvg.render().asPng();

  result.satori = {
    svg: Buffer.from(svg).toString("base64"),
    png: Buffer.from(png).toString("base64"),
    ms: Math.round(satoriMs * 100) / 100,
  };
} catch (e: any) {
  result.satori = { error: e.message };
}

if (takumiRenderer) {
  try {
    const tn = vnodeToTakumi(vnode);
    const t0 = performance.now();
    const png = await takumiRenderer.render(tn, { width: WIDTH, height: HEIGHT, format: "png" as const });
    const takumiMs = performance.now() - t0;

    result.takumi = {
      png: Buffer.from(png).toString("base64"),
      ms: Math.round(takumiMs * 100) / 100,
    };
  } catch (e: any) {
    result.takumi = { error: e.message };
  }
}

process.stdout.write(JSON.stringify(result));
