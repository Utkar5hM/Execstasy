"use client";

import { useEffect, useState } from "react";
import { z } from "zod";

import { columns } from "./components/columns";
import { DataTable } from "./components/data-table";
import { taskSchema } from "./data/schema";
import { apiClient } from "@/utils/apiClient";
import { Button } from "@/components/ui/button";
import { IconDeviceLaptop } from "@tabler/icons-react";
import { Skeleton } from "@/components/ui/skeleton"
import Link from "next/link";

import {decodeJwt} from "@/utils/userToken"
const decodedToken = decodeJwt();

export default function TaskPage() {
  const [tasks, setTasks] = useState<z.infer<typeof taskSchema>[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchTasks() {
      try {
        type Task = {
          ID: number;
          name: string;
          status: string;
        };
        const response = await apiClient.get<Task[]>("/api/instances");
        // const transformedData = response.data.map((task: any) => {
        //   const { ID, ...rest } = task; // Extract "ID" and the rest of the object
        //   return { id: ID, ...rest }; // Rename "ID" to "id"
        // });
  2
        const tasks = z.array(taskSchema).parse(response.data);
        setTasks(tasks);
      } catch (err) {
        console.error("Failed to fetch tasks:", err);
        setError("Failed to load tasks.");
      } finally {
        setLoading(false);
      }
    }

    fetchTasks();
  }, []);

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
    return <div>Error: {error}</div>;
  }

  return (
    <div className="hidden h-full flex-1 flex-col space-y-8 p-8 md:flex">
      <div className="flex items-center justify-between space-y-2">
        <div>
          <h2 className="text-2xl font-bold tracking-tight">Instances</h2>
        </div>
        {decodedToken?.role === "admin" && (
        <div className="flex items-center space-x-2">
          <Link href="/instances/add">
        <Button size="sm">
      <IconDeviceLaptop /> Add Instance
    </Button>
    </Link>
        </div>)}
      </div>
      <DataTable data={tasks} columns={columns} />
    </div>
  );
}