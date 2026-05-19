(() => {
    const API = "http://localhost:8080/api/v1";

    const $ = (id) => document.getElementById(id);

    const el = {
        authButton: $("authButton"),
        authModal: $("authModal"),
        closeModal: $("closeModal"),

        loginForm: $("loginForm"),
        registerForm: $("registerForm"),

        loader: $("loader"),
        notification: $("notification"),

        requestButton: $("requestButton"),
        downloadButton: $("downloadButton"),

        formatModal: $("formatModal"),

        city: $("cityInput"),
        from: $("startDate"),
        to: $("endDate"),
        dimensions: $("dimensionsInput"),
        scale: $("scaleInput"),

        bands: $("bandsContainer"),
        closeFormatModal: $("closeFormatModal"),

        image: $("satelliteImage"),
        placeholder: $("imagePlaceholder"),

        viewer: $("imageViewer"),
        viewerImage: $("viewerImage"),

        history: $("historyList"),

        tabs: document.querySelectorAll(".tab-button"),
    };

    let zoom = 1;

    const storage = {
        get token() {
            return localStorage.getItem("token");
        },

        set token(value) {
            localStorage.setItem("token", value);
        },

        clear() {
            localStorage.removeItem("token");
        },
    };

    const notify = (text, error = false) => {
        el.notification.textContent = text;

        el.notification.style.background =
            error ? "#dc2626" : "#16a34a";

        el.notification.classList.remove("hidden");

        clearTimeout(window.notifyTimer);

        window.notifyTimer = setTimeout(() => {
            el.notification.classList.add("hidden");
        }, 3500);
    };

    const toggleLoader = (show) =>
        el.loader.classList.toggle("hidden", !show);

    const headers = () => ({
        "Content-Type": "application/json",

        ...(storage.token && {
            Authorization: `Bearer ${storage.token}`,
        }),
    });

    async function api(path, options = {}) {

        const response = await fetch(
            `${API}${path}`,
            {
                ...options,

                headers: {
                    ...headers(),
                    ...options.headers,
                },
            }
        );

        if (!response.ok) {

            let message = "Ошибка";

            try {

                const data = await response.json();

                message = data.error || message;

            } catch {}

            throw new Error(message);
        }

        return response;
    }

    function getBands() {
        return [
            ...el.bands.querySelectorAll(
                'input[type="checkbox"]:checked'
            ),
        ]
            .map((x) => x.value)
            .join(",");
    }

    function getPayload() {

        const payload = {
            city: el.city.value.trim(),
            from: el.from.value,
            to: el.to.value,
            bands: getBands(),
            dimensions: Number(el.dimensions.value),
            scale: Number(el.scale.value),
        };

        const valid = Object.values(payload)
            .every(Boolean);

        return valid ? payload : null;
    }

    function renderImage(src) {

        el.image.src = src;

        el.image.classList.remove("hidden");

        el.placeholder.classList.add("hidden");
    }

    function clearImage() {

        el.image.src = "";

        el.image.classList.add("hidden");

        el.placeholder.classList.remove("hidden");
    }

    async function previewImage() {

        const payload = getPayload();

        if (!payload) {
            notify("Заполните все поля", true);
            return;
        }

        try {

            toggleLoader(true);

            const response = await api(
                "/sentinel/imagery/preview",
                {
                    method: "POST",

                    body: JSON.stringify(payload),
                }
            );

            const blob = await response.blob();

            renderImage(
                URL.createObjectURL(blob)
            );

            await loadHistory();

            notify("Превью получено");

        } catch (error) {

            notify(error.message, true);

        } finally {

            toggleLoader(false);
        }
    }

    async function downloadImage(format) {

        const payload = getPayload();

        if (!payload) {
            notify("Заполните все поля", true);
            return;
        }

        try {

            toggleLoader(true);

            const response = await api(
                "/sentinel/imagery/download",
                {
                    method: "POST",

                    body: JSON.stringify({
                        ...payload,
                        output_format: format,
                    }),
                }
            );

            const blob = await response.blob();

            const url =
                URL.createObjectURL(blob);

            const link =
                document.createElement("a");

            link.href = url;

            link.download =
                `satellite-image.${
                    format === "png"
                        ? "png"
                        : "tif"
                }`;

            document.body.appendChild(link);

            link.click();

            link.remove();

            URL.revokeObjectURL(url);

            notify("Снимок скачан");

        } catch (error) {

            notify(error.message, true);

        } finally {

            toggleLoader(false);

            el.formatModal.classList.add(
                "hidden"
            );
        }
    }

    async function loadHistory() {

        try {

            const response = await api(
                "/sentinel/requests"
            );

            const history =
                await response.json();

            el.history.innerHTML = "";

            history
                .sort(
                    (a, b) =>
                        new Date(b.updatedAt) -
                        new Date(a.updatedAt)
                )
                .forEach(addHistoryItem);

        } catch (error) {

            console.error(error);
        }
    }

    async function deleteHistory(id, item) {

        try {

            toggleLoader(true);

            await api(
                `/sentinel/requests/delete/${id}`,
                {
                    method: "DELETE",
                }
            );

            item.remove();

            notify("Запрос удалён");

        } catch (error) {

            notify(error.message, true);

        } finally {

            toggleLoader(false);
        }
    }

    function fillForm(data) {

        el.city.value = data.city;

        el.from.value = data.from;

        el.to.value = data.to;

        el.dimensions.value =
            data.dimensions;

        el.scale.value =
            data.scale;

        const bands =
            (data.bands || "").split(",");

        el.bands
            .querySelectorAll(
                'input[type="checkbox"]'
            )
            .forEach((checkbox) => {

                checkbox.checked =
                    bands.includes(
                        checkbox.value
                    );
            });
    }

    function addHistoryItem(item) {

        const div =
            document.createElement("div");

        div.className = "history-item";

        div.innerHTML = `
            <div class="history-header">

                <span>${item.city}</span>

                <div class="history-actions">

                    <span>
                        ${new Date(
                            item.updatedAt
                        ).toLocaleString("ru-RU")}
                    </span>

                    <button
                        class="delete-history-button"
                        data-id="${item.id}"
                        type="button"
                    >
                        ×
                    </button>

                </div>

            </div>

            <div class="history-body hidden">

                <div>
                    <strong>Период:</strong>
                    ${item.from} — ${item.to}
                </div>

                <div>
                    <strong>Слои:</strong>
                    ${item.bands}
                </div>

                <div>
                    <strong>Разрешение превью:</strong>
                    ${item.dimensions}
                </div>

                <div>
                    <strong>Пространственное разрешение:</strong>
                    ${item.scale} м
                </div>

                <button
                    class="repeat-button"
                    type="button"

                    data-city="${item.city}"
                    data-from="${item.from}"
                    data-to="${item.to}"
                    data-bands="${item.bands}"
                    data-dimensions="${item.dimensions}"
                    data-scale="${item.scale}"
                >
                    Заполнить параметры
                </button>

            </div>
        `;

        el.history.appendChild(div);
    }

    async function login(event) {

        event.preventDefault();

        const [email, password] =
            el.loginForm.querySelectorAll(
                "input"
            );

        try {

            toggleLoader(true);

            const response = await api(
                "/auth/login",
                {
                    method: "POST",

                    body: JSON.stringify({
                        email: email.value.trim(),
                        password:
                            password.value.trim(),
                    }),
                }
            );

            const data =
                await response.json();

            storage.token = data.token;

            el.authModal.classList.add(
                "hidden"
            );

            await checkAuth();

            notify("Успешный вход");

        } catch (error) {

            notify(error.message, true);

        } finally {

            toggleLoader(false);
        }
    }

    async function register(event) {

        event.preventDefault();

        const [name, email, password] =
            el.registerForm.querySelectorAll(
                "input"
            );

        try {

            toggleLoader(true);

            await api(
                "/auth/register",
                {
                    method: "POST",

                    body: JSON.stringify({
                        name: name.value.trim(),
                        email: email.value.trim(),
                        password:
                            password.value.trim(),
                    }),
                }
            );

            notify("Регистрация успешна");

            el.tabs[0].click();

        } catch (error) {

            notify(error.message, true);

        } finally {

            toggleLoader(false);
        }
    }

    async function checkAuth() {

        if (!storage.token) return;

        try {

            const response = await api(
                "/auth/me"
            );

            const user =
                await response.json();

            el.authButton.textContent =
                `${user.name} | Выйти`;

            await loadHistory();

        } catch {

            storage.clear();
        }
    }

    function logout() {

        storage.clear();

        clearImage();

        el.history.innerHTML = "";

        el.authButton.textContent =
            "Войти / Регистрация";

        notify("Вы вышли");
    }

    function openViewer() {

        if (!el.image.src) return;

        zoom = 1;

        el.viewerImage.src = el.image.src;

        el.viewerImage.style.transform =
            "scale(1)";

        el.viewer.classList.remove(
            "hidden"
        );
    }

    function closeViewer() {

        el.viewer.classList.add(
            "hidden"
        );
    }

    function bindEvents() {

        el.authButton.addEventListener(
            "click",
            () => {

                if (storage.token) {

                    logout();

                } else {

                    el.authModal.classList.remove(
                        "hidden"
                    );
                }
            }
        );

        el.closeModal.addEventListener(
            "click",
            () =>
                el.authModal.classList.add(
                    "hidden"
                )
        );

        el.loginForm.addEventListener(
            "submit",
            login
        );

        el.registerForm.addEventListener(
            "submit",
            register
        );

        el.requestButton.addEventListener(
            "click",
            previewImage
        );

        el.downloadButton.addEventListener(
            "click",
            () => {
                el.formatModal.classList.remove(
                    "hidden"
                );
            }
        );

        document
            .querySelectorAll(".format-button")
            .forEach((button) => {

                button.addEventListener(
                    "click",
                    () =>
                        downloadImage(
                            button.dataset.format
                        )
                );
            });

        el.tabs.forEach((tab) => {

            tab.addEventListener(
                "click",
                () => {

                    el.tabs.forEach((x) =>
                        x.classList.remove(
                            "active"
                        )
                    );

                    tab.classList.add("active");

                    const login =
                        tab.dataset.tab ===
                        "login";

                    el.loginForm.classList.toggle(
                        "hidden",
                        !login
                    );

                    el.registerForm.classList.toggle(
                        "hidden",
                        login
                    );
                }
            );
        });

        el.history.addEventListener(
            "click",
            async (event) => {

                const deleteButton =
                    event.target.closest(
                        ".delete-history-button"
                    );

                if (deleteButton) {

                    await deleteHistory(
                        deleteButton.dataset.id,

                        deleteButton.closest(
                            ".history-item"
                        )
                    );

                    return;
                }

                const repeatButton =
                    event.target.closest(
                        ".repeat-button"
                    );

                if (repeatButton) {

                    fillForm(
                        repeatButton.dataset
                    );

                    notify(
                        "Параметры заполнены"
                    );

                    return;
                }

                const item =
                    event.target.closest(
                        ".history-item"
                    );

                if (!item) return;

                item
                    .querySelector(
                        ".history-body"
                    )
                    ?.classList.toggle(
                        "hidden"
                    );
            }
        );

        el.image.addEventListener(
            "click",
            openViewer
        );

        el.viewer.addEventListener(
            "click",
            closeViewer
        );

        el.viewerImage.addEventListener(
            "wheel",
            (event) => {

                event.preventDefault();

                zoom +=
                    event.deltaY * -0.001;

                zoom = Math.min(
                    Math.max(1, zoom),
                    5
                );

                el.viewerImage.style.transform =
                    `scale(${zoom})`;
            }
        );

        el.bands.addEventListener(
            "change",
            () => {

                const checked = [
                    ...el.bands.querySelectorAll(
                        'input[type="checkbox"]:checked'
                    ),
                ];

                if (checked.length > 3) {

                    checked.at(-1).checked = false;

                    notify(
                        "Максимум 3 слоя",
                        true
                    );
                }
            }
        );

        el.closeFormatModal.addEventListener(
            "click",
            () => {
                el.formatModal.classList.add(
                    "hidden"
                );
            }
        );

        el.formatModal.addEventListener(
            "click",
            (event) => {

                if (
                    event.target === el.formatModal
                ) {

                    el.formatModal.classList.add(
                        "hidden"
                    );
                }
            }
        );
    }

    async function init() {

        bindEvents();

        await checkAuth();
    }

    init();
})();