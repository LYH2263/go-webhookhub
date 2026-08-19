const $ = (id) => document.getElementById(id);

async function api(path, opts) {
  const res = await fetch(path, {
    headers: { "Content-Type": "application/json", ...(opts && opts.headers) },
    ...opts,
  });
  const text = await res.text();
  let data = {};
  try { data = text ? JSON.parse(text) : {}; } catch (_) { data = { error: text }; }
  if (!res.ok) throw new Error(data.error || res.statusText);
  return data;
}

async function refreshEndpoints() {
  const data = await api("/api/endpoints");
  const ul = $("ep-list");
  ul.innerHTML = "";
  (data.endpoints || []).forEach((e) => {
    const li = document.createElement("li");
    const meta = document.createElement("div");
    meta.innerHTML = `<strong>${e.id}</strong><div class="muted">${e.url}</div>
      <div class="muted">${(e.events || []).join(", ")} · secret ${e.secret_len}B · ${e.enabled ? "启用" : "禁用"}</div>`;
    const actions = document.createElement("div");
    const tog = document.createElement("button");
    tog.className = "ghost";
    tog.textContent = e.enabled ? "禁用" : "启用";
    tog.onclick = async (ev) => {
      ev.stopPropagation();
      const path = e.enabled ? `/api/endpoints/${encodeURIComponent(e.id)}/disable` : `/api/endpoints/${encodeURIComponent(e.id)}/enable`;
      await api(path, { method: "POST", body: "{}" });
      await refreshAll();
    };
    const del = document.createElement("button");
    del.className = "ghost";
    del.textContent = "删除";
    del.onclick = async (ev) => {
      ev.stopPropagation();
      await api(`/api/endpoints/${encodeURIComponent(e.id)}`, { method: "DELETE" });
      await refreshAll();
    };
    actions.append(tog, del);
    li.append(meta, actions);
    ul.appendChild(li);
  });
}

async function refreshDeliveries() {
  const data = await api("/api/deliveries?n=30");
  const ul = $("del-list");
  ul.innerHTML = "";
  (data.deliveries || []).slice().reverse().forEach((d) => {
    const li = document.createElement("li");
    const cls = d.ok ? "ok" : "bad";
    li.innerHTML = `<span class="${cls}">${d.ok ? "OK" : "FAIL"}</span>
      <span class="muted">${d.event}</span> → ${d.endpoint_id}
      <div class="muted">HTTP ${d.status_code} · ${d.attempts} 次 · ${d.error || ""}</div>`;
    ul.appendChild(li);
  });
}

async function refreshStats() {
  const s = await api("/api/stats");
  $("stats-line").textContent = `端点 ${s.endpoints}（启用 ${s.enabled}）· 投递 ${s.deliveries} · 成功 ${s.successes} · 失败 ${s.failures}`;
}

async function refreshAll() {
  await refreshEndpoints();
  await refreshDeliveries();
  await refreshStats();
}

$("ep-form").addEventListener("submit", async (e) => {
  e.preventDefault();
  const events = $("ep-events").value.split(",").map((x) => x.trim()).filter(Boolean);
  await api("/api/endpoints", {
    method: "POST",
    body: JSON.stringify({
      url: $("ep-url").value.trim(),
      secret: $("ep-secret").value,
      events,
      enabled: $("ep-enabled").checked,
    }),
  });
  $("ep-url").value = "";
  await refreshAll();
});

$("dispatch-form").addEventListener("submit", async (e) => {
  e.preventDefault();
  const data = await api("/api/dispatch", {
    method: "POST",
    body: JSON.stringify({
      event: $("ev-name").value.trim(),
      body: $("ev-body").value,
    }),
  });
  $("dispatch-out").textContent = JSON.stringify(data, null, 2);
  await refreshAll();
});

refreshAll().catch((err) => {
  $("stats-line").textContent = "无法连接 API：" + err.message;
});
