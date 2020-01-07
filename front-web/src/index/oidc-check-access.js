import { vuexOidcCreateRouterMiddleware } from "vuex-oidc"
import router from "./router"
import store from "./store"

router.beforeEach(vuexOidcCreateRouterMiddleware(store, "oidc"))
