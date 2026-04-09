import satori from "satori";
import { Resvg } from "@resvg/resvg-js";
import { readdir, readFile, writeFile, mkdir } from "fs/promises";
import { join, basename } from "path";

import { Renderer } from "@takumi-rs/core";
import { container, text } from "@takumi-rs/helpers";

let takumiRenderer: InstanceType<typeof Renderer> | null = null;

const FIXTURES_DIR = join(import.meta.dir, "..", "fixtures");
const REFERENCE_DIR = join(import.meta.dir, "..", "reference");
const WIDTH = 1200;
const HEIGHT = 630;

const fontRegular = await fetch(
  "https://cdn.jsdelivr.net/fontsource/fonts/inter@latest/latin-400-normal.ttf"
).then((r) => r.arrayBuffer());

const fontBold = await fetch(
  "https://cdn.jsdelivr.net/fontsource/fonts/inter@latest/latin-700-normal.ttf"
).then((r) => r.arrayBuffer());

await mkdir(REFERENCE_DIR, { recursive: true });

try {
  takumiRenderer = new Renderer({
    fonts: [
      { name: "sans-serif", data: Buffer.from(fontRegular), weight: 400, style: "normal" },
      { name: "sans-serif", data: Buffer.from(fontBold), weight: 700, style: "normal" },
    ],
    loadDefaultFonts: false,
  });
  console.log("Takumi renderer initialized");
} catch (e: any) {
  console.log("Takumi not available:", e.message);
}

function vnodeToTakumi(n: any): any {
  if (typeof n === "string") return text({ style: {} }, n);
  const s = n.props?.style || {};
  const ch = (n.props?.children || []).map(vnodeToTakumi);
  return container({ style: s }, ch);
}

function parseStyle(styleStr: string): Record<string, string | number> {
  const style: Record<string, string | number> = {};
  for (const decl of styleStr.split(";")) {
    const idx = decl.indexOf(":");
    if (idx === -1) continue;
    const prop = decl.slice(0, idx).trim();
    const val = decl.slice(idx + 1).trim();
    const camel = prop.replace(/-([a-z])/g, (_, c) => c.toUpperCase());
    const num = Number(val);
    if (!isNaN(num) && val !== "") {
      style[camel] = num;
    } else if (val.endsWith("px")) {
      const n = Number(val.slice(0, -2));
      if (!isNaN(n)) {
        style[camel] = n;
      } else {
        style[camel] = val;
      }
    } else {
      style[camel] = val;
    }
  }
  return style;
}

type VNode = {
  type: string;
  props: {
    style?: Record<string, string | number>;
    children?: (VNode | string)[];
    [key: string]: any;
  };
};

function htmlToVNode(html: string): VNode {
  html = html.trim();
  const nodes = parseNodes(html);
  if (nodes.length === 1) return nodes[0] as VNode;
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
      const selfClosing = html[tagEnd - 1] === "/";
      const tagContent = html.slice(i + 1, selfClosing ? tagEnd - 1 : tagEnd).trim();
      const spaceIdx = tagContent.indexOf(" ");
      const tag = spaceIdx === -1 ? tagContent : tagContent.slice(0, spaceIdx);
      const attrsStr = spaceIdx === -1 ? "" : tagContent.slice(spaceIdx + 1);

      const props: Record<string, any> = {};
      const styleMatch = attrsStr.match(/style="([^"]*)"/);
      if (styleMatch) {
        props.style = parseStyle(styleMatch[1]);
      }

      if (selfClosing || ["br", "img", "hr", "input"].includes(tag)) {
        nodes.push({ type: tag, props });
        i = tagEnd + 1;
      } else {
        const innerStart = tagEnd + 1;
        let depth = 1;
        let j = innerStart;
        while (j < html.length && depth > 0) {
          if (html[j] === "<") {
            if (html[j + 1] === "/") {
              const closeEnd = html.indexOf(">", j);
              const closeTag = html.slice(j + 2, closeEnd).trim();
              if (closeTag === tag) depth--;
              if (depth === 0) {
                const inner = html.slice(innerStart, j);
                const children = parseNodes(inner);
                if (children.length > 0) {
                  props.children = children;
                }
                nodes.push({ type: tag, props });
                i = closeEnd + 1;
                break;
              }
              j = closeEnd + 1;
            } else {
              const nextClose = html.indexOf(">", j);
              const isSelfClose = html[nextClose - 1] === "/";
              if (!isSelfClose) {
                const tc = html.slice(j + 1, nextClose).trim();
                const si = tc.indexOf(" ");
                const nt = si === -1 ? tc : tc.slice(0, si);
                if (!["br", "img", "hr", "input"].includes(nt)) {
                  depth++;
                }
              }
              j = nextClose + 1;
            }
          } else {
            j++;
          }
        }
        if (depth > 0) {
          i = html.length;
        }
      }
    } else {
      let end = html.indexOf("<", i);
      if (end === -1) end = html.length;
      const text = html.slice(i, end).trim();
      if (text) nodes.push(text);
      i = end;
    }
  }
  return nodes;
}

function vnodeToSatori(node: VNode | string): any {
  if (typeof node === "string") return node;
  const { type, props } = node;
  const { children, ...rest } = props;
  const satoriChildren = children?.map(vnodeToSatori);
  return {
    type,
    props: {
      ...rest,
      children: satoriChildren,
    },
  };
}

const files = (await readdir(FIXTURES_DIR)).filter((f) => f.endsWith(".html")).sort();

console.log(`Generating ${files.length} reference images...`);

for (const file of files) {
  const name = basename(file, ".html");
  const htmlContent = await readFile(join(FIXTURES_DIR, file), "utf-8");

  const vnode = htmlToVNode(htmlContent);

  try {
    const element = vnodeToSatori(vnode);

    const svg = await satori(element, {
      width: WIDTH,
      height: HEIGHT,
      fonts: [
        { name: "sans-serif", data: fontRegular, weight: 400, style: "normal" as const },
        { name: "sans-serif", data: fontBold, weight: 700, style: "normal" as const },
      ],
    });

    await writeFile(join(REFERENCE_DIR, `${name}.svg`), svg);

    const resvg = new Resvg(svg, {
      fitTo: { mode: "width" as const, value: WIDTH },
    });
    const png = resvg.render().asPng();
    await writeFile(join(REFERENCE_DIR, `${name}.png`), png);

    console.log(`  ✓ ${name} (satori)`);
  } catch (err: any) {
    console.error(`  ✗ ${name} (satori): ${err.message}`);
  }

  if (takumiRenderer) {
    try {
      const takumiNode = vnodeToTakumi(vnode);
      const takumiPng = await takumiRenderer.render(takumiNode, {
        width: WIDTH,
        height: HEIGHT,
        format: "png" as const,
      });
      await writeFile(join(REFERENCE_DIR, `${name}.takumi.png`), Buffer.from(takumiPng));
      console.log(`  ✓ ${name} (takumi)`);
    } catch (err: any) {
      console.error(`  ✗ ${name} (takumi): ${err.message}`);
    }
  }
}

console.log("Done.");
