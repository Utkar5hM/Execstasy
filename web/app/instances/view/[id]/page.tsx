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

export default function InstanceViewPage() {
	const params = useParams(); // Access the params object
	const id = params?.id; // Extract the `id` parameter
  const [instance, setInstance] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  useEffect(() => {
    async function fetchInstance() {
      try {
        const response = await apiClient.get(`/api/instances/${id}`);
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
  if (loading) {
    return (
      <div className="flex h-screen items-center justify-center">
        <div>Loading...</div>
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
  return (
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
    
  );
}