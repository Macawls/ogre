export const templates = {
  blog: `<div style="display:flex;flex-direction:column;width:100%;height:100%;padding:64px;background-image:linear-gradient(160deg,#0f172a 0%,#1e293b 50%,#0f172a 100%)">
  <div style="display:flex;align-items:center;justify-content:space-between">
    <div style="display:flex;gap:12px">
      <div style="background-color:#e11d48;color:white;font-size:14px;padding:6px 14px;border-radius:6px;font-weight:600">Engineering</div>
      <div style="background-color:rgba(255,255,255,0.08);color:#94a3b8;font-size:14px;padding:6px 14px;border-radius:6px">12 min read</div>
    </div>
    <div style="display:flex;align-items:center;gap:8px">
      <img src="data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Ccircle fill='none' cx='12' cy='12' r='10' stroke='%23475569' stroke-width='2'/%3E%3Cpath fill='%23475569' d='M12 6v6l4 2'/%3E%3C/svg%3E" style="width:18px;height:18px" />
      <div style="font-size:14px;color:#475569">March 15, 2026</div>
    </div>
  </div>
  <div style="display:flex;flex-direction:column;margin-top:auto">
    <div style="font-size:54px;font-weight:700;color:white;line-height:1.1">Why We Migrated to Postgres</div>
    <div style="font-size:20px;color:#64748b;margin-top:14px;line-height:1.4">Lessons from moving 2TB of data with zero downtime</div>
    <div style="display:flex;align-items:center;gap:16px;margin-top:32px">
      <img src="https://api.dicebear.com/9.x/lorelei/svg?seed=Sarah" style="width:48px;height:48px;border-radius:24px" />
      <div style="display:flex;flex-direction:column">
        <div style="color:white;font-size:16px;font-weight:600">Sarah Chen</div>
        <div style="color:#475569;font-size:14px">Staff Engineer at Veritas</div>
      </div>
    </div>
  </div>
</div>`,

  event: `<div style="display:flex;flex-direction:column;width:100%;height:100%;padding:64px;justify-content:space-between;background-image:linear-gradient(135deg,#0c0a1a 0%,#1a0533 40%,#2d0a4e 100%)">
  <div style="display:flex;flex-direction:column">
    <div style="display:flex;align-items:center;gap:10px">
      <img src="data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='%23c084fc' d='M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5a2.5 2.5 0 1 1 0-5 2.5 2.5 0 0 1 0 5z'/%3E%3C/svg%3E" style="width:18px;height:18px" />
      <div style="color:#c084fc;font-size:13px;font-weight:700;letter-spacing:5px">JUNE 12–14, 2026 · BERLIN</div>
    </div>
    <div style="color:white;font-size:72px;font-weight:800;margin-top:24px;line-height:1">Systems Conf</div>
    <div style="color:#a78bfa;font-size:22px;margin-top:16px;line-height:1.4">Infrastructure, distributed systems, and platform engineering</div>
  </div>
  <div style="display:flex;align-items:center;justify-content:space-between">
    <div style="display:flex;align-items:center;gap:12px">
      <div style="background-color:#c084fc;color:#0c0a1a;font-size:15px;font-weight:700;padding:12px 28px;border-radius:8px">Get Tickets →</div>
      <div style="border:1px solid rgba(192,132,252,0.3);color:#c084fc;font-size:15px;font-weight:600;padding:12px 28px;border-radius:8px">View Speakers</div>
    </div>
    <div style="display:flex;align-items:center;gap:16px">
      <div style="display:flex;align-items:center;gap:6px">
        <img src="data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='%239333ea' d='M16 11c1.66 0 2.99-1.34 2.99-3S17.66 5 16 5c-1.66 0-3 1.34-3 3s1.34 3 3 3zm-8 0c1.66 0 2.99-1.34 2.99-3S9.66 5 8 5C6.34 5 5 6.34 5 8s1.34 3 3 3zm0 2c-2.33 0-7 1.17-7 3.5V19h14v-2.5c0-2.33-4.67-3.5-7-3.5zm8 0c-.29 0-.62.02-.97.05 1.16.84 1.97 1.97 1.97 3.45V19h6v-2.5c0-2.33-4.67-3.5-7-3.5z'/%3E%3C/svg%3E" style="width:18px;height:18px" />
        <div style="color:#a78bfa;font-size:14px;font-weight:600">32 speakers</div>
      </div>
      <div style="display:flex;align-items:center;gap:6px">
        <img src="data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='%239333ea' d='M11.99 2C6.47 2 2 6.48 2 12s4.47 10 9.99 10C17.52 22 22 17.52 22 12S17.52 2 11.99 2zM12 20c-4.42 0-8-3.58-8-8s3.58-8 8-8 8 3.59 8 8-3.58 8-8 8zm.5-13H11v6l5.25 3.15.75-1.23-4.5-2.67z'/%3E%3C/svg%3E" style="width:18px;height:18px" />
        <div style="color:#a78bfa;font-size:14px;font-weight:600">3 days</div>
      </div>
    </div>
  </div>
</div>`,

  product: `<div style="display:flex;width:100%;height:100%;background-color:#ffffff">
  <div style="display:flex;flex-direction:column;width:50%;padding:56px 64px;justify-content:space-between">
    <div style="display:flex;align-items:center;gap:12px">
      <div style="width:40px;height:40px;border-radius:10px;background-image:linear-gradient(135deg,#3b82f6,#8b5cf6)"></div>
      <div style="font-size:18px;font-weight:700;color:#111827;letter-spacing:-0.3px">Acme Analytics</div>
    </div>
    <div style="display:flex;flex-direction:column">
      <div style="font-size:44px;font-weight:800;color:#111827;line-height:1.08;letter-spacing:-1px">Real-time dashboards for your entire team</div>
      <div style="font-size:17px;color:#6b7280;margin-top:14px;line-height:1.5">Track metrics, set alerts, and share insights across your organization.</div>
    </div>
    <div style="display:flex;align-items:center;gap:12px">
      <div style="background-image:linear-gradient(135deg,#3b82f6,#8b5cf6);color:white;font-size:15px;font-weight:600;padding:12px 28px;border-radius:8px">Start Free Trial</div>
      <div style="color:#6b7280;font-size:13px">No credit card required</div>
    </div>
  </div>
  <div style="display:flex;flex-direction:column;width:50%;background-image:linear-gradient(135deg,#eff6ff,#f5f3ff);padding:48px 40px;justify-content:center;gap:14px">
    <div style="display:flex;align-items:center;gap:14px;background-color:white;padding:18px 22px;border-radius:12px;box-shadow:0 1px 3px rgba(0,0,0,0.08)">
      <div style="width:36px;height:36px;border-radius:8px;background-color:#eff6ff;display:flex;align-items:center;justify-content:center">
        <img src="data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Crect fill='%233b82f6' x='4' y='13' width='4' height='8' rx='1'/%3E%3Crect fill='%233b82f6' x='10' y='5' width='4' height='16' rx='1'/%3E%3Crect fill='%233b82f6' x='16' y='9' width='4' height='12' rx='1'/%3E%3C/svg%3E" style="width:22px;height:22px" />
      </div>
      <div style="display:flex;flex-direction:column">
        <div style="font-size:12px;color:#6b7280;letter-spacing:0.5px">REVENUE</div>
        <div style="display:flex;align-items:center;gap:8px">
          <div style="font-size:22px;font-weight:700;color:#111827">$2.4M</div>
          <div style="font-size:12px;font-weight:600;color:#16a34a;background-color:#f0fdf4;padding:2px 8px;border-radius:4px">+12.5%</div>
        </div>
      </div>
    </div>
    <div style="display:flex;align-items:center;gap:14px;background-color:white;padding:18px 22px;border-radius:12px;box-shadow:0 1px 3px rgba(0,0,0,0.08)">
      <div style="width:36px;height:36px;border-radius:8px;background-color:#f5f3ff;display:flex;align-items:center;justify-content:center">
        <img src="data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='%238b5cf6' d='M16 11c1.66 0 2.99-1.34 2.99-3S17.66 5 16 5c-1.66 0-3 1.34-3 3s1.34 3 3 3zm-8 0c1.66 0 2.99-1.34 2.99-3S9.66 5 8 5C6.34 5 5 6.34 5 8s1.34 3 3 3zm0 2c-2.33 0-7 1.17-7 3.5V19h14v-2.5c0-2.33-4.67-3.5-7-3.5zm8 0c-.29 0-.62.02-.97.05 1.16.84 1.97 1.97 1.97 3.45V19h6v-2.5c0-2.33-4.67-3.5-7-3.5z'/%3E%3C/svg%3E" style="width:22px;height:22px" />
      </div>
      <div style="display:flex;flex-direction:column">
        <div style="font-size:12px;color:#6b7280;letter-spacing:0.5px">ACTIVE USERS</div>
        <div style="display:flex;align-items:center;gap:8px">
          <div style="font-size:22px;font-weight:700;color:#111827">14,892</div>
          <div style="font-size:12px;font-weight:600;color:#16a34a;background-color:#f0fdf4;padding:2px 8px;border-radius:4px">+8.1%</div>
        </div>
      </div>
    </div>
    <div style="display:flex;align-items:center;gap:14px;background-color:white;padding:18px 22px;border-radius:12px;box-shadow:0 1px 3px rgba(0,0,0,0.08)">
      <div style="width:36px;height:36px;border-radius:8px;background-color:#fefce8;display:flex;align-items:center;justify-content:center">
        <img src="data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='%23ca8a04' d='M13 2L3 14h9l-1 8 10-12h-9l1-8z'/%3E%3C/svg%3E" style="width:22px;height:22px" />
      </div>
      <div style="display:flex;flex-direction:column">
        <div style="font-size:12px;color:#6b7280;letter-spacing:0.5px">UPTIME</div>
        <div style="display:flex;align-items:center;gap:8px">
          <div style="font-size:22px;font-weight:700;color:#111827">99.98%</div>
          <div style="font-size:12px;font-weight:600;color:#16a34a;background-color:#f0fdf4;padding:2px 8px;border-radius:4px">+0.02%</div>
        </div>
      </div>
    </div>
  </div>
</div>`,

  repo: `<div style="display:flex;flex-direction:column;width:100%;height:100%;background-color:#0d1117;padding:56px 64px;justify-content:space-between">
  <div style="display:flex;align-items:center;justify-content:space-between">
    <div style="display:flex;align-items:center;gap:16px">
      <img src="https://github.com/vercel.png" style="width:44px;height:44px;border-radius:22px" />
      <div style="display:flex;align-items:center;gap:8px">
        <div style="font-size:20px;color:#58a6ff">vercel</div>
        <div style="font-size:20px;color:#8b949e">/</div>
        <div style="font-size:20px;font-weight:700;color:#f0f6fc">next.js</div>
      </div>
    </div>
    <div style="display:flex;align-items:center;gap:12px">
      <div style="display:flex;align-items:center;gap:6px;background-color:#21262d;padding:8px 14px;border-radius:6px;border:1px solid #30363d">
        <img src="data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='%23e3b341' d='M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01z'/%3E%3C/svg%3E" style="width:16px;height:16px" />
        <div style="color:#c9d1d9;font-size:14px;font-weight:500">128k</div>
      </div>
      <div style="display:flex;align-items:center;gap:6px;background-color:#21262d;padding:8px 14px;border-radius:6px;border:1px solid #30363d">
        <img src="data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='%238b949e' d='M6 3a3 3 0 1 0 0 6 3 3 0 0 0 0-6zm12 0a3 3 0 1 0 0 6 3 3 0 0 0 0-6zm-6 12a3 3 0 1 0 0 6 3 3 0 0 0 0-6zM6 9v1.5A4.5 4.5 0 0 0 10.5 15h3a4.5 4.5 0 0 0 4.5-4.5V9'/%3E%3C/svg%3E" style="width:16px;height:16px" />
        <div style="color:#c9d1d9;font-size:14px;font-weight:500">14.2k</div>
      </div>
    </div>
  </div>
  <div style="display:flex;flex-direction:column">
    <div style="font-size:44px;font-weight:700;color:#f0f6fc;line-height:1.15">The React Framework for the Web</div>
    <div style="font-size:18px;color:#8b949e;margin-top:14px;line-height:1.5">Used by some of the world's largest companies, Next.js enables you to create full-stack web applications.</div>
  </div>
  <div style="display:flex;align-items:center;gap:24px">
    <div style="display:flex;align-items:center;gap:8px">
      <div style="width:14px;height:14px;border-radius:7px;background-color:#3178c6"></div>
      <div style="color:#8b949e;font-size:15px">TypeScript</div>
    </div>
    <div style="color:#8b949e;font-size:15px">MIT License</div>
    <div style="display:flex;align-items:center;gap:6px">
      <div style="width:8px;height:8px;border-radius:4px;background-color:#3fb950"></div>
      <div style="color:#8b949e;font-size:15px">Updated 2 hours ago</div>
    </div>
  </div>
</div>`,

  emoji: `<div style="display:flex;flex-direction:column;width:100%;height:100%;background-image:linear-gradient(135deg,#18181b 0%,#27272a 100%);padding:64px;justify-content:space-between">
  <div style="display:flex;flex-direction:column;gap:8px">
    <div style="color:white;font-size:56px;font-weight:800;line-height:1.1">Ship faster 🚀</div>
    <div style="color:#a1a1aa;font-size:22px;margin-top:8px;line-height:1.5">Deploy with confidence ✨ Monitor in real-time 📊</div>
  </div>
  <div style="display:flex;gap:16px">
    <div style="display:flex;align-items:center;gap:10px;background-color:rgba(255,255,255,0.06);padding:16px 24px;border-radius:12px;border:1px solid rgba(255,255,255,0.08)">
      <div style="font-size:16px">🌍</div>
      <div style="display:flex;flex-direction:column">
        <div style="color:white;font-size:20px;font-weight:700">42 regions</div>
        <div style="color:#71717a;font-size:14px">Global edge network</div>
      </div>
    </div>
    <div style="display:flex;align-items:center;gap:10px;background-color:rgba(255,255,255,0.06);padding:16px 24px;border-radius:12px;border:1px solid rgba(255,255,255,0.08)">
      <div style="font-size:16px">⚡</div>
      <div style="display:flex;flex-direction:column">
        <div style="color:white;font-size:20px;font-weight:700">&lt;50ms</div>
        <div style="color:#71717a;font-size:14px">Average latency</div>
      </div>
    </div>
    <div style="display:flex;align-items:center;gap:10px;background-color:rgba(255,255,255,0.06);padding:16px 24px;border-radius:12px;border:1px solid rgba(255,255,255,0.08)">
      <div style="font-size:16px">🔒</div>
      <div style="display:flex;flex-direction:column">
        <div style="color:white;font-size:20px;font-weight:700">SOC 2</div>
        <div style="color:#71717a;font-size:14px">Certified</div>
      </div>
    </div>
  </div>
</div>`,

  rtl: `<div style="display:flex;flex-direction:column;width:100%;height:100%;padding:64px;justify-content:space-between;background-image:linear-gradient(135deg,#1a1a2e 0%,#16213e 50%,#0f3460 100%);direction:rtl;font-family:'Noto Sans Arabic',sans-serif">
  <div style="display:flex;align-items:center;justify-content:space-between">
    <div style="display:flex;align-items:center;gap:12px">
      <div style="background-color:#e2b714;color:#1a1a2e;font-size:14px;font-weight:700;padding:6px 16px;border-radius:6px">تقنية</div>
      <div style="background-color:rgba(255,255,255,0.08);color:#94a3b8;font-size:14px;padding:6px 14px;border-radius:6px">8 دقائق للقراءة</div>
    </div>
    <div style="font-size:14px;color:#475569">15 مارس 2026</div>
  </div>
  <div style="display:flex;flex-direction:column">
    <div style="font-size:52px;font-weight:800;color:white;line-height:1.2">لماذا انتقلنا إلى Kubernetes</div>
    <div style="font-size:20px;color:#64748b;margin-top:14px;line-height:1.5">دروس مستفادة من ترحيل 50 خدمة مصغرة إلى بنية تحتية جديدة</div>
  </div>
  <div style="display:flex;align-items:center;gap:24px">
    <div style="display:flex;align-items:center;gap:8px">
      <div style="width:14px;height:14px;border-radius:7px;background-color:#e2b714"></div>
      <div style="color:#94a3b8;font-size:15px">DevOps</div>
    </div>
    <div style="color:#64748b;font-size:15px">مصدر مفتوح</div>
    <div style="display:flex;align-items:center;gap:6px">
      <div style="width:8px;height:8px;border-radius:4px;background-color:#3fb950"></div>
      <div style="color:#64748b;font-size:15px">تم التحديث مؤخرًا</div>
    </div>
  </div>
</div>`,
  tailwind: `<div class="flex flex-col w-full h-full bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900 p-16 justify-between">
  <div class="flex items-center gap-4">
    <div class="w-12 h-12 rounded-xl bg-violet-500"></div>
    <div class="text-white text-2xl font-bold">Nebula</div>
    <div class="text-slate-500 text-base">Developer Platform</div>
  </div>
  <div class="flex flex-col gap-3">
    <div class="text-5xl font-extrabold text-white">Deploy anywhere in seconds</div>
    <div class="text-xl text-slate-400">Push your code. We handle the infrastructure, scaling, and monitoring so you can focus on shipping.</div>
  </div>
  <div class="flex items-center gap-6">
    <div class="bg-violet-500 text-white text-base font-semibold px-7 py-3 rounded-lg">Start Building</div>
    <div class="border border-slate-600 text-slate-300 text-base font-medium px-7 py-3 rounded-lg">View Pricing</div>
    <div class="flex items-center gap-2">
      <div class="w-2 h-2 rounded-full bg-emerald-400"></div>
      <div class="text-sm text-slate-500">99.99% uptime SLA</div>
    </div>
  </div>
</div>`,
} as const;

export type TemplateName = keyof typeof templates;
export const templateNames = Object.keys(templates) as TemplateName[];
