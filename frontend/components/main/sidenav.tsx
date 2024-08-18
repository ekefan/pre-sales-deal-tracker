import NavLinks from "./navLinks";

export default function SideNav() {
  return (
    <>
      <div className="bg-sky-100 h-14 md:h-16 rounded flex items-end justify-start pt-3 pl-3 pb-1 font-medium text-sky-600 text-base md:text-lg">Vas Deal Tracker</div>
      <div className="bg-slate-50 h-14 grow flex space-x-2 md:space-x-0 md:space-y-2 justify-between items-center md:flex-col p-2 rounded overflow-y-auto"><NavLinks user=""/></div>
    </>
  );
}
