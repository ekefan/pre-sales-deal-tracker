import Link from "next/link";

export default async function Page() {
  return (
    <div className="flex flex-col text-sm md:text-base p-2 h-full w-full xl:w-11/12 relative">
      <div className="p-3">Pitch Requests</div>
      <div className="px-3 py-1 text-green-300">
        Oh great, no pending requests
      </div>
      <div className="flex flex-col gap-3  p-3 w-full grow h-auto">
        <section className="flex flex-col w-full h-full gap-3 ">
          <div className="bg-green-200 w-full h-32 rounded-lg border flex justify-end items-center">
            <Link
              href="/dashboard/pitch-requests/update-pr"
              className="border p-2 rounded-xl h-full flex items-center justify-center bg-green-500"
            >
              <p>update</p>
            </Link>
          </div>
          <div className="bg-yellow-200 w-full h-32 rounded-lg border flex justify-end items-center">
            <Link
              href="/dashboard/pitch-requests/create"
              className="border p-2 rounded-xl h-full flex flex-col items-center justify-center bg-yellow-500"
            >
              <p>create</p>
              <p>pitch r</p>
            </Link>
          </div>
          <div className="bg-indigo-300 w-full h-32 rounded-lg border flex justify-end items-center">
            <Link
              href="/dashboard/pitch-requests/create-deal"
              className=" p-2 rounded-xl h-full flex items-center justify-center bg-indigo-500 flex-col"
            >
              <p>create</p>
              <p>deal</p>
            </Link>
          </div>
          <div className="bg-pink-300 w-full h-32 rounded-lg border"></div>
        </section>
      </div>
    </div>
  );
}
