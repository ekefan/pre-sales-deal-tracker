import Link from "next/link";
export function UserCard({
  userId,
  username,
  fullname,
  email,
}: {
  userId: number;
  username: string;
  fullname: string;
  email: string;
}) {

  const handleUpdateClick = () => {
    console.log("waiting to be set", username, userId)
  };
  return (
    <div className="w-full h-32 rounded-xl bg-slate-100 flex gap-2 p-1 text-sm md:text-base">
      <div className="flex flex-col gap-3 w-full justify-center pl-2">
        <p>Full name: {fullname}</p>
        <p>Username: {username}</p>
        <p>Email: {email}</p>
      </div>
      <div className="flex  text-sm gap-2 justify-end items-center w-1/6">
        <button className="bg-slate-200 h-full w-full p-1 rounded-xl ">
          <p>reset</p>
          <p>password</p>
        </button>
        <Link
          href={`/dashboard/users/update/?user_id=${userId}`}
          onClick={handleUpdateClick}
          className="p-2 rounded-xl h-full w-full flex items-center border border-slate-300 justify-center bg-slate-200"
        >
          <p>update</p>
        </Link>
      </div>
    </div>
  );
}
