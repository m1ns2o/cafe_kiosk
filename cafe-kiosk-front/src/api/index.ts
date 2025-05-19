// src/api/index.ts
import axios from 'axios';
import type { AxiosInstance, AxiosRequestConfig } from 'axios';

// 기본 설정 옵션
const config: AxiosRequestConfig = {
  baseURL: '/api',
  timeout: 100000,
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json'
  }
};

// Axios 인스턴스 생성
const apiClient: AxiosInstance = axios.create(config);

export default apiClient;