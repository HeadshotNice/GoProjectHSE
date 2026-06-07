package httpapi

func indexHTML() string {
	return `<!doctype html>
<html lang="ru">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Система проверки документов</title>
  <style>
    :root {
      --bg: #f4efe8;
      --bg-deep: #e7ddcf;
      --surface: rgba(255, 255, 255, 0.84);
      --surface-strong: rgba(255, 255, 255, 0.94);
      --text: #211c17;
      --muted: #716556;
      --line: rgba(74, 59, 43, 0.12);
      --accent: #ba5a31;
      --accent-2: #2f7555;
      --accent-soft: rgba(186, 90, 49, 0.12);
      --success-soft: rgba(47, 117, 85, 0.12);
      --danger-soft: rgba(152, 56, 56, 0.14);
      --shadow: 0 26px 70px rgba(50, 37, 23, 0.14);
      --radius-xl: 30px;
      --radius-lg: 20px;
      --radius-md: 16px;
      --radius-sm: 12px;
      --transition: 220ms ease;
    }

    * { box-sizing: border-box; }

    @keyframes fadeUp {
      from {
        opacity: 0;
        transform: translateY(18px);
      }
      to {
        opacity: 1;
        transform: translateY(0);
      }
    }

    @keyframes fadeIn {
      from { opacity: 0; }
      to { opacity: 1; }
    }

    body {
      margin: 0;
      min-height: 100vh;
      font-family: "Trebuchet MS", "Segoe UI", sans-serif;
      color: var(--text);
      background:
        radial-gradient(circle at top left, rgba(47, 117, 85, 0.14), transparent 26%),
        radial-gradient(circle at 85% 12%, rgba(186, 90, 49, 0.18), transparent 24%),
        linear-gradient(135deg, var(--bg) 0%, var(--bg-deep) 100%);
    }

    .page {
      max-width: 1180px;
      margin: 0 auto;
      padding: 24px 18px 40px;
    }

    .auth-layout,
    .shell-card,
    .panel {
      border: 1px solid var(--line);
      border-radius: var(--radius-xl);
      background: var(--surface);
      backdrop-filter: blur(18px);
      box-shadow: var(--shadow);
    }

    .auth-layout {
      display: grid;
      grid-template-columns: minmax(0, 1.15fr) minmax(330px, 430px);
      overflow: hidden;
      min-height: 560px;
      animation: fadeIn 420ms ease;
    }

    .hero {
      position: relative;
      display: flex;
      flex-direction: column;
      justify-content: center;
      padding: 40px 44px;
      color: #faf3eb;
      background:
        linear-gradient(180deg, rgba(255,255,255,0.14), rgba(255,255,255,0)),
        linear-gradient(145deg, #285d45 0%, #1d3127 50%, #171412 100%);
    }

    .hero::after {
      content: "";
      position: absolute;
      right: -70px;
      bottom: -70px;
      width: 280px;
      height: 280px;
      border-radius: 50%;
      background: radial-gradient(circle, rgba(186, 90, 49, 0.48), rgba(186, 90, 49, 0));
      pointer-events: none;
    }

    .badge {
      display: inline-flex;
      align-items: center;
      gap: 10px;
      padding: 10px 14px;
      border-radius: 999px;
      background: rgba(255, 255, 255, 0.1);
      font-size: 13px;
      letter-spacing: 0.08em;
      text-transform: uppercase;
    }

    h1 {
      margin: 22px 0 0;
      font-size: clamp(40px, 5.8vw, 66px);
      line-height: 0.96;
      letter-spacing: 0;
      max-width: 560px;
    }

    .hero p {
      max-width: 520px;
      margin: 0;
      color: rgba(250, 243, 235, 0.82);
      font-size: 18px;
      line-height: 1.65;
    }

    .story {
      display: grid;
      gap: 14px;
      margin-top: 34px;
    }

    .story-card {
      padding: 18px;
      border-radius: var(--radius-lg);
      background: rgba(255, 255, 255, 0.08);
      border: 1px solid rgba(255, 255, 255, 0.08);
      animation: fadeUp 420ms ease backwards;
    }

    .story-card:nth-child(2) { animation-delay: 70ms; }
    .story-card:nth-child(3) { animation-delay: 140ms; }

    .story-card strong {
      display: block;
      margin-bottom: 6px;
      font-size: 16px;
    }

    .story-card span {
      color: rgba(250, 243, 235, 0.74);
      line-height: 1.5;
      font-size: 14px;
    }

    .auth-side {
      display: flex;
      flex-direction: column;
      justify-content: center;
      padding: 26px;
      background: linear-gradient(180deg, rgba(255,255,255,0.92), rgba(255,255,255,0.76));
    }

    .tabs {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 8px;
      padding: 8px;
      border-radius: 18px;
      background: rgba(33, 28, 23, 0.06);
      margin-bottom: 24px;
    }

    button,
    input,
    textarea {
      font: inherit;
    }

    .tab-btn,
    .primary-btn,
    .secondary-btn,
    .logout-btn {
      appearance: none;
      border: 0;
      transition: transform var(--transition), box-shadow var(--transition), background var(--transition), opacity var(--transition);
    }

    .tab-btn {
      cursor: pointer;
      border-radius: 14px;
      padding: 14px 16px;
      background: transparent;
      color: var(--muted);
      font-weight: 700;
    }

    .tab-btn.active {
      background: var(--surface-strong);
      color: var(--text);
      box-shadow: 0 10px 24px rgba(33, 28, 23, 0.08);
    }

    .tab-btn:hover,
    .primary-btn:hover,
    .secondary-btn:hover,
    .logout-btn:hover {
      transform: translateY(-1px);
    }

    .auth-card,
    .cabinet[hidden] {
      display: none !important;
    }

    .auth-card.active,
    .cabinet.active {
      display: block !important;
    }

    .auth-card {
      animation: fadeUp 240ms ease;
    }

    .auth-head {
      margin-bottom: 20px;
    }

    .auth-head h2,
    .cabinet-header h2 {
      margin: 0 0 8px;
      font-size: 32px;
      line-height: 1.08;
    }

    .auth-head p,
    .soft-text {
      margin: 0;
      color: var(--muted);
      line-height: 1.56;
    }

    .field {
      margin-bottom: 16px;
    }

    label {
      display: block;
      margin-bottom: 8px;
      font-size: 14px;
      font-weight: 700;
    }

    input,
    textarea {
      width: 100%;
      border: 1px solid rgba(74, 59, 43, 0.14);
      border-radius: var(--radius-md);
      background: rgba(255, 255, 255, 0.92);
      color: var(--text);
      padding: 14px 16px;
      outline: none;
      transition: border-color var(--transition), box-shadow var(--transition), background var(--transition);
    }

    input:focus,
    textarea:focus {
      border-color: rgba(186, 90, 49, 0.72);
      box-shadow: 0 0 0 4px rgba(186, 90, 49, 0.12);
      background: #fff;
    }

    textarea {
      min-height: 160px;
      resize: vertical;
      line-height: 1.5;
    }

    .hint {
      margin-top: 8px;
      color: var(--muted);
      font-size: 13px;
      line-height: 1.45;
    }

    .primary-btn,
    .secondary-btn,
    .logout-btn {
      cursor: pointer;
      border-radius: var(--radius-md);
      padding: 15px 18px;
      font-weight: 700;
    }

    .primary-btn {
      width: 100%;
      color: #fffaf5;
      background: linear-gradient(135deg, var(--accent) 0%, #da7e3a 100%);
      box-shadow: 0 16px 30px rgba(186, 90, 49, 0.24);
    }

    .secondary-btn,
    .logout-btn {
      background: rgba(33, 28, 23, 0.07);
      color: var(--text);
    }

    .notice {
      min-height: 22px;
      margin-top: 14px;
      font-size: 14px;
      color: var(--muted);
      transition: color var(--transition);
    }

    .notice.success { color: var(--accent-2); }
    .notice.error { color: #8e3535; }

    .cabinet {
      animation: fadeIn 280ms ease;
    }

    .shell-card {
      padding: 24px;
    }

    .cabinet-top {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      gap: 18px;
      margin-bottom: 20px;
    }

    .identity {
      display: inline-flex;
      align-items: center;
      gap: 8px;
      padding: 8px 12px;
      border-radius: 999px;
      background: var(--accent-soft);
      color: var(--accent);
      font-size: 13px;
      font-weight: 700;
      margin-bottom: 12px;
    }

    .cabinet-grid {
      display: grid;
      grid-template-columns: minmax(0, 0.94fr) minmax(0, 1.14fr);
      gap: 18px;
    }

    .panel {
      padding: 22px;
      background: var(--surface-strong);
      animation: fadeUp 300ms ease backwards;
    }

    .panel:nth-child(2) { animation-delay: 70ms; }

    .panel h3 {
      margin: 0 0 8px;
      font-size: 22px;
      line-height: 1.1;
    }

    .actions {
      display: flex;
      gap: 12px;
      margin: 18px 0 0;
    }

    .actions .secondary-btn {
      min-width: 190px;
    }

    .table-wrap {
      margin-top: 18px;
      border-radius: 18px;
      overflow: hidden;
      border: 1px solid rgba(74, 59, 43, 0.08);
      background: rgba(255, 255, 255, 0.8);
    }

    table {
      width: 100%;
      border-collapse: collapse;
    }

    th,
    td {
      text-align: left;
      padding: 14px 12px;
      border-bottom: 1px solid rgba(74, 59, 43, 0.08);
      vertical-align: top;
      font-size: 14px;
    }

    th {
      background: rgba(33, 28, 23, 0.04);
      color: var(--muted);
      font-size: 12px;
      text-transform: uppercase;
      letter-spacing: 0.08em;
    }

    tr:last-child td {
      border-bottom: 0;
    }

    .doc-title {
      margin-bottom: 4px;
      font-weight: 700;
    }

    .doc-body {
      color: var(--muted);
      line-height: 1.48;
      white-space: pre-wrap;
      word-break: break-word;
    }

    .pill {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      padding: 8px 12px;
      border-radius: 999px;
      font-size: 12px;
      font-weight: 700;
      letter-spacing: 0.06em;
      text-transform: uppercase;
      background: rgba(33, 28, 23, 0.08);
    }

    .pill.pending_review {
      background: var(--accent-soft);
      color: var(--accent);
    }

    .pill.in_review {
      background: rgba(52, 121, 90, 0.14);
      color: #22573d;
    }

    .pill.approved {
      background: var(--success-soft);
      color: var(--accent-2);
    }

    .pill.rejected {
      background: var(--danger-soft);
      color: #8c3232;
    }

    .empty {
      padding: 22px;
      border-radius: var(--radius-lg);
      background: rgba(33, 28, 23, 0.04);
      color: var(--muted);
      text-align: center;
      margin-top: 18px;
      animation: fadeIn 240ms ease;
    }

    @media (max-width: 980px) {
      .auth-layout {
        grid-template-columns: 1fr;
      }

      .hero,
      .auth-side {
        padding: 28px 22px;
      }

      .cabinet-grid {
        grid-template-columns: 1fr;
      }
    }

    @media (max-width: 640px) {
      .page {
        padding: 14px 12px 28px;
      }

      .auth-layout,
      .shell-card,
      .panel {
        border-radius: 24px;
      }

      .cabinet-top,
      .actions {
        flex-direction: column;
        align-items: stretch;
      }

      .actions .secondary-btn,
      .logout-btn {
        width: 100%;
      }

      h1 {
        font-size: 38px;
      }

      .auth-head h2,
      .cabinet-header h2 {
        font-size: 28px;
      }
    }
  </style>
</head>
<body>
  <div class="page">
    <section class="auth-layout" id="authView">
      <div class="hero">
        <div class="badge">Secure Review Workspace</div>
        <h1>Проверка документов</h1>
      </div>

      <div class="auth-side">
        <div class="tabs">
          <button type="button" class="tab-btn active" id="loginTabBtn">Вход</button>
          <button type="button" class="tab-btn" id="registerTabBtn">Регистрация</button>
        </div>

        <section class="auth-card active" id="loginCard">
          <div class="auth-head">
            <h2>Войти в кабинет</h2>
            <p>Введите email и пароль, чтобы продолжить работу с документами.</p>
          </div>

          <div class="field">
            <label for="loginEmail">Email</label>
            <input id="loginEmail" type="email" placeholder="you@example.com" autocomplete="email" />
          </div>

          <div class="field">
            <label for="loginPass">Пароль</label>
            <input id="loginPass" type="password" placeholder="Минимум 6 символов" autocomplete="current-password" />
          </div>

          <button type="button" class="primary-btn" id="loginBtn">Войти</button>
          <div class="notice" id="loginMsg"></div>
        </section>

        <section class="auth-card" id="registerCard">
          <div class="auth-head">
            <h2>Создать аккаунт</h2>
            <p>Используйте email и пароль не короче 6 символов. Можно использовать буквы, цифры и символы.</p>
          </div>

          <div class="field">
            <label for="regEmail">Email</label>
            <input id="regEmail" type="email" placeholder="you@example.com" autocomplete="email" />
          </div>

          <div class="field">
            <label for="regPass">Пароль</label>
            <input id="regPass" type="password" placeholder="Например: Secure123" autocomplete="new-password" />
            <div class="hint">Минимум 6 символов. Подойдёт любой удобный надёжный пароль.</div>
          </div>

          <button type="button" class="primary-btn" id="regBtn">Зарегистрироваться</button>
          <div class="notice" id="regMsg"></div>
        </section>
      </div>
    </section>

    <section class="shell-card cabinet" id="cabinetView" hidden>
      <div class="cabinet-top">
        <div class="cabinet-header">
          <div class="identity" id="userBadge">Аккаунт</div>
          <h2>Ваш кабинет</h2>
          <p class="soft-text">Здесь вы отправляете документы на проверку и следите за их текущим состоянием.</p>
        </div>
        <button type="button" class="logout-btn" id="logoutBtn">Выйти</button>
      </div>

      <div class="cabinet-grid">
        <section class="panel">
          <h3>Отправка документа</h3>
          <p class="soft-text">После отправки документ появится в списке справа и начнёт проходить этапы проверки.</p>

          <div class="field">
            <label for="docTitle">Название документа</label>
            <input id="docTitle" type="text" placeholder="Например: Паспорт, справка, договор" />
          </div>

          <div class="field">
            <label for="docContent">Описание или содержимое</label>
            <textarea id="docContent" placeholder="Кратко опишите документ или вставьте текст, который нужно проверить"></textarea>
          </div>

          <button type="button" class="primary-btn" id="sendDocBtn">Отправить на проверку</button>
          <div class="notice" id="docMsg"></div>
        </section>

        <section class="panel">
          <h3>Статусы проверки</h3>
          <p class="soft-text">Список ниже показывает только ваши документы и обновляется по запросу.</p>

          <div class="actions">
            <button type="button" class="secondary-btn" id="refreshBtn">Обновить список</button>
          </div>

          <div class="empty" id="emptyState" hidden>Документов пока нет. Отправьте первый документ, и он появится здесь.</div>

          <div class="table-wrap" id="tableWrap" hidden>
            <table>
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Документ</th>
                  <th>Статус</th>
                  <th>Обновлено</th>
                </tr>
              </thead>
              <tbody id="docRows"></tbody>
            </table>
          </div>
        </section>
      </div>
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
    const emailKey = "hse_email";

    const authView = document.getElementById("authView");
    const cabinetView = document.getElementById("cabinetView");
    const loginCard = document.getElementById("loginCard");
    const registerCard = document.getElementById("registerCard");
    const loginTabBtn = document.getElementById("loginTabBtn");
    const registerTabBtn = document.getElementById("registerTabBtn");
    const regMsg = document.getElementById("regMsg");
    const loginMsg = document.getElementById("loginMsg");
    const docMsg = document.getElementById("docMsg");
    const userBadge = document.getElementById("userBadge");
    const docRows = document.getElementById("docRows");
    const tableWrap = document.getElementById("tableWrap");
    const emptyState = document.getElementById("emptyState");

    function getToken() {
      return localStorage.getItem(tokenKey) || "";
    }

    function getEmail() {
      return localStorage.getItem(emailKey) || "";
    }

    function setSession(token, email) {
      if (token) {
        localStorage.setItem(tokenKey, token);
      } else {
        localStorage.removeItem(tokenKey);
      }

      if (email) {
        localStorage.setItem(emailKey, email);
      } else {
        localStorage.removeItem(emailKey);
      }
    }

    function setNotice(node, text, kind) {
      node.textContent = text || "";
      node.className = "notice" + (kind ? " " + kind : "");
    }

    function setMode(mode) {
      const loginActive = mode === "login";
      loginCard.classList.toggle("active", loginActive);
      registerCard.classList.toggle("active", !loginActive);
      loginTabBtn.classList.toggle("active", loginActive);
      registerTabBtn.classList.toggle("active", !loginActive);
      setNotice(regMsg, "", "");
      setNotice(loginMsg, "", "");
    }

    function renderLayout() {
      const hasToken = Boolean(getToken());
      authView.hidden = hasToken;
      cabinetView.hidden = !hasToken;
      cabinetView.classList.toggle("active", hasToken);
      userBadge.textContent = getEmail() ? ("Аккаунт: " + getEmail()) : "Аккаунт";
      if (hasToken) {
        refreshDocs();
      }
    }

    function isPasswordValid(password) {
      return password.trim().length >= 6;
    }

    function statusLabel(status) {
      const map = {
        pending_review: "Ожидает проверки",
        in_review: "На проверке",
        approved: "Одобрен",
        rejected: "Отклонён"
      };
      return map[status] || status;
    }

    function escapeHtml(value) {
      return String(value).replace(/[&<>"']/g, function (m) {
        return ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", "\"": "&quot;", "'": "&#39;" })[m];
      });
    }

    async function refreshDocs() {
      docRows.innerHTML = "";
      tableWrap.hidden = true;
      emptyState.hidden = true;

      try {
        const data = await api.get("/documents", getToken());
        const docs = data.documents || [];

        if (!docs.length) {
          emptyState.hidden = false;
          return;
        }

        for (const doc of docs) {
          const tr = document.createElement("tr");
          tr.innerHTML =
            "<td>" + doc.id + "</td>" +
            "<td><div class='doc-title'>" + escapeHtml(doc.title) + "</div><div class='doc-body'>" + escapeHtml(doc.content) + "</div></td>" +
            "<td><span class='pill " + escapeHtml(doc.status) + "'>" + escapeHtml(statusLabel(doc.status)) + "</span></td>" +
            "<td>" + new Date(doc.updated_at || doc.created_at).toLocaleString() + "</td>";
          docRows.appendChild(tr);
        }

        tableWrap.hidden = false;
      } catch (e) {
        setNotice(docMsg, e.message, "error");
      }
    }

    loginTabBtn.onclick = function () { setMode("login"); };
    registerTabBtn.onclick = function () { setMode("register"); };

    document.getElementById("regBtn").onclick = async function () {
      const email = document.getElementById("regEmail").value.trim();
      const password = document.getElementById("regPass").value;
      setNotice(regMsg, "", "");

      if (!isPasswordValid(password)) {
        setNotice(regMsg, "Пароль должен содержать минимум 6 символов.", "error");
        return;
      }

      try {
        await api.post("/auth/register", { email: email, password: password });
        setNotice(regMsg, "Аккаунт создан. Теперь можно войти.", "success");
        document.getElementById("loginEmail").value = email;
        document.getElementById("loginPass").value = password;
        setMode("login");
      } catch (e) {
        setNotice(regMsg, e.message, "error");
      }
    };

    document.getElementById("loginBtn").onclick = async function () {
      const email = document.getElementById("loginEmail").value.trim();
      const password = document.getElementById("loginPass").value;
      setNotice(loginMsg, "", "");

      try {
        const data = await api.post("/auth/login", { email: email, password: password });
        setSession(data.token || "", email);
        renderLayout();
      } catch (e) {
        setNotice(loginMsg, e.message, "error");
      }
    };

    document.getElementById("sendDocBtn").onclick = async function () {
      const title = document.getElementById("docTitle").value.trim();
      const content = document.getElementById("docContent").value.trim();
      setNotice(docMsg, "", "");

      try {
        const data = await api.post("/documents", { title: title, content: content }, getToken());
        setNotice(docMsg, "Документ отправлен на проверку. ID: " + data.document_id, "success");
        document.getElementById("docTitle").value = "";
        document.getElementById("docContent").value = "";
        refreshDocs();
      } catch (e) {
        setNotice(docMsg, e.message, "error");
      }
    };

    document.getElementById("refreshBtn").onclick = refreshDocs;

    document.getElementById("logoutBtn").onclick = function () {
      setSession("", "");
      setNotice(docMsg, "", "");
      docRows.innerHTML = "";
      tableWrap.hidden = true;
      emptyState.hidden = false;
      setMode("login");
      renderLayout();
    };

    setMode("login");
    renderLayout();
  </script>
</body>
</html>`
}
