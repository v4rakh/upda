import { UpdateMetadataResponse } from '../../types';
import { Typography } from 'antd';
import parse from 'html-react-parser';
import linkifyHtml from 'linkify-html';
import { FC } from 'react';

export interface UpdateMetadataProps {
	metadata: UpdateMetadataResponse;
}

const { Paragraph } = Typography;

const UpdateMetadata: FC<UpdateMetadataProps> = ({ metadata }): JSX.Element => {
	return (
		<Paragraph>
			<pre>
				{parse(
					linkifyHtml(JSON.stringify(metadata, null, 2), {
						target: '_blank',
						ignoreTags: ['script', 'style'],
						defaultProtocol: 'https',
						rel: 'noopener',
						attributes: {
							referrerpolicy: 'no-referrer'
						}
					})
				)}
			</pre>
		</Paragraph>
	);
};

export default UpdateMetadata;
