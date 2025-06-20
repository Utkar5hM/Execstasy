export type DefaultStatusResponse = {
	status: string;
	message: string;
	error?: string;
	error_description?: string;
  };

  export type AuthJwtPayload = {
	name?: string;
	username?: string;
	email?: string;
	// add other custom claims if needed
  };
  export type APIRole = {
	ID: number;
	Name: string;
	Description: string;
	CreatedBy: string;
	CreatedAt: string;
	UpdatedAt: string;
  }

  export type APIInstances = {
			ID: number;
			name: string;
			status: string;
		  } 
export type APIInstanceView = {
	ID: number;
	Name: string;
	Status: string;
	CreatedBy: string;
	Description: string;
	HostAddress: string;
	HostUsers: string[] | null;
	ClientID: string;
}
  export type APIRoles = {
	created_at: string;
	created_by: string;
	description: string;
	id: number;
	name: string;
	updated_at: string;
	users: string;
  };

  export type APIInstanceUsers = {
	username: string;
	id: number;
	role: string;
	name: string;
	host_username: string;
  }

  export type APIInstanceRoles = {
	id: number;
	name: string;
	host_username: string;
  }

  export type CreateInstanceResponse = {
	client_id: string;
	id: number;
	message: string;
	status: string;
  }

  export type APIRolesAccess = {
	access: boolean;
	message: string;
	status: string;
  }

  export type APIRoleUsers = {
	id: number;
	name: string;
	role: string;
	username: string;
  }

  export type APIMyProfile = {
	email: string;
	id: number;
	name: string;
	role: string;
	username: string;
  }

export type APIMyInstance = {
	ID: number;
	Name: string;
	HostAddress: string;
	Status: string;
	CreatedBy: string;
}

export type APIMyRoles = {
	created_at: string;
	description: string;
	id: number;
	name: string;
	role: string;
	updated_at: string;
}