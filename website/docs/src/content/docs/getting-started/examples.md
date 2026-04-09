---
title: Examples
description: Real-world input/output examples — blog posts, events, products, repos — rendered in SVG, PNG, and JPEG.
---

## Blog post

```html
<div class="flex flex-col w-full h-full bg-slate-950 p-16 justify-between">
  <div class="flex gap-3">
    <div class="bg-rose-500 text-white text-sm px-3 py-1 rounded-md font-semibold">Engineering</div>
    <div class="bg-slate-800 text-slate-300 text-sm px-3 py-1 rounded-md">12 min read</div>
  </div>
  <div class="flex flex-col mt-6">
    <div class="text-6xl font-extrabold text-white">Why We Migrated to Postgres</div>
    <div class="text-2xl text-slate-400 mt-4">Lessons from moving 2TB of data with zero downtime</div>
  </div>
  <div class="flex items-center gap-4 mt-auto">
    <div class="w-12 h-12 rounded-full bg-rose-600"></div>
    <div class="flex flex-col">
      <div class="text-white font-semibold">Sarah Chen</div>
      <div class="text-slate-500 text-sm">March 15, 2026</div>
    </div>
  </div>
</div>
```

| SVG | PNG | JPEG |
|-----|-----|------|
| ![SVG](/examples/ex-blog.svg) | ![PNG](/examples/ex-blog.png) | ![JPEG](/examples/ex-blog.jpg) |

## Conference

```html
<div style="display:flex;flex-direction:column;width:100%;height:100%;padding:64px;
  justify-content:space-between;
  background-image:linear-gradient(135deg,#0c0a1a 0%,#1a0533 40%,#2d0a4e 100%)">
  <div style="display:flex;flex-direction:column">
    <div style="color:#c084fc;font-size:14px;font-weight:700;letter-spacing:4px">
      JUNE 12–14, 2026 · BERLIN</div>
    <div style="color:white;font-size:68px;font-weight:800;margin-top:24px">
      Systems Conf Europe</div>
    <div style="color:#a78bfa;font-size:24px;margin-top:12px">
      Infrastructure, distributed systems, and platform engineering</div>
  </div>
  <div style="display:flex;align-items:center;gap:16px">
    <div style="background-color:#c084fc;color:#0c0a1a;font-size:16px;
      font-weight:700;padding:10px 24px;border-radius:8px">Get Tickets</div>
    <div style="color:#a78bfa;font-size:16px">Early bird pricing ends May 1</div>
  </div>
</div>
```

| SVG | PNG | JPEG |
|-----|-----|------|
| ![SVG](/examples/ex-event.svg) | ![PNG](/examples/ex-event.png) | ![JPEG](/examples/ex-event.jpg) |

## Product

```html
<div style="display:flex;flex-direction:column;width:100%;height:100%;
  background-color:#fafafa;padding:64px;justify-content:space-between">
  <div style="display:flex;align-items:center;gap:12px">
    <div style="width:44px;height:44px;border-radius:10px;
      background-image:linear-gradient(135deg,#3b82f6,#8b5cf6)"></div>
    <div style="font-size:22px;font-weight:700;color:#111827">Acme Analytics</div>
  </div>
  <div style="display:flex;flex-direction:column;margin-top:32px">
    <div style="font-size:52px;font-weight:800;color:#111827">
      Real-time dashboards for your entire team</div>
    <div style="font-size:22px;color:#6b7280;margin-top:16px">
      Track metrics, set alerts, and share insights. No SQL required.</div>
  </div>
  <div style="display:flex;align-items:center;gap:16px;margin-top:auto">
    <div style="background-color:#3b82f6;color:white;font-size:16px;
      font-weight:600;padding:10px 24px;border-radius:8px">Start Free Trial</div>
    <div style="color:#6b7280;font-size:16px">No credit card required</div>
  </div>
</div>
```

| SVG | PNG | JPEG |
|-----|-----|------|
| ![SVG](/examples/ex-product.svg) | ![PNG](/examples/ex-product.png) | ![JPEG](/examples/ex-product.jpg) |

## Repository

```html
<div style="display:flex;flex-direction:column;width:100%;height:100%;
  background-color:#0d1117;padding:64px;justify-content:space-between">
  <div style="display:flex;align-items:center;gap:14px">
    <div style="width:48px;height:48px;border-radius:24px;background-color:#238636"></div>
    <div style="font-size:24px;font-weight:600;color:#f0f6fc">vercel / next.js</div>
  </div>
  <div style="display:flex;flex-direction:column;margin-top:32px">
    <div style="font-size:48px;font-weight:700;color:#f0f6fc">
      The React Framework for the Web</div>
    <div style="font-size:20px;color:#8b949e;margin-top:16px;line-height:1.5">
      Used by some of the world's largest companies, Next.js enables you
      to create full-stack web applications.</div>
  </div>
  <div style="display:flex;align-items:center;gap:24px;margin-top:auto">
    <div style="color:#8b949e;font-size:15px">TypeScript</div>
    <div style="color:#8b949e;font-size:15px">MIT License</div>
    <div style="color:#f0f6fc;font-size:15px;font-weight:600">★ 128k</div>
    <div style="color:#8b949e;font-size:15px">Updated 2 hours ago</div>
  </div>
</div>
```

| SVG | PNG | JPEG |
|-----|-----|------|
| ![SVG](/examples/ex-repo.svg) | ![PNG](/examples/ex-repo.png) | ![JPEG](/examples/ex-repo.jpg) |

## Emoji

```html
<div style="display:flex;flex-direction:column;width:100%;height:100%;
  background-color:#18181b;padding:64px;justify-content:center">
  <div style="color:white;font-size:52px;font-weight:800">Ship faster 🚀</div>
  <div style="color:#a1a1aa;font-size:24px;margin-top:16px">
    Deploy with confidence ✨ Monitor in real-time 📊 Scale automatically ⚡</div>
</div>
```

| SVG | PNG | JPEG |
|-----|-----|------|
| ![SVG](/examples/ex-emoji.svg) | ![PNG](/examples/ex-emoji.png) | ![JPEG](/examples/ex-emoji.jpg) |

---

## Rendering

```bash
ogre --render card.html --output card.png --format png
```

```go
result, _ := ogre.Render(html, ogre.Options{
    Width:  1200,
    Height: 630,
    Format: ogre.FormatPNG,
})
os.WriteFile("card.png", result.Data, 0644)
```
