import { MoveRight, MoveLeft } from "lucide-react";
export default async function Page() {
  const newDate = new Date(Date.now());
  const dateString = newDate.toDateString();

  return (
    <div className="flex flex-col gap-4 pb-2 px-2 h-full w-full relative">
      <div className="flex flex-col p-2 bg-slate-50 gap-4">
        <div>
          <p>Welcome {`<first name>`}</p>
          <p>{dateString}</p>
        </div>
        <p>Ongoing Deals</p>
      </div>
      <div className="flex flex-col gap-3  p-2 w-full grow h-auto">
        <section className="flex flex-col w-full h-full gap-3 ">
          <div className="bg-green-200 w-full h-32 rounded-lg border"></div>
          <div className="bg-yellow-200 w-full h-32 rounded-lg border"></div>
          <div className="bg-indigo-300 w-full h-32 rounded-lg border"></div>
          <div className="bg-pink-300 w-full h-32 rounded-lg border"></div>
          <div className="bg-green-200 w-full h-32 rounded-lg border"></div>
          <div className="bg-yellow-200 w-full h-32 rounded-lg border"></div>
          <div className="bg-indigo-300 w-full h-32 rounded-lg border"></div>
          <div className="bg-pink-300 w-full h-32 rounded-lg border"></div>
        </section>
      </div>
    </div>
  );
}

/*
main
  div
    div
    div
      section
        div
          div
*/

/*
 <div className="flex gap-5 justify-end items-center mt-5 bg-sky-400 p-2">
            <button className="border rounded">
              <MoveLeft />
            </button>
            <button className="border rounded">
              <MoveRight />
            </button>
          </div>
*/
