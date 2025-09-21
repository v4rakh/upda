import { useAuth } from '../../auth/AuthContext';
import AppPaths from '../../constants/appPaths';
import { AuthTypeSessionLoginRequest } from '../../types';
import { getPageFullPath } from '../../utils/urlHelper';
import AppBreadcrumb from '../common/AppBreadcrumb';
import { Button, Card, Flex, Form, Input, Space } from 'antd';
import { useCallback, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router';

const SessionLogin = () => {
	const [t] = useTranslation('session_login');
	const { login, isAuthenticated } = useAuth();
	const [form] = Form.useForm();
	const [isLoggingIn, setIsLoggingIn] = useState(false);
	const navigate = useNavigate();

	useEffect(() => {
		if (isAuthenticated) {
			navigate(getPageFullPath(AppPaths.HOME));
		}
	}, [isAuthenticated, navigate]);

	const preLogin = useCallback(() => {
		setIsLoggingIn(true);
	}, []);

	const postLogin = useCallback(() => {
		setIsLoggingIn(false);
	}, []);

	const sendLogin = useCallback(
		(values: AuthTypeSessionLoginRequest) => {
			login({ credentials: values, preLogin: preLogin, postLogin: postLogin });
		},
		[login, postLogin, preLogin]
	);

	const onFinish = useCallback(
		(values: { username: string; password: string }) => {
			sendLogin({ username: values.username, password: values.password } as AuthTypeSessionLoginRequest);
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
								<Button type="primary" htmlType="submit" loading={isLoggingIn}>
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

export default SessionLogin;
