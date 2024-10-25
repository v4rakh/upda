import { getApiErrorMessage } from '../../utils/apiHelper';
import { notification } from 'antd';
import { ReactNode } from 'react';

const key = 'updatable';

export const apiNotification = {
	error: (props: {
		i18n: {
			default: ReactNode;
			badRequest?: ReactNode;
			notFound?: ReactNode;
			unAuthorized?: ReactNode;
			forbidden?: ReactNode;
		};
		error: any;
	}) =>
		notification.error({
			placement: 'top',
			duration: null,
			message: '',
			description: <span style={{ paddingRight: '10px' }}>{getApiErrorMessage(props)}</span>,
			key
		}),
	simpleError: ({ message, title }: { message: ReactNode; title?: ReactNode }) => {
		return notification.error({
			placement: 'top',
			duration: null,
			message: title || '',
			description: <span style={{ paddingRight: '5px' }}>{message}</span>,
			key
		});
	},
	simpleInfo: ({
		message,
		title,
		duration = 3
	}: {
		message: ReactNode;
		title?: ReactNode;
		duration?: number | null;
	}) => {
		return notification.info({
			placement: 'top',
			duration: duration,
			message: title || '',
			description: <span style={{ paddingRight: '5px' }}>{message}</span>,
			key
		});
	}
};
