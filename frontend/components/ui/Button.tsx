"use client";

import { motion } from "framer-motion";
import { ReactNode } from "react";

interface ButtonProps {
  children: ReactNode;
  onClick?: () => void;
  type?: "button" | "submit" | "reset";
  variant?: "primary" | "secondary" | "outline" | "danger" | "likert";
  className?: string;
  disabled?: boolean;
  isActive?: boolean;
}

export function Button({
  children,
  onClick,
  type = "button",
  variant = "primary",
  className = "",
  disabled = false,
  isActive = false,
}: ButtonProps) {
  let baseVariants = "";

  switch (variant) {
    case "primary":
      baseVariants = "bg-gradient-to-r from-palette4 to-palette5 text-white shadow-[0_8px_20px_-6px_rgba(74,143,231,0.5)] hover:shadow-[0_12px_25px_-6px_rgba(74,143,231,0.6)] border-none";
      break;
    case "secondary":
      baseVariants = "bg-white text-slate-500 hover:bg-slate-50 border-2 border-slate-100";
      break;
    case "outline":
      baseVariants = "bg-transparent border-2 border-slate-200 text-slate-600 hover:border-slate-300 hover:bg-slate-50";
      break;
    case "likert":
      baseVariants = isActive 
        ? "bg-palette4/10 border-palette4 text-palette4 shadow-sm" 
        : "bg-white border-slate-100 text-slate-500 hover:border-slate-200 hover:bg-slate-50 shadow-sm";
      break;
  }

  const defaultClasses = `px-6 py-3.5 rounded-xl font-bold transition-all flex items-center justify-center gap-2 ${baseVariants} border-2 ${disabled ? 'opacity-50 cursor-not-allowed' : ''}`;

  return (
    <motion.button
      type={type}
      onClick={onClick}
      disabled={disabled}
      className={`${defaultClasses} ${className}`}
      whileHover={disabled ? {} : { scale: 1.02, translateY: -2 }}
      whileTap={disabled ? {} : { scale: 0.98 }}
    >
      {children}
    </motion.button>
  );
}
