"use client";

import { useEffect, useState } from "react";
import Particles, { initParticlesEngine } from "@tsparticles/react";
import { loadSlim } from "@tsparticles/slim";
import type { Container } from "@tsparticles/engine";
import { IconTerminal2 } from "@tabler/icons-react"
import { LoginForm } from "@/components/login-form";
import particlesConfig from "@/lib/particles-config";

export default function LoginPage() {
  const [init, setInit] = useState(false);

  useEffect(() => {
    initParticlesEngine(async (engine) => {
      await loadSlim(engine);
    }).then(() => {
      setInit(true);
    });
  }, []);

  const particlesLoaded = async (container?: Container): Promise<void> => {
    console.log(container);
  };

  if (init){
  return (
    <div className="relative min-h-svh overflow-hidden">
      {/* Particles background */}
      <Particles
        id="tsparticles"
        particlesLoaded={particlesLoaded}
        options={particlesConfig}
        className="absolute inset-0 z-0"
      />

      {/* Foreground content */}
      <div className="relative z-10 flex min-h-svh flex-col items-center justify-center gap-6 p-6 md:p-10">
        <div className="flex w-full max-w-sm flex-col gap-6">
          <a href="#" className="flex items-center gap-2 self-center font-medium">
            <div className="flex size-6 items-center justify-center rounded-md">
              <IconTerminal2 />
            </div>
            Execstasy
          </a>
          <LoginForm />
        </div>
      </div>
    </div>
  );
} else {
  return(
    <div className="relative min-h-svh overflow-hidden">
      <div className="relative z-10 flex min-h-svh flex-col items-center justify-center gap-6 p-6 md:p-10">
        <div className="flex w-full max-w-sm flex-col gap-6">
          <a href="#" className="flex items-center gap-2 self-center font-medium">
            <div className="bg-primary text-primary-foreground flex size-6 items-center justify-center rounded-md">
              <IconTerminal2 />
            </div>
            Execstasy
          </a>
          <LoginForm />
        </div>
      </div>
    </div>
  )
}
}
