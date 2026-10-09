// Runs in <head> so the saved theme applies before the page paints.
(function () {
  var theme = null;
  try {
    theme = localStorage.getItem("theme");
  } catch (e) {}
  if (theme !== "light" && theme !== "dark") {
    theme = window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light";
  }
  document.documentElement.setAttribute("data-theme", theme);
})();
