export default function Layout({ children }: { children: React.ReactNode }) {
  return (
    <>
      <div className="bg-slate-700 z-10 w-full h-screen md:w-5/6 md:rounded-r-2xl">
        <div className="flex p-4 gap-2 bg-white flex-col md:flex-row md:overflow-hidden">
          <div className="border p-2">
            <p>SideNav</p>
          </div>
          <div className="border flex-grow md:overflow-y-auto">{children}</div>
        </div>
      </div>
    </>
  );
}
