vi.mock('react-i18next', () => ({
	useTranslation: (): [(key: string) => string] => [(key: string): string => key]
}));
export {};
