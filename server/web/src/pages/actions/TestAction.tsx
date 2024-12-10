import { useTestActionMutation } from '../../api/actionsApi';
import { apiNotification } from '../common/apiNotification';
import { PlayCircleTwoTone } from '@ant-design/icons';
import { Button } from 'antd';
import { FC, ReactNode, useCallback, useEffect } from 'react';
import { useTranslation } from 'react-i18next';

export interface TestActionProps {
	id: string;
}

const TestAction: FC<TestActionProps> = ({ id }): ReactNode => {
	const [t] = useTranslation('action_test');

	const [test, { data, isLoading, isSuccess, isError, error }] = useTestActionMutation();

	const onClick = useCallback(() => {
		test({
			id,
			body: {
				application: t('application'),
				host: t('host'),
				provider: t('provider'),
				version: t('version'),
				state: t('state')
			}
		});
	}, [test, id, t]);

	useEffect(() => {
		if (isSuccess) {
			if (data.data.success) {
				apiNotification.simpleInfo({
					title: t('tested_title_success'),
					message: t('tested_message_success'),
					duration: 5
				});
			} else {
				apiNotification.simpleError({
					title: t('tested_title_error'),
					message: t('tested_message_error', { reason: data?.data.message })
				});
			}
		}
	}, [data?.data.message, data?.data.success, isSuccess, t]);

	useEffect(() => {
		if (isError) {
			apiNotification.error({
				i18n: {
					notFound: t('error_unable'),
					unAuthorized: t('error_unauthorized'),
					forbidden: t('error_forbidden'),
					badRequest: t('error_bad_request'),
					default: t('error_default')
				},
				error: error
			});
		}
	}, [isError, error, t]);

	return (
		<Button
			loading={isLoading}
			key="test"
			icon={<PlayCircleTwoTone twoToneColor={'green'} />}
			type="text"
			onClick={onClick}>
			{t('test')}
		</Button>
	);
};

export default TestAction;
