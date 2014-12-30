var skella = skella || {};
skella.api = skella.api || {};
skella.schema = skella.schema || {};
skella.events = skella.events || {};

// Used by the authentication mechanism
skella.api.sessionCookie = "skella_auth";
skella.api.emailCookie = "skella_email";

skella.events.SchemaPopulated = 'populated';
skella.events.LoggedIn = 'logged-in';
skella.events.LoggedOut = 'logged-out';

skella.schema.pathVariablesRegex = new RegExp('{[^{]+}', 'g');
skella.schema.acceptFormat = "application/vnd.api+json; version="

skella.schema.generateURL = function(path, attributes){
		var tokens = path.match(skella.schema.pathVariablesRegex);
		if(tokens == null || tokens.length == 0) {
			return path;
		}
		var result = "";
		var index = 0;
		for(var i=0; i < tokens.length; i++){
			var tokenIndex = path.indexOf(tokens[i]);
			result += path.substring(index, tokenIndex);
			index = tokenIndex + tokens[i].length;
			var name = tokens[i].substring(1, tokens[i].length - 1).split(':')[0];
			result += attributes[name]
		}
		if(index < path.length){
			result += path.substring(index);
		}
		return result;	
}

// Add the API version to the XHR headers when syncing models or collections
skella.schema.versionedSync = function(method, model, options){
	var beforeSend = options.beforeSend;
	var version = this.version;
	options.beforeSend = function(xhr) {
		xhr.setRequestHeader('Accept', skella.schema.acceptFormat + version);
		if (beforeSend) return beforeSend.apply(this, arguments);
	};
	Backbone.Model.prototype.sync.apply(this, arguments);
}

skella.schema.Collection = Backbone.Collection.extend({
	initialize: function(options){
		this.options = options;
	},
	parse: function(response){
		this.offset = response.offset;
		this.limit = response.limit;
		return response.objects;
	},
	url: function(){
		return skella.schema.generateURL(this.schema.path, this.options);
	},
	sync: skella.schema.versionedSync
});

skella.schema.Model = Backbone.Model.extend({
	initialize: function(options){
		this.options = options;
	},
	url: function(){
		return skella.schema.generateURL(this.schema.path, this.attributes);
	},
	sync: skella.schema.versionedSync
});

skella.schema.Schema = Backbone.Model.extend({
	initialize: function(options){
		_.bindAll(this, 'populate', 'hasProperties');
		this.options = options;
		this.user = null; // Will be set to schema.api.User if the session is authenticated
		this.api = {}; // This is where we will put the Backbone Models and Collections populated from the schema
		this.populated = false;
		if(!this.options.url){
			throw 'Schema requires you to pass in a "url" option';
		}
		this.on('sync', this.populate);
	},
	url: function(){
		return this.options.url;
	},
	populate: function(){
		this.version = this.get('api').version;
		for(var i in this.attributes.endpoints){
			var endpoint = this.attributes.endpoints[i];
			if(this.hasProperties(endpoint['properties'], ['offset', 'limit', 'objects']) == true){
				continue;
			}
			var name = skella.schema.objectifyEndpointName(endpoint['name']);
			this.api[name] = skella.schema.Model.extend({
				'schema':endpoint,
				'version':this.version
			});
		}

		for(var i in this.attributes.endpoints){
			var endpoint = this.attributes.endpoints[i];
			if(this.hasProperties(endpoint['properties'], ['offset', 'limit', 'objects']) == false){
				continue;
			}

			var model = null;
			var objectsProperty = this.getProperty(endpoint.properties, 'objects');
			if(objectsProperty && objectsProperty['children-type']){
				var childName = skella.schema.objectifyEndpointName(objectsProperty['children-type']);
				if(this.api[childName]){
					model = this.api[childName];
				}
			}

			var name = skella.schema.objectifyEndpointName(endpoint['name']);
			this.api[name] = skella.schema.Collection.extend({
				'schema':endpoint,
				'model':model,
				'version':this.version
			});
		}
		this.populated = true;
		this.trigger(skella.events.SchemaPopulated, this);
	},
	getProperty: function(properties, name){
		for(var i=0; i < properties.length; i++){
			if(properties[i].name == name){
				return properties[i];
			}
		}
		return null;
	},
	hasProperties: function(properties, names) {
		for(var i=0; i < names.length; i++){
			var found = false;
			for(var j=0; j < properties.length; j++){
				if(properties[j].name == names[i]){
					found = true;
					break;
				}
			}
			if(!found) return false;
		}
		return true;
	}
});

skella.schema.objectifyEndpointName = function(name){
	if(!name) return null;
	var tokens = name.split('-');
	result = "";
	for(var i=0; i < tokens.length; i++){
		result += skella.schema.initialCap(tokens[i]);
	}
	return result;
}

skella.schema.initialCap = function(val){
	return val[0].toUpperCase() + val.substring(1);
}

// TODO stop hard coding the API version number here
window.API_VERSION = "0.1.0";

$(document).ready(function(){
	window.schema = new skella.schema.Schema({'url':'/api/' + window.API_VERSION + '/schema'});
	window.schema.on(skella.events.SchemaPopulated, function(){
		if(localStorage.user){
			window.schema.user = new window.schema.api.User(JSON.parse(localStorage.user));
			// Update the localStorage
			window.schema.user.on('sync', function(){
				localStorage.user = JSON.stringify(window.schema.user.attributes);
			});
		} else {
			window.schema.user = null;
		}
	});
	window.schema.fetch();
})

/*
Returns true if the session cookie exists
This depends on the jquery.cookie plugin: https://github.com/carhartl/jquery-cookie

*/
skella.api.loggedIn = function(){
	return !!$.cookie(skella.api.sessionCookie);
}

/*
	Connect to the API and authenticate
*/
skella.api.login = function(email, password, successCallback, errorCallback){
	$.ajax({
		url: '/api/' + window.API_VERSION + '/user/current',
		method: 'post',
		contentType: 'application/json',
		data: JSON.stringify({
			'email': email,
			'password': password
		}),
		headers :  {
			'Accept': skella.schema.acceptFormat + window.API_VERSION
		},
		error: function(jqXHR, status, error) {
			if (errorCallback) {
				errorCallback.apply(this, arguments);
			}
		},
		success: function(data, status, jqXHR) {
			localStorage.user = JSON.stringify(data); // Used to populate window.schema.user
			if(window.schema){
				if(window.schema.user){
					window.schema.user.set(data);
				} else {
					window.schema.user = new window.schema.api.User(data);
					window.schema.user.on('sync', function(){
						localStorage.user = JSON.stringify(window.schema.user.attributes);
					});
				}
				window.schema.trigger(skella.events.LoggedIn);
			}
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
		url: '/api/' + window.API_VERSION + '/user/current',
		method: 'delete',
		headers :  {
			'Accept': skella.schema.acceptFormat + window.API_VERSION
		},
		error: function(jqXHR, status, error) {
			if (errorCallback) {
				errorCallback.apply(this, arguments);
			}
		},
		success: function(data, status, jqXHR) {
			// Delete the localStorage
			localStorage.removeItem('user');
			if(window.schema){
				window.schema.user = null;
				window.schema.trigger(skella.events.LoggedOut);
			}
			if (successCallback) {
				successCallback.apply(this, arguments);
			}
		}
	});
}