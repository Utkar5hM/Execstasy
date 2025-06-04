import { z } from "zod"

// We're keeping a simple non-relational schema here.
// IRL, you will have a schema for your data models.
export const taskSchema = z.object({
  ID: z.number(),
  Name: z.string(),
  HostAddress: z.string(),
  Status: z.string(),
  CreatedBy: z.string(),
})

export type Task = z.infer<typeof taskSchema>
