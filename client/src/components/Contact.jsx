export default function Contact(props) {
  return (
      <div class="p-3 flex items-center rounded-2xl cursor-pointer hover:bg-(--primary-hover) mx-1" onClick={props.handler}>
        <img src="assets/pfp.jpg" class="rounded-full" width="40" />
        <div class="ml-2">
          <p class="font-semibold">{props.name}</p>
          <p class="text-sm">Start the chat</p>
        </div>
      </div>
  )
}
