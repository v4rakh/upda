import { getApiErrorMessage } from '../utils/apiHelper';
import { App } from 'antd';
import { ReactNode, useCallback } from 'react';

const key = 'updatable';

export interface NotificationErrorProps {
	i18n: {
		default: ReactNode;
		conflict?: ReactNode;
		badRequest?: ReactNode;
		notFound?: ReactNode;
		unAuthorized?: ReactNode;
		forbidden?: ReactNode;
	};
	error: any;
}

export interface NotificationSimpleErrorProps {
	message: ReactNode;
	title?: ReactNode;
}

export interface NotificationSimpleInfoProps {
	message: ReactNode;
	title?: ReactNode;
	duration?: number | null;
}

export interface NotificationHook {
	apiError: (props: NotificationErrorProps) => void;
	simpleError: (props: NotificationSimpleErrorProps) => void;
	simpleInfo: (props: NotificationSimpleInfoProps) => void;
}

export const useNotification = (): NotificationHook => {
	const { notification } = App.useApp();

	const apiError = useCallback(
		(props: NotificationErrorProps) => {
			notification.error({
				placement: 'top',
				duration: null,
				message: '',
				description: <span style={{ paddingRight: '10px' }}>{getApiErrorMessage(props)}</span>,
				key
			});
		},
		[notification]
	);

	const simpleError = useCallback(
		({ message, title }: NotificationSimpleErrorProps) => {
			notification.error({
				placement: 'top',
				duration: null,
				message: title || '',
				description: <span style={{ paddingRight: '5px' }}>{message}</span>,
				key
			});
		},
		[notification]
	);

	const simpleInfo = useCallback(
		({ message, title, duration }: NotificationSimpleInfoProps) => {
			notification.info({
				placement: 'top',
				duration: duration,
				message: title || '',
				description: <span style={{ paddingRight: '5px' }}>{message}</span>,
				key
			});
		},
		[notification]
	);

	return {
		apiError,
		simpleError,
		simpleInfo
	};
};
