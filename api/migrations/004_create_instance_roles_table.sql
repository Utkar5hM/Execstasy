CREATE TABLE instance_roles (
    instance_id INT NOT NULL,
    role_id INT NOT NULL,
    instance_host_username VARCHAR(255)  NOT NULL,
    PRIMARY KEY (instance_id, role_id, instance_host_username),
    FOREIGN KEY (instance_id) REFERENCES instances(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (instance_id, instance_host_username) REFERENCES instance_host_users(instance_id, username) ON DELETE CASCADE
);