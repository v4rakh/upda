import { injectEndpoints } from './index';
import ApiTags from '../constants/apiTags';
import {
	CommentSingleResponse,
	CommentsRequestParams,
	CommentsResponse,
	CreateCommentRequest,
	ModifyCommentContentRequest
} from '../types';
import { FetchBaseQueryError } from '@reduxjs/toolkit/query';
import CommentFilterQueryParamNames from '../constants/api/commentFilterQueryParamNames';

const TAG_LIST_ID = 'LIST';

const invalidatesTags = (results?: CommentsResponse | CommentSingleResponse | void, error?: FetchBaseQueryError) => {
	if (error) {
		return [];
	}
	return [ApiTags.Comments] as any;
};

export const commentsApi = injectEndpoints({
	endpoints: (build) => {
		return {
			getComments: build.query<CommentsResponse, CommentsRequestParams>({
				query: ({ ...args }) => {
					const { updateId, page, pageSize } = args;

					const params = new URLSearchParams();
					if (page) {
						params.append(CommentFilterQueryParamNames.PAGE, `${page}`);
					}
					if (pageSize) {
						params.append(CommentFilterQueryParamNames.PAGE_SIZE, `${pageSize}`);
					}

					return { url: `comments/${updateId}?${params.toString()}` };
				},
				providesTags: (result, error) => {
					if (!error && result?.data.content) {
						return [
							{ type: ApiTags.Comments, id: TAG_LIST_ID },
							...result.data.content.map(({ id }) => ({ type: ApiTags.Comments, id }))
						];
					}
					return [];
				}
			}),
			createComment: build.mutation<CommentSingleResponse, { updateId: string; body: CreateCommentRequest }>({
				query: ({ updateId, body }) => ({ url: `comments/${updateId}`, method: 'POST', body }),
				invalidatesTags
			}),
			modifyCommentContent: build.mutation<
				CommentSingleResponse,
				{ id: string; body: ModifyCommentContentRequest }
			>({
				query: ({ id, body }) => ({ url: `comments/${id}/content`, method: 'PATCH', body }),
				invalidatesTags: (result, error, arg) => {
					if (error) {
						return [];
					}
					return [{ type: ApiTags.Comments, id: arg.id }];
				}
			}),
			deleteComment: build.mutation<void, { id: string }>({
				query: ({ id }) => ({ url: `comments/${id}`, method: 'DELETE' }),
				invalidatesTags
			})
		};
	}
});

export const {
	useLazyGetCommentsQuery,
	useCreateCommentMutation,
	useModifyCommentContentMutation,
	useDeleteCommentMutation
} = commentsApi;
