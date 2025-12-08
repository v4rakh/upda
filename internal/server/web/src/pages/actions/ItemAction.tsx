import DeleteAction from './DeleteAction';
import TestAction from './TestAction';
import UpdateMatchApplicationAction from './UpdateMatchApplicationAction';
import UpdateMatchEventAction from './UpdateMatchEventAction';
import UpdateMatchHostAction from './UpdateMatchHostAction';
import UpdateMatchProviderAction from './UpdateMatchProviderAction';
import UpdatePayloadAction from './UpdatePayloadAction';
import { ActionResponse } from '../../types';
import { Descriptions, Space } from 'antd';
import { FC, ReactNode } from 'react';
import { useTranslation } from 'react-i18next';

const { Item } = Descriptions;

export interface ItemActionProps {
	e: ActionResponse;
}

const ItemActionInvocation: FC<ItemActionProps> = ({ e }): ReactNode => {
	const [t] = useTranslation('action_item');

	return (
		<Space orientation="vertical">
			<Descriptions
				size="small"
				layout="vertical"
				column={{ xs: 1, sm: 2, md: 4, lg: 4, xl: 4, xxl: 4 }}
				colon={false}>
				<Item label={t('match_event')}>
					<UpdateMatchEventAction id={e.id} matchEvent={e.matchEvent} />
				</Item>
				<Item label={t('match_application')}>
					<UpdateMatchApplicationAction id={e.id} matchApplication={e.matchApplication} />
				</Item>
				<Item label={t('match_host')}>
					<UpdateMatchHostAction id={e.id} matchHost={e.matchHost} />
				</Item>
				<Item label={t('match_provider')}>
					<UpdateMatchProviderAction id={e.id} matchProvider={e.matchProvider} />
				</Item>
			</Descriptions>
			<UpdatePayloadAction id={e.id} type={e.type} payload={e.payload} />
			<TestAction id={e.id} />
			<DeleteAction id={e.id} />
		</Space>
	);
};

export default ItemActionInvocation;
