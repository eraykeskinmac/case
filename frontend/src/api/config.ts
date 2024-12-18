import axios from "axios";
import type { ApiResponse } from "../types";

const API_URL = "http://localhost:3000/api/v1";

export const api = axios.create({
  baseURL: API_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

api.interceptors.response.use(
  (response) => {
    return {
      ...response,
      data: {
        data: response.data,
        message: response.data.message,
        meta: response.data.meta,
      } as ApiResponse<unknown>,
    };
  },
  (error) => {
    const message = error.response?.data?.message || "An error occurred";
    return Promise.reject(new Error(message));
  }
);
