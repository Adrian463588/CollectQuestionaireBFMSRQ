"use client";

import { motion } from "framer-motion";

interface ProgressBarProps {
  progress: number; // 0 to 100
}

export function ProgressBar({ progress }: ProgressBarProps) {
  return (
    <div className="w-full bg-slate-100 rounded-full h-2.5 overflow-hidden">
      <motion.div
        className="bg-gradient-to-r from-palette4 to-palette5 h-full rounded-full"
        initial={{ width: 0 }}
        animate={{ width: `${progress}%` }}
        transition={{ duration: 0.4, ease: "easeOut" }}
      />
    </div>
  );
}
