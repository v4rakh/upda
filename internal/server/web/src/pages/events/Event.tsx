import EventText from './EventText';
import { useDeleteEventMutation } from '../../api/eventsApi';
import DateTimeStyle from '../../constants/dateTimeStyle';
import { useLocaleProviderContext } from '../../providers/LocaleContextProvider';
import { EventResponse } from '../../types/event';
import { useNotification } from '../../use/useNotification';
import { formatDateTimeWithTimeZone } from '../../utils/datetimeHelper';
import { DeleteOutlined } from '@ant-design/icons';
import { Button, Card, Popconfirm, Tooltip, Typography } from 'antd';
import { FC, ReactNode, useEffect } from 'react';
import { useTranslation } from 'react-i18next';

const { Text } = Typography;

export interface EventProps {
	entity: EventResponse;
	onDeleteSuccess?: () => void;
}

const Event: FC<EventProps> = ({ entity, onDeleteSuccess }): ReactNode => {
	const [t] = useTranslation('event');
	const { apiError } = useNotification();
	const { locale } = useLocaleProviderContext();

	const [deleteEvent, { isSuccess, isLoading, isError, error }] = useDeleteEventMutation();

	useEffect(() => {
		if (isError) {
			apiError({
				i18n: {
					notFound: t('error_unable_delete'),
					unAuthorized: t('error_unauthorized_delete'),
					forbidden: t('error_forbidden_delete'),
					default: t('error_default_delete')
				},
				error: error
			});
		}

		if (isSuccess && onDeleteSuccess) {
			onDeleteSuccess();
		}
	}, [isError, error, isSuccess, onDeleteSuccess, t, apiError]);

	return (
		<Card
			loading={isLoading}
			size="small"
			actions={[
				<Text key={`${entity.id}_created`} italic type="secondary">
					{formatDateTimeWithTimeZone(entity.createdAt, DateTimeStyle.LONG, DateTimeStyle.MEDIUM, locale)}
				</Text>,
				<Popconfirm
					key={`${entity.id}_del_confirm`}
					title={t('delete_title')}
					onConfirm={() => deleteEvent({ id: entity.id })}
					okText={t('delete')}
					placement="bottom"
					cancelText={t('cancel')}
					okButtonProps={{ icon: <DeleteOutlined />, type: 'primary', danger: true }}>
					<Tooltip title={t('help_delete')} placement="bottom">
						<Button key="del" size="small" icon={<DeleteOutlined />} type="text" danger />
					</Tooltip>
				</Popconfirm>
			]}>
			<EventText name={entity.name} payload={entity.payload} />
		</Card>
	);
};

export default Event;
