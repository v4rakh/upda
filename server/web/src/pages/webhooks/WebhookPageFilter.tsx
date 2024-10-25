import WebhookFilterQueryParamNames from '../../constants/api/webhookFilterQueryParamNames';
import WebhookOrder from '../../constants/api/webhookOrder';
import WebhookOrderBy from '../../constants/api/webhookOrderBy';
import useWebhooksFilterQueryParams from '../../use/useWebhooksFilterQueryParams';
import { Form, Select } from 'antd';
import { FC, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { useSearchParams } from 'react-router-dom';

const WebhookPageFilter: FC = () => {
	const [t] = useTranslation('webhooks_filters');
	const [form] = Form.useForm();

	const [queryParams, setSearchQueryParams] = useSearchParams();
	const { orderBy, order } = useWebhooksFilterQueryParams();

	const onOrderByChange = useCallback(
		(value: WebhookOrderBy | undefined) => {
			if (!value) {
				queryParams.delete(WebhookFilterQueryParamNames.ORDER_BY);
			} else {
				queryParams.set(WebhookFilterQueryParamNames.ORDER_BY, value ?? undefined);
			}
			setSearchQueryParams(queryParams);
		},
		[queryParams, setSearchQueryParams]
	);

	const onOrderChange = useCallback(
		(value: WebhookOrder | undefined) => {
			if (!value) {
				queryParams.delete(WebhookFilterQueryParamNames.ORDER);
			} else {
				queryParams.set(WebhookFilterQueryParamNames.ORDER, value);
			}
			setSearchQueryParams(queryParams);
		},
		[queryParams, setSearchQueryParams]
	);

	const initialValues = { orderBy: orderBy ?? WebhookOrderBy.LABEL, order: order ?? WebhookOrder.ASC };

	return (
		<Form layout="inline" form={form} initialValues={initialValues}>
			<Form.Item label={t('order_by')} name="orderBy">
				<Select
					style={{ width: 120 }}
					variant="filled"
					onChange={onOrderByChange}
					options={[
						{ value: WebhookOrderBy.ID, label: t(`order_by_${WebhookOrderBy.ID.toLowerCase()}`) },
						{
							value: WebhookOrderBy.CREATED_AT,
							label: t(`order_by_${WebhookOrderBy.CREATED_AT.toLowerCase()}`)
						},
						{
							value: WebhookOrderBy.UPDATED_AT,
							label: t(`order_by_${WebhookOrderBy.UPDATED_AT.toLowerCase()}`)
						},
						{
							value: WebhookOrderBy.LABEL,
							label: t(`order_by_${WebhookOrderBy.LABEL.toLowerCase()}`)
						},
						{
							value: WebhookOrderBy.TYPE,
							label: t(`order_by_${WebhookOrderBy.TYPE.toLowerCase()}`)
						}
					]}
				/>
			</Form.Item>
			<Form.Item name="order">
				<Select
					variant="filled"
					style={{ width: 120 }}
					onChange={onOrderChange}
					options={[
						{ value: WebhookOrder.DESC, label: t(`order_${WebhookOrder.DESC.toLowerCase()}`) },
						{ value: WebhookOrder.ASC, label: t(`order_${WebhookOrder.ASC.toLowerCase()}`) }
					]}
				/>
			</Form.Item>
		</Form>
	);
};

export default WebhookPageFilter;
