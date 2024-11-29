import i18next from 'i18next';
import LanguageDetector from 'i18next-browser-languagedetector';
import { initReactI18next } from 'react-i18next';

import en from './translations/en.json';
import { isDevelopment } from '../utils/envHelper';

const resources = {
	en: {
		...en
	}
};

i18next
	.use(initReactI18next)
	.use(LanguageDetector)
	.init(
		{
			resources,
			lng: 'en',
			fallbackLng: 'en',
			debug: isDevelopment(),
			interpolation: {
				escapeValue: false // not needed for react as it escapes by default
			},
			defaultNS: 'common',
			fallbackNS: 'common',
			detection: {
				order: ['navigator']
			}
		},
		undefined
	);

export default i18next;
