"use client"
import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { apiClient } from "@/utils/apiClient";

export default function InstanceViewPage() {
	const params = useParams(); // Access the params object
	const id = params?.id; // Extract the `id` parameter
  const [instance, setInstance] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

//   useEffect(() => {
//     async function fetchInstance() {
//       try {
//         const response = await apiClient.get(`/api/instances/${id}`);
//         setInstance(response.data);
//       } catch (err) {
//         console.error("Failed to fetch instance:", err);
//         setError("Failed to load instance.");
//       } finally {
//         setLoading(false);
//       }
//     }

//     fetchInstance();
//   }, [id]);

//   if (loading) {
//     return <div>Loading...</div>;
//   }

//   if (error) {
//     return <div>Error: {error}</div>;
//   }

  return (
    <div className="p-8">
      <h1 className="text-2xl font-bold">Instance Details</h1>
      <p><strong>ID: </strong> {id}</p>
      {/* <p><strong>Name:</strong> {instance.name}</p>
      <p><strong>Status:</strong> {instance.status}</p>
      <p><strong>Created By:</strong> {instance.createdBy}</p> */}
    </div>
  );
}