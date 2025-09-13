import { useLazyGetConstantsQuery } from '../api/constantsApi';
import { useLazyGetSecretsQuery } from '../api/secretsApi';
import { CalculatorOutlined, FolderOutlined, LockOutlined } from '@ant-design/icons';
import { Space } from 'antd';
import { MentionsOptionProps } from 'antd/lib/mentions';
import { concat, map } from 'lodash';
import { useCallback, useEffect, useState } from 'react';

const VAR_OPTIONS = map(['APPLICATION', 'PROVIDER', 'HOST', 'VERSION', 'STATE'], (e) => {
	return {
		value: `VAR>${e}</VAR>`,
		label: (
			<Space>
				<CalculatorOutlined style={{ color: 'blue' }} />
				{e}
			</Space>
		)
	};
}) as MentionsOptionProps[];

export interface AutoSuggestionHook {
	mentionOptions: MentionsOptionProps[];
	reloadMentionOptions: () => void;
}

const useActionAutoSuggestion = (): AutoSuggestionHook => {
	const [suggestions, setSuggestions] = useState<MentionsOptionProps[]>([]);
	const [getConstants, constantsResult] = useLazyGetConstantsQuery();
	const [getSecrets, secretsResult] = useLazyGetSecretsQuery();
	const [optionsReady, setOptionsReady] = useState(false);

	useEffect(() => {
		if (!optionsReady) {
			getConstants();
			getSecrets();
		}
	}, [getConstants, getSecrets, optionsReady]);

	useEffect(() => {
		if (constantsResult?.data?.data.content && secretsResult?.data?.data.content) {
			const constantOptions = map(constantsResult?.data?.data.content, (e) => {
				return {
					value: `CONST>${e.key}</CONST>`,
					label: (
						<Space>
							<FolderOutlined style={{ color: 'gray' }} />
							{e.key}
						</Space>
					)
				};
			}) as MentionsOptionProps[];

			const secretOptions = map(secretsResult?.data?.data.content, (e) => {
				return {
					value: `SECRET>${e.key}</SECRET>`,
					label: (
						<Space>
							<LockOutlined style={{ color: 'red' }} />
							{e.key}
						</Space>
					)
				};
			}) as MentionsOptionProps[];

			setSuggestions(concat(VAR_OPTIONS, constantOptions, secretOptions));
			setOptionsReady(true);
		}
	}, [constantsResult, secretsResult]);

	const onRefresh = useCallback(() => {
		setOptionsReady(false);
	}, []);

	return {
		mentionOptions: suggestions,
		reloadMentionOptions: onRefresh
	};
};

export default useActionAutoSuggestion;
