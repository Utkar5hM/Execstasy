import { z } from "zod"

import { columns } from "./components/columns"
import { DataTable } from "./components/data-table"
import { taskSchema } from "./data/schema"
import { apiClient } from "@/utils/apiClient"

// Simulate a database read for tasks.
async function getTasks() {
  // const response = await apiClient.get("/api/instances")
  // const tasks = response.data

  return z.array(taskSchema).parse([{
    "id": "7",
    "Name": "Staging Server",
    "HostAddress": "myinstancia.lfg",
    "Status": "active",
    "CreatedBy": "utkarshrm568@gmail.com"
},
{
  "id": "2",
  "Name": "Execstasy",
  "HostAddress": "testtt.com",
  "Status": "disabled",
  "CreatedBy": "utkarshrm568@gmail.com"
},
{
  "id": "3",
  "Name": "Olaola uberaaa",
  "HostAddress": "10.32.45.12",
  "Status": "disabled",
  "CreatedBy": "utkarshrm568@gmail.com"
}])
}

// Server Component
export default async function TaskPage() {
  const tasks = await getTasks()

  return (
    <TaskPageClient tasks={tasks} />
  )
}

// Client Component
function TaskPageClient({ tasks }: { tasks: z.infer<typeof taskSchema>[] }) {
  return (
    <>
      <div className="md:hidden">
      </div>
      <div className="hidden h-full flex-1 flex-col space-y-8 p-8 md:flex">
        <div className="flex items-center justify-between space-y-2">
          <div>
            <h2 className="text-2xl font-bold tracking-tight">Your Instances</h2>
          </div>
          <div className="flex items-center space-x-2">
            <button className="btn btn-primary">
              Create Instance
            </button>
          </div>
        </div>
        <DataTable data={tasks} columns={columns} />
      </div>
    </>
  )
}