module.exports = {
	page: ["index"],
	port: 9090,
	//pc端将其置为null
	px2rem: null,
	imports: "iview",
	esLint: "fix",
	checkVersion: false,
	publicPath(env) {
		return {
			dev: "../",
			daily: "http://zerodb.2dfire-daily.com/",
			// pre: 'http://work.2dfire-pre.com/',
			publish: "http://zerodb.2dfire.net/"
		}[env]
	},
	proxy: {
		"/app": {
			target: "https://zerodb.app.2dfire-daily.com",
			pathRewrite: {
				"/app": ""
			}
		}
	}
}
