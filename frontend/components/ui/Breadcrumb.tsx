"use client";

import Link from "next/link";
import { motion } from "framer-motion";

export interface BreadcrumbItem {
  label: string;
  href?: string;
}

interface BreadcrumbProps {
  items: BreadcrumbItem[];
}

/**
 * Reusable Breadcrumb component.
 * Renders a trail of navigation links separated by "/".
 * The last item (active page) is shown as plain text without a link.
 */
export function Breadcrumb({ items }: BreadcrumbProps) {
  return (
    <motion.nav
      aria-label="Breadcrumb"
      initial={{ opacity: 0, y: -6 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.35, ease: "easeOut" }}
      className="flex items-center flex-wrap gap-1.5 text-xs font-semibold text-slate-400 mb-6"
    >
      {items.map((item, index) => {
        const isLast = index === items.length - 1;

        return (
          <span key={item.label} className="flex items-center gap-1.5">
            {isLast ? (
              <span className="text-slate-600 font-bold" aria-current="page">
                {item.label}
              </span>
            ) : (
              <Link
                href={item.href ?? "/"}
                className="hover:text-palette4 transition-colors duration-200"
              >
                {item.label}
              </Link>
            )}
            {!isLast && (
              <span className="text-slate-300 select-none">/</span>
            )}
          </span>
        );
      })}
    </motion.nav>
  );
}
