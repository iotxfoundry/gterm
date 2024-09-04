import "@xterm/xterm";
import "@xterm/addon-attach";
import "@xterm/addon-clipboard";
import "@xterm/addon-fit";
import "@xterm/addon-image";
// import { LigaturesAddon } from "@xterm/addon-ligatures";
// import { SearchAddon, ISearchOptions } from "@xterm/addon-search";
// import { SerializeAddon } from "@xterm/addon-serialize";
import "@xterm/addon-web-links";
import "@xterm/addon-webgl";
//import "@xterm/addon-unicode11";

const isWindows = navigator.userAgentData.platform == "Windows";

const xtermjsTheme = {
  foreground: "#F8F8F8",
  background: "#2D2E2C",
  selectionBackground: "#5DA5D533",
  selectionInactiveBackground: "#555555AA",
  black: "#1E1E1D",
  brightBlack: "#262625",
  red: "#CE5C5C",
  brightRed: "#FF7272",
  green: "#5BCC5B",
  brightGreen: "#72FF72",
  yellow: "#CCCC5B",
  brightYellow: "#FFFF72",
  blue: "#5D5DD3",
  brightBlue: "#7279FF",
  magenta: "#BC5ED1",
  brightMagenta: "#E572FF",
  cyan: "#5DA5D5",
  brightCyan: "#72F0FF",
  white: "#F8F8F8",
  brightWhite: "#FFFFFF",
};

const imageSettings = {
  enableSizeReports: true, // whether to enable CSI t reports (see below)
  pixelLimit: 16777216, // max. pixel size of a single image
  sixelSupport: true, // enable sixel support
  sixelScrolling: true, // whether to scroll on image output
  sixelPaletteLimit: 256, // initial sixel palette size
  sixelSizeLimit: 25000000, // size limit of a single sixel sequence
  storageLimit: 128, // FIFO storage limit in MB
  showPlaceholder: true, // whether to show a placeholder for evicted images
  iipSupport: true, // enable iTerm IIP support
  iipSizeLimit: 20000000, // size limit of a single IIP sequence
};

const termOptions = {
  allowProposedApi: true,
  windowsPty: isWindows
    ? {
        // In a real scenario, these values should be verified on the backend
        backend: "conpty",
        buildNumber: 22621,
      }
    : undefined,
  fontFamily:
    '"Fira Code", courier-new, courier, monospace, "Powerline Extra Symbols"',
  theme: xtermjsTheme,
};

const isWebGL2Available = () => {
  try {
    const canvas = document.createElement("canvas");
    return !!(window.WebGL2RenderingContext && canvas.getContext("webgl2"));
  } catch (e) {
    return false;
  }
};

const webSocket = () => {
  const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
  const path = window.location.pathname.replace(/[\/]+$/, "");
  const wsUrl = [
    protocol,
    "//",
    window.location.host,
    path,
    "/v1/ws",
    window.location.search,
  ].join("");
  const ws = new WebSocket(wsUrl, ["tty"]);
  return ws;
};

(function () {
  "use strict";
  const term = new Terminal(termOptions);

  const fitAddon = new FitAddon.FitAddon();
  const webglAddon = new WebglAddon.WebglAddon();
  const attachAddon = new AttachAddon.AttachAddon(webSocket());
  const unicode11Addon = new Unicode11Addon.Unicode11Addon();
  const webLinksAddon = new WebLinksAddon.WebLinksAddon();
  const clipboardAddon = new ClipboardAddon.ClipboardAddon();
  const imageAddon = new ImageAddon.ImageAddon(imageSettings);

  term.loadAddon(fitAddon);
  term.loadAddon(webglAddon);
  term.loadAddon(attachAddon);
  term.loadAddon(unicode11Addon);
  term.loadAddon(webLinksAddon);
  term.loadAddon(clipboardAddon);
  term.loadAddon(imageAddon);

  term.open(document.getElementById("terminal"));

  window.addEventListener("resize", () => {
    fitAddon.fit();
  });
  fitAddon.fit();

  // The browser may drop WebGL contexts for various reasons like OOM or after the system has been suspended.
  // There is an API exposed that fires the webglcontextlost event fired on the canvas so embedders can handle it however they wish.
  // An easy, but suboptimal way, to handle this is by disposing of WebglAddon when the event fires:
  webglAddon.onContextLoss((e) => {
    webglAddon.dispose();
  });

  // activate the new version
  term.unicode.activeVersion = "11";
})();
