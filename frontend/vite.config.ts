import { defineConfig } from "vite";
import { sveltekit } from '@sveltejs/kit/vite';
import tailwindcss from "@tailwindcss/vite";

// https://vitejs.dev/config/
export default defineConfig({
	plugins: [sveltekit(), tailwindcss()],
	server: {
		fs: {
			allow: ['/Users/alexeyzoz/Desktop/Projects/WranglerLocalDev/myproject/frontend']
		}
	}
});
