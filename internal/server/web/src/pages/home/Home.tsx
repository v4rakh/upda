import AppPaths from '../../constants/appPaths';
import { getPageFullPath } from '../../utils/urlHelper';
import { Navigate } from 'react-router';

const Home = () => {
	return <Navigate to={getPageFullPath(AppPaths.UPDATES)} replace={true} />;
};

export default Home;
