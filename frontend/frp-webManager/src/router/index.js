	import { createRouter, createWebHashHistory, createWebHistory } from "vue-router"

	const routes = [
	    {
	        path: "/",
	        alias: ["/home", "/index"],
	        component: () => import("@/App.vue")
        }
	]

	const router = createRouter({
	    history: createWebHistory(),
	    routes
	})

	export default router