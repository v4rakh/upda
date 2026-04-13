import CreateStateDefinition from './CreateStateDefinition';
import DeleteStateDefinition from './DeleteStateDefinition';
import {
	useGetUpdateStateDefinitionsQuery,
	useReorderUpdateStateDefinitionsMutation
} from '../../api/updateStateDefinitionsApi';
import AppPaths from '../../constants/appPaths';
import DateTimeStyle from '../../constants/dateTimeStyle';
import { useLocaleProviderContext } from '../../providers/LocaleContextProvider';
import { UpdateStateDefinition } from '../../types';
import { useNotification } from '../../use/useNotification';
import { formatDateTimeWithTimeZone } from '../../utils/datetimeHelper';
import { renderIcon } from '../../utils/iconHelper';
import { getPageFullPath } from '../../utils/urlHelper';
import AppBreadcrumb from '../common/AppBreadcrumb';
import { CheckOutlined, EditOutlined, HolderOutlined, QuestionCircleOutlined, ReloadOutlined } from '@ant-design/icons';
import { PageHeader } from '@ant-design/pro-layout';
import { DndContext, DragEndEvent, PointerSensor, useSensor, useSensors } from '@dnd-kit/core';
import { restrictToVerticalAxis } from '@dnd-kit/modifiers';
import { SortableContext, useSortable, verticalListSortingStrategy } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import { Button, Col, Result, Row, Skeleton, Space, Table, Tag, Tooltip, Typography } from 'antd';
import { ColumnsType } from 'antd/es/table';
import parse from 'html-react-parser';
import { FC, useCallback, useEffect, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Link } from 'react-router';

const { Text } = Typography;

interface DraggableRowProps extends React.HTMLAttributes<HTMLTableRowElement> {
	'data-row-key': string;
}

const DraggableRow: FC<DraggableRowProps> = (props) => {
	const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({
		id: props['data-row-key']
	});

	const style: React.CSSProperties = {
		...props.style,
		transform: CSS.Translate.toString(transform),
		transition,
		cursor: 'move',
		...(isDragging ? { position: 'relative', zIndex: 9999, background: '#fafafa' } : {})
	};

	return <tr {...props} ref={setNodeRef} style={style} {...attributes} {...listeners} />;
};

