export type AuthProfile = {
	preferredUsername: string;
};

export type AuthProfileResponse = {
	data: AuthProfile;
};

export interface AuthTypeSessionLoginRequest {
	username: string;
	password: string;
}
