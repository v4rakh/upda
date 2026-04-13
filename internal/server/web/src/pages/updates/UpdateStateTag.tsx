import { useGetUpdateStateDefinitionsQuery } from '../../api/updateStateDefinitionsApi';
import { UpdateStateValue } from '../../types';
import { renderIcon } from '../../utils/iconHelper';
import {
	getUpdateStateColorFromDefinitions,
	getUpdateStateIconFromDefinitions,
	getUpdateStateLabelFromDefinitions
} from '../../utils/updateHelper';
import { Tag } from 'antd';
import { FC, ReactNode } from 'react';

export interface UpdateStateProps {
	state: UpdateStateValue;
}

const UpdateStateTag: FC<UpdateStateProps> = ({ state }): ReactNode => {
	const { data: statesData } = useGetUpdateStateDefinitionsQuery();
	const definitions = statesData?.data?.content;

	const color = getUpdateStateColorFromDefinitions(state, definitions);
	const label = getUpdateStateLabelFromDefinitions(state, definitions);
	const icon = getUpdateStateIconFromDefinitions(state, definitions);

	return (
		<Tag color={color} icon={renderIcon(icon, { marginRight: 4 })}>
			{label}
		</Tag>
	);
};
export default UpdateStateTag;
