const mockConfig = {};

Object.defineProperty(window, 'runtime_config', {
	writable: true,
	value: mockConfig
});

export { mockConfig };
