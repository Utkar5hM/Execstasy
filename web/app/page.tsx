import { IconTerminal2 } from "@tabler/icons-react";

export default function LandingPage() {
  return (
<div className="flex h-screen flex-col items-center justify-center">
  {/* Title with Icon */}
  <div className="flex items-center space-x-4">
    <IconTerminal2 size={64} className="" />
    <h1 className="scroll-m-20 text-6xl font-extrabold tracking-tight">
      ExecStasy
    </h1>
  </div>

  {/* Caption */}
  <h4 className="scroll-m-20 mt-4 text-center text-l font-semibold tracking-tight">
    Effortless RBAC and SSH access for your Linux Instances.
  </h4>
</div>
  );
}