const termOptions = {
    fontSize: 13,
    fontFamily: 'Menlo For Powerline,Consolas,Liberation Mono,Menlo,Courier,monospace',
    theme: {
        foreground: '#d2d2d2',
        background: '#2b2b2b',
        cursor: '#adadad',
        black: '#000000',
        red: '#d81e00',
        green: '#5ea702',
        yellow: '#cfae00',
        blue: '#427ab3',
        magenta: '#89658e',
        cyan: '#00a7aa',
        white: '#dbded8',
        brightBlack: '#686a66',
        brightRed: '#f54235',
        brightGreen: '#99e343',
        brightYellow: '#fdeb61',
        brightBlue: '#84b0d8',
        brightMagenta: '#bc94b7',
        brightCyan: '#37e6e8',
        brightWhite: '#f1f1f0',
    },
};

const isWebGL2Available = () => {
    try {
        const canvas = document.createElement('canvas');
        return !!(window.WebGL2RenderingContext && canvas.getContext('webgl2'));
    } catch (e) {
        return false;
    }
};

(function () {
    'use strict'
    const term = new Terminal(termOptions);
    const fitAddon = new FitAddon.FitAddon();
    term.loadAddon(fitAddon);
    window.addEventListener('resize', () => {
        fitAddon.fit();
    });

    term.open(document.getElementById('terminal'));
    fitAddon.fit();

    if (isWebGL2Available()) {
        term.loadAddon(new WebglAddon.WebglAddon());
    }

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const path = window.location.pathname.replace(/[\/]+$/, '');
    const wsUrl = [protocol, '//', window.location.host, path, '/v1/ws', window.location.search].join('');
    const ws = new WebSocket(wsUrl, ['tty'])
    const attachAddon = new AttachAddon.AttachAddon(ws);
    term.loadAddon(attachAddon);
    // term.write('\x1B[1;3;31mPress Enter to Continue\x1B[0m')
})();