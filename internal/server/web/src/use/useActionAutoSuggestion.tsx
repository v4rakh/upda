import { useLazyGetConstantsQuery } from '../api/constantsApi';
import { useLazyGetSecretsQuery } from '../api/secretsApi';
import { EventName } from '../types/event';
import { CalculatorOutlined, FolderOutlined, LockOutlined } from '@ant-design/icons';
import { Space } from 'antd';
import { MentionsOptionProps } from 'antd/lib/mentions';
import { concat, map } from 'lodash';
import { useCallback, useEffect, useMemo, useState } from 'react';

const BASE_VAR_OPTIONS = ['APPLICATION', 'PROVIDER', 'HOST', 'VERSION', 'STATE'];
const COMMENT_VAR_OPTIONS = ['COMMENT_AUTHOR', 'COMMENT_CONTENT'];

const createVarOptions = (vars: string[]): MentionsOptionProps[] =>
	map(vars, (e) => ({
		value: `VAR>${e}</VAR>`,
		label: (
			<Space>
				<CalculatorOutlined style={{ color: 'blue' }} />
				{e}
			</Space>
		)
	})) as MentionsOptionProps[];

export interface AutoSuggestionHook {
	mentionOptions: MentionsOptionProps[];
	reloadMentionOptions: () => void;
}

export interface UseActionAutoSuggestionProps {
	matchEvent?: EventName | string;
}

const useActionAutoSuggestion = (props?: UseActionAutoSuggestionProps): AutoSuggestionHook => {
	const { matchEvent } = props || {};
	const [suggestions, setSuggestions] = useState<MentionsOptionProps[]>([]);
	const [getConstants, constantsResult] = useLazyGetConstantsQuery();
	const [getSecrets, secretsResult] = useLazyGetSecretsQuery();
	const [optionsReady, setOptionsReady] = useState(false);

	const varOptions = useMemo(() => {
		const vars =
			!matchEvent || matchEvent === EventName.COMMENT_CREATED
				? concat(BASE_VAR_OPTIONS, COMMENT_VAR_OPTIONS)
				: BASE_VAR_OPTIONS;
		return createVarOptions(vars);
	}, [matchEvent]);

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

			setSuggestions(concat(varOptions, constantOptions, secretOptions));
			setOptionsReady(true);
		}
	}, [constantsResult, secretsResult, varOptions]);

	const onRefresh = useCallback(() => {
		setOptionsReady(false);
	}, []);

	return {
		mentionOptions: suggestions,
		reloadMentionOptions: onRefresh
	};
};

export default useActionAutoSuggestion;
