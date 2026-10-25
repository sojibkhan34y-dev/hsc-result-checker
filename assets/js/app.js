// HSC Result Checker - Theme Controller

const themeButton = document.createElement("button");

themeButton.className = "theme-btn";
themeButton.type = "button";
themeButton.setAttribute("aria-label", "Toggle dark mode");

document.body.appendChild(themeButton);

function applyTheme(theme) {
    const isDark = theme === "dark";

    document.body.classList.toggle("dark", isDark);

    themeButton.textContent = isDark ? "☀️" : "🌙";

    themeButton.setAttribute(
        "aria-label",
        isDark ? "Enable light mode" : "Enable dark mode"
    );
}

const savedTheme = localStorage.getItem("hsc-theme");

applyTheme(
    savedTheme === "dark" ? "dark" : "light"
);

themeButton.addEventListener("click", () => {
    const nextTheme =
        document.body.classList.contains("dark")
            ? "light"
            : "dark";

    applyTheme(nextTheme);

    localStorage.setItem("hsc-theme", nextTheme);
});