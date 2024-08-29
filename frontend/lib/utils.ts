import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"
import axios from "axios"
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}


export const DEFAULT_PAGE_SIZE: number = 10;
export const ADMIN_ROLE: string = "admin";
export const BASE_URL: string = "http://localhost:8080/"