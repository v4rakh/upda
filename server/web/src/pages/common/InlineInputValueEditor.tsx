import { CheckOutlined } from '@ant-design/icons';
import { Button, Input, Space } from 'antd';
import type { InputStatus } from 'antd/es/_util/statusUtils';
import { Variant } from 'antd/es/config-provider';
import { FC, FormEvent, useEffect, useState } from 'react';

const { Compact } = Space;

export interface InlineInputValueEditorProps {
	initialValue?: string;
	placeholder?: string;
	allowBlank?: boolean;
	variant?: Variant;
	isLoading: boolean;
	isError: boolean;
	isSuccess: boolean;
	resetOnSuccess?: boolean;
	resetOnError?: boolean;
	onSubmit: (val?: string) => void;
}

interface InlineInputValueEditorState {
	currentValue?: string;
	currentStatus: InputStatus;
	submitDisabled: boolean;
}

const InlineInputValueEditor: FC<InlineInputValueEditorProps> = ({
	initialValue,
	placeholder = undefined,
	allowBlank = true,
	resetOnSuccess = false,
	resetOnError = false,
	variant = 'filled',
	isLoading,
	isError,
	isSuccess,
	onSubmit
}) => {
	const [state, setState] = useState<InlineInputValueEditorState>({
		currentValue: initialValue,
		currentStatus: '',
		submitDisabled: true
	});

	useEffect(() => {
		if (isError || isSuccess) {
			const resetValueError = resetOnError ? '' : state.currentValue;
			const resetValueSuccess = resetOnSuccess ? '' : state.currentValue;

			setState({
				currentValue: isSuccess ? resetValueSuccess : resetValueError,
				currentStatus: '',
				submitDisabled: true
			});
		}
	}, [setState, isError, isSuccess, resetOnError, state.currentValue, resetOnSuccess]);

	const onChange = (e: FormEvent<HTMLInputElement>) => {
		const newVal = e.currentTarget.value;

		const isBlank = newVal === '' || newVal == undefined;

		let inputStatus = '' as InputStatus;
		if (isBlank && !allowBlank) {
			inputStatus = 'error';
		}

		setState({
			currentValue: newVal,
			currentStatus: inputStatus,
			submitDisabled: isLoading || newVal === initialValue || inputStatus === 'error'
		});
	};

	const submit = () => {
		if (!state.submitDisabled) {
			onSubmit(state.currentValue);
		}
	};

	return (
		<Compact>
			<Input
				placeholder={placeholder}
				variant={variant}
				status={state.currentStatus}
				defaultValue={state.currentValue}
				value={state.currentValue}
				onChange={onChange}
				allowClear={allowBlank}
				onPressEnter={submit}
			/>
			<Button
				loading={isLoading}
				disabled={state.submitDisabled}
				type="primary"
				onClick={submit}
				icon={<CheckOutlined />}
			/>
		</Compact>
	);
};

export default InlineInputValueEditor;
