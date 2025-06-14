type ApiResponse<T> = {
	status: number;
	ok: boolean;
	data: T;
  };

export const apiClient = {
	async request<T>(
	  endpoint: string,
	  options: RequestInit = {}
	): Promise<ApiResponse<T>> {
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
  
	  const json = await response.json().catch(() => {
		throw new Error("Failed to parse JSON response");
	  });
	  
	  return {
		status: response.status, // Include the response status
		ok: response.ok,         // Include the `ok` status
		data: json,              // Include the parsed JSON data
	  };
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
  
	delete<T>(endpoint: string, body: any, options: RequestInit = {}) {
	  return this.request<T>(endpoint, { ...options, 
		method: "DELETE",
		body: JSON.stringify(body),
	   });
	},
  };