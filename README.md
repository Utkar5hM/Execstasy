# Execstasy

Effortless RBAC and SSH access management for your Linux instances.

## Overview

Execstasy is a modern platform for managing Role-Based Access Control (RBAC) and SSH access to Linux servers. It provides a secure, centralized way to control which users can access which instances, and integrates easily with your existing infrastructure.

- **Centralized user and role management**
- **OAuth2 Device flows for secure authentication**

## Getting Started

1. **Clone this repository:**

```sh
   git clone https://github.com/Utkar5hM/Execstasy.git
   cd Execstasy
```

2. **configure app.env file**

copy app.env.example to app.env and setup all the required environment variables.

3.  **Start Dev Environment**

```sh
make dev
```

4.  **Start Prod Environment**

```sh
make prod
```
-----------

**Linux PAM Module:**

To enforce Execstasy access control at the SSH, use the companion PAM module: 👉 [ExecStasy-PAM](https://github.com/Utkar5hM/execstasy-pam)