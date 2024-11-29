import { useGetProbeLoginMutation } from '../../api/loginApi';
import AppPaths from '../../constants/appPaths';
import { updateAuth } from '../../slices/authSlice';
import { useAppDispatch } from '../../store';
import { LoginRequest } from '../../types';
import { getPageFullPath } from '../../utils/urlHelper';
import { apiNotification } from '../common/apiNotification';
import AppBreadcrumb from '../common/AppBreadcrumb';
import { Button, Card, Flex, Form, Input, Space } from 'antd';
import { useCallback, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router';

const Login = () => {
	const [t] = useTranslation('login');
	const dispatch = useAppDispatch();
	const navigate = useNavigate();

	const [form] = Form.useForm();
	const [getProbeLogin, { isLoading, isSuccess, isError, error }] = useGetProbeLoginMutation();

	useEffect(() => {
		if (isError) {
			apiNotification.error({
				i18n: {
					unAuthorized: t('unauthorized'),
					default: t('default_message')
				},
				error: error
			});
		}
	}, [error, isError, t]);

	useEffect(() => {
		if (isSuccess) {
			const auth = { username: form.getFieldValue('username'), password: form.getFieldValue('password') };
			dispatch(updateAuth(auth));
			navigate(getPageFullPath(AppPaths.HOME));
		}
	}, [isSuccess, t, form, dispatch, navigate]);

	const sendLogin = useCallback(
		(values: LoginRequest) => {
			getProbeLogin(values);
		},
		[getProbeLogin]
	);

	const onFinish = useCallback(
		(values: any) => {
			sendLogin(values);
		},
		[sendLogin]
	);

	const onReset = useCallback(() => {
		form.resetFields();
	}, [form]);

	return (
		<>
			<AppBreadcrumb items={[{ label: t('title'), active: true, path: '' }]} />
			<Flex justify="center" align="center">
				<Card type="inner" title={t('title')}>
					<Form form={form} layout="vertical" name="login_form" onFinish={onFinish}>
						<Form.Item
							name={['username']}
							label={t('username')}
							rules={[{ required: true, message: t('username_required') }]}>
							<Input variant="filled" autoFocus />
						</Form.Item>
						<Form.Item
							name={['password']}
							label={t('password')}
							rules={[{ required: true, message: t('password_required') }]}>
							<Input.Password variant="filled" />
						</Form.Item>
						<Form.Item>
							<Space>
								<Button type="primary" htmlType="submit" loading={isLoading}>
									{t('submit')}
								</Button>
								<Button htmlType="button" onClick={onReset}>
									{t('reset')}
								</Button>
							</Space>
						</Form.Item>
					</Form>
				</Card>
			</Flex>
		</>
	);
};

export default Login;
