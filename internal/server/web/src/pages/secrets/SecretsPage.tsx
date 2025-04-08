import CreateSecret from './CreateSecret';
import DeleteSecret from './DeleteSecret';
import UpdateValueSecret from './UpdateValueSecret';
import { useGetSecretsQuery } from '../../api/secretsApi';
import DateTimeStyle from '../../constants/dateTimeStyle';
import { useLocaleProviderContext } from '../../providers/LocaleContextProvider';
import { SecretResponse } from '../../types';
import { formatDateTimeWithTimeZone } from '../../utils/datetimeHelper';
import { sortAlphaIgnoringCase } from '../../utils/sortHelper';
import AppBreadcrumb from '../common/AppBreadcrumb';
import { QuestionCircleOutlined, ReloadOutlined } from '@ant-design/icons';
import { PageHeader } from '@ant-design/pro-layout';
import { Button, Col, Result, Row, Skeleton, Table, Tooltip, Typography } from 'antd';
import { ColumnsType } from 'antd/es/table';
import parse from 'html-react-parser';
import { FC, useCallback, useMemo } from 'react';
import { useTranslation } from 'react-i18next';

const { Text } = Typography;

const SecretsPage: FC = () => {
	const [t] = useTranslation('secrets');
	const { locale } = useLocaleProviderContext();

	const { isLoading, isError, refetch, isFetching, isSuccess, data } = useGetSecretsQuery();

	const invokeReload = useCallback(() => {
		refetch();
	}, [refetch]);

	const columns: ColumnsType<SecretResponse> = useMemo(() => {
		return [
			{
				title: t('col_key'),
				dataIndex: 'key',
				key: 'key',
				ellipsis: true,
				sorter: (a, b) => sortAlphaIgnoringCase(a.key, b.key),
				render: (_: string, entity: SecretResponse) => {
					return <Text copyable>{entity.key}</Text>;
				}
			},
			{
				title: t('col_value'),
				dataIndex: 'id',
				key: 'value',
				render: (id: string, entity: SecretResponse) => {
					return <UpdateValueSecret id={id} entityValue={entity.value} />;
				}
			},
			{
				title: t('col_created_at'),
				dataIndex: 'createdAt',
				key: 'createdAt',
				ellipsis: true,
				responsive: ['sm', 'md', 'lg', 'xl', 'xxl'],
				sorter: (a, b) => sortAlphaIgnoringCase(a.createdAt, b.createdAt),
				render: (value) => formatDateTimeWithTimeZone(value, DateTimeStyle.LONG, DateTimeStyle.MEDIUM, locale)
			},
			{
				title: t('col_updated_at'),
				dataIndex: 'updatedAt',
				key: 'updatedAt',
				ellipsis: true,
				responsive: ['sm', 'md', 'lg', 'xl', 'xxl'],
				sorter: (a, b) => sortAlphaIgnoringCase(a.updatedAt, b.updatedAt),
				render: (value) => formatDateTimeWithTimeZone(value, DateTimeStyle.LONG, DateTimeStyle.MEDIUM, locale)
			},
			{
				title: t('actions'),
				dataIndex: 'id',
				key: 'actions',
				ellipsis: false,
				render: (id: string) => {
					return <DeleteSecret id={id} />;
				}
			}
		];
	}, [locale, t]);

	return (
		<>
			<AppBreadcrumb items={[{ label: t('title'), active: true, path: '' }]} />
			<PageHeader
				className="pl-0"
				title={
					<Typography.Title level={4} ellipsis>
						{t('title')}
						<Tooltip placement="bottom" title={parse(t('help'))}>
							<Button icon={<QuestionCircleOutlined />} type="link" />
						</Tooltip>
					</Typography.Title>
				}
				extra={
					<Tooltip title={t('reload_tooltip')} placement="bottom">
						<Button
							icon={<ReloadOutlined />}
							type="link"
							onClick={invokeReload}
							loading={isFetching}
							disabled={isFetching || isLoading}
						/>
					</Tooltip>
				}
			/>
			<CreateSecret />
			{isLoading && <Skeleton loading={isLoading} active={isLoading} />}
			{isError && <Result status="error" title={t('error_default_loading')} />}
			{isSuccess && data.data.content.length === 0 && <Result status={404} title={t('no_secrets')} />}
			{isSuccess && data.data.content.length > 0 && (
				<Row justify="center" align="middle">
					<Col xs={24} lg={24}>
						<Table
							rowKey="id"
							columns={columns}
							loading={isLoading}
							dataSource={data.data.content}
							pagination={false}
						/>
					</Col>
				</Row>
			)}
		</>
	);
};

export default SecretsPage;
