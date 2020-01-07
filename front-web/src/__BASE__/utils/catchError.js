import Vue from "vue";
import Cookie from '@2dfire/share/Cookie'
// let ERR_INFO = {
//     "0001": "网络暂时开小差了，请稍后重试" /*系统错误*/
// };

function catchError(e) {
		console.log(e, "catchError");
		if (e && e.data && +e.data.code === 302) {
				console.log("[登录过期，请重新登录]");
				window.location.href = e.data.location;
				// Vue.prototype.$Modal.error({
				// 		content: "登录过期，请重新登录",
				// 		onOk() {
				// 				Cookie.clear();
				// 				sessionStorage.clear();
				// 				console.log(2222);
				// 				window.location.href = e.data.location;
				// 		}
				// });
		} else if (e && e.data) {
		    let error = e.data;
		    let errorMSg = error.message;
		    Vue.prototype.$Message.error(errorMSg);
		} else {
		    Vue.prototype.$Modal.error({
		        title: "请注意",
		        content: "网络错误，请刷新页面",
		        okText: "刷新",
		        onOk() {
		            window.location.reload();
		        }
		    });
		}
    // if (e && +e.status === 401) {
    //     console.log("[401错误]");
    //     console.warn("Token过期");
    //     Vue.prototype.$Modal.error({
    //         content: "登录过期，请重新登录",
    //         onOk() {
    //             Cookie.clear();
    //             sessionStorage.clear();
    //             window.location.href = "./middle.html";
    //         }
    //     });
    // } else if (e && e.data) {
    //     let error = e.data;
    //     let errorMSg = error.message;
    //     Vue.prototype.$Message.error(errorMSg);
    // } else {
    //     Vue.prototype.$Modal.error({
    //         title: "请注意",
    //         content: "网络错误，请刷新页面",
    //         okText: "刷新",
    //         onOk() {
    //             window.location.reload();
    //         }
    //     });
    // }
}
export default catchError;
