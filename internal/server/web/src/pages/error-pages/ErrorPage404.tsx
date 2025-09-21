import AppPaths from '../../constants/appPaths';
import { getPageFullPath } from '../../utils/urlHelper';
import { Button, Result } from 'antd';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router';

const ErrorPage404 = () => {
	const [t] = useTranslation('error');
	const navigate = useNavigate();

	return (
		<Result
			status="404"
			title={t('404_title')}
			subTitle={t('404_description')}
			extra={
				<Button type="primary" onClick={() => navigate(getPageFullPath(AppPaths.HOME))}>
					{t('back_home')}
				</Button>
			}
		/>
	);
};
export default ErrorPage404;
