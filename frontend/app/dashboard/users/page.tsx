import { Plus } from 'lucide-react';
import Link from 'next/link';
export default async function Page() {
  return (
    <div className="flex text-sm md:text-base flex-col p-2 h-full w-full xl:w-11/12 relative">
      <div className="p-3">Users</div>
      <div className="flex px-3 py-1 justify-end">
        <Link href="/dashboard/users/create" className="bg-slate-800 rounded-lg p-2 text-sm md:text-base text-white flex items-center gap-2 hover:bg-sky-200 hover:text-slate-900">
            <p>create user</p>
            <Plus size={20}/>
        </Link>
      </div>
      <div className="flex flex-col gap-3  p-3 w-full grow h-auto">
        <section className="flex flex-col w-full h-full gap-3 ">
        <div className="bg-green-200 w-full h-32 rounded-lg border flex justify-end items-center">
              <Link href="/dashboard/users/update" className="border p-2 rounded-xl h-full flex items-center justify-center bg-green-500"><p>update</p></Link>
            </div>
          <div className="bg-yellow-200 w-full h-32 rounded-lg border"></div>
          <div className="bg-indigo-300 w-full h-32 rounded-lg border"></div>
          <div className="bg-pink-300 w-full h-32 rounded-lg border"></div>
          
        </section>
      </div>
    </div>
  );
}