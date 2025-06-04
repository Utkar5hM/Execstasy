export const apiClient = {
	async request<T>(
	  endpoint: string,
	  options: RequestInit = {}
	): Promise<T> {
	  const token = document.cookie
		.split("; ")
		.find((row) => row.startsWith("jwt="))
		?.split("=")[1];
  
	  // Construct headers
	  const headers: HeadersInit = {
		"Content-Type": "application/json",
		...(token && { Authorization: `Bearer ${token}` }), // Add Authorization header only if token exists
		...options.headers,
	  };
  
	  const response = await fetch(endpoint, {
		...options,
		headers,
	  });
  
	  if (!response.ok) {
		const error = await response.json();
		throw new Error(error.message || "An error occurred");
	  }
  
	  return response.json();
	},
  
	get<T>(endpoint: string, options: RequestInit = {}) {
	  return this.request<T>(endpoint, { ...options, method: "GET" });
	},
  
	post<T>(endpoint: string, body: any, options: RequestInit = {}) {
	  return this.request<T>(endpoint, {
		...options,
		method: "POST",
		body: JSON.stringify(body),
	  });
	},
  
	put<T>(endpoint: string, body: any, options: RequestInit = {}) {
	  return this.request<T>(endpoint, {
		...options,
		method: "PUT",
		body: JSON.stringify(body),
	  });
	},
  
	delete<T>(endpoint: string, options: RequestInit = {}) {
	  return this.request<T>(endpoint, { ...options, method: "DELETE" });
	},
  };