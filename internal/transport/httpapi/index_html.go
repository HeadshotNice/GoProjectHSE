package httpapi

func indexHTML() string {
	return `<!doctype html>
<html lang="ru">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Кабинет документов</title>
  <style>
    :root { color-scheme: light; }
    body { margin: 0; font-family: Segoe UI, sans-serif; background: #f6f8fb; color: #1f2937; }
    .wrap { max-width: 960px; margin: 0 auto; padding: 24px; }
    h1 { margin: 0 0 16px; }
    .grid { display: grid; gap: 16px; grid-template-columns: repeat(auto-fit, minmax(280px, 1fr)); }
    .card { background: #fff; border: 1px solid #dbe3f0; border-radius: 14px; padding: 16px; box-shadow: 0 8px 20px rgba(0,0,0,.04); }
    label { display: block; margin: 10px 0 6px; font-size: 14px; font-weight: 600; }
    input, textarea, button { width: 100%; box-sizing: border-box; border-radius: 10px; border: 1px solid #cfd8e6; padding: 10px 12px; font-size: 14px; }
    textarea { min-height: 110px; resize: vertical; }
    button { margin-top: 12px; background: #1d4ed8; color: #fff; border: none; cursor: pointer; font-weight: 600; }
    button:hover { background: #1e40af; }
    .secondary { background: #0f766e; }
    .secondary:hover { background: #115e59; }
    .muted { color: #6b7280; font-size: 13px; margin-top: 8px; min-height: 18px; }
    .token { word-break: break-all; font-family: Consolas, monospace; background: #eef2ff; padding: 8px; border-radius: 8px; font-size: 12px; }
    table { width: 100%; border-collapse: collapse; margin-top: 12px; background: #fff; }
    th, td { text-align: left; padding: 10px; border-bottom: 1px solid #e5e7eb; font-size: 14px; vertical-align: top; }
    th { background: #f3f4f6; }
  </style>
</head>
<body>
  <div class="wrap">
    <h1>Портал документов</h1>
    <div class="grid">
      <section class="card">
        <h3>Регистрация</h3>
        <label for="regEmail">Email</label>
        <input id="regEmail" type="email" placeholder="you@example.com" />
        <label for="regPass">Пароль</label>
        <input id="regPass" type="password" placeholder="Минимум 1 символ" />
        <button id="regBtn">Зарегистрироваться</button>
        <div id="regMsg" class="muted"></div>
      </section>
      <section class="card">
        <h3>Логин</h3>
        <label for="loginEmail">Email</label>
        <input id="loginEmail" type="email" />
        <label for="loginPass">Пароль</label>
        <input id="loginPass" type="password" />
        <button id="loginBtn">Войти</button>
        <div id="loginMsg" class="muted"></div>
        <div id="tokenBox" class="token" hidden></div>
      </section>
    </div>

    <section class="card" style="margin-top:16px;">
      <h3>Отправка документа</h3>
      <label for="docTitle">Название документа</label>
      <input id="docTitle" type="text" placeholder="Например: Паспорт" />
      <label for="docContent">Содержимое / комментарий</label>
      <textarea id="docContent" placeholder="Текст или описание документа"></textarea>
      <button id="sendDocBtn" class="secondary">Отправить на проверку</button>
      <div id="docMsg" class="muted"></div>
    </section>

    <section class="card" style="margin-top:16px;">
      <h3>Статус проверки</h3>
      <button id="refreshBtn">Обновить список</button>
      <table>
        <thead><tr><th>ID</th><th>Документ</th><th>Статус</th><th>Дата</th></tr></thead>
        <tbody id="docRows"></tbody>
      </table>
    </section>
  </div>
  <script>
    const api = {
      async post(url, body, token) {
        const res = await fetch(url, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            ...(token ? { "Authorization": "Bearer " + token } : {})
          },
          body: JSON.stringify(body)
        });
        return parseJSON(res);
      },
      async get(url, token) {
        const res = await fetch(url, {
          headers: token ? { "Authorization": "Bearer " + token } : {}
        });
        return parseJSON(res);
      }
    };

    async function parseJSON(res) {
      const data = await res.json().catch(() => ({}));
      if (!res.ok) throw new Error(data.error || ("HTTP " + res.status));
      return data;
    }

    const tokenKey = "hse_token";
    const regMsg = document.getElementById("regMsg");
    const loginMsg = document.getElementById("loginMsg");
    const docMsg = document.getElementById("docMsg");
    const tokenBox = document.getElementById("tokenBox");
    const docRows = document.getElementById("docRows");

    function getToken() { return localStorage.getItem(tokenKey) || ""; }
    function setToken(token) {
      localStorage.setItem(tokenKey, token);
      tokenBox.hidden = !token;
      tokenBox.textContent = token ? ("JWT: " + token) : "";
    }
    setToken(getToken());

    document.getElementById("regBtn").onclick = async () => {
      regMsg.textContent = "";
      try {
        const email = document.getElementById("regEmail").value.trim();
        const password = document.getElementById("regPass").value;
        const data = await api.post("/auth/register", { email, password });
        regMsg.textContent = "Пользователь создан, id: " + data.user_id;
      } catch (e) {
        regMsg.textContent = e.message;
      }
    };

    document.getElementById("loginBtn").onclick = async () => {
      loginMsg.textContent = "";
      try {
        const email = document.getElementById("loginEmail").value.trim();
        const password = document.getElementById("loginPass").value;
        const data = await api.post("/auth/login", { email, password });
        setToken(data.token || "");
        loginMsg.textContent = "Успешный вход";
      } catch (e) {
        loginMsg.textContent = e.message;
      }
    };

    document.getElementById("sendDocBtn").onclick = async () => {
      docMsg.textContent = "";
      try {
        const title = document.getElementById("docTitle").value.trim();
        const content = document.getElementById("docContent").value.trim();
        const token = getToken();
        const data = await api.post("/documents", { title, content }, token);
        docMsg.textContent = "Документ отправлен, id: " + data.document_id;
        await refreshDocs();
      } catch (e) {
        docMsg.textContent = e.message;
      }
    };

    async function refreshDocs() {
      docRows.innerHTML = "";
      try {
        const token = getToken();
        const data = await api.get("/documents", token);
        for (const doc of data.documents || []) {
          const tr = document.createElement("tr");
          tr.innerHTML = "<td>" + doc.id + "</td><td><b>" + escapeHtml(doc.title) + "</b><br>" + escapeHtml(doc.content) + "</td><td>" + escapeHtml(doc.status) + "</td><td>" + new Date(doc.created_at).toLocaleString() + "</td>";
          docRows.appendChild(tr);
        }
      } catch (e) {
        docMsg.textContent = e.message;
      }
    }

    function escapeHtml(value) {
      return String(value).replace(/[&<>"']/g, (m) => ({ "&":"&amp;","<":"&lt;",">":"&gt;","\"":"&quot;","'":"&#39;" }[m]));
    }

    document.getElementById("refreshBtn").onclick = refreshDocs;
    refreshDocs();
  </script>
</body>
</html>`
}
