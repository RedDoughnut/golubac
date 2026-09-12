import Contact from "./components/Contact";

export default function App() {
  return (
        <div class="grid grid-cols-1 md:grid-cols-[30%_70%]">
        <div class="bg-(--background) text-white h-screen">
            <h1 class="p-2 text-xl">Golubac</h1>
            <p class="p-2 text-sm">E2E Encrypted</p>
            <Contact name="Vuk Lazic"/>
            <Contact name="Komplet Lepinja"/>
            <Contact name="David Mitic Mitke"/>
          </div>

          <div class="bg-(--background) text-white h-screen border-l-2">
            Your messages will appear here!
          </div>
        </div>
  );
}
