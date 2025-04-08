// Note that Object.freeze is NOT recursive
const runtime_config = Object.freeze({});

Object.defineProperty(window, 'runtime_config', {
    value: runtime_config,
    writable: false
});
