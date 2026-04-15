"use client";

import { useState, useEffect, ReactNode } from "react";
import { motion } from "framer-motion";
import { BentoCard } from "@/components/ui/BentoCard";
import { Button } from "@/components/ui/Button";
import { AnimatedBackground } from "@/components/ui/AnimatedBackground";

export function AdminGuard({ children }: { children: ReactNode }) {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [isMounted, setIsMounted] = useState(false);

  useEffect(() => {
    setIsMounted(true);
    const authStatus = sessionStorage.getItem("adminAuth");
    if (authStatus === "true") {
      setIsAuthenticated(true);
    }
  }, []);

  const handleLogin = (e: React.FormEvent) => {
    e.preventDefault();
    const correctPassword = process.env.NEXT_PUBLIC_ADMIN_PASSWORD || "Admin123";

    if (password === correctPassword) {
      sessionStorage.setItem("adminAuth", "true");
      setIsAuthenticated(true);
      setError("");
    } else {
      setError("Password salah.");
    }
  };

  if (!isMounted) return null; // Prevent hydration mismatch

  if (isAuthenticated) {
    return <>{children}</>;
  }

  return (
    <div className="min-h-screen bg-bgLight relative flex items-center justify-center p-4 overflow-hidden">
      <AnimatedBackground />

      <BentoCard className="max-w-md w-full z-10">
        <div className="mb-6 text-center">
          <div className="w-14 h-14 bg-palette4/10 rounded-2xl mx-auto mb-4 flex items-center justify-center text-palette4">
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"></path>
            </svg>
          </div>
          <h2 className="text-2xl font-extrabold text-textMain tracking-tight">Login Admin</h2>
          <p className="text-slate-500 font-medium text-sm mt-1">Masukkan password untuk melihat data partisipan</p>
        </div>

        <form onSubmit={handleLogin} className="space-y-4">
          <div>
            <input
              type="password"
              placeholder="Password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full px-4 py-3.5 rounded-xl border border-slate-200 text-textMain placeholder-slate-400 focus:outline-none focus:border-palette4 focus:ring-4 focus:ring-palette4/15 transition-all shadow-sm font-medium"
              required
            />
          </div>
          
          {error && (
            <motion.p 
              initial={{ opacity: 0, y: -5 }} animate={{ opacity: 1, y: 0 }}
              className="text-red-500 text-xs font-bold"
            >
              {error}
            </motion.p>
          )}

          <Button type="submit" variant="primary" className="w-full mt-2">
            Masuk
          </Button>

          <div className="text-center mt-4 pt-2">
             <a href="/" className="text-xs font-bold text-slate-400 hover:text-palette5 transition-colors">
               ← Kembali ke Home
             </a>
          </div>
        </form>
      </BentoCard>
    </div>
  );
}
