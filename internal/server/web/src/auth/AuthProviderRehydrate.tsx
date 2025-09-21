import { useAuth } from './AuthContext';
import { LoadingOutlined } from '@ant-design/icons';
import { Col, Row, Spin } from 'antd';
import React, { FC, ReactNode, useCallback, useEffect, useState } from 'react';

type AuthProviderRehydrateState = {
	initialLoad: boolean;
	isRehydrating: boolean;
};

const initialAuthProviderRehydrateState = {
	initialLoad: true,
	isRehydrating: true
} as AuthProviderRehydrateState;

export const AuthProviderRehydrate: FC<{ children: ReactNode }> = ({ children }) => {
	const { rehydrate } = useAuth();

	const [rehydrateState, setRehydrateState] = useState<AuthProviderRehydrateState>(initialAuthProviderRehydrateState);

	const postRehydrate = useCallback(() => {
		setRehydrateState({ initialLoad: false, isRehydrating: false });
	}, []);

	useEffect(() => {
		if (rehydrateState.initialLoad) {
			rehydrate({ postRehydrate: postRehydrate });
		}
	}, [rehydrate, postRehydrate, rehydrateState.initialLoad]);

	if (rehydrateState.isRehydrating) {
		return (
			<Row justify="center" align="middle" style={{ height: '100vh' }}>
				<Col>
					<Spin fullscreen size="large" indicator={<LoadingOutlined />} />
				</Col>
			</Row>
		);
	}

	return <>{children}</>;
};
