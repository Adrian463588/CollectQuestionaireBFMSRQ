"use client";

import { motion } from "framer-motion";
import { ReactNode } from "react";

interface BentoCardProps {
  children: ReactNode;
  className?: string;
  onClick?: () => void;
  hoverEffect?: boolean;
  delay?: number;
}

export function BentoCard({ children, className = "", onClick, hoverEffect = false, delay = 0 }: BentoCardProps) {
  const baseClasses = "bg-white/90 backdrop-blur-xl rounded-3xl p-6 md:p-8 shadow-[0_20px_50px_-12px_rgba(0,0,0,0.1)] border border-white z-10 block";
  
  const hoverProps = hoverEffect ? {
    whileHover: { scale: 1.02, translateY: -4 },
    whileTap: { scale: 0.98 },
  } : {};

  return (
    <motion.div
      initial={{ opacity: 0, scale: 0.95, y: 10 }}
      animate={{ opacity: 1, scale: 1, y: 0 }}
      transition={{ duration: 0.5, ease: "easeOut", delay }}
      {...hoverProps}
      onClick={onClick}
      className={`${baseClasses} ${hoverEffect ? 'cursor-pointer hover:shadow-[0_20px_50px_-12px_rgba(74,143,231,0.2)] transition-all group' : ''} ${className}`}
    >
      {children}
    </motion.div>
  );
}
