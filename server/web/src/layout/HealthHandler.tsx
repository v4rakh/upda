import classes from './style/HealthHandler.module.less';

import { useGetHealthQuery } from '../api/healthApi';
import { App, ModalFuncProps, Skeleton } from 'antd';
import { FC, ReactNode, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

const HealthHandler: FC<{ children: ReactNode | ReactNode[] }> = ({ children }): ReactNode => {
	const [t] = useTranslation('health');
	const { modal: antModal } = App.useApp();

	const { isLoading, isSuccess, isError, data } = useGetHealthQuery(undefined);

	const [showChildren, setShowChildren] = useState<boolean>(false);
	const [showErrorModal, setShowErrorModal] = useState<boolean>(false);

	const [modal, setModal] = useState<{ update: (config: ModalFuncProps) => void; destroy: () => void }>();

	useEffect(() => {
		if (!isLoading && isSuccess && data?.data.healthy) {
			setShowChildren(true);
			setShowErrorModal(false);
		} else if (!isLoading && (isError || (!isError && !data?.data.healthy))) {
			setShowChildren(false);
			setShowErrorModal(true);
		}
	}, [isLoading, isSuccess, isError, data]);

	useEffect(() => {
		if (showErrorModal) {
			const title = <strong>{t('generic_error_title')}</strong>;
			const content = t('generic_error_content');
			const okText = t('reload');

			if (modal) {
				modal.update({ title, content });
			} else {
				const props: ModalFuncProps = {
					title,
					content,
					okButtonProps: { className: classes.okBtnHidden },
					okText,
					onOk: () => {
						window.location.reload();
					}
				};
				const error = antModal.error(props);
				setModal(error);
			}
		} else {
			modal?.destroy();
		}
	}, [t, modal, showErrorModal, antModal]);

	if (isLoading) {
		return <Skeleton loading={isLoading} active={isLoading} />;
	} else {
		if (showChildren) {
			return <>{children}</>;
		} else {
			return <></>;
		}
	}
};

export default HealthHandler;
