import axios from "axios";
import { store } from "@/lib/store";

export const http = axios.create({ baseURL: "/api" });

http.interceptors.request.use((config) => {
	const token = store.getState().auth.token;
	if (token) config.headers.Authorization = `Bearer ${token}`;
	return config;
});
