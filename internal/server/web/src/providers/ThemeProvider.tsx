import LocalStorageKeys from '../constants/localStorageKeys';
import { useLocalStorage } from '../use/useLocalStorage';
import { ConfigProvider, theme } from 'antd';
import { createContext, FC, ReactNode, useCallback, useContext, useEffect } from 'react';

interface ThemeContextType {
	isDarkTheme: boolean;
	toggleTheme: () => void;
}

const ThemeContext = createContext<ThemeContextType | undefined>(undefined);

interface ThemeProviderProps {
	children: ReactNode | ReactNode[];
}

const ThemeProvider: FC<ThemeProviderProps> = ({ children }): ReactNode => {
	const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
	const [isDarkTheme, setIsDarkTheme] = useLocalStorage<boolean>(LocalStorageKeys.DARK_THEME_ENABLED, prefersDark);

	const toggleTheme = useCallback(() => {
		setIsDarkTheme((prev) => !prev);
	}, [setIsDarkTheme]);

	useEffect(() => {
		document.documentElement.setAttribute('data-color-mode', isDarkTheme ? 'dark' : 'light');
	}, [isDarkTheme]);

	const algorithm = isDarkTheme ? theme.darkAlgorithm : theme.defaultAlgorithm;

	return (
		<ThemeContext.Provider value={{ isDarkTheme, toggleTheme }}>
			<ConfigProvider theme={{ algorithm }}>{children}</ConfigProvider>
		</ThemeContext.Provider>
	);
};

export const useTheme = (): ThemeContextType => {
	const context = useContext(ThemeContext);
	if (context === undefined) {
		throw new Error('useTheme must be used within ThemeProvider');
	}
	return context;
};

export default ThemeProvider;
