import Languages from '../constants/languages';
import Locales from '../constants/locales';
import { Col, ConfigProvider, Layout, Row, Spin } from 'antd';
import { Locale } from 'antd/lib/locale';
import en_US from 'antd/lib/locale/en_US';
// vite cannot bundle locales of moment, see https://github.com/moment/moment/issues/5926
import moment from 'moment-timezone';
import React, { createContext, FC, ReactNode, useCallback, useContext, useEffect, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';

export type LocaleProviderContextType = {
	language: string;
	locale: string;
};

type LocaleProviderProviderProps = {
	children: ReactNode | ReactNode[];
};

const LocaleContext = createContext<LocaleProviderContextType | undefined>(undefined);

const LocaleContextProvider: FC<LocaleProviderProviderProps> = ({ children }): React.JSX.Element => {
	const { i18n, t } = useTranslation();
	const [antLocale, setAntLocale] = useState<Locale>();
	const [language, setLanguage] = useState<string>(Languages.DEFAULT);
	const [locale, setLocale] = useState<string>(Locales.DEFAULT);

	const isLocale = useCallback((i18nLanguage: string) => {
		return i18nLanguage.includes('-');
	}, []);

	const getLocaleFromLanguage = useCallback((language: string) => {
		if (Languages.ENGLISH === language) {
			return Locales.ENGLISH_US;
		}

		return Locales.DEFAULT;
	}, []);

	const getLanguageFromLocale = useCallback((locale: string) => {
		const language = locale.split('-')[0];

		if (Languages.ENGLISH === language) {
			return Languages.ENGLISH;
		}
		return Languages.DEFAULT;
	}, []);

	const setLocales = useCallback(() => {
		let i18nLanguage = i18n.resolvedLanguage;
		let locale: string;
		let language: string;

		if (!i18nLanguage) {
			i18nLanguage = Languages.DEFAULT;
		}

		if (isLocale(i18nLanguage)) {
			language = getLanguageFromLocale(i18nLanguage);
			locale = getLocaleFromLanguage(language);
		} else {
			locale = getLocaleFromLanguage(i18nLanguage);
			language = getLanguageFromLocale(locale);
		}

		if (Languages.ENGLISH === language) {
			setAntLocale(en_US);
			moment.updateLocale(language, null);
		} else {
			setAntLocale(en_US);
		}

		setLocale(locale);
		setLanguage(language);
	}, [getLanguageFromLocale, getLocaleFromLanguage, i18n.resolvedLanguage, isLocale]);
	const renderProvider = useMemo(() => {
		return (
			<LocaleContext.Provider value={{ locale: locale, language: language }}>
				<ConfigProvider locale={antLocale}>{children}</ConfigProvider>
			</LocaleContext.Provider>
		);
	}, [antLocale, children, language, locale]);

	useEffect(() => {
		setLocales();
	}, [setLocales, i18n]);

	if (language && locale) {
		return renderProvider;
	}
	return (
		<Layout>
			<Layout.Content>
				<Row align="middle">
					<Col span={24}>
						<Spin tip={t('loading')} spinning />
					</Col>
				</Row>
			</Layout.Content>
		</Layout>
	);
};
export const useLocaleProviderContext = (): LocaleProviderContextType => {
	const context = useContext(LocaleContext);
	if (context === undefined) {
		throw new Error('useLocaleProviderContext must be used within LocaleProvider');
	}
	return context;
};

export default LocaleContextProvider;
