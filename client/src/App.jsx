import Contact from "./components/Contact";
import Message from "./components/Message";
import { createSignal } from "solid-js";
import Setting from "./components/Setting";

export default function App() {
  const [settingsOpen, setSettingsOpen] = createSignal(false);
  const handleClick = (event) => {
      console.log("Div clicked!", event.currentTarget);
    };
  return (
        <div class="grid grid-cols-1 lg:grid-cols-[25%_75%] overflow-x-hidden">
          <div class="bg-(--background) text-white h-screen" classList={{
              "hidden": settingsOpen(),
            }}>
            <div class="flex justify-between m-3 mb-0 items-center">
              <h1 class="text-2xl font-bold">Golubac</h1>
              <button onClick={() => setSettingsOpen(true)} class="h-[3em] w-[3em] rounded-md hover:bg-(--primary-hover) flex justify-center items-center cursor-pointer"><svg class="w-[2.5em] h-[2.5em]" viewBox="-2.4 -2.4 28.80 28.80" fill="none" xmlns="http://www.w3.org/2000/svg" transform="rotate(0)"><g id="SVGRepo_bgCarrier" stroke-width="0"></g><g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round" stroke="#CCCCCC" stroke-width="0.096"></g><g id="SVGRepo_iconCarrier"> <path d="M9.65202 4.56614C9.85537 3.65106 10.667 3 11.6044 3H12.3957C13.3331 3 14.1447 3.65106 14.3481 4.56614L14.551 5.47935C15.2121 5.73819 15.8243 6.09467 16.3697 6.53105L17.2639 6.24961C18.1581 5.96818 19.1277 6.34554 19.5964 7.15735L19.9921 7.84264C20.4608 8.65445 20.3028 9.68287 19.612 10.3165L18.9218 10.9496C18.9733 11.2922 19.0001 11.643 19.0001 12C19.0001 12.357 18.9733 12.7078 18.9218 13.0504L19.612 13.6835C20.3028 14.3171 20.4608 15.3455 19.9921 16.1574L19.5965 16.8426C19.1278 17.6545 18.1581 18.0318 17.2639 17.7504L16.3698 17.4689C15.8243 17.9053 15.2121 18.2618 14.551 18.5206L14.3481 19.4339C14.1447 20.3489 13.3331 21 12.3957 21H11.6044C10.667 21 9.85537 20.3489 9.65202 19.4339L9.44909 18.5206C8.78796 18.2618 8.17579 17.9053 7.63034 17.4689L6.73614 17.7504C5.84199 18.0318 4.87234 17.6545 4.40364 16.8426L4.00798 16.1573C3.53928 15.3455 3.69731 14.3171 4.38811 13.6835L5.07833 13.0504C5.02678 12.7077 5.00005 12.357 5.00005 12C5.00005 11.643 5.02678 11.2922 5.07833 10.9496L4.38813 10.3165C3.69732 9.68288 3.5393 8.65446 4.008 7.84265L4.40365 7.15735C4.87235 6.34554 5.842 5.96818 6.73616 6.24962L7.63035 6.53106C8.1758 6.09467 8.78796 5.73819 9.44909 5.47935L9.65202 4.56614Z" stroke="#ffffff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"></path> <path d="M13 12C13 12.5523 12.5523 13 12 13C11.4477 13 11 12.5523 11 12C11 11.4477 11.4477 11 12 11C12.5523 11 13 11.4477 13 12Z" stroke="#ffffff" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round"></path> </g></svg></button>
            </div>
              <div class="flex items-center ml-3">
                <svg class="h-3.5" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><g id="SVGRepo_bgCarrier" stroke-width="0"></g><g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g><g id="SVGRepo_iconCarrier"> <path d="M12 14V16M8 9V6C8 3.79086 9.79086 2 12 2C14.2091 2 16 3.79086 16 6V9M7 21H17C18.1046 21 19 20.1046 19 19V11C19 9.89543 18.1046 9 17 9H7C5.89543 9 5 9.89543 5 11V19C5 20.1046 5.89543 21 7 21Z" stroke="#ffffff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"></path> </g></svg>
                <p class="text-sm ml-0.5">E2E Encrypted</p>
              </div>
            <Contact name="Vuk Lazic" handler={handleClick} />
            <Contact name="Komplet Lepinja" handler={handleClick} />
            <Contact name="David Mitic Mitke" handler={handleClick} />
          </div>
          <div class="bg-(--background) text-white h-screen" classList={{
              "hidden": !settingsOpen(),
            }}>
            <div class="flex justify-between m-3 mb-0 items-center">
              <h1 class="text-2xl font-bold">Settings</h1>
              <button onClick={() => setSettingsOpen(false)} class="h-[3em] w-[3em] rounded-md hover:bg-(--primary-hover) flex justify-center items-center cursor-pointer"><svg class="w-[2.5em] h-[2.5em]" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><g id="SVGRepo_bgCarrier" stroke-width="0"></g><g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g><g id="SVGRepo_iconCarrier"> <path d="M6 6L18 18M18 6L6 18" stroke="#ffffff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"></path> </g></svg></button>
            </div>
            <Setting name="User settings" />
          </div>
          <div class="bg-(--background) text-white h-screen border-l-2 flex flex-col">
            <div class="overflow-scroll overflow-x-hidden scrollbar-thin scrollbar-thumb-(--primary-hover) scrollbar-track-(--primary)">
              <Message rounded={false} yours={true} />
              <Message rounded={true} yours={true} />
              <Message rounded={false} yours={false} />
              <Message rounded={false} yours={true} />
              <Message rounded={true} yours={true} />
              <Message rounded={false} yours={false} />
            </div>
            <div class="mt-auto p-5 max-h-20">
              <form class="flex items-center gap-2">
                <input type="text" placeholder="Write a message" class="h-10 rounded-md w-full text-md bg-(--primary-hover) p-1 caret-fuchsia-300 focus:outline-none focus:ring-0" />
                <button type="submit" class="h-10 w-10 cursor-pointer inline-flex items-center justify-center rounded-md hover:bg-(--primary-hover)">
                <svg class="h-full max-h-8" viewBox="0 0 24.00 24.00" fill="none" xmlns="http://www.w3.org/2000/svg" stroke="#ffffff"><g id="SVGRepo_bgCarrier" stroke-width="0"></g><g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g><g id="SVGRepo_iconCarrier"> <path d="M10 14L14 21L21 3M10 14L3 10L21 3M10 14L21 3" stroke="#ffffff" stroke-width="1.2" stroke-linejoin="round"></path> </g></svg>
                </button>
                </form>
            </div>
          </div>
      </div>
  );
}
