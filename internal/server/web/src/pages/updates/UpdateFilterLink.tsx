import UpdateFilterQueryParamNames from '../../constants/api/updateFilterQueryParamNames';
import UpdateSearchIn from '../../constants/api/updateSearchIn';
import AppPaths from '../../constants/appPaths';
import { getPageFullPathWithQueryParameters } from '../../utils/urlHelper';
import { Typography } from 'antd';
import React, { ReactNode, useCallback } from 'react';
import { useNavigate } from 'react-router';

interface UpdateFilterLinkProps {
	label: string;
	searchIn: UpdateSearchIn;
}

const { Link } = Typography;

const UpdateFilterLink = ({ label, searchIn }: UpdateFilterLinkProps): ReactNode => {
	const navigate = useNavigate();

	const redirect = useCallback(() => {
		const redirectTo = getPageFullPathWithQueryParameters(AppPaths.UPDATES, {
			[UpdateFilterQueryParamNames.SEARCH_TERM]: label,
			[UpdateFilterQueryParamNames.SEARCH_IN]: searchIn
		});

		navigate(redirectTo);
	}, [label, navigate, searchIn]);

	return <Link onClick={redirect}>{label}</Link>;
};

export default UpdateFilterLink;
