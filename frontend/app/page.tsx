import Image from "next/image";
import LoginForm from "@/components/main/loginform";
export default function Home() {
  return (
    <main className="flex min-h-screen flex-col items-center p-24 gap-4">
      <p>Login</p>
      <LoginForm/>
    </main>
  );
}
