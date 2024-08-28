import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"
import axios from "axios"
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}


export const DEFAULT_PAGE_SIZE: number = 20;
export const ADMIN_ROLE: string = "admin";