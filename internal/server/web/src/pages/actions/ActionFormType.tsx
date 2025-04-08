import { ActionType } from '../../types';
import { Select, Form } from 'antd';
import { FC, ReactNode } from 'react';
import { useTranslation } from 'react-i18next';

export interface ActionFormTypeProps {
	isLoading: boolean;
	initialValue?: ActionType;
}

const ActionFormType: FC<ActionFormTypeProps> = ({ isLoading, initialValue }): ReactNode => {
	const [t] = useTranslation('action_form_type');

	return (
		<Form.Item name="type" label={t('type_label')} tooltip={t('type_help')} required={true}>
			<Select
				value={initialValue}
				loading={isLoading}
				disabled={isLoading}
				style={{ width: 150 }}
				variant="filled"
				options={[{ value: ActionType.SHOUTRRR, label: `${t(ActionType.SHOUTRRR.toLowerCase())}` }]}
			/>
		</Form.Item>
	);
};

export default ActionFormType;
