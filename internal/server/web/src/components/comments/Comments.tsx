import CreateComment from './CreateComment';
import DeleteComment from './DeleteComment';
import EditableComment from './EditableComment';
import { useLazyGetCommentsQuery } from '../../api/commentsApi';
import DateTimeStyle from '../../constants/dateTimeStyle';
import { LIST_PAGE_DEFAULT, LIST_PAGE_SIZE_DEFAULT } from '../../constants/pagination';
import { useLocaleProviderContext } from '../../providers/LocaleContextProvider';
import { CommentResponse, CommentSingleResponse } from '../../types';
import { useNotification } from '../../use/useNotification';
import { formatDateTimeWithTimeZone } from '../../utils/datetimeHelper';
import { DownOutlined, ReloadOutlined } from '@ant-design/icons';
import { Avatar, Button, Flex, List, Result, Skeleton, Tooltip } from 'antd';
import { filter, map, unionBy } from 'lodash';
import React, { ReactNode, useCallback, useEffect, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';

interface CommentProps {
	updateId: string;
}

const Comments = ({ updateId }: CommentProps): ReactNode => {
	const [t] = useTranslation('comments');
	const { apiError } = useNotification();
	const { locale } = useLocaleProviderContext();

	const [trigger, result] = useLazyGetCommentsQuery();

	const [page, setPage] = useState<number>(LIST_PAGE_DEFAULT);
	const [comments, setComments] = useState<CommentResponse[] | undefined>(undefined);
	const [hasMore, setHasMore] = useState<boolean>(true);

	const fetchData = useCallback(
		async (updateId: string, page: number) => {
			const res = await trigger({ page: page, pageSize: LIST_PAGE_SIZE_DEFAULT, updateId: updateId }, false);
			if (res.isSuccess && res.data && res.data.data.content) {
				const merged = unionBy(comments, res.data.data.content, 'id');
				setComments(merged);
				setHasMore(merged.length < res.data.data.totalElements);
			}
			if (res.isError) {
				apiError({
					i18n: {
						unAuthorized: t('error_unauthorized_get'),
						forbidden: t('error_forbidden_get'),
						badRequest: t('error_bad_request_get'),
						default: t('error_default_get')
					},
					error: res.error
				});
			}
		},
		[apiError, comments, t, trigger]
	);

	useEffect(() => {
		if (!comments) {
			fetchData(updateId, page);
		}
	}, [comments, fetchData, page, updateId]);

	const removeComment = useCallback((id: string) => {
		setComments((prevComments) => filter(prevComments, (e) => e.id !== id));
	}, []);

	const editComment = useCallback((response: CommentSingleResponse) => {
		setComments((prevComments) =>
			map(prevComments, (e) => (e.id === response.data.id ? { ...e, ...response.data } : e))
		);
	}, []);

	const invokeReload = useCallback(() => {
		setPage(LIST_PAGE_DEFAULT);
		setHasMore(false);
		setComments(undefined);
	}, []);

	const onLoadMore = useCallback(async () => {
		if (hasMore) {
			const nextPage = page + 1;
			setPage(nextPage);
			await fetchData(updateId, nextPage);
		}
	}, [fetchData, hasMore, page, updateId]);

	const loadMore = useMemo(() => {
		return hasMore ? (
			<div
				style={{
					textAlign: 'center'
				}}>
				<Button
					size="small"
					type="link"
					onClick={onLoadMore}
					icon={<DownOutlined />}
					loading={result.isLoading || result.isFetching}
					disabled={result.isLoading || result.isFetching}>
					{t('load_more')}
				</Button>
			</div>
		) : null;
	}, [hasMore, onLoadMore, result.isFetching, result.isLoading, t]);

	return (
		<>
			<CreateComment updateId={updateId} onCreateSuccess={invokeReload} />
			<Flex justify="end" align="center">
				<Tooltip title={t('reload_tooltip')} placement="bottom">
					<Button
						icon={<ReloadOutlined />}
						type="link"
						onClick={invokeReload}
						loading={result.isFetching}
						disabled={result.isLoading || result.isFetching}
					/>
				</Tooltip>
			</Flex>
			{result.isLoading ||
				(result.isFetching && (
					<Skeleton
						loading={result.isLoading || result.isFetching}
						active={result.isLoading || result.isFetching}
					/>
				))}
			{result.isError && <Result status="error" title={t('error_default_loading')} />}
			{result.isSuccess && (
				<List
					itemLayout="vertical"
					loading={result.isLoading || result.isFetching}
					locale={{ emptyText: t('no_data') }}
					size="large"
					loadMore={result.isSuccess && loadMore}
					dataSource={comments}
					renderItem={(comment) => (
						<List.Item
							extra={
								<DeleteComment
									onDeleteSuccess={() => removeComment(comment.id)}
									key={`delete_${comment.id}`}
									id={comment.id}
								/>
							}>
							<List.Item.Meta
								avatar={
									<Avatar draggable={false} size="large">
										{comment.author}
									</Avatar>
								}
								title={formatDateTimeWithTimeZone(
									comment.updatedAt,
									DateTimeStyle.LONG,
									DateTimeStyle.MEDIUM,
									locale
								)}
							/>
							<EditableComment comment={comment} onEditSuccess={editComment} />
						</List.Item>
					)}
				/>
			)}
		</>
	);
};

export default Comments;