const StateDefinitionsPage: FC = () => {
	const [t] = useTranslation('state_definitions');
	const { locale } = useLocaleProviderContext();
	const { apiError } = useNotification();

	const { isLoading, isError, refetch, isFetching, isSuccess, data } = useGetUpdateStateDefinitionsQuery();
	const [reorder, { isError: isReorderError, error: reorderError }] = useReorderUpdateStateDefinitionsMutation();

	const [localData, setLocalData] = useState<UpdateStateDefinition[]>([]);

	useEffect(() => {
		if (data?.data?.content) {
			setLocalData([...data.data.content]);
		}
	}, [data]);

	useEffect(() => {
		if (isReorderError) {
			apiError({
				i18n: {
					notFound: t('error_reorder'),
					unAuthorized: t('error_unauthorized'),
					forbidden: t('error_forbidden'),
					badRequest: t('error_reorder'),
					default: t('error_reorder')
				},
				error: reorderError
			});
			// Revert to server data on error
			if (data?.data?.content) {
				setLocalData([...data.data.content]);
			}
		}
	}, [isReorderError, reorderError, t, apiError, data]);

	const sensors = useSensors(
		useSensor(PointerSensor, {
			activationConstraint: {
				distance: 1
			}
		})
	);

	const invokeReload = useCallback(() => {
		refetch();
	}, [refetch]);

	const handleDragEnd = useCallback(
		(event: DragEndEvent) => {
			const { active, over } = event;

			if (over && active.id !== over.id) {
				const oldIndex = localData.findIndex((item) => item.id === active.id);
				const newIndex = localData.findIndex((item) => item.id === over.id);

				if (oldIndex !== -1 && newIndex !== -1) {
					// Create new array with reordered items
					const newData = [...localData];
					const [movedItem] = newData.splice(oldIndex, 1);
					newData.splice(newIndex, 0, movedItem);

					// Update local state optimistically
					setLocalData(newData);

					// Prepare reorder request with new sort orders
					const items = newData.map((item, index) => ({
						id: item.id,
						sortOrder: index
					}));

					// Send reorder request
					reorder({ items });
				}
			}
		},
		[localData, reorder]
	);

	const columns: ColumnsType<UpdateStateDefinition> = useMemo(() => {
		return [
			{
				key: 'drag',
				width: '2%',
				render: () => <HolderOutlined style={{ cursor: 'grab', color: 'lightblue' }} />
			},
			{
				title: t('col_name'),
				dataIndex: 'name',
				key: 'name',
				responsive: ['sm', 'md', 'lg', 'xl', 'xxl']
			},
			{
				title: t('col_label'),
				dataIndex: 'label',
				key: 'label',
				ellipsis: true
			},
			{
				title: t('col_description'),
				dataIndex: 'description',
				key: 'description',
				ellipsis: true,
				responsive: ['sm', 'md', 'lg', 'xl', 'xxl'],
				render: (description: string | undefined) => {
					return description ? <Text>{description}</Text> : <Text type="secondary">-</Text>;
				}
			},
			{
				title: t('col_color'),
				dataIndex: 'color',
				key: 'color',
				responsive: ['sm', 'md', 'lg', 'xl', 'xxl'],
				render: (color: string) => {
					return <Tag color={color}>{color}</Tag>;
				}
			},
			{
				title: t('col_icon'),
				dataIndex: 'icon',
				key: 'icon',
				render: (icon: string | undefined, r) => {
					return renderIcon(icon, { fontSize: 16, color: r.color });
				}
			},
			{
				title: t('col_is_initial'),
				dataIndex: 'isInitial',
				key: 'isInitial',
				width: '10%',
				responsive: ['sm', 'md', 'lg', 'xl', 'xxl'],
				render: (isInitial: boolean) => {
					return isInitial ? (
						<Tag color="green" icon={<CheckOutlined />}>
							{t('yes')}
						</Tag>
					) : null;
				}
			},
			{
				title: t('col_skip_on_new_version'),
				dataIndex: 'skipOnNewVersion',
				key: 'skipOnNewVersion',
				width: '10%',
				responsive: ['sm', 'md', 'lg', 'xl', 'xxl'],
				render: (skipOnNewVersion: boolean) => {
					return skipOnNewVersion ? (
						<Tag color="green" icon={<CheckOutlined />}>
							{t('yes')}
						</Tag>
					) : null;
				}
			},
			{
				title: t('col_created_at'),
				dataIndex: 'createdAt',
				key: 'createdAt',
				ellipsis: true,
				responsive: ['sm', 'md', 'lg', 'xl', 'xxl'],
				render: (value) => formatDateTimeWithTimeZone(value, DateTimeStyle.LONG, DateTimeStyle.MEDIUM, locale)
			},
			{
				title: t('col_updated_at'),
				dataIndex: 'updatedAt',
				key: 'updatedAt',
				ellipsis: true,
				responsive: ['sm', 'md', 'lg', 'xl', 'xxl'],
				render: (value) => formatDateTimeWithTimeZone(value, DateTimeStyle.LONG, DateTimeStyle.MEDIUM, locale)
			},
			{
				title: t('actions'),
				dataIndex: 'id',
				key: 'actions',
				ellipsis: false,
				render: (_: string, entity: UpdateStateDefinition) => {
					return (
						<Space>
							<Link to={getPageFullPath(`${AppPaths.STATE_DEFINITIONS}/${entity.id}`)}>
								<Button type="link" icon={<EditOutlined />} />
							</Link>
							<DeleteStateDefinition id={entity.id} />
						</Space>
					);
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
			<CreateStateDefinition />
			{isLoading && <Skeleton loading={isLoading} active={isLoading} />}
			{isError && <Result status="error" title={t('error_default_loading')} />}
			{isSuccess && localData.length === 0 && <Result status={404} title={t('no_state_definitions')} />}
			{isSuccess && localData.length > 0 && (
				<Row justify="center" align="middle">
					<Col xs={24} lg={24}>
						<DndContext sensors={sensors} modifiers={[restrictToVerticalAxis]} onDragEnd={handleDragEnd}>
							<SortableContext
								items={localData.map((item) => item.id)}
								strategy={verticalListSortingStrategy}>
								<Table
									rowKey="id"
									columns={columns}
									loading={isLoading}
									dataSource={localData}
									pagination={false}
									components={{
										body: {
											row: DraggableRow
										}
									}}
								/>
							</SortableContext>
						</DndContext>
					</Col>
				</Row>
			)}
		</>
	);
};

export default StateDefinitionsPage;
