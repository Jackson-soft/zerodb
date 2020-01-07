import "@styles/common.less"
import "@styles/init.less"
import "@styles/iviewreset.less"
import {
	Breadcrumb,
	BreadcrumbItem,
	Button,
	Card,
	Checkbox,
	CheckboxGroup,
	Col,
	DatePicker,
	Dropdown,
	DropdownItem,
	DropdownMenu,
	Form,
	FormItem,
	Icon,
	Input,
	InputNumber,
	Menu,
	MenuItem,
	Message,
	Modal,
	Option,
	Page,
	Radio,
	RadioGroup,
	Row,
	Select,
	Submenu,
	Switch,
	Table,
	TabPane,
	Tabs,
	Timeline,
	TimelineItem,
	Upload
} from "iview"
import Vue from "vue"
import { mapGetters } from "vuex"
import App from "./App"
import "./oidc-check-access"
import router from "./router"
import store from "./store"

const componentsArr = {
	Tabs,
	TabPane,
	Row,
	Col,
	Card,
	Modal,
	Table,
	DropdownMenu,
	Dropdown,
	DatePicker,
	Menu,
	Submenu,
	MenuItem,
	DropdownItem,
	CheckboxGroup,
	Checkbox,
	Form,
	Page,
	Icon,
	Message,
	FormItem,
	Button,
	Select,
	Option,
	Input,
	RadioGroup,
	Radio,
	Breadcrumb,
	BreadcrumbItem,
	Timeline,
	TimelineItem,
	InputNumber,
	Upload,
	iSwitch: Switch
}
// Object.keys(filters).map((key) => {
// 	Vue.filter(key, filters[key])
// })
Object.keys(componentsArr).forEach(item => {
	Vue.component(item, componentsArr[item])
})
Vue.prototype.$Message = Message
Vue.prototype.$Modal = Modal
Vue.config.productionTip = false
new Vue({
	el: "#app",
	router,
	store,
	computed: {
		...mapGetters("oidc", ["oidcIsAuthenticated"]),
		hasAccess() {
			return this.oidcIsAuthenticated || this.$route.meta.isPublic
		}
	},
	mounted() {
		window.addEventListener("vuexoidc:userLoaded", this.userLoaded)
	},
	destroyed() {
		window.removeEventListener("vuexoidc:userLoaded", this.userLoaded)
	},
	methods: {
		userLoaded(user) {
			console.log("I am listening to the user loaded event in vuex-oidc" + JSON.stringify(user))
		}
	},
	template: "<App/>",
	components: {
		App
	}
})
