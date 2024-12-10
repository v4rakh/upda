import { UpdateState } from '../../types';
import { getUpdateStateColor } from '../../utils/updateHelper';
import { Tag } from 'antd';
import { FC, ReactNode } from 'react';
import { useTranslation } from 'react-i18next';

export interface UpdateStateProps {
	state: UpdateState;
}

const UpdateStateTag: FC<UpdateStateProps> = ({ state }): ReactNode => {
	const [t] = useTranslation('update_state_tag');

	return <Tag color={getUpdateStateColor(state)}>{t(`state_${state.toLocaleLowerCase()}`)}</Tag>;
};
export default UpdateStateTag;
