import Link from "next/link";

export default function LoginForm() {
  return (
    <form action="">
    <div className="p-3 border rounded-xl flex flex-col gap-4">
      <div className="border rounded p-1 flex gap-2">
        <label htmlFor="email">email:</label>
        <input
          type="text"
          placeholder="example@gmail.com"
          className="outline-none"
        />
      </div>
      <div>
        <div className="border rounded p-1 flex gap-2">
          <label htmlFor="password">password:</label>
          <input
            type="password"
            className="outline-none"
            name=""
            id=""
            aria-describedby="helpId"
            placeholder="password"
          />
        </div>
      </div>
      <div className="w-full flex justify-center items-center">
        <Link href="/dashboard" className="border flex items-center justify-center rounded-md w-1/4 ">login</Link>
      </div>
    </div>
  </form>
  );
}
