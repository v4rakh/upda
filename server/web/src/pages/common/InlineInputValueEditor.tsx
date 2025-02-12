import { CheckOutlined } from '@ant-design/icons';
import { Button, Input, Space } from 'antd';
import type { InputStatus } from 'antd/es/_util/statusUtils';
import { Variant } from 'antd/es/config-provider';
import { FC, FormEvent, useCallback, useEffect, useState } from 'react';

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

const DEFAULT_VALIDATION_STATUS = '';
const DEFAULT_SUBMIT_DISABLED = true;

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
	const [submitDisabled, setSubmitDisabled] = useState(DEFAULT_SUBMIT_DISABLED);
	const [value, setValue] = useState(initialValue);
	const [validationStatus, setValidationStatus] = useState<InputStatus>(DEFAULT_VALIDATION_STATUS);

	useEffect(() => {
		if (isError || isSuccess) {
			const resetValueError = resetOnError ? '' : value;
			const resetValueSuccess = resetOnSuccess ? '' : value;

			setValue(isSuccess ? resetValueSuccess : resetValueError);
			setValidationStatus(DEFAULT_VALIDATION_STATUS);
			setSubmitDisabled(DEFAULT_SUBMIT_DISABLED);
		}
	}, [isError, isSuccess, resetOnError, resetOnSuccess, value]);

	useEffect(() => {
		if (isLoading || initialValue === value || validationStatus === 'error') {
			setSubmitDisabled(true);
		} else {
			setSubmitDisabled(false);
		}
	}, [initialValue, isLoading, validationStatus, value]);

	const onChange = useCallback(
		(e: FormEvent<HTMLInputElement>) => {
			const newVal = e.currentTarget.value;

			const isBlank = newVal === '' || newVal == undefined;

			let inputStatus = '' as InputStatus;
			if (isBlank && !allowBlank) {
				inputStatus = 'error';
			}

			setValue(newVal);
			setValidationStatus(inputStatus);
		},
		[allowBlank]
	);

	const submit = () => {
		if (!submitDisabled) {
			onSubmit(value);
		}
	};

	return (
		<Compact>
			<Input
				placeholder={placeholder}
				variant={variant}
				status={validationStatus}
				value={value}
				onChange={onChange}
				allowClear={allowBlank}
				onPressEnter={submit}
			/>
			<Button
				loading={isLoading}
				disabled={submitDisabled}
				type="primary"
				onClick={submit}
				icon={<CheckOutlined />}
			/>
		</Compact>
	);
};

export default InlineInputValueEditor;
