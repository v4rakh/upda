import { EventName } from '../../types/event';
import { Select } from 'antd';
import { FC, ReactNode } from 'react';
import { useTranslation } from 'react-i18next';

export interface ActionSelectEventProps {
	name?: EventName;
	onChange?: (selected?: string) => void;
	loading?: boolean;
}

const noop = () => {
	return;
};

const ActionSelectEvent: FC<ActionSelectEventProps> = ({ name, onChange = noop, loading }): ReactNode => {
	const [t] = useTranslation('action_select_event');

	return (
		<Select
			loading={loading}
			disabled={loading}
			allowClear
			style={{ width: 200 }}
			placeholder={t('all')}
			defaultValue={name}
			variant="filled"
			onChange={onChange}
			options={[
				{ value: EventName.UPDATE_CREATED, label: `${t(EventName.UPDATE_CREATED.toLowerCase())}` },
				{ value: EventName.UPDATE_UPDATED, label: `${t(EventName.UPDATE_UPDATED.toLowerCase())}` },
				{ value: EventName.UPDATE_UPDATED_STATE, label: `${t(EventName.UPDATE_UPDATED_STATE.toLowerCase())}` },
				{
					value: EventName.UPDATE_UPDATED_VERSION,
					label: `${t(EventName.UPDATE_UPDATED_VERSION.toLowerCase())}`
				},
				{ value: EventName.UPDATE_DELETED, label: `${t(EventName.UPDATE_DELETED.toLowerCase())}` },
				{ value: EventName.COMMENT_CREATED, label: `${t(EventName.COMMENT_CREATED.toLowerCase())}` }
			]}
		/>
	);
};

export default ActionSelectEvent;
