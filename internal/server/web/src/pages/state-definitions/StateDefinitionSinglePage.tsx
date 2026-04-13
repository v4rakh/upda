import UpdateColorStateDefinition from './UpdateColorStateDefinition';
import UpdateDescriptionStateDefinition from './UpdateDescriptionStateDefinition';
import UpdateIconStateDefinition from './UpdateIconStateDefinition';
import UpdateIsInitialStateDefinition from './UpdateIsInitialStateDefinition';
import UpdateLabelStateDefinition from './UpdateLabelStateDefinition';
import UpdateNameStateDefinition from './UpdateNameStateDefinition';
import UpdateSkipOnNewVersionStateDefinition from './UpdateSkipOnNewVersionStateDefinition';
import {
	useDeleteUpdateStateDefinitionMutation,
	useGetUpdateStateDefinitionByIdQuery
} from '../../api/updateStateDefinitionsApi';
import ApiErrorCodes from '../../constants/apiErrorCodes';
import AppPathParamNames from '../../constants/appPathParamNames';
import AppPaths from '../../constants/appPaths';
import DateTimeStyle from '../../constants/dateTimeStyle';
import { useLocaleProviderContext } from '../../providers/LocaleContextProvider';
import { useNotification } from '../../use/useNotification';
import { formatDateTimeWithTimeZone } from '../../utils/datetimeHelper';
import { getPageFullPath } from '../../utils/urlHelper';
import AppBreadcrumb from '../common/AppBreadcrumb';
import { DeleteOutlined } from '@ant-design/icons';
import { PageHeader } from '@ant-design/pro-layout';
import { Button, Descriptions, Popconfirm, Result, Skeleton, Tooltip, Typography } from 'antd';
import { FC, ReactNode, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate, useParams } from 'react-router';

const { Text } = Typography;
const { Item } = Descriptions;

const StateDefinitionSinglePage: FC = (): ReactNode => {
	const [t] = useTranslation('state_definition_single');
	const [tDelete] = useTranslation('state_definition_delete');
	const { locale } = useLocaleProviderContext();
	const { apiError } = useNotification();
	const navigate = useNavigate();

	const { [AppPathParamNames.STATE_DEFINITION_ID]: stateDefinitionId } = useParams();

	const { isLoading, isError, isSuccess, data, error } = useGetUpdateStateDefinitionByIdQuery(
		{
			id: stateDefinitionId || ''
		},
		{
			skip: stateDefinitionId === undefined
		}
	);

	const [
		deleteStateDefinition,
		{ isLoading: isDeleteLoading, isError: isDeleteError, error: deleteError, isSuccess: isDeleteSuccess }
	] = useDeleteUpdateStateDefinitionMutation();

	useEffect(() => {
		if (isDeleteError) {
			apiError({
				i18n: {
					notFound: tDelete('error_unable_delete'),
					unAuthorized: tDelete('error_unauthorized_delete'),
					forbidden: tDelete('error_forbidden_delete'),
					badRequest: tDelete('error_bad_request_delete'),
					default: tDelete('error_default_delete')
				},
				error: deleteError
			});
		}
	}, [isDeleteError, deleteError, tDelete, apiError]);

	useEffect(() => {
		if (isDeleteSuccess) {
			navigate(getPageFullPath(AppPaths.STATE_DEFINITIONS));
		}
	}, [isDeleteSuccess, navigate]);

	return (
		<>
			<AppBreadcrumb
				items={[
					{
						label: t('title_parent'),
						path: getPageFullPath(AppPaths.STATE_DEFINITIONS)
					},
					{
						label: data?.data.name || '',
						active: true,
						path: ''
					}
				]}
			/>
			<PageHeader
				className="pl-0"
				title={
					<Typography.Title level={4} ellipsis>
						{t('title', { name: data?.data.name })}
					</Typography.Title>
				}
				extra={
					data && (
						<Popconfirm
							title={t('delete_title')}
							onConfirm={() => deleteStateDefinition({ id: data.data.id })}
							okText={t('delete')}
							placement="bottom"
							cancelText={t('cancel')}
							okButtonProps={{ icon: <DeleteOutlined />, type: 'primary', danger: true }}>
							<Tooltip title={t('help_delete')} placement="bottom">
								<Button loading={isDeleteLoading} icon={<DeleteOutlined />} type="text" danger />
							</Tooltip>
						</Popconfirm>
					)
				}
			/>
			{isLoading && <Skeleton loading={isLoading} active={isLoading} />}
			{isError && (error as any)?.data?.status !== ApiErrorCodes.NOT_FOUND && (
				<Result status={500} title={t('error_default_loading')} />
			)}
			{isError && (error as any)?.data?.status === ApiErrorCodes.NOT_FOUND && (
				<Result status={404} title={t('no_state')} />
			)}
			{isSuccess && (
				<Descriptions
					layout="vertical"
					size="small"
					column={{ xs: 1, sm: 2, md: 3, lg: 3, xl: 4, xxl: 4 }}
					colon={false}>
					<Item label={t('name')}>
						<UpdateNameStateDefinition entity={data.data} />
					</Item>
					<Item label={t('label')}>
						<UpdateLabelStateDefinition entity={data.data} />
					</Item>
					<Item label={t('description')}>
						<UpdateDescriptionStateDefinition entity={data.data} />
					</Item>
					<Item label={t('color')}>
						<UpdateColorStateDefinition entity={data.data} />
					</Item>
					<Item label={t('icon')}>
						<UpdateIconStateDefinition entity={data.data} />
					</Item>
					<Item label={t('is_initial')}>
						<UpdateIsInitialStateDefinition entity={data.data} />
					</Item>
					<Item label={t('skip_on_new_version')}>
						<UpdateSkipOnNewVersionStateDefinition entity={data.data} />
					</Item>
					<Item label={t('created')}>
						<Text>
							{formatDateTimeWithTimeZone(
								data.data.createdAt,
								DateTimeStyle.LONG,
								DateTimeStyle.MEDIUM,
								locale
							)}
						</Text>
					</Item>
					<Item label={t('updated')}>
						<Text>
							{formatDateTimeWithTimeZone(
								data.data.updatedAt,
								DateTimeStyle.LONG,
								DateTimeStyle.MEDIUM,
								locale
							)}
						</Text>
					</Item>
				</Descriptions>
			)}
		</>
	);
};

export default StateDefinitionSinglePage;
