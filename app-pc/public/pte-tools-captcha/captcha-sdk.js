/* PTE Captcha SDK: browser/WebView entry. Mini-program and native clients use the same REST contract. */
(function (global) {
  "use strict";
  const text = { "zh-CN": { retry: "验证失败，请重新操作", refreshed: "验证失败，请重新操作", refreshing: "正在刷新验证…", refresh: "换一张", slider: "拖动滑块完成拼图", restore: "拖动图片条完成还原", rotate: "拖动滑块旋转圆形拼图", image: "输入图中的字符", click: "请依次点击", submit: "确认", close: "关闭", secure: "安全验证", secureHint: "请完成验证以继续操作" }, en: { retry: "Verification failed. Please try again.", refreshed: "Verification failed. Please try again.", refreshing: "Refreshing challenge…", refresh: "Refresh", slider: "Drag the slider to complete the puzzle", restore: "Drag the image strip to restore the picture", rotate: "Drag the slider to rotate the centre disc", image: "Enter the characters", click: "Click in order", submit: "Verify", close: "Close", secure: "Security check", secureHint: "Complete verification to continue" } };
  const lang = (locale) => text[String(locale || "").toLowerCase().startsWith("zh") ? "zh-CN" : "en"];
  const svgToURL = (svg) => `data:image/svg+xml;charset=utf-8,${encodeURIComponent(svg)}`;
  // PNG is emitted by the service for every visual challenge. Prefer it so a
  // WebView follows the same rendering contract as mini programs and App SDKs.
  const captchaImageURL = (pngBase64, svg) => pngBase64 ? `data:image/png;base64,${pngBase64}` : svgToURL(svg || "");

  class CaptchaSDK {
    constructor(options) { if (!options || !options.endpoint || !options.action) throw new Error("CaptchaSDK requires endpoint and action"); this.options = Object.assign({ locale: navigator.language, theme: "light", timeout: 10000, preferredMode: "invisible", forceHighRisk: false }, options); this.startedAt = 0; this.events = 0; this.fingerprint = this.options.fingerprint || this.makeFingerprint(); this.lastRiskScore = undefined; }
    async verify() {
      this.startedAt = Date.now(); this.events = 0;
      const events = ["pointermove", "keydown", "touchstart"], record = () => { this.events++; };
      events.forEach((e) => document.addEventListener(e, record, { passive: true }));
      try {
        const challenge = await this.createChallenge();
        const verified = challenge.status === "verified" ? challenge : await this.showChallenge(challenge);
        if (verified.status !== "verified") throw new Error(verified.message || lang(this.options.locale).retry);
        return this.success(verified);
      } finally { events.forEach((e) => document.removeEventListener(e, record)); }
    }
    async createChallenge() { const challenge = await this.request("/api/v1/challenges", { preferred_mode: this.options.preferredMode, action: this.options.action, locale: this.options.locale, theme: this.options.theme, client: this.signals() }); this.lastRiskScore = challenge.risk_score; return challenge; }
    success(response) { const result = { verificationToken: response.verification_token, mode: response.mode, riskScore: response.risk_score ?? this.lastRiskScore }; if (typeof this.options.onSuccess === "function") this.options.onSuccess(result); return result; }
    signals() { if (this.options.forceHighRisk) return { fingerprint: "", event_count: 0, duration_ms: 0, language: this.options.locale }; return { fingerprint: this.fingerprint, event_count: this.events, duration_ms: Date.now() - this.startedAt, language: this.options.locale }; }
    async request(path, body) {
      const controller = new AbortController(), timeout = setTimeout(() => controller.abort(), this.options.timeout);
      try { const response = await fetch(this.options.endpoint.replace(/\/$/, "") + path, { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify(body), signal: controller.signal }); const data = await response.json(); if (!response.ok) throw new Error(data.error || "Captcha request failed"); return data; } finally { clearTimeout(timeout); }
    }
    showChallenge(challenge) {
      return new Promise((resolve, reject) => {
        const t = lang(this.options.locale), root = document.createElement("div");
        root.className = `pte-captcha pte-captcha--${this.options.theme}${this.options.className ? ` ${this.options.className}` : ""}`;
        Object.entries(this.options.styles || {}).forEach(([name, value]) => root.style.setProperty(name, value));
        root.innerHTML = `<div class="pte-captcha__dialog" role="dialog" aria-modal="true"><button class="pte-captcha__close" aria-label="${t.close}">×</button><div class="pte-captcha__header"><span class="pte-captcha__shield">✓</span><div><div class="pte-captcha__title">${t.secure}</div><div class="pte-captcha__subtitle">${t.secureHint}</div></div></div><div class="pte-captcha__message"></div><div class="pte-captcha__content"></div><div class="pte-captcha__feedback"><div class="pte-captcha__error" role="alert"></div><button class="pte-captcha__refresh" type="button" aria-label="${t.refresh}" title="${t.refresh}"><span aria-hidden="true">↻</span></button></div></div>`;
        document.body.appendChild(root);
        const close = (error) => { root.remove(); if (error) reject(error); };
        root.querySelector(".pte-captcha__close").onclick = () => close(new Error("已取消安全验证"));
        const content = root.querySelector(".pte-captcha__content"), message = root.querySelector(".pte-captcha__message"), error = root.querySelector(".pte-captcha__error"), refresh = root.querySelector(".pte-captcha__refresh");
        let refreshing = false, renderChallenge;
        const refreshChallenge = async (notice = "") => {
          if (refreshing) return;
          refreshing = true; refresh.disabled = true; content.setAttribute("aria-busy", "true"); error.textContent = notice || t.refreshing;
          try { renderChallenge(await this.createChallenge(), notice); } catch (err) { content.removeAttribute("aria-busy"); error.textContent = err.message; } finally { refreshing = false; refresh.disabled = false; }
        };
        refresh.onclick = () => refreshChallenge();
        renderChallenge = (current, notice = "") => {
          message.hidden = true; message.textContent = current.message; error.textContent = notice; content.removeAttribute("aria-busy");
          const verify = async (payload) => { try { const result = await this.request(`/api/v1/challenges/${current.challenge_id}/verify`, Object.assign(payload, { client: this.signals() })); if (result.status === "verified") { close(); resolve(result); return; } error.textContent = result.message || t.retry; content.setAttribute("aria-busy", "true"); await new Promise((done) => setTimeout(done, 550)); await refreshChallenge(t.refreshed); } catch (err) { content.removeAttribute("aria-busy"); error.textContent = err.message; } };
          if (current.mode === "image") this.renderImage(content, current, t, verify); else if (current.mode === "click") this.renderClick(content, current, t, verify); else if (current.mode === "rotate") this.renderRotate(content, current, t, verify); else if (current.mode === "restore") this.renderRestore(content, current, t, verify); else this.renderPuzzle(content, current, t, verify);
        };
        renderChallenge(challenge);
      });
    }
    renderImage(node, challenge, t, verify) { node.innerHTML = `<img class="pte-captcha__image" alt="captcha" src="${captchaImageURL(challenge.image_png_base64, challenge.image_svg)}"><label>${t.image}<input class="pte-captcha__input" autocomplete="off" maxlength="8"></label><button class="pte-captcha__button">${t.submit}</button>`; const submit = () => verify({ answer: node.querySelector("input").value }); node.querySelector("button").onclick = submit; node.querySelector("input").onkeydown = (e) => { if (e.key === "Enter") submit(); }; }
    renderClick(node, challenge, t, verify) {
      const targets = Array.isArray(challenge.click_targets) ? challenge.click_targets : [], targetSVGs = Array.isArray(challenge.click_target_svgs) ? challenge.click_target_svgs : [], targetPNGs = Array.isArray(challenge.click_target_png_base64s) ? challenge.click_target_png_base64s : [];
      node.innerHTML = `<div class="pte-captcha__click-prompt">${t.click}<span class="pte-captcha__click-targets">${targets.map((target, index) => `<span class="pte-captcha__click-target" aria-label="${target}" title="${target}">${targetPNGs[index] || targetSVGs[index] ? `<img alt="${target}" src="${captchaImageURL(targetPNGs[index], targetSVGs[index])}">` : target}</span>`).join("")}</span></div><div class="pte-captcha__click-scene" role="application" aria-label="${t.click}"><img alt="" src="${captchaImageURL(challenge.click_png_base64, challenge.click_svg)}"><div class="pte-captcha__click-markers"></div></div><button class="pte-captcha__button" disabled>${t.submit}</button>`;
      const scene = node.querySelector(".pte-captcha__click-scene"), markers = node.querySelector(".pte-captcha__click-markers"), submit = node.querySelector("button"), clicks = [];
      scene.onpointerup = (event) => {
        if (clicks.length >= targets.length) return;
        const box = scene.getBoundingClientRect(), x = Math.round((event.clientX - box.left) * 320 / box.width), y = Math.round((event.clientY - box.top) * 160 / box.height);
        if (x < 0 || x > 320 || y < 0 || y > 160) return;
        clicks.push({ x, y });
        const marker = document.createElement("span"); marker.className = "pte-captcha__click-marker"; marker.textContent = String(clicks.length); marker.style.left = `${(x / 320) * 100}%`; marker.style.top = `${(y / 160) * 100}%`; markers.appendChild(marker);
        submit.disabled = clicks.length !== targets.length;
      };
      submit.onclick = () => verify({ clicks });
    }
    renderPuzzle(node, challenge, t, verify) {
      node.innerHTML = `<div class="pte-captcha__puzzle"><img class="pte-captcha__background" alt="puzzle" src="${captchaImageURL(challenge.background_png_base64, challenge.background_svg)}"><img class="pte-captcha__piece" alt="" src="${captchaImageURL(challenge.piece_png_base64, challenge.piece_svg)}"></div><p class="pte-captcha__slider-label">${t.slider}</p><div class="pte-captcha__slider" role="slider" tabindex="0" aria-label="${t.slider}" aria-valuemin="0" aria-valuemax="${challenge.puzzle_max_x}" aria-valuenow="0"><div class="pte-captcha__slider-fill"></div><button class="pte-captcha__slider-handle" type="button" aria-label="drag"><span>»</span></button></div>`;
      const puzzle = node.querySelector(".pte-captcha__puzzle"), slider = node.querySelector(".pte-captcha__slider"), fill = node.querySelector(".pte-captcha__slider-fill"), handle = node.querySelector(".pte-captcha__slider-handle"), piece = node.querySelector(".pte-captcha__piece"), track = [];
      let active = false, value = 0;
      const positionPiece = () => { const scale = puzzle.clientWidth / 320; piece.style.width = `${46 * scale}px`; piece.style.height = `${46 * scale}px`; piece.style.top = `${challenge.puzzle_y * scale}px`; piece.style.left = `${value * scale}px`; };
      const point = () => track.push({ x: Math.round(value), y: challenge.puzzle_y, t: Date.now() - this.startedAt });
      const setValue = (next) => { value = Math.max(0, Math.min(challenge.puzzle_max_x, next)); const handleWidth = handle.offsetWidth || 48, travel = Math.max(1, slider.clientWidth - handleWidth), left = value / challenge.puzzle_max_x * travel; positionPiece(); fill.style.width = `${left + handleWidth}px`; handle.style.left = `${left}px`; slider.setAttribute("aria-valuenow", String(Math.round(value))); };
      const fromPointer = (event) => { const box = slider.getBoundingClientRect(), handleWidth = handle.offsetWidth || 48, travel = Math.max(1, box.width - handleWidth); setValue((event.clientX - box.left - handleWidth / 2) / travel * challenge.puzzle_max_x); };
      slider.onpointerdown = (event) => { active = true; slider.classList.add("is-dragging"); slider.setPointerCapture(event.pointerId); fromPointer(event); point(); event.preventDefault(); };
      slider.onpointermove = (event) => { if (!active) return; fromPointer(event); point(); };
      slider.onpointerup = (event) => { if (!active) return; fromPointer(event); point(); active = false; slider.classList.remove("is-dragging"); if (slider.hasPointerCapture(event.pointerId)) slider.releasePointerCapture(event.pointerId); verify({ offset_x: Math.round(value), trajectory: track }); };
      slider.onkeydown = (event) => { if (event.key !== "ArrowLeft" && event.key !== "ArrowRight") return; setValue(value + (event.key === "ArrowRight" ? 8 : -8)); point(); event.preventDefault(); };
      handle.onpointerdown = (event) => event.preventDefault();
      if (typeof ResizeObserver !== "undefined") new ResizeObserver(positionPiece).observe(puzzle);
      setValue(0);
    }
    renderRestore(node, challenge, t, verify) {
      node.innerHTML = `<div class="pte-captcha__restore"><img class="pte-captcha__background" alt="restore" src="${captchaImageURL(challenge.restore_background_png_base64)}"><img class="pte-captcha__restore-piece" alt="" src="${captchaImageURL(challenge.restore_piece_png_base64)}"></div><p class="pte-captcha__slider-label">${t.restore}</p><div class="pte-captcha__slider" role="slider" tabindex="0" aria-label="${t.restore}" aria-valuemin="0" aria-valuemax="${challenge.restore_max_x}" aria-valuenow="0"><div class="pte-captcha__slider-fill"></div><button class="pte-captcha__slider-handle" type="button" aria-label="drag"><span>»</span></button></div>`;
      const stage = node.querySelector(".pte-captcha__restore"), slider = node.querySelector(".pte-captcha__slider"), fill = node.querySelector(".pte-captcha__slider-fill"), handle = node.querySelector(".pte-captcha__slider-handle"), piece = node.querySelector(".pte-captcha__restore-piece"), track = [];
      const max = challenge.restore_max_x || 266, pieceWidth = challenge.restore_piece_width || 54;
      let active = false, value = 0;
      const positionPiece = () => { const scale = stage.clientWidth / 320; piece.style.width = `${pieceWidth * scale}px`; piece.style.height = `${stage.clientHeight}px`; piece.style.left = `${value * scale}px`; };
      const point = () => track.push({ x: Math.round(value), y: 80, t: Date.now() - this.startedAt });
      const setValue = (next) => { value = Math.max(0, Math.min(max, next)); const handleWidth = handle.offsetWidth || 48, travel = Math.max(1, slider.clientWidth - handleWidth), left = value / max * travel; positionPiece(); fill.style.width = `${left + handleWidth}px`; handle.style.left = `${left}px`; slider.setAttribute("aria-valuenow", String(Math.round(value))); };
      const fromPointer = (event) => { const box = slider.getBoundingClientRect(), handleWidth = handle.offsetWidth || 48, travel = Math.max(1, box.width - handleWidth); setValue((event.clientX - box.left - handleWidth / 2) / travel * max); };
      slider.onpointerdown = (event) => { active = true; slider.classList.add("is-dragging"); slider.setPointerCapture(event.pointerId); fromPointer(event); point(); event.preventDefault(); };
      slider.onpointermove = (event) => { if (!active) return; fromPointer(event); point(); };
      slider.onpointerup = (event) => { if (!active) return; fromPointer(event); point(); active = false; slider.classList.remove("is-dragging"); if (slider.hasPointerCapture(event.pointerId)) slider.releasePointerCapture(event.pointerId); verify({ offset_x: Math.round(value), trajectory: track }); };
      slider.onkeydown = (event) => { if (event.key !== "ArrowLeft" && event.key !== "ArrowRight") return; setValue(value + (event.key === "ArrowRight" ? 8 : -8)); point(); event.preventDefault(); };
      handle.onpointerdown = (event) => event.preventDefault();
      if (typeof ResizeObserver !== "undefined") new ResizeObserver(positionPiece).observe(stage);
      setValue(0);
    }
    renderRotate(node, challenge, t, verify) {
      const min = Number.isFinite(challenge.rotate_min) ? challenge.rotate_min : 0, max = Number.isFinite(challenge.rotate_max) ? challenge.rotate_max : 180;
      node.innerHTML = `<div class="pte-captcha__rotate"><img class="pte-captcha__background" alt="rotate" src="${captchaImageURL(challenge.rotate_background_png_base64, challenge.rotate_png_base64)}"><div class="pte-captcha__rotate-disc"><img alt="" src="${captchaImageURL(challenge.rotate_piece_png_base64)}"></div></div><p class="pte-captcha__slider-label">${t.rotate}</p><div class="pte-captcha__slider" role="slider" tabindex="0" aria-label="${t.rotate}" aria-valuemin="${min}" aria-valuemax="${max}" aria-valuenow="${min}"><div class="pte-captcha__slider-fill"></div><button class="pte-captcha__slider-handle" type="button" aria-label="drag"><span>»</span></button></div>`;
      const slider = node.querySelector(".pte-captcha__slider"), fill = node.querySelector(".pte-captcha__slider-fill"), handle = node.querySelector(".pte-captcha__slider-handle"), disc = node.querySelector(".pte-captcha__rotate-disc img"), track = [];
      let active = false, value = min;
      const point = () => track.push({ x: Math.round(value), y: 0, t: Date.now() - this.startedAt });
      const setValue = (next) => { value = Math.max(min, Math.min(max, next)); const handleWidth = handle.offsetWidth || 48, travel = Math.max(1, slider.clientWidth - handleWidth), left = (value - min) / (max - min) * travel; disc.style.transform = `rotate(${value}deg)`; fill.style.width = `${left + handleWidth}px`; handle.style.left = `${left}px`; slider.setAttribute("aria-valuenow", String(Math.round(value))); };
      const fromPointer = (event) => { const box = slider.getBoundingClientRect(), handleWidth = handle.offsetWidth || 48, travel = Math.max(1, box.width - handleWidth); setValue(min + (event.clientX - box.left - handleWidth / 2) / travel * (max - min)); };
      slider.onpointerdown = (event) => { active = true; slider.classList.add("is-dragging"); slider.setPointerCapture(event.pointerId); fromPointer(event); point(); event.preventDefault(); };
      slider.onpointermove = (event) => { if (!active) return; fromPointer(event); point(); };
      slider.onpointerup = (event) => { if (!active) return; fromPointer(event); point(); active = false; slider.classList.remove("is-dragging"); if (slider.hasPointerCapture(event.pointerId)) slider.releasePointerCapture(event.pointerId); verify({ rotate_angle: Math.round(value), rotate_trajectory: track }); };
      slider.onkeydown = (event) => { if (event.key !== "ArrowLeft" && event.key !== "ArrowRight") return; setValue(value + (event.key === "ArrowRight" ? 8 : -8)); point(); event.preventDefault(); };
      handle.onpointerdown = (event) => event.preventDefault();
      setValue(min);
    }
    makeFingerprint() { const raw = [navigator.userAgent, navigator.language, screen.width, screen.height, Intl.DateTimeFormat().resolvedOptions().timeZone].join("|"); let n = 2166136261; for (let i = 0; i < raw.length; i++) n = Math.imul(n ^ raw.charCodeAt(i), 16777619); return `web-${(n >>> 0).toString(36)}`; }
  }
  global.CaptchaSDK = CaptchaSDK;
})(window);
