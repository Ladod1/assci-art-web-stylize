// Progressive enhancements: the site works without JavaScript, this file
// only adds live feedback and shortcuts on top of the server-rendered pages.
(function () {
  "use strict";

  var root = document.documentElement;

  function store(key, value) {
    try {
      localStorage.setItem(key, value);
    } catch (e) {}
  }

  function load(key) {
    try {
      return localStorage.getItem(key);
    } catch (e) {
      return null;
    }
  }

  // ---------- Toast ----------
  var toast = document.getElementById("toast");
  var toastTimer;

  function showToast(message) {
    if (!toast) return;
    toast.textContent = message;
    toast.classList.add("is-visible");
    clearTimeout(toastTimer);
    toastTimer = setTimeout(function () {
      toast.classList.remove("is-visible");
    }, 2200);
  }

  // ---------- Theme toggle ----------
  var themeToggle = document.getElementById("theme-toggle");
  if (themeToggle) {
    themeToggle.hidden = false;
    var syncThemeLabel = function () {
      var dark = root.getAttribute("data-theme") === "dark";
      var label = dark ? "Switch to light mode" : "Switch to dark mode";
      themeToggle.setAttribute("aria-label", label);
      themeToggle.title = label;
    };
    syncThemeLabel();
    themeToggle.addEventListener("click", function () {
      var next = root.getAttribute("data-theme") === "dark" ? "light" : "dark";
      root.setAttribute("data-theme", next);
      store("theme", next);
      syncThemeLabel();
    });
  }

  // ---------- Form ----------
  var form = document.getElementById("ascii-form");
  if (form) {
    var textarea = document.getElementById("text");
    var counter = document.getElementById("counter");
    var count = document.getElementById("count");
    var errorBox = document.getElementById("text-error");
    var errorMsg = document.getElementById("text-error-msg");
    var generate = document.getElementById("generate");
    var clear = document.getElementById("clear");
    var max = parseInt(textarea.getAttribute("maxlength"), 10) || 200;
    // Printable ASCII plus line breaks, matching the server-side check.
    var allowed = /^[\x20-\x7E\r\n]*$/;

    var setError = function (message) {
      errorMsg.textContent = message;
      errorBox.hidden = !message;
      if (message) {
        textarea.setAttribute("aria-invalid", "true");
        textarea.setAttribute("aria-describedby", "text-hint text-error");
      } else {
        textarea.removeAttribute("aria-invalid");
        textarea.setAttribute("aria-describedby", "text-hint");
      }
    };

    var validate = function (submitting) {
      var value = textarea.value;
      if (!allowed.test(value)) {
        var bad = value.match(/[^\x20-\x7E\r\n]/);
        setError("“" + bad[0] + "” is not supported. Only printable ASCII characters (letters, digits, spaces and symbols like !?#) can be used.");
        return false;
      }
      if (submitting && value === "") {
        setError("Please enter some text to convert.");
        return false;
      }
      setError("");
      return true;
    };

    var update = function () {
      var length = textarea.value.length;
      count.textContent = length;
      counter.classList.toggle("is-near", length >= max * 0.85 && length < max);
      counter.classList.toggle("is-full", length >= max);
      clear.hidden = length === 0;
    };

    textarea.addEventListener("input", function () {
      update();
      validate(false);
    });

    textarea.addEventListener("keydown", function (event) {
      if (event.key === "Enter" && (event.ctrlKey || event.metaKey)) {
        event.preventDefault();
        if (typeof form.requestSubmit === "function") {
          form.requestSubmit();
        } else if (validate(true)) {
          form.submit();
        }
      }
    });

    clear.addEventListener("click", function () {
      textarea.value = "";
      update();
      setError("");
      textarea.focus();
    });

    form.addEventListener("submit", function (event) {
      if (!validate(true)) {
        event.preventDefault();
        textarea.focus();
        return;
      }
      generate.disabled = true;
      generate.classList.add("is-loading");
      generate.querySelector(".btn-label").textContent = "Generating…";
    });

    // Re-enable the button if the page is restored from the back/forward cache.
    window.addEventListener("pageshow", function () {
      generate.disabled = false;
      generate.classList.remove("is-loading");
      generate.querySelector(".btn-label").textContent = "Generate";
    });

    update();
  }

  // ---------- Result ----------
  var result = document.getElementById("result");
  if (result) {
    var art = document.getElementById("art");
    var tools = document.getElementById("result-tools");
    var zoom = document.getElementById("zoom");
    tools.hidden = false;

    var setSize = function (size) {
      root.style.setProperty("--art-size", size + "px");
    };
    var savedSize = load("art-size");
    if (savedSize) {
      zoom.value = savedSize;
    }
    setSize(zoom.value);
    zoom.addEventListener("input", function () {
      setSize(zoom.value);
      store("art-size", zoom.value);
    });

    // Older browsers and some privacy settings block the Clipboard API,
    // so fall back to selecting a hidden textarea and copying it.
    var legacyCopy = function (text) {
      var helper = document.createElement("textarea");
      helper.value = text;
      helper.setAttribute("readonly", "");
      helper.style.position = "fixed";
      helper.style.opacity = "0";
      document.body.appendChild(helper);
      helper.select();
      var copied = false;
      try {
        copied = document.execCommand("copy");
      } catch (e) {}
      document.body.removeChild(helper);
      return copied;
    };

    var reportCopy = function (copied) {
      showToast(copied ? "Copied to clipboard" : "Could not copy, please select the text manually");
    };

    document.getElementById("copy").addEventListener("click", function () {
      var text = art.textContent;
      if (navigator.clipboard && window.isSecureContext) {
        navigator.clipboard.writeText(text).then(function () {
          reportCopy(true);
        }, function () {
          reportCopy(legacyCopy(text));
        });
        return;
      }
      reportCopy(legacyCopy(text));
    });

    document.getElementById("download").addEventListener("click", function () {
      var blob = new Blob([art.textContent], { type: "text/plain" });
      var link = document.createElement("a");
      link.href = URL.createObjectURL(blob);
      link.download = "ascii-art.txt";
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      URL.revokeObjectURL(link.href);
      showToast("Download started");
    });

    result.scrollIntoView({ behavior: "smooth", block: "start" });
    result.focus({ preventScroll: true });
  }
})();
