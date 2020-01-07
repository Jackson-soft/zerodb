import jwt_decode from 'jwt-decode'
export const getNickName = () => {
	let pageData = JSON.parse(sessionStorage.getItem('baseData')) || '';
	const user = {};
	if (pageData.access_token) {
		user.Info = jwt_decode(pageData.access_token) || ''
		if (user.Info.name) {
			return user.Info.name
		}
	}
}
export const getUsername = () => {
	let pageData = JSON.parse(sessionStorage.getItem('baseData')) || '';
	const user = {};
	if (pageData.access_token) {
		user.Info = jwt_decode(pageData.access_token) || ''
		if (user.Info.preferred_username) {
			return user.Info.preferred_username
		}
	}
}
export const getUserInfo = () => {
	let pageData = JSON.parse(sessionStorage.getItem('baseData')) || '';
	const user = {};
	if (pageData.access_token) {
		user.access_token = pageData.access_token || '';
		user.token_type = pageData.token_type || '';
		user.Info = jwt_decode(pageData.access_token) || ''
	}
	return user
}

export const clearSession = () => {
	return sessionStorage.removeItem('baseData')
}
