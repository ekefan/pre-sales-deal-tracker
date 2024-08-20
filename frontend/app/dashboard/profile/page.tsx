import Link from "next/link";
export default async function Page() {
  return (
    <div className="flex text-sm md:text-base flex-col p-2 h-full w-full sm:w-5/6 md:w-10/12 lg:w-8/12 xl:w-6/12 relative">
      <div className="p-3">Profile</div>
      <div className="px-3 py-1 flex flex-col gap-3 w-full">
        <div className="flex justify-between items-center">
          <p className="font-medium">Personal Information</p>
          <Link href="/dashboard/profile/update" className="rounded-lg items-center bg-slate-800 hover:bg-sky-200 hover:text-slate-900 text-white p-2">Update profile</Link>
        </div>
        <div className="flex border rounded-xl p-2 gap-1">
          <p>username:</p>
          <p>{"<username>"}</p>
        </div>
        <div className="flex border rounded-xl p-2 gap-1">
          <p>Full name:</p>
          <p>{"<firstname>"}</p>
        </div>
        <div className="flex border gap-1 rounded-xl p-2">
          <p>Email:</p>
          <p>{"<email>"}</p>
        </div>
        <div>
          <Link href="/dashboard/profile/password" className="rounded-lg border p-2 font-medium hover:bg-green-100 hover:text-green-400">Update password</Link>
        </div>
      </div>
    </div>
  );
}
