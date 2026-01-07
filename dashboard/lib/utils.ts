import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

// XSS対策のためのサニタイズ関数
export function sanitizeInput(input: string): string {
  // HTMLタグを除去
  const sanitized = input
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#39;")
    .replace(/&/g, "&amp;")
    .replace(/\\/g, "&#92;")
    .replace(/\//g, "&#47;")

  return sanitized
}
