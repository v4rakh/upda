import classes from './style/HealthHandler.module.less';

import { useGetHealthQuery } from '../api/healthApi';
import { Modal, ModalFuncProps, Skeleton } from 'antd';
import { FC, ReactNode, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

const HealthHandler: FC<{ children: ReactNode | ReactNode[] }> = ({ children }): JSX.Element => {
	const [t] = useTranslation('health');
	const { isLoading, isSuccess, isError, data } = useGetHealthQuery(undefined);

	const [showChildren, setShowChildren] = useState<boolean>(false);
	const [showModal, setShowModal] = useState<boolean>(false);

	const [modal, setModal] = useState<{ update: (config: ModalFuncProps) => void; destroy: () => void }>();

	useEffect(() => {
		if (!isLoading && isSuccess && data?.data.healthy) {
			setShowChildren(true);
			setShowModal(false);
		} else if (!isLoading && (isError || (!isError && !data?.data.healthy))) {
			setShowChildren(false);
			setShowModal(true);
		}
	}, [isLoading, isSuccess, isError, data]);

	useEffect(() => {
		if (showModal) {
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
				const error = Modal.error(props);
				setModal(error);
			}
		} else {
			modal?.destroy();
		}
	}, [t, modal, showModal]);

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
