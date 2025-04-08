import HttpStatusCode from '../constants/httpStatusCode';
import { FetchBaseQueryError } from '@reduxjs/toolkit/query';
import { ReactNode } from 'react';

export const getApiErrorMessage = ({
	i18n,
	error
}: {
	i18n: {
		default: ReactNode;
		conflict?: ReactNode;
		badRequest?: ReactNode;
		notFound?: ReactNode;
		unAuthorized?: ReactNode;
		forbidden?: ReactNode;
	};
	error: any;
}) => {
	let message = i18n.default;
	if ((error as FetchBaseQueryError)?.status) {
		if (error.status === HttpStatusCode.STATUS_409 || error.originalStatus === HttpStatusCode.STATUS_409) {
			message = i18n.conflict || i18n.default;
		} else if (error.status === HttpStatusCode.STATUS_404 || error.originalStatus === HttpStatusCode.STATUS_404) {
			message = i18n.notFound || i18n.default;
		} else if (error.status === HttpStatusCode.STATUS_401 && i18n.forbidden) {
			message = i18n.forbidden;
		} else if (error.status === HttpStatusCode.STATUS_401 || error.status === HttpStatusCode.STATUS_403) {
			message = i18n.unAuthorized || i18n.default;
		} else if (error.status === HttpStatusCode.STATUS_400) {
			message = i18n.badRequest || i18n.default;
		}
	}
	return message;
};

export const convertToLowerCaseUnderscore = (str: string): string =>
	str
		.replace(/([a-z])([A-Z])/g, '$1_$2')
		.replace(/[\s_]+/g, '_')
		.toLowerCase();
