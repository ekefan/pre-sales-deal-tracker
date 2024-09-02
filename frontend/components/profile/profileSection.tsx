"use client"
import Link from "next/link";
import { useUser } from "@/context/userContext";


export function ProfileSection() {
    const {usr} = useUser()
  return (
    <section className="px-3 py-1 flex flex-col gap-3 w-full">
      <div className="flex justify-between items-center">
        <p className="font-medium">Personal Information</p>
        <Link
          href={`/dashboard/profile/update?user_id=${usr?.user.user_id}/`}
          className="rounded-lg items-center bg-slate-800 hover:bg-sky-200 hover:text-slate-900 text-white p-2"
        >
          Update profile
        </Link>
      </div>
      <div className="flex border rounded-xl p-2 gap-1">
        <p>username:</p>
        <p>{usr?.user?.username}</p>
      </div>
      <div className="flex border rounded-xl p-2 gap-1">
        <p>Full name:</p>
        <p>{usr?.user?.fullname}</p>
      </div>
      <div className="flex border gap-1 rounded-xl p-2">
        <p>Email:</p>
        <p>{usr?.user?.email}</p>
      </div>
      <div>
        <Link
          href={`/dashboard/profile/password?user_id=${usr?.user.user_id}`}
          className="rounded-lg border p-2 font-medium hover:bg-green-100 hover:text-green-400"
        >
          Update password
        </Link>
      </div>
    </section>
  );
}
