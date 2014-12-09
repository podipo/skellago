var skella = skella || {};
skella.api = skella.api || {};

/*
	Connect to the API and authenticate
*/
skella.api.login = function(email, password, successCallback, errorCallback){
	$.ajax({
		url: '/api/user/current',
		method: 'post',
		contentType: 'application/json',
		data: JSON.stringify({
			'email': email,
			'password': password
		}),
		error: function(jqXHR, status, error) {
			if (errorCallback) {
				errorCallback.apply(this, arguments);
			}
		},
		success: function(data, status, jqXHR) {
			if (successCallback) {
				successCallback.apply(this, arguments);
			}
		}
	});
}

/*
	Deauthenticate 
*/
skella.api.logout = function(successCallback, errorCallback){
	$.ajax({
		url: '/api/user/current',
		method: 'delete',
		error: function(jqXHR, status, error) {
			if (errorCallback) {
				errorCallback.apply(this, arguments);
			}
		},
		success: function(data, status, jqXHR) {
			if (successCallback) {
				successCallback.apply(this, arguments);
			}
		}
	});
}