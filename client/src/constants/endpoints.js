// export const API_BASE_URL = "http://localhost:8080/api/v1";
export const API_BASE_URL =
  process.env.NODE_ENV === "production"
    ? window.location.origin + "/api/v1"
    : "http://localhost:8080/api/v1";
export const API_BASE_URL_HISTORY = `${API_BASE_URL}/history`;
export const API_BASE_URL_REALTIME = `${API_BASE_URL}/ws`;
export const API_BASE_URL_STATS = `${API_BASE_URL}/stat`;
