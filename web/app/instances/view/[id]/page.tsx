"use client"
import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { apiClient } from "@/utils/apiClient";
import {
  CheckCircle,
  CircleOff,
} from "lucide-react"
import { Separator } from "@/components/ui/separator"
import { Button } from "@/components/ui/button"
import Link from "next/link"
import { ExecTable } from "@/components/execTable";
import {
  ColumnDef,
} from "@tanstack/react-table"
import { Badge } from "@/components/ui/badge"
import { ArrowUpDown, BadgeCheckIcon, ChevronDown,  MoreHorizontal } from "lucide-react"
import { Skeleton } from "@/components/ui/skeleton";
export default function InstanceViewPage() {
	const params = useParams(); // Access the params object
	const id = params?.id; // Extract the `id` parameter
  const [instance, setInstance] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [usersData, setUsersData] = useState([]);
  const [rolesData, setRolesData] = useState([]);
  useEffect(() => {
    async function fetchInstance() {
      try {
        const response = await apiClient.get(`/api/instances/view/${id}`);
        setInstance(response.data);
      } catch (err) {
        console.error("Failed to fetch instance:", err);
        setError("Failed to load instance.");
      } finally {
        setLoading(false);
      }
    }

    fetchInstance();
  }, [id]);

  const data = [
    { id: 1, name: "John Doe", email: "john@example.com", role: "admin" },
    { id: 2, name: "Jane Smith", email: "jane@example.com", role: "user" },
  ];
  
  const columns = [
    {
      accessorKey: "name",
      header: "Name",
    },
    {
      accessorKey: "email",
      header: "Email",
    },
    {
      accessorKey: "role",
      header: "Role",
    },
  ];
  useEffect(() => {
    async function fetchInstance() {
      try {
        const response = await apiClient.get(`/api/instances/users/${id}`);
        setUsersData(response.data);
      } catch (err) {
        console.error("Failed to fetch instance:", err);
        setError("Failed to load instance.");
      } finally {
        setLoading(false);
      }
    }

    fetchInstance();
  }, [id]);
  
  useEffect(() => {
    async function fetchInstance() {
      try {
        const response = await apiClient.get(`/api/instances/roles/${id}`);
        setRolesData(response.data);
      } catch (err) {
        console.error("Failed to fetch instance:", err);
        setError("Failed to load instance.");
      } finally {
        setLoading(false);
      }
    }

    fetchInstance();
  }, [id]);

  const RolesColumns = [
    {
      accessorKey: "name",
      header: "Name",
    },
    {
      accessorKey: "host_username",
      header: "Host Username",
    },
  ]
  const Userscolumns = [
    {
      accessorKey: "name",
      header: "Name",
    },
    {
      accessorKey: "username",
      header: "username",
    },
    {
      accessorKey: "host_username",
      header: "Host Username",
    },
  ];

  if (loading) {
    return (
      <div className="flex h-screen items-center justify-center">
        <div className="flex items-center space-x-4">
          <Skeleton className="h-12 w-12 rounded-full" />
          <div className="space-y-2">
            <Skeleton className="h-4 w-[250px]" />
            <Skeleton className="h-4 w-[200px]" />
          </div>
        </div>
      </div>
        );
  }

  if (error) {
    return (
      <div className="flex h-screen items-center justify-center">
        <div className="text-red-500">{error}</div>
      </div>
    );
  }
  return (<>
    <div className="p-8">
          <h1 className="scroll-m-20 text-4xl font-extrabold tracking-tight text-balance">
          {instance.Name}</h1>

      <p className="pt-4  text-muted-foreground text-xl">{instance.Description}</p>
      <p className="pt-4 flex items-center space-x-2"><strong>Host Address:</strong>     <code className="bg-muted relative rounded px-[0.3rem] py-[0.2rem] font-mono text-sm font-semibold">
      {instance.HostAddress}
    </code>
      </p>
      <p className="pt-4 flex items-center space-x-2 ">
      <strong>Status:</strong>
  {instance.Status === "active" ? (
    <span className="text-muted-foreground flex items-center space-x-1">
      <CheckCircle className="w-5 h-5" />
      <span>Active</span>
    </span>
  ) : (
    <span className="flex items-center space-x-1">
      <CircleOff className="w-5 h-5" />
      <span>Disabled</span>
    </span>
  )}
</p>
      <p className="pt-4 flex items-center space-x-2"><strong>Created By:&nbsp;</strong> {instance.CreatedBy}</p>
      <p className="pt-4 flex items-center space-x-2">
<Button
        size="sm"
        className="hidden h-8 lg:flex w-[80px]"
      >Edit
      </Button>
      <Button
        variant="destructive"
        size="sm"
        className="hidden h-8 lg:flex w-[80px]"
      >Delete
      </Button>
      </p>
      <Separator className="my-6" />
    </div>
    <div className="px-8">
      <h1 className="text-2xl font-bold mb-4">Users</h1>
      <ExecTable data={usersData}
        columns={Userscolumns}
        filterColumn="name"
        headerContent={
          <>
            <Button>Add User</Button>
          </>
        } />
      </div>
    <div className="px-8">
      <h1 className="text-2xl font-bold mb-4">Roles</h1>
      <ExecTable data={rolesData} 
      columns={RolesColumns} 
      filterColumn="name" 
        headerContent={
          <>
            <Button>Add Role</Button>
          </>
        } />
      </div>
    </>
  );
}