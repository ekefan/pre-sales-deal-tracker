import NavLinks from "./navLinks";

export default function SideNav() {
  return (
    <>
      <div className="bg-slate-50 h-14 md:h-16 rounded flex items-end p-2">VDT</div>
      <div className="bg-slate-50 h-8 grow flex space-x-2 md:space-x-0 md:space-y-2 justify-between items-center md:flex-col p-2 rounded"><NavLinks user=""/></div>
    </>
  );
}
