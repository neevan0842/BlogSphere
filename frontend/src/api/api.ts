import axios, { type AxiosInstance, type AxiosRequestHeaders } from "axios";

const API_URL: string = import.meta.env.VITE_API_URL;

const api: AxiosInstance = axios.create({
  baseURL: API_URL,
  withCredentials: true,
});

// Interceptor: sets the Authorization header before each request.
api.interceptors.request.use(
  (config) => {
    const accessToken = localStorage.getItem("access-token");
    if (accessToken) {
      config.headers = {
        ...config.headers,
        Authorization: `Bearer ${accessToken}`,
      } as AxiosRequestHeaders;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  },
);

export { api };
