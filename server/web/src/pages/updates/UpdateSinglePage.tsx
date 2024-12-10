import UpdateMetadata from './UpdateMetadata';
import UpdateStateTag from './UpdateStateTag';
import { useGetUpdateByIdQuery } from '../../api/updatesApi';
import ApiErrorCodes from '../../constants/apiErrorCodes';
import AppPathParamNames from '../../constants/appPathParamNames';
import AppPaths from '../../constants/appPaths';
import DateTimeStyle from '../../constants/dateTimeStyle';
import { useLocaleProviderContext } from '../../providers/LocaleContextProvider';
import { formatDateTimeWithTimeZone } from '../../utils/datetimeHelper';
import { getPageFullPath } from '../../utils/urlHelper';
import AppBreadcrumb from '../common/AppBreadcrumb';
import { PageHeader } from '@ant-design/pro-layout';
import { Descriptions, Result, Skeleton, Typography } from 'antd';
import { FC, ReactNode } from 'react';
import { useTranslation } from 'react-i18next';
import { useParams } from 'react-router';

const { Text } = Typography;

const UpdateSinglePage: FC = (): ReactNode => {
	const [t] = useTranslation('updates_single');
	const { locale } = useLocaleProviderContext();

	const { [AppPathParamNames.UPDATE_ID]: updateId } = useParams();

	const { isLoading, isError, isSuccess, data, error } = useGetUpdateByIdQuery(
		{
			id: updateId || ''
		},
		{
			skip: updateId === undefined
		}
	);

	return (
		<>
			<AppBreadcrumb
				items={[
					{
						label: t('title_parent'),
						path: getPageFullPath(AppPaths.UPDATES)
					},
					{
						label: data?.data.application || '',
						active: true,
						path: ''
					}
				]}
			/>
			<PageHeader
				className={'pl-0'}
				title={
					<Typography.Title level={4} ellipsis>
						{t('title')}
					</Typography.Title>
				}
			/>
			{isLoading && <Skeleton loading={isLoading} active={isLoading} />}
			{isError && (error as any)?.data?.status !== ApiErrorCodes.NOT_FOUND && (
				<Result status={500} title={t('default_loading_error_message')} />
			)}
			{isError && (error as any)?.data?.status === ApiErrorCodes.NOT_FOUND && (
				<Result status={404} title={t('no_update')} />
			)}
			{isSuccess && (
				<>
					<Descriptions
						layout="vertical"
						size="small"
						column={{ xs: 1, sm: 2, md: 4, lg: 4, xl: 4, xxl: 4 }}
						colon={false}>
						<Descriptions.Item label={t('state')}>
							<UpdateStateTag state={data.data.state} />
						</Descriptions.Item>
						<Descriptions.Item label={t('application')}>{data.data.application}</Descriptions.Item>
						<Descriptions.Item label={t('version')}>
							<Text code>{data.data.version}</Text>
						</Descriptions.Item>
						<Descriptions.Item label={t('host')}>{data.data.host}</Descriptions.Item>
						<Descriptions.Item label={t('provider')}>{data.data.provider}</Descriptions.Item>
						<Descriptions.Item label={t('created')}>
							{formatDateTimeWithTimeZone(
								data.data.createdAt,
								DateTimeStyle.LONG,
								DateTimeStyle.MEDIUM,
								locale
							)}
						</Descriptions.Item>
						<Descriptions.Item label={t('updated')}>
							{formatDateTimeWithTimeZone(
								data.data.updatedAt,
								DateTimeStyle.LONG,
								DateTimeStyle.MEDIUM,
								locale
							)}
						</Descriptions.Item>
					</Descriptions>
					<Descriptions
						layout="vertical"
						size="small"
						column={{ xs: 1, sm: 2, md: 4, lg: 4, xl: 4, xxl: 4 }}
						colon={false}>
						<Descriptions.Item label={t('metadata')}>
							<UpdateMetadata metadata={data.data.metadata} />
						</Descriptions.Item>
					</Descriptions>
				</>
			)}
		</>
	);
};

export default UpdateSinglePage;
