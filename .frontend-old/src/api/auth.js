const API_URL = "http://localhost:8080";

export const loginWithGoogle = async () => {
  window.location.href = `${API_URL}/auth/google`;
};

export const logout = () => {
  localStorage.removeItem("token");
  window.location.reload();
};

export const getToken = () => {
  return localStorage.getItem("token");
};
