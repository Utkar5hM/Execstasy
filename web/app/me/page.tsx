"use client"
import { useEffect, useState } from "react";
import { apiClient } from "@/utils/apiClient";
import {
  CheckCircle,
  CircleOff,
} from "lucide-react"
import { Separator } from "@/components/ui/separator"
import { Button } from "@/components/ui/button"
import { ExecTable } from "@/components/execTable";
import { Skeleton } from "@/components/ui/skeleton";
import * as React from "react"
import Link from "next/link";
import { DataTableColumnHeader } from "../instances/components/data-table-column-header";
import { DefaultStatusResponse, APIMyInstance, APIMyProfile, APIMyRoles } from "@/utils/ResponseTypes";
import { ColumnDef } from "@tanstack/react-table";


export default function ProfileViewPage() {
  const [userProfile, setUserProfile] = useState<APIMyProfile | undefined>(undefined);
  const [loading, setLoading] = useState(true);
  const [loading2, setLoading2] = useState(true);
  const [loading3, setLoading3] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [instancesData, setInstancesData] = useState<APIMyInstance[]>([]);
  const [rolesData, setRolesData] = useState<APIMyRoles[]>([]);


  useEffect(() => {
    async function fetchInstance() {
      try {
        const response = await apiClient.get<APIMyProfile | DefaultStatusResponse>(`/api/users/me`);
        if (response.status =200){
        setUserProfile(response.data  as APIMyProfile);
        } else {
          setError("Failed to load user profile.");
        }
      } catch (err) {
        console.error("Failed to fetch instance:", err);
        setError("Failed to load instance.");
      } finally {
        setLoading(false);
      }
    }

    fetchInstance();
  },  []);


  useEffect(() => {
    async function fetchInstance() {
      try {
        const response = await apiClient.get<APIMyInstance[] | DefaultStatusResponse>(`/api/instances/me`);
        if (response.status =200){
          setInstancesData(response.data as APIMyInstance[]);
        } else {
          setError("Failed to load instances.");
        }
      } catch (err) {
        console.error("Failed to fetch instance:", err);
        setError("Failed to load instance.");
      } finally {
        setLoading2(false);
      }
    }

    fetchInstance();
  }, []);
  
  useEffect(() => {
    async function fetchInstance() {
      try {
        const response = await apiClient.get<APIMyRoles[] | DefaultStatusResponse>(`/api/roles/me`);
        if (response.status =200){
        setRolesData(response.data as APIMyRoles[]);
        } else {
          setError("Failed to load roles.");
        }
      } catch (err) {
        console.error("Failed to fetch instance:", err);
        setError("Failed to load instance.");
      } finally {
        setLoading3(false);
      }
    }

    fetchInstance();
  }, []);
  const statuses = [
    {
      value: "active",
      label: "active",
      icon: CheckCircle,
    },
    {
      value: "disabled",
      label: "disabled",
      icon: CircleOff,
    },
  ]
  const RolesColumns: ColumnDef<APIMyRoles, unknown>[] = [
    {
      accessorKey: "name",
      header: "Name",
    },
    {
      accessorKey: "role",
      header: "Role",
    },
    {
    accessorKey: "ID",
    header: "Actions",
    enableSorting: false,
    enableHiding: false,
    cell: ({ row }) => (
      <Link href={`/roles/view/${row.original.id}`}
      className="hover:underline"
      >
        <Button variant="outline">
          View
        </Button>
      </Link>
    ),
    }
  ]
  const Instancecolumns: ColumnDef<APIMyInstance, unknown>[] = [
    {
      accessorKey: "Name",
      header: "Name",
    },
      {
        accessorKey: "HostAddress",
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title="Host Address" />
        ),
        cell: ({ row }) => <div className="w-[80px]">{row.getValue("HostAddress")}</div>,
        enableSorting: false,
        enableHiding: false,
      },
      {
        accessorKey: "Status",
        header: ({ column }) => (
          <DataTableColumnHeader column={column} title="Status" />
        ),
        cell: ({ row }) => {
          const status = statuses.find(
            (s) => s.value === row.getValue("Status")
          )
          if (!status) {
            return <div className="w-[40px]">Cannot find</div>
          }
    
          return (
            <div className="flex w-[100px] items-center">
              {status.icon && (
                <status.icon className="mr-2 h-4 w-4 text-muted-foreground" />
              )}
              <span>{status.value}</span>
            </div>
          )
        },
      },
      {
      accessorKey: "ID",
      header: "Actions",
      enableSorting: false,
      enableHiding: false,
      cell: ({ row }) => (
        <Link href={`/instances/view/${row.original.ID}`}
        className="hover:underline"
        >
          <Button variant="outline">
            View
          </Button>
        </Link>
      ),
      }
  ];


  if (loading || loading2 || loading3) {
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

  if (error || !userProfile) {
    return (
      <div className="flex h-screen items-center justify-center">
        <div className="text-red-500">{error}</div>
      </div>
    );
  }


  return (<>
    <div className="p-8">
          <h1 className="scroll-m-20 text-4xl font-extrabold tracking-tight text-balance">
          {userProfile.name}</h1>
      <p className="pt-4 flex items-center space-x-2"><strong>Role:</strong>     <code className="bg-muted relative rounded px-[0.3rem] py-[0.2rem] font-mono text-sm font-semibold">
      {userProfile.role}
    </code>
      </p>
      <p className="pt-4 flex items-center space-x-2"><strong>Username:&nbsp;</strong> {userProfile.username}</p>
      <p className="pt-4 flex items-center space-x-2"><strong>Email:&nbsp;</strong> {userProfile.email}</p>
      <p className="pt-4 flex items-center space-x-2">
<Link href={`/me/edit`}>
<Button
        size="sm"
        className="hidden h-8 lg:flex w-[80px]"
      >Edit
      </Button>
      </Link>
      </p>
      <Separator className="my-6" />
    </div>
    <div className="px-8">
      <h1 className="text-2xl font-bold mb-4">My Instances</h1>
      <ExecTable data={instancesData}
        columns={Instancecolumns}
        filterColumn="Name" />
      </div>
    <div className="px-8">
      <h1 className="text-2xl font-bold mb-4">My Roles</h1>
      <ExecTable data={rolesData} 
      columns={RolesColumns} 
      filterColumn="name" 
         />
      </div>
    </>
  );
}